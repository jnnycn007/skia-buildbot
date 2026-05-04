package task_scheduler

import (
	"fmt"

	"go.skia.org/infra/mcp/common"
	"go.skia.org/infra/task_scheduler/go/db"
	"go.skia.org/infra/task_scheduler/go/types"
)

const (
	argStartTime     = "start_time"
	argEndTime       = "end_time"
	argIssue         = "issue"
	argPatchset      = "patchset"
	argTaskStatus    = "status"
	argRepo          = "repo"
	argRevision      = "revision"
	argTaskName      = "name"
	argLimit         = "limit"
	argIncludeStable = "include_stable"
	argTaskId        = "task_id"

	taskStatusPending = "PENDING"
)

func GetTools(c *TaskSchedulerClient) []common.Tool {
	return []common.Tool{
		{
			Name:        "search_tasks",
			Description: `Retrieve a list of matching tasks from the database.`,
			Arguments: []common.ToolArgument{
				{
					Name: argStartTime,
					Description: `
[Optional] The start of the time range to search for tasks.
The input should be in the RFC 3339 format and GMT should be
used as the default timezone, eg. "2025-07-12T14:30:00-00:00".`,
				},
				{
					Name: argEndTime,
					Description: `
[Optional] The end of the time range to search for tasks.
The input should be in the RFC 3339 format and GMT should be
used as the default timezone, eg. "2025-07-12T14:30:00-00:00".
If not provided, the current time is used.`,
				},
				{
					Name:        argIssue,
					Description: "[Optional] CL issue ID. If not provided, try jobs are excluded from results.",
				},
				{
					Name:        argPatchset,
					Description: "[Optional] CL patchset ID. If not provided, try jobs are excluded from results.",
				},
				{
					Name: argTaskStatus,
					Description: fmt.Sprintf("[Optional] Task status, one of %v", []string{
						taskStatusPending,
						string(types.TASK_STATUS_RUNNING),
						string(types.TASK_STATUS_SUCCESS),
						string(types.TASK_STATUS_FAILURE),
						string(types.TASK_STATUS_MISHAP),
					}),
				},
				{
					Name:        argRepo,
					Description: `[Optional] Git repository URL of the task, eg. "https://skia.googlesource.com/skia.git"`,
				},
				{
					Name:        argRevision,
					Description: "[Optional] Full git commit hash at which the task ran.",
				},
				{
					Name:        argTaskName,
					Description: "[Optional] Name of the task.",
				},
				{
					Name:        argLimit,
					Description: fmt.Sprintf("[Optional] Maximum number of tasks to return. Default %d", db.SearchResultLimit),
				},
			},
			Handler: c.SearchTasksHandler,
		},
		{
			Name:        "get_task_health_report",
			Description: "Retrieve a summary of task health over a series of commits.",
			Arguments: []common.ToolArgument{
				{
					Name:        argRepo,
					Description: `Git repository URL of the task, eg. "https://skia.googlesource.com/skia.git"`,
					Required:    true,
				},
				{
					Name:        argRevision,
					Description: "Git commit hash or branch name to start at.",
					Required:    true,
				},
				{
					Name:        argLimit,
					Description: "Number of commits to trace backward in history.",
					Required:    true,
				},
				{
					Name:         argIncludeStable,
					Description:  "If set, include results for tasks which are succeeding or failing consistently within the given window.",
					ArgumentType: common.BooleanArgument,
				},
			},
			Handler: c.GetTaskHealthReportHandler,
		},
		{
			Name:        "get_task",
			Description: "Retrieve the full details for a task.",
			Arguments: []common.ToolArgument{
				{
					Name:        argTaskId,
					Description: "ID of the task.",
					Required:    true,
				},
			},
			Handler: c.GetTaskHandler,
		},
	}
}
