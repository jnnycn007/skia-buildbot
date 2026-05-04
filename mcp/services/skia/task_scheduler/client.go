package task_scheduler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/mark3labs/mcp-go/mcp"
	"go.skia.org/infra/go/gitiles"
	"go.skia.org/infra/go/httputils"
	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/task_scheduler/go/db"
	"go.skia.org/infra/task_scheduler/go/db/firestore"
	"go.skia.org/infra/task_scheduler/go/types"
	"golang.org/x/oauth2/google"
)

type TaskSchedulerClient struct {
	client *http.Client
	db     db.DBCloser
}

func NewClient(ctx context.Context, firestoreInstance string) (*TaskSchedulerClient, error) {
	ts, err := google.DefaultTokenSource(ctx, datastore.ScopeDatastore)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	client := httputils.DefaultClientConfig().WithTokenSource(ts).Client()
	db, err := firestore.NewDBWithParams(ctx, firestore.FIRESTORE_PROJECT, firestoreInstance, ts)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	return &TaskSchedulerClient{
		client: client,
		db:     db,
	}, nil
}

func (c *TaskSchedulerClient) Close() error {
	return skerr.Wrap(c.db.Close())
}

func (c *TaskSchedulerClient) SearchTasksHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	startTime, err := parseTimeOrNil(req, argStartTime)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	endTime, err := parseTimeOrNil(req, argEndTime)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// If issue and patchset aren't provided, assume the caller doesn't want try jobs included.
	issue := req.GetString(argIssue, "")
	patchset := req.GetString(argPatchset, "")

	status := getStringOrNil(req, argTaskStatus)
	if status != nil && *status == "PENDING" {
		*status = string(types.TASK_STATUS_PENDING)
	}

	limit := req.GetInt(argLimit, db.SearchResultLimit)

	searchParams := &db.TaskSearchParams{
		Status:    (*types.TaskStatus)(status),
		Issue:     &issue,
		Name:      getStringOrNil(req, argTaskName),
		Patchset:  &patchset,
		Repo:      getStringOrNil(req, argRepo),
		Revision:  getStringOrNil(req, argRevision),
		TimeStart: startTime,
		TimeEnd:   endTime,
		Limit:     &limit,
	}
	tasks, err := c.db.SearchTasks(ctx, searchParams)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	b, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}

type TaskHealthReport struct {
	CommitGraph []string                       `json:"commit_graph"`
	Tasks       map[string][]TaskResultSummary `json:"tasks"`
}

type TaskResultSummary struct {
	ID        string           `json:"id"`
	Revision  string           `json:"rev"`
	Status    types.TaskStatus `json:"status"`
	BlameSize int              `json:"blame_size"`
}

func (c *TaskSchedulerClient) GetTaskHealthReportHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	repoUrl, err := req.RequireString(argRepo)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	revision, err := req.RequireString(argRevision)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	limit, err := req.RequireInt(argLimit)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	includeStable := req.GetBool(argIncludeStable, false)

	repo := gitiles.NewRepo(repoUrl, c.client)

	commits, err := repo.Log(ctx, revision, gitiles.LogLimit(limit))
	if err != nil {
		return nil, err
	}

	commitHashes := make([]string, 0, len(commits))
	for _, c := range commits {
		commitHashes = append(commitHashes, c.Hash)
	}

	taskHistory := make(map[string][]TaskResultSummary)
	seenTasks := make(map[string]bool)

	for _, commit := range commits {
		tasks, err := c.db.SearchTasks(ctx, &db.TaskSearchParams{
			Repo:              &repoUrl,
			BlamelistContains: &commit.Hash,
		})
		if err != nil {
			return nil, err
		}

		for _, t := range tasks {
			if seenTasks[t.Id] {
				continue
			}
			seenTasks[t.Id] = true

			taskHistory[t.Name] = append(taskHistory[t.Name], TaskResultSummary{
				ID:        t.Id,
				Revision:  t.Revision,
				Status:    t.Status,
				BlameSize: len(t.Commits),
			})
		}
	}

	report := &TaskHealthReport{
		CommitGraph: commitHashes,
		Tasks:       make(map[string][]TaskResultSummary),
	}

	for name, history := range taskHistory {
		if !includeStable {
			if len(history) > 0 {
				firstStatus := history[0].Status
				stable := true
				for _, res := range history {
					if res.Status != firstStatus {
						stable = false
						break
					}
				}
				if stable {
					continue
				}
			}
		}
		report.Tasks[name] = history
	}
	return encodeJSONResponse(report)
}

func (c *TaskSchedulerClient) GetTaskHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskID, err := req.RequireString(argTaskId)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	task, err := c.db.GetTaskById(ctx, taskID)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	return encodeJSONResponse(task)
}

func getStringOrNil(req mcp.CallToolRequest, arg string) *string {
	// Using RequireString is cleaner than choosing a placeholder to use when
	// the arg is not provided, even if we don't really require it.
	str, err := req.RequireString(arg)
	if err != nil {
		return nil
	}
	return &str
}

func parseTimeOrNil(req mcp.CallToolRequest, arg string) (*time.Time, error) {
	str, err := req.RequireString(arg)
	if err != nil {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	return &parsed, nil
}

func encodeJSONResponse(resp interface{}) (*mcp.CallToolResult, error) {
	b, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	return mcp.NewToolResultText(string(b)), nil
}
