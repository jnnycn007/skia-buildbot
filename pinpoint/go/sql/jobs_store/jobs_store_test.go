package jobstore

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	buildbucketpb "go.chromium.org/luci/buildbucket/proto"
	"go.skia.org/infra/perf/go/sql/sqltest"
	"go.skia.org/infra/pinpoint/go/common"
	"go.skia.org/infra/pinpoint/go/sql/schema/spanner"
	"go.skia.org/infra/pinpoint/go/workflows"
	pinpointpb "go.skia.org/infra/pinpoint/proto/v1"
)

// setupTestDB creates a new Spanner test database and a JobStore instance.
func setupTestDB(t *testing.T) JobStore {
	db := sqltest.NewSpannerDBForTests(t, "jobstore_test_db")
	js := NewJobStore(db)

	// The schema is defined in pinpoint/go/sql/schema/spanner/spanner.go
	_, err := db.Exec(context.Background(), spanner.Schema)
	require.NoError(t, err)

	return js
}

var startCommit = pinpointpb.CombinedCommit{
	Main: &pinpointpb.Commit{
		GitHash: "start_hash",
		Url:     "https://chromium.googlesource.com/chromium/src.git+/start_hash",
		Author:  "chrome@chromium.org",
	},
}

var endCommit = pinpointpb.CombinedCommit{
	Main: &pinpointpb.Commit{
		GitHash: "end_hash",
		Url:     "https://chromium.googlesource.com/chromium/src.git+/end_hash",
		Author:  "chrome@chromium.org",
	},
}

// Helper struct to unmarshal the JSON stored in additional_request_parameters["commit_runs"]
type storedCommitRuns struct {
	Left  *CommitRunData `json:"left"`
	Right *CommitRunData `json:"right"`
}

var leftData = CommitRunData{
	Build: &workflows.Build{
		ID:     123,
		Status: buildbucketpb.Status_SUCCESS,
		BuildParams: workflows.BuildParams{
			Commit: common.NewCombinedCommit(common.NewChromiumCommit("left_build_hash")),
		},
	},
	Runs: []*workflows.TestRun{
		{
			TaskID: "left_task_1",
			Status: "COMPLETED",
			Values: map[string][]float64{"chart_a": {1.0, 2.0}},
		},
	},
}

var rightData = CommitRunData{
	Build: &workflows.Build{
		ID:     456,
		Status: buildbucketpb.Status_SUCCESS, // Success
		BuildParams: workflows.BuildParams{
			Commit: common.NewCombinedCommit(common.NewChromiumCommit("right_build_hash")),
		},
	},
	Runs: []*workflows.TestRun{
		{
			TaskID: "right_task_1",
			Status: "COMPLETED",
			Values: map[string][]float64{"chart_b": {3.0, 4.0}},
		},
	},
}

func TestAddInitialJob(t *testing.T) {
	js := setupTestDB(t)

	ctx := context.Background()
	jobID := uuid.New().String()

	req := &pinpointpb.SchedulePairwiseRequest{
		Benchmark:           "speedometer3",
		Story:               "story1",
		Chart:               "Score",
		Configuration:       "mac-m1_mini_2020-perf",
		InitialAttemptCount: "50",
		StartCommit:         &startCommit,
		EndCommit:           &endCommit,
		StoryTags:           "tag1,tag2",
		AggregationMethod:   "mean",
		Target:              "target_val",
		Project:             "project_val",
		BugId:               "bug_val",
	}

	err := js.AddInitialJob(ctx, req, jobID)
	require.NoError(t, err)

	// Verify the job was added
	retrievedJob, err := js.GetJob(ctx, jobID)
	require.NoError(t, err)
	assert.Equal(t, jobID, retrievedJob.JobID)
	assert.Equal(t, "default", retrievedJob.JobName)
	assert.Equal(t, JobType, retrievedJob.JobType)
	assert.Equal(t, "default", retrievedJob.SubmittedBy)
	assert.Equal(t, req.Benchmark, retrievedJob.Benchmark)
	assert.Equal(t, req.Configuration, retrievedJob.BotName)

	expectedAdditionalParams := map[string]string{
		"start_commit_githash":  req.StartCommit.Main.GitHash,
		"end_commit_githash":    req.EndCommit.Main.GitHash,
		"story":                 req.Story,
		"story_tags":            req.StoryTags,
		"initial_attempt_count": req.InitialAttemptCount,
		"aggregation_method":    req.AggregationMethod,
		"target":                req.Target,
		"project":               req.Project,
		"bug_id":                req.BugId,
		"chart":                 req.Chart,
	}
	assert.Equal(t, expectedAdditionalParams, retrievedJob.AdditionalRequestParameters)
	assert.Empty(t, retrievedJob.MetricSummary)
	assert.Empty(t, retrievedJob.ErrorMessage)
}

