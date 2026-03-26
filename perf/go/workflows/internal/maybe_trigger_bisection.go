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
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	ctx = workflow.WithActivityOptions(ctx, regularActivityOptions)
	logger := workflow.GetLogger(ctx)
	var err error
	var agsa AnomalyGroupServiceActivity
	var gsa GerritServiceActivity
	var csa CulpritServiceActivity

	if err = waitForAnomalyClusteringWindow(ctx); err != nil {
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

	logger.Info(
		"MaybeTriggerBisectionWorkflow",
		"WorkflowID",
		workflow.GetInfo(ctx).WorkflowExecution.ID,
		"AnomalyGroup",
		input.AnomalyGroupId,
		"GroupAction",
		anomalyGroupResponse.AnomalyGroup.GroupAction,
	)

	// Temporary code checking whether the rate limiter works as expected on production.
	bisectionAllowed, err := isBisectionAllowed(ctx, agsa)
	if err != nil {
		logger.Error(fmt.Sprintf("Rate limiter error: %s", skerr.Wrap(err)))
	}
	logger.Info("MaybeTriggerBisectionWorkflow", "Bisection allowed", bisectionAllowed)

	if anomalyGroupResponse.AnomalyGroup.GroupAction == ag_pb.GroupActionType_BISECT {
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

		child_wf_id, err := invokeBisection(ctx, input, topAnomaly, startHash, endHash)
		if err != nil {
			return nil, skerr.Wrap(err)
		}

		// Update the anomaly group with the bisection id.
		updateRequest := ag_pb.UpdateAnomalyGroupRequest{
			AnomalyGroupId: input.AnomalyGroupId,
			BisectionId:    child_wf_id,
		}
		if err = updateAnomalyGroup(ctx, agsa, input.AnomalyGroupServiceUrl, &updateRequest); err != nil {
			return nil, skerr.Wrap(err)
		}
		metrics2.GetCounter("anomalygroup_bisected").Inc(1)
		return &workflows.MaybeTriggerBisectionResult{
			JobId: child_wf_id,
		}, nil
	} else if anomalyGroupResponse.AnomalyGroup.GroupAction == ag_pb.GroupActionType_REPORT {
		// Step 3. Load Anomalies data
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
		topAnomalies := make([]*c_pb.Anomaly, len(topAnomaliesResponse.Anomalies))
		// Currently the protos in culprit service and anomaly service are having two identical
		// copies of definition on Anomaly. We should merge them into one.
		for i, anomaly := range topAnomaliesResponse.Anomalies {
			topAnomalies[i] = &c_pb.Anomaly{
				StartCommit:          anomaly.StartCommit,
				EndCommit:            anomaly.EndCommit,
				Paramset:             anomaly.Paramset,
				ImprovementDirection: anomaly.ImprovementDirection,
				MedianBefore:         anomaly.MedianBefore,
				MedianAfter:          anomaly.MedianAfter,
			}
		}
		// Step 4. Notify the user of the top anomalies
		var notifyUserOfAnomalyResponse *c_pb.NotifyUserOfAnomalyResponse
		if err = workflow.ExecuteActivity(ctx, csa.NotifyUserOfAnomaly, input.CulpritServiceUrl, &c_pb.NotifyUserOfAnomalyRequest{
			AnomalyGroupId: input.AnomalyGroupId,
			Anomaly:        topAnomalies,
		}).Get(ctx, &notifyUserOfAnomalyResponse); err != nil {
			return nil, err
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

	return nil, skerr.Fmt(
		"Unhandled GroupAction type %s",
		anomalyGroupResponse.AnomalyGroup.GroupAction,
	)
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
	if topAnomaliesResponse != nil && len(topAnomaliesResponse.Anomalies) == 0 {
		return nil, skerr.Fmt("No anomalies found for anomalygroup %s", anomalyGroupID)
	}
	return topAnomaliesResponse, nil
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

func invokeBisection(
	ctx workflow.Context,
	input *workflows.MaybeTriggerBisectionParam,
	anomaly *ag_pb.Anomaly,
	startHash, endHash string,
) (string, error) {
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

	chart, stat := parseStatisticNameFromChart(anomaly.Paramset["measurement"])

	benchmark := anomaly.Paramset["benchmark"]
	story := anomaly.Paramset["story"]
	if benchmarkStoriesNeedUpdate(benchmark) {
		story = updateStoryDescriptorName(story)
	}
	find_culprit_wf := workflow.ExecuteChildWorkflow(c_ctx, pinpoint.CulpritFinderWorkflow,
		&pinpoint.CulpritFinderParams{
			Request: &pp_pb.ScheduleCulpritFinderRequest{
				StartGitHash:         startHash,
				EndGitHash:           endHash,
				Configuration:        anomaly.Paramset["bot"],
				Benchmark:            benchmark,
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
