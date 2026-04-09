package internal

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.skia.org/infra/go/metrics2"
	"go.skia.org/infra/go/skerr"
	ag_pb "go.skia.org/infra/perf/go/anomalygroup/proto/v1"
	c_pb "go.skia.org/infra/perf/go/culprit/proto/v1"

	"go.skia.org/infra/perf/go/types"
	"go.skia.org/infra/perf/go/workflows"

	// TODO(b/500974820): Replace `legacyPinpoint` with `pinpoint`.
	legacyPinpoint "go.skia.org/infra/pinpoint/go/pinpoint"

	// TODO(b/500974820): Remove the new `pinpoint` backend.
	pinpoint "go.skia.org/infra/pinpoint/go/workflows"
	pp_pb "go.skia.org/infra/pinpoint/proto/v1"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

const (
	_WAIT_TIME_FOR_ANOMALIES = 30 * time.Minute
)

func agsaToken() *AnomalyGroupServiceActivity {
	return &AnomalyGroupServiceActivity{}
}

func gsaToken() *GerritServiceActivity {
	return &GerritServiceActivity{}
}

func csaToken() *CulpritServiceActivity {
	return &CulpritServiceActivity{}
}

// MaybeTriggerBisectionWorkflow is the entry point for the workflow which handles anomaly group
// processing. It is responsible for triggering a bisection if the anomalygroup's
// group action = BISECT. If group action = REPORT, files a bug notifying user of the anomalies.
func MaybeTriggerBisectionWorkflow(
	ctx workflow.Context,
	input *workflows.MaybeTriggerBisectionParam,
) (*workflows.MaybeTriggerBisectionResult, error) {
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	ctx = workflow.WithActivityOptions(ctx, regularActivityOptions)

	if err := waitForAnomalyClusteringWindow(ctx); err != nil {
		return nil, skerr.Wrap(err)
	}

	anomalyGroupResponse, err := loadAnomalyGroupByID(
		ctx,
		input.AnomalyGroupServiceUrl,
		input.AnomalyGroupId,
	)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	workflow.GetLogger(ctx).Info(
		"MaybeTriggerBisectionWorkflow",
		"WorkflowID",
		workflow.GetInfo(ctx).WorkflowExecution.ID,
		"AnomalyGroup",
		input.AnomalyGroupId,
		"GroupAction",
		anomalyGroupResponse.AnomalyGroup.GroupAction,
	)

	switch anomalyGroupResponse.AnomalyGroup.GroupAction {
	case ag_pb.GroupActionType_BISECT:
		bisectionAllowed, err := isBisectionAllowed(ctx)
		if err != nil {
			return nil, skerr.Wrap(err)
		}
		if bisectionAllowed {
			return processAnomaliesAsBisection(ctx, input)
		} else {
			// Fallback to reporting if the rate limiter prevents creating bisect jobs.
			return processAnomaliesAsReporting(ctx, input)
		}
	case ag_pb.GroupActionType_REPORT:
		return processAnomaliesAsReporting(ctx, input)
	case ag_pb.GroupActionType_NOACTION:
		metrics2.GetCounter("anomalygroup_ignored").Inc(1)
		return nil, nil
	default:
		return nil, skerr.Fmt(
			"Unhandled GroupAction type %s",
			anomalyGroupResponse.AnomalyGroup.GroupAction,
		)
	}
}

func processAnomaliesAsBisection(
	ctx workflow.Context,
	input *workflows.MaybeTriggerBisectionParam,
) (*workflows.MaybeTriggerBisectionResult, error) {
	anomaliesCount := 1
	topAnomaliesResponse, err := findTopAnomalies(
		ctx,
		input.AnomalyGroupServiceUrl,
		input.AnomalyGroupId,
		anomaliesCount,
	)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	topAnomaly := topAnomaliesResponse.Anomalies[0]
	startHash, endHash, err := getCommitHashes(
		ctx,
		topAnomaly.StartCommit,
		topAnomaly.EndCommit,
	)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	jobId, err := createBisectJob(
		ctx,
		input,
		topAnomaly,
		startHash,
		endHash,
	)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	workflow.GetLogger(ctx).Info("Pinpoint Job created", "jobId", jobId)

	// Update the anomaly group with the bisection id.
	updateRequest := ag_pb.UpdateAnomalyGroupRequest{
		AnomalyGroupId: input.AnomalyGroupId,
		BisectionId:    jobId,
	}
	if err = updateAnomalyGroup(ctx, input.AnomalyGroupServiceUrl, &updateRequest); err != nil {
		return nil, skerr.Wrap(err)
	}
	metrics2.GetCounter("anomalygroup_bisected").Inc(1)
	return &workflows.MaybeTriggerBisectionResult{
		JobId: jobId,
	}, nil
}

