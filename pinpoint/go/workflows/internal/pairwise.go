package internal

import (
	"time"

	"github.com/google/uuid"
	swarming_pb "go.chromium.org/luci/swarming/proto/api_v2"
	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/pinpoint/go/common"
	"go.skia.org/infra/pinpoint/go/compare"
	"go.skia.org/infra/pinpoint/go/workflows"
	"go.temporal.io/sdk/workflow"

	pinpoint_proto "go.skia.org/infra/pinpoint/proto/v1"
)

func PairwiseWorkflow(ctx workflow.Context, p *workflows.PairwiseParams) (*pinpoint_proto.PairwiseExecution, error) {
	if p.Request.StartBuild == nil && p.Request.StartCommit == nil {
		return nil, skerr.Fmt("Base build and commit are empty.")
	}
	if p.Request.EndBuild == nil && p.Request.EndCommit == nil {
		return nil, skerr.Fmt("Experiment build and commit are empty.")
	}

	leftCas, err := convertCas(p.Request.StartBuild)
	if err != nil {
		return nil, skerr.Wrapf(err, "start build is invalid")
	}

	rightCas, err := convertCas(p.Request.EndBuild)
	if err != nil {
		return nil, skerr.Wrapf(err, "end build is invalid")
	}

	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	ctx = workflow.WithActivityOptions(ctx, regularActivityOptions)

	jobID := uuid.New().String()
	wkStartTime := time.Now().UnixNano()

	// Benchmark runs can sometimes generate an inconsistent number of data points.
	// So even if all benchmark runs were successful, the number of data values
	// generated by commit A vs B will be inconsistent. This can fail the balancing
	// requirement of the statistical test. It can pair up data incorrectly i.e.
	// commit A: [1], [2, 3]
	// commit B: [4, 5], [6]
	// So the analysis will pair up [2,5] together, which are from different runs,
	// violating pairwise analysis
	if p.Request.AggregationMethod == "" {
		p.Request.AggregationMethod = "mean"
	}

	pairwiseRunnerParams := PairwiseCommitsRunnerParams{
		SingleCommitRunnerParams: SingleCommitRunnerParams{
			PinpointJobID:     jobID,
			BotConfig:         p.Request.Configuration,
			Benchmark:         p.Request.Benchmark,
			Chart:             p.Request.Chart,
			Story:             p.Request.Story,
			StoryTags:         p.Request.StoryTags,
			AggregationMethod: p.Request.AggregationMethod,
			Iterations:        p.GetInitialAttempt(),
		},
		LeftCAS:     leftCas,
		RightCAS:    rightCas,
		LeftCommit:  (*common.CombinedCommit)(p.Request.StartCommit),
		RightCommit: (*common.CombinedCommit)(p.Request.EndCommit),
	}

	mh := workflow.GetMetricsHandler(ctx).WithTags(map[string]string{
		"job_id":    jobID,
		"benchmark": p.Request.Benchmark,
		"config":    p.Request.Configuration,
		"story":     p.Request.Story,
	})

	defer func() {
		duration := time.Now().UnixNano() - wkStartTime
		mh.Timer("pairwise_duration").Record(time.Duration(duration))
	}()

	var pr *PairwiseRun
	if err := workflow.ExecuteChildWorkflow(ctx, workflows.PairwiseCommitsRunner, pairwiseRunnerParams).Get(ctx, &pr); err != nil {
		return nil, skerr.Wrap(err)
	}

	results, err := comparePairwiseRuns(ctx, pr, compare.UnknownDir)
	if err != nil {
		return nil, skerr.Wrapf(err, "failed to compare pairwise runs")
	}

	// we only return the "culprit" if this is a culprit-verification run.
	// Culprit verification workflows are run in parallel and returned via
	// channels. So the parent culprit_finder workflow is unable to determine
	// which completed child workflow belongs to which set of inputs. We
	// explicitly return this commit to work around this limitation.
	var culpritCandidate *pinpoint_proto.CombinedCommit
	if p.CulpritVerify {
		culpritCandidate = (*pinpoint_proto.CombinedCommit)(pairwiseRunnerParams.RightCommit)
	}

	protoResults := map[string]*pinpoint_proto.PairwiseExecution_WilcoxonResult{}
	for chart, res := range results {
		protoResults[chart] = &pinpoint_proto.PairwiseExecution_WilcoxonResult{
			// Significant is used in CulpritFinder to determine whether to bisect.
			// Significant = true means that there's indeed a regression and it should
			// be investigated. If significant is not explicitly set to False, we see
			// Temporal workflows with Significant and Culprit omitted because the resolve
			// to nil.
			Significant:              res.Verdict == compare.Different,
			PValue:                   res.PValue,
			ConfidenceIntervalLower:  res.LowerCi,
			ConfidenceIntervalHigher: res.UpperCi,
			ControlMedian:            res.XMedian,
			TreatmentMedian:          res.YMedian,
		}
	}

	return &pinpoint_proto.PairwiseExecution{
		JobId:               jobID,
		CulpritCandidate:    culpritCandidate,
		Results:             protoResults,
		LeftSwarmingStatus:  pr.Left.GetSwarmingStatus(),
		RightSwarmingStatus: pr.Right.GetSwarmingStatus(),
	}, nil
}

func convertCas(pinpointCas *pinpoint_proto.CASReference) (*swarming_pb.CASReference, error) {
	if pinpointCas == nil {
		return nil, nil
	}
	if pinpointCas.Digest == nil {
		return nil, skerr.Fmt("cas digest cannot be nil: %v", pinpointCas)
	}
	return &swarming_pb.CASReference{
		CasInstance: pinpointCas.CasInstance,
		Digest: &swarming_pb.Digest{
			Hash:      pinpointCas.Digest.Hash,
			SizeBytes: pinpointCas.Digest.SizeBytes,
		},
	}, nil
}