func TestUpdateJobStatus(t *testing.T) {
	js := setupTestDB(t)

	ctx := context.Background()
	jobID := uuid.New().String()
	initialReq := &pinpointpb.SchedulePairwiseRequest{
		Benchmark:     "test_benchmark",
		Configuration: "test_config",
		Story:         "test_story",
		StartCommit:   &startCommit,
		EndCommit:     &endCommit,
	}

	err := js.AddInitialJob(ctx, initialReq, jobID)
	require.NoError(t, err)

	newStatus := "Completed"
	durationInNanoseconds := int64(10 * time.Minute)
	err = js.UpdateJobStatus(ctx, jobID, newStatus, durationInNanoseconds)
	require.NoError(t, err)

	retrievedJob, err := js.GetJob(ctx, jobID)
	require.NoError(t, err)
	assert.Equal(t, newStatus, retrievedJob.JobStatus)
	assert.Equal(t, "10", retrievedJob.AdditionalRequestParameters["duration"])
}

func TestAddResults(t *testing.T) {
	js := setupTestDB(t)

	ctx := context.Background()
	jobID := uuid.New().String()
	initialReq := &pinpointpb.SchedulePairwiseRequest{
		Benchmark:     "test_benchmark",
		Configuration: "test_config",
		Story:         "test_story",
		StartCommit:   &startCommit,
		EndCommit:     &endCommit,
	}

	err := js.AddInitialJob(ctx, initialReq, jobID)
	require.NoError(t, err)

	results := map[string]*pinpointpb.PairwiseExecution_WilcoxonResult{
		"chart1": {
			Significant:              true,
			PValue:                   0.01,
			ConfidenceIntervalLower:  10.0,
			ConfidenceIntervalHigher: 20.0,
			ControlMedian:            15.0,
			TreatmentMedian:          25.0,
		},
		"chart2": {
			Significant:              false,
			PValue:                   0.5,
			ConfidenceIntervalLower:  5.0,
			ConfidenceIntervalHigher: 15.0,
			ControlMedian:            10.0,
			TreatmentMedian:          11.0,
		},
	}

	err = js.AddResults(ctx, jobID, results)
	require.NoError(t, err)

	retrievedJob, err := js.GetJob(ctx, jobID)
	require.NoError(t, err)
	assert.Equal(t, results, retrievedJob.MetricSummary)
}

func TestAddErrors(t *testing.T) {
	js := setupTestDB(t)

	ctx := context.Background()
	jobID := uuid.New().String()
	initialReq := &pinpointpb.SchedulePairwiseRequest{
		Benchmark:     "test_benchmark",
		Configuration: "test_config",
		Story:         "test_story",
		StartCommit:   &startCommit,
		EndCommit:     &endCommit,
	}

	err := js.AddInitialJob(ctx, initialReq, jobID)
	require.NoError(t, err)

	testError := errors.New("something went wrong during execution")
	err = js.SetErrors(ctx, jobID, testError)
	require.NoError(t, err)

	retrievedJob, err := js.GetJob(ctx, jobID)
	require.NoError(t, err)
	assert.Equal(t, testError.Error(), retrievedJob.ErrorMessage)

	// Test with nil error
	err = js.SetErrors(ctx, jobID, nil)
	require.NoError(t, err)
	retrievedJob, err = js.GetJob(ctx, jobID)
	require.NoError(t, err)
	assert.Empty(t, retrievedJob.ErrorMessage)
}