func processAnomaliesAsReporting(
	ctx workflow.Context,
	input *workflows.MaybeTriggerBisectionParam,
) (*workflows.MaybeTriggerBisectionResult, error) {
	// Load Anomalies data
	anomaliesCount := 10
	topAnomaliesResponse, err := findTopAnomalies(
		ctx,
		input.AnomalyGroupServiceUrl,
		input.AnomalyGroupId,
		anomaliesCount,
	)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	topAnomalies := convertToCulpritAnomalies(topAnomaliesResponse.Anomalies)

	notifyUserOfAnomalyResponse, err := notifyUserOfAnomalies(
		ctx,
		topAnomalies,
		input.CulpritServiceUrl,
		input.AnomalyGroupId,
	)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	// Update the anomaly group with the reported issue id.
	if notifyUserOfAnomalyResponse != nil && notifyUserOfAnomalyResponse.IssueId != "" {
		updateRequest := ag_pb.UpdateAnomalyGroupRequest{
			AnomalyGroupId: input.AnomalyGroupId,
			IssueId:        notifyUserOfAnomalyResponse.IssueId,
		}
		if err = updateAnomalyGroup(ctx, input.AnomalyGroupServiceUrl, &updateRequest); err != nil {
			return nil, skerr.Wrap(err)
		}
	}

	metrics2.GetCounter("anomalygroup_reported").Inc(1)
	return &workflows.MaybeTriggerBisectionResult{}, nil
}

// Mimic the story name update in the legacy descriptor logic.
// The original source in catapult/dashboard/dashboard/common/descriptor.py
func benchmarkStoriesNeedUpdate(b string) bool {
	system_health_benchmark_prefix := "system_health"
	legacy_complex_cases_benchmarks := []string{
		"tab_switching.typical_25",
		"v8.browsing_desktop",
		"v8.browsing_desktop-future",
		"v8.browsing_mobile",
		"v8.browsing_mobile-future",
		"heap_profiling.mobile.disabled",
	}
	if strings.HasPrefix(b, system_health_benchmark_prefix) {
		return true
	}
	for _, benchmark := range legacy_complex_cases_benchmarks {
		if benchmark == b {
			return true
		}
	}
	return false
}

func updateStoryDescriptorName(s string) string {
	return strings.Replace(s, "_", ":", -1)
}

func parseStatisticNameFromChart(chart_name string) (string, string) {
	parts := strings.Split(chart_name, "_")
	part_count := len(parts)
	if part_count < 1 {
		return chart_name, ""
	}
	for _, stat := range types.AllMeasurementStats {
		if parts[part_count-1] == stat {
			return strings.Join(parts[:part_count-1], "_"), parts[part_count-1]
		}
	}
	return chart_name, ""
}

// waitForAnomalyClusteringWindow waits for some time so that more anomalies can
// be detected and grouped.
func waitForAnomalyClusteringWindow(ctx workflow.Context) error {
	if err := workflow.Sleep(ctx, _WAIT_TIME_FOR_ANOMALIES); err != nil {
		return skerr.Wrap(err)
	}
	return nil
}

func loadAnomalyGroupByID(
	ctx workflow.Context,
	url string,
	anomalyGroupID string,
) (*ag_pb.LoadAnomalyGroupByIDResponse, error) {
	var anomalyGroupResponse *ag_pb.LoadAnomalyGroupByIDResponse
	err := workflow.ExecuteActivity(ctx, agsaToken().LoadAnomalyGroupByID, url,
		&ag_pb.LoadAnomalyGroupByIDRequest{
			AnomalyGroupId: anomalyGroupID,
		}).
		Get(ctx, &anomalyGroupResponse)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	return anomalyGroupResponse, nil
}

