package task_details

import "go.skia.org/infra/mcp/common"

const (
	argTaskID         = "task_id"
	argStepID         = "step_id"
	argLogID          = "log_id"
	argSwarmingTaskID = "swarming_task_id"
	argLogPath        = "log_path"
	argStartIndex     = "start_index"
	argLimit          = "limit"
	argCursor         = "cursor"
)

func GetTools(c *TaskDetailsClient) []common.Tool {
	return []common.Tool{
		{
			Name:        "get_task_steps",
			Description: "Retrieve the full step listing for a task. Depending on the task, this may return a Task Driver, a Recipe, or a raw Swarming task log.",
			Arguments: []common.ToolArgument{
				{
					Name:        argTaskID,
					Description: "ID of the task.",
					Required:    true,
				},
			},
			Handler: c.GetTaskStepsHandler,
		},
		{
			Name:        "get_task_driver_step_logs",
			Description: "Retrieve log entries for a task driver step.",
			Arguments: []common.ToolArgument{
				{
					Name:        argTaskID,
					Description: "ID of the task.",
					Required:    true,
				},
				{
					Name:        argStepID,
					Description: "ID of the step.",
					Required:    true,
				},
				{
					Name:        argLogID,
					Description: "ID of the log.",
					Required:    true,
				},
			},
			Handler: c.GetTaskDriverLogsHandler,
		},
		{
			Name:        "get_recipe_step_logs",
			Description: "Retrieve log lines for a recipe step.",
			Arguments: []common.ToolArgument{
				{
					Name:        argSwarmingTaskID,
					Description: "ID of the Swarming task.",
					Required:    true,
				},
				{
					Name:        argLogPath,
					Description: "Path of the step log as specified in Recipe step result data.",
					Required:    true,
				},
				{
					Name:        argStartIndex,
					Description: "Starting log message index.",
					Required:    true,
				},
				{
					Name:        argLimit,
					Description: "Maximum number of entries to load. Logs may be very large, so providing a limit is recommended.",
				},
				{
					Name:        argCursor,
					Description: "Starting point for the next page of results, when paginating.",
				},
			},
			Handler: c.GetRecipeStepLogsHandler,
		},
	}
}
