package task_scheduler

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/mark3labs/mcp-go/mcp"
	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/task_scheduler/go/db"
	"go.skia.org/infra/task_scheduler/go/db/firestore"
	"go.skia.org/infra/task_scheduler/go/types"
	"golang.org/x/oauth2/google"
)

type TaskSchedulerClient struct {
	db db.DBCloser
}

func NewClient(ctx context.Context, firestoreInstance string) (*TaskSchedulerClient, error) {
	ts, err := google.DefaultTokenSource(ctx, datastore.ScopeDatastore)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	db, err := firestore.NewDBWithParams(ctx, firestore.FIRESTORE_PROJECT, firestoreInstance, ts)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	return &TaskSchedulerClient{
		db: db,
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