func isBisectionAllowed(ctx workflow.Context) (bool, error) {
	var bisectionAllowed bool
	err := workflow.ExecuteActivity(ctx, agsaToken().CheckBisectionAllowed).
		Get(ctx, &bisectionAllowed)
	if err != nil {
		return false, skerr.Wrap(err)
	}
	workflow.GetLogger(ctx).Info(
		"MaybeTriggerBisectionWorkflow",
		"Bisection allowed",
		bisectionAllowed,
	)
	return bisectionAllowed, nil
}

func findTopAnomalies(
	ctx workflow.Context,
	url string,
	anomalyGroupID string,
	limit int,
) (*ag_pb.FindTopAnomaliesResponse, error) {
	var topAnomaliesResponse *ag_pb.FindTopAnomaliesResponse
	if err := workflow.ExecuteActivity(ctx, agsaToken().FindTopAnomalies, url, &ag_pb.FindTopAnomaliesRequest{
		AnomalyGroupId: anomalyGroupID,
		Limit:          int64(limit),
	}).Get(ctx, &topAnomaliesResponse); err != nil {
		return nil, skerr.Wrap(err)
	}
	if topAnomaliesResponse == nil || len(topAnomaliesResponse.Anomalies) == 0 {
		return nil, skerr.Fmt("No anomalies found for anomalygroup %s", anomalyGroupID)
	}
	return topAnomaliesResponse, nil
}

// Currently the protos in culprit service and anomaly service are having two identical
// copies of definition on Anomaly. We should merge them into one.
func convertToCulpritAnomalies(anomalies []*ag_pb.Anomaly) []*c_pb.Anomaly {
	result := make([]*c_pb.Anomaly, len(anomalies))
	for i, anomaly := range anomalies {
		result[i] = &c_pb.Anomaly{
			StartCommit:          anomaly.StartCommit,
			EndCommit:            anomaly.EndCommit,
			Paramset:             anomaly.Paramset,
			ImprovementDirection: anomaly.ImprovementDirection,
			MedianBefore:         anomaly.MedianBefore,
			MedianAfter:          anomaly.MedianAfter,
		}
	}
	return result
}

// getCommitHashes converts start and end commit postions to commit hash.
func getCommitHashes(
	ctx workflow.Context,
	startCommit int64,
	endCommit int64,
) (string, string, error) {
	var startHash, endHash string
	if err := workflow.ExecuteActivity(ctx, gsaToken().GetCommitRevision, startCommit).
		Get(ctx, &startHash); err != nil {
		return "", "", skerr.Wrap(err)
	}
	if err := workflow.ExecuteActivity(ctx, gsaToken().GetCommitRevision, endCommit).
		Get(ctx, &endHash); err != nil {
		return "", "", skerr.Wrap(err)
	}
	return startHash, endHash, nil
}