func TestAddCommitRuns(t *testing.T) {
	js := setupTestDB(t)

	ctx := context.Background()
	jobID := uuid.New().String()
	initialReq := &pinpointpb.SchedulePairwiseRequest{
		Benchmark:           "test_benchmark",
		Configuration:       "test_config",
		Story:               "test_story",
		InitialAttemptCount: "50",
		StartCommit:         &startCommit,
		EndCommit:           &endCommit,
	}

	err := js.AddInitialJob(ctx, initialReq, jobID)
	require.NoError(t, err)

	err = js.AddCommitRuns(ctx, jobID, &leftData, &rightData)
	require.NoError(t, err)

	retrievedJob, err := js.GetJob(ctx, jobID)
	require.NoError(t, err)

	// we need to unmarshal the "commit_runs" value from the map.
	commitRunsJSON, ok := retrievedJob.AdditionalRequestParameters["commit_runs"]
	assert.True(t, ok)

	var actualStoredCommitRuns storedCommitRuns
	err = json.NewDecoder(strings.NewReader(commitRunsJSON)).Decode(&actualStoredCommitRuns)
	require.NoError(t, err)

	assert.Equal(t, &leftData, actualStoredCommitRuns.Left)
	assert.Equal(t, &rightData, actualStoredCommitRuns.Right)
	assert.Equal(t, initialReq.Story, retrievedJob.AdditionalRequestParameters["story"])
}

func insertTestListJob(t *testing.T, js JobStore, jobID, jobName, jobType, benchmark, status, botName, user string) {
	jsi := js.(*jobStoreImpl)
	query := `
       INSERT INTO jobs (
           job_id,
           job_name,
           job_status,
           job_type,
           submitted_by,
           benchmark,
           bot_name,
           additional_request_parameters,
           error_message
       ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
       `
	_, err := jsi.db.Exec(
		context.Background(),
		query,
		jobID,
		jobName,
		status,
		jobType,
		user,
		benchmark,
		botName,
		map[string]string{},
		"",
	)
	require.NoError(t, err)
}

func TestListJobs(t *testing.T) {
	js := setupTestDB(t)
	ctx := context.Background()

	job1ID := uuid.New().String()
	job2ID := uuid.New().String()
	job3ID := uuid.New().String()

	insertTestListJob(t, js, job1ID, "Job A", "Pairwise", "speedometer", "Completed", "linux-perf", "user1@google.com")
	insertTestListJob(t, js, job2ID, "Job B", "Bisect", "jetstream", "Pending", "windows-perf", "user2@google.com")
	insertTestListJob(t, js, job3ID, "Another Job C", "Pairwise", "speedometer", "Running", "mac-perf", "user1@google.com")

	t.Run("Default behavior without options", func(t *testing.T) {
		jobs, err := js.ListJobs(ctx, ListJobsOptions{})
		require.NoError(t, err)
		require.Len(t, jobs, 3)
		// Default sort is by createdat DESC
		assert.Equal(t, job3ID, jobs[0].JobID)
		assert.Equal(t, job2ID, jobs[1].JobID)
		assert.Equal(t, job1ID, jobs[2].JobID)
	})

	t.Run("Search by term case-sensitive", func(t *testing.T) {
		// Search for a specific job, case-insensitive
		opts := ListJobsOptions{SearchTerm: "Job A"}
		jobs, err := js.ListJobs(ctx, opts)
		require.NoError(t, err)
		require.Len(t, jobs, 1)
		assert.Equal(t, job1ID, jobs[0].JobID)

		// Search for a broader term that matches all jobs
		opts = ListJobsOptions{SearchTerm: "Job"}
		jobs, err = js.ListJobs(ctx, opts)
		require.NoError(t, err)
		require.Len(t, jobs, 3)
	})

	t.Run("With limit", func(t *testing.T) {
		opts := ListJobsOptions{Limit: 2}
		jobs, err := js.ListJobs(ctx, opts)
		require.NoError(t, err)
		require.Len(t, jobs, 2)
		assert.Equal(t, job3ID, jobs[0].JobID)
		assert.Equal(t, job2ID, jobs[1].JobID)
	})

	t.Run("With limit and search", func(t *testing.T) {
		opts := ListJobsOptions{SearchTerm: "Job", Limit: 2}
		jobs, err := js.ListJobs(ctx, opts)
		require.NoError(t, err)
		require.Len(t, jobs, 2)
		assert.Equal(t, job3ID, jobs[0].JobID)
		assert.Equal(t, job2ID, jobs[1].JobID)
	})

	t.Run("With limit and offset for pagination", func(t *testing.T) {
		// Get the second page, which should only have one job.
		// Jobs are ordered by creation date DESC, so job3, job2, job1.
		opts := ListJobsOptions{Limit: 2, Offset: 2}
		jobs, err := js.ListJobs(ctx, opts)
		require.NoError(t, err)
		require.Len(t, jobs, 1)
		assert.Equal(t, job1ID, jobs[0].JobID)
	})
}
