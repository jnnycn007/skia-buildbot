package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.skia.org/infra/go/metrics2"
	"go.skia.org/infra/go/skerr"
	ag_pb "go.skia.org/infra/perf/go/anomalygroup/proto/v1"
	c_pb "go.skia.org/infra/perf/go/culprit/proto/v1"
	legacyPinpoint "go.skia.org/infra/perf/go/pinpoint"
	"go.skia.org/infra/perf/go/types"
	"go.skia.org/infra/perf/go/workflows"
	pinpoint "go.skia.org/infra/pinpoint/go/workflows"
	pp_pb "go.skia.org/infra/pinpoint/proto/v1"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

const (
	_WAIT_TIME_FOR_ANOMALIES = 30 * time.Minute
)

// MaybeTriggerBisectionWorkflow is the entry point for the workflow which handles anomaly group
// processing. It is responsible for triggering a bisection if the anomalygroup's
// group action = BISECT. If group action = REPORT, files a bug notifying user of the anomalies.
func MaybeTriggerBisectionWorkflow(
	ctx workflow.Context,
	input *workflows.MaybeTriggerBisectionParam,
) (*workflows.MaybeTriggerBisectionResult, error) {
	var agsa AnomalyGroupServiceActivity

	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	ctx = workflow.WithActivityOptions(ctx, regularActivityOptions)

	if err := waitForAnomalyClusteringWindow(ctx); err != nil {
		return nil, skerr.Wrap(err)
	}

	anomalyGroupResponse, err := loadAnomalyGroupByID(
		ctx,
		agsa,
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
		bisectionAllowed, err := isBisectionAllowed(ctx, agsa)
		if err != nil {
			return nil, skerr.Wrap(err)
		}
		if bisectionAllowed {
			return processAnomaliesAsBisection(ctx, agsa, input)
		} else {
			// Fallback to reporting if the rate limiter prevents creating bisect jobs.
			return processAnomaliesAsReporting(ctx, agsa, input)
		}
	case ag_pb.GroupActionType_REPORT:
		return processAnomaliesAsReporting(ctx, agsa, input)
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
	agsa AnomalyGroupServiceActivity,
	input *workflows.MaybeTriggerBisectionParam,
) (*workflows.MaybeTriggerBisectionResult, error) {
	var gsa GerritServiceActivity
	anomaliesCount := 1
	topAnomaliesResponse, err := findTopAnomalies(
		ctx,
		agsa,
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
		gsa,
		topAnomaly.StartCommit,
		topAnomaly.EndCommit,
	)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	jobId, err := createBisectJob(
		ctx,
		agsa,
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
	if err = updateAnomalyGroup(ctx, agsa, input.AnomalyGroupServiceUrl, &updateRequest); err != nil {
		return nil, skerr.Wrap(err)
	}
	metrics2.GetCounter("anomalygroup_bisected").Inc(1)
	return &workflows.MaybeTriggerBisectionResult{
		JobId: jobId,
	}, nil
}

func processAnomaliesAsReporting(
	ctx workflow.Context,
	agsa AnomalyGroupServiceActivity,
	input *workflows.MaybeTriggerBisectionParam,
) (*workflows.MaybeTriggerBisectionResult, error) {
	var csa CulpritServiceActivity
	// Load Anomalies data
	anomaliesCount := 10
	topAnomaliesResponse, err := findTopAnomalies(
		ctx,
		agsa,
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
		csa,
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
		if err = updateAnomalyGroup(ctx, agsa, input.AnomalyGroupServiceUrl, &updateRequest); err != nil {
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
	agsa AnomalyGroupServiceActivity,
	url string,
	anomalyGroupID string,
) (*ag_pb.LoadAnomalyGroupByIDResponse, error) {
	var anomalyGroupResponse *ag_pb.LoadAnomalyGroupByIDResponse
	err := workflow.ExecuteActivity(ctx, agsa.LoadAnomalyGroupByID, url,
		&ag_pb.LoadAnomalyGroupByIDRequest{
			AnomalyGroupId: anomalyGroupID,
		}).
		Get(ctx, &anomalyGroupResponse)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	return anomalyGroupResponse, nil
}

func isBisectionAllowed(ctx workflow.Context, agsa AnomalyGroupServiceActivity) (bool, error) {
	var bisectionAllowed bool
	err := workflow.ExecuteActivity(ctx, agsa.CheckBisectionAllowed).Get(ctx, &bisectionAllowed)
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
	agsa AnomalyGroupServiceActivity,
	url string,
	anomalyGroupID string,
	limit int,
) (*ag_pb.FindTopAnomaliesResponse, error) {
	var topAnomaliesResponse *ag_pb.FindTopAnomaliesResponse
	if err := workflow.ExecuteActivity(ctx, agsa.FindTopAnomalies, url, &ag_pb.FindTopAnomaliesRequest{
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
	gsa GerritServiceActivity,
	startCommit int64,
	endCommit int64,
) (string, string, error) {
	var startHash, endHash string
	if err := workflow.ExecuteActivity(ctx, gsa.GetCommitRevision, startCommit).Get(ctx, &startHash); err != nil {
		return "", "", skerr.Wrap(err)
	}
	if err := workflow.ExecuteActivity(ctx, gsa.GetCommitRevision, endCommit).Get(ctx, &endHash); err != nil {
		return "", "", skerr.Wrap(err)
	}
	return startHash, endHash, nil
}

func createBisectJob(
	ctx workflow.Context,
	agsa AnomalyGroupServiceActivity,
	input *workflows.MaybeTriggerBisectionParam,
	anomaly *ag_pb.Anomaly,
	startHash, endHash string,
) (string, error) {
	var isLegacyPinpointEnabled bool
	if err := workflow.ExecuteActivity(ctx, agsa.ShouldUseLegacyPinpoint).
		Get(ctx, &isLegacyPinpointEnabled); err != nil {
		return "", skerr.Wrap(err)
	}
	if isLegacyPinpointEnabled {
		return createLegacyBisectJob(ctx, agsa, anomaly, startHash, endHash)
	}
	child_wf_id := uuid.New().String()
	// Childworkflow options includes:
	//   WorkflowID: 		The UUID to be used as the Pinpoint job id. We pre-assigne it
	//				 		here to avoid extra calls to get it from the spawned workflow.
	//	 TaskQueue:  		Assign the cihld workflow to the correct task queue. If this is
	//				 		empty, it will be assigned to the current grouping queue.
	//   ParentClosePolicy: Using _ABANDON option to ensure the child workflow will
	//    					continue even if the parent workflow exits.
	child_wf_options := workflow.ChildWorkflowOptions{
		WorkflowID:        child_wf_id,
		TaskQueue:         input.PinpointTaskQueue,
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}
	c_ctx := workflow.WithChildOptions(ctx, child_wf_options)

	story, chart, stat := parseStoryChartStat(anomaly)
	find_culprit_wf := workflow.ExecuteChildWorkflow(c_ctx, pinpoint.CulpritFinderWorkflow,
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
		})
	// This Get() call will wait for the child workflow to start.
	if err := find_culprit_wf.GetChildWorkflowExecution().Get(ctx, nil); err != nil {
		return "", skerr.Wrapf(err, "Child workflow failed to start.")
	}
	return child_wf_id, nil
}

func createLegacyBisectJob(ctx workflow.Context,
	agsa AnomalyGroupServiceActivity,
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
		// TODO(b/495782839): Remove this workaround by providing a test path as a
		// paramset parameter.
		TestPath: getAnomalyTestPath(anomaly),
	}
	var resp *legacyPinpoint.CreatePinpointResponse
	err := workflow.ExecuteActivity(ctx, agsa.CreateLegacyBisectJob, &req).Get(ctx, &resp)
	if err != nil {
		return "", skerr.Wrap(err)
	}
	return resp.JobID, nil
}

func getAnomalyTestPath(anomaly *ag_pb.Anomaly) string {
	if testPath, ok := anomaly.Paramset["test_path"]; ok && testPath != "" {
		return testPath
	}

	bot := anomaly.Paramset["bot"]
	benchmark := anomaly.Paramset["benchmark"]
	measurement := anomaly.Paramset["measurement"]
	story := anomaly.Paramset["story"]
	return fmt.Sprintf("ChromiumPerf/%s/%s/%s/%s", bot, benchmark, measurement, story)
}

func updateAnomalyGroup(
	ctx workflow.Context,
	agsa AnomalyGroupServiceActivity,
	url string,
	req *ag_pb.UpdateAnomalyGroupRequest,
) error {
	var updateAnomalyGroupResponse *ag_pb.UpdateAnomalyGroupResponse
	future := workflow.ExecuteActivity(ctx, agsa.UpdateAnomalyGroup, url, req)
	if err := future.Get(ctx, &updateAnomalyGroupResponse); err != nil {
		return skerr.Wrap(err)
	}
	return nil
}

func notifyUserOfAnomalies(
	ctx workflow.Context,
	csa CulpritServiceActivity,
	anomalies []*c_pb.Anomaly,
	culpritServiceUrl, anomalyGroupId string,
) (*c_pb.NotifyUserOfAnomalyResponse, error) {
	var notifyUserOfAnomalyResponse *c_pb.NotifyUserOfAnomalyResponse
	request := c_pb.NotifyUserOfAnomalyRequest{
		AnomalyGroupId: anomalyGroupId,
		Anomaly:        anomalies,
	}
	future := workflow.ExecuteActivity(ctx, csa.NotifyUserOfAnomaly, culpritServiceUrl, &request)
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