func createBisectJob(
	ctx workflow.Context,
	input *workflows.MaybeTriggerBisectionParam,
	anomaly *ag_pb.Anomaly,
	startHash, endHash string,
) (string, error) {
	var isLegacyPinpointEnabled bool
	if err := workflow.ExecuteActivity(ctx, agsaToken().ShouldUseLegacyPinpoint).
		Get(ctx, &isLegacyPinpointEnabled); err != nil {
		return "", skerr.Wrap(err)
	}
	if isLegacyPinpointEnabled {
		return createLegacyBisectJob(ctx, anomaly, startHash, endHash)
	}
	jobId := uuid.New().String()
	// Childworkflow options includes:
	//   WorkflowID: 		The UUID to be used as the Pinpoint job id. We pre-assigne it
	//				 		here to avoid extra calls to get it from the spawned workflow.
	//	 TaskQueue:  		Assign the child workflow to the correct task queue. If this is
	//				 		empty, it will be assigned to the current grouping queue.
	//   ParentClosePolicy: Using _ABANDON option to ensure the child workflow will
	//    					continue even if the parent workflow exits.
	options := workflow.ChildWorkflowOptions{
		WorkflowID:        jobId,
		TaskQueue:         input.PinpointTaskQueue,
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}
	story, chart, stat := parseStoryChartStat(anomaly)
	wf := workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(ctx, options),
		pinpoint.CulpritFinderWorkflow,
		&pinpoint.CulpritFinderParams{
			Request: &pp_pb.ScheduleCulpritFinderRequest{
				StartGitHash:         startHash,
				EndGitHash:           endHash,
				Configuration:        anomaly.Paramset["bot"],
				Benchmark:            anomaly.Paramset["benchmark"],
				Story:                story,
				Chart:                chart,
				Statistic:            stat,
				ImprovementDirection: anomaly.ImprovementDirection,
			},
			CallbackParams: &pp_pb.CulpritProcessingCallbackParams{
				AnomalyGroupId:        input.AnomalyGroupId,
				CulpritServiceUrl:     input.CulpritServiceUrl,
				TemporalTaskQueueName: input.GroupingTaskQueue,
			},
		},
	)
	// This Get() call will wait for the child workflow to start.
	if err := wf.GetChildWorkflowExecution().Get(ctx, nil); err != nil {
		return "", skerr.Wrapf(err, "Child workflow failed to start.")
	}
	return jobId, nil
}

func createLegacyBisectJob(ctx workflow.Context,
	anomaly *ag_pb.Anomaly,
	startHash, endHash string,
) (string, error) {
	story, chart, stat := parseStoryChartStat(anomaly)
	req := legacyPinpoint.BisectJobCreateRequest{
		ComparisonMode: "performance",
		StartGitHash:   startHash,
		EndGitHash:     endHash,
		Configuration:  anomaly.Paramset["bot"],
		Benchmark:      anomaly.Paramset["benchmark"],
		Story:          story,
		Chart:          chart,
		Statistic:      stat,
		TestPath:       anomaly.Paramset["test_path"],
	}
	var resp *legacyPinpoint.CreatePinpointResponse
	err := workflow.ExecuteActivity(ctx, agsaToken().CreateLegacyBisectJob, &req).Get(ctx, &resp)
	if err != nil {
		return "", skerr.Wrap(err)
	}
	if resp.JobID == "" {
		return "", skerr.Wrap(errors.New("Chromeperf failed to create a new job"))
	}
	return resp.JobID, nil
}

func updateAnomalyGroup(
	ctx workflow.Context,
	url string,
	req *ag_pb.UpdateAnomalyGroupRequest,
) error {
	var updateAnomalyGroupResponse *ag_pb.UpdateAnomalyGroupResponse
	future := workflow.ExecuteActivity(ctx, agsaToken().UpdateAnomalyGroup, url, req)
	if err := future.Get(ctx, &updateAnomalyGroupResponse); err != nil {
		return skerr.Wrap(err)
	}
	return nil
}

func notifyUserOfAnomalies(
	ctx workflow.Context,
	anomalies []*c_pb.Anomaly,
	culpritServiceUrl, anomalyGroupId string,
) (*c_pb.NotifyUserOfAnomalyResponse, error) {
	var notifyUserOfAnomalyResponse *c_pb.NotifyUserOfAnomalyResponse
	request := c_pb.NotifyUserOfAnomalyRequest{
		AnomalyGroupId: anomalyGroupId,
		Anomaly:        anomalies,
	}
	future := workflow.ExecuteActivity(
		ctx,
		csaToken().NotifyUserOfAnomaly,
		culpritServiceUrl,
		&request,
	)
	if err := future.Get(ctx, &notifyUserOfAnomalyResponse); err != nil {
		return nil, skerr.Wrap(err)
	}
	return notifyUserOfAnomalyResponse, nil
}

func parseStoryChartStat(anomaly *ag_pb.Anomaly) (string, string, string) {
	chart, stat := parseStatisticNameFromChart(anomaly.Paramset["measurement"])

	story := anomaly.Paramset["story"]
	if benchmarkStoriesNeedUpdate(anomaly.Paramset["benchmark"]) {
		story = updateStoryDescriptorName(story)
	}
	return story, chart, stat
}
