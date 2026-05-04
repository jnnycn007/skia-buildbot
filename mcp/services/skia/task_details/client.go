package task_details

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/bigtable"
	"cloud.google.com/go/datastore"
	"github.com/mark3labs/mcp-go/mcp"
	"go.chromium.org/luci/grpc/prpc"
	"go.chromium.org/luci/logdog/client/coordinator"
	"go.chromium.org/luci/logdog/common/fetcher"
	"go.chromium.org/luci/logdog/common/types"
	annopb "go.chromium.org/luci/luciexe/legacy/annotee/proto"
	"go.skia.org/infra/go/auth"
	"go.skia.org/infra/go/httputils"
	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/go/swarming"
	td_db "go.skia.org/infra/task_driver/go/db"
	td_bigtable "go.skia.org/infra/task_driver/go/db/bigtable"
	"go.skia.org/infra/task_driver/go/logs"
	ts_db "go.skia.org/infra/task_scheduler/go/db"
	"go.skia.org/infra/task_scheduler/go/db/firestore"
	"golang.org/x/oauth2/google"
	"google.golang.org/protobuf/proto"
)

const (
	logdogProject = "skia"
	logdogHost    = "logs.chromium.org"

	logdogPathTmplRun      = "%s/+/annotations"
	logdogPathTmplStepLogs = "%s/+/%s"
)

type TaskDetailsClient struct {
	swarm  swarming.ApiClient
	td     td_db.DB
	tdLogs *logs.LogsManager
	ts     ts_db.DBCloser
	logdog *coordinator.Client
}

func NewClient(ctx context.Context, btProject, btInstance, firestoreInstance string) (*TaskDetailsClient, error) {
	ts, err := google.DefaultTokenSource(ctx, auth.ScopeUserinfoEmail, bigtable.Scope, datastore.ScopeDatastore)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	tdDB, err := td_bigtable.NewBigTableDB(ctx, btProject, btInstance, ts)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	tsDB, err := firestore.NewDBWithParams(ctx, firestore.FIRESTORE_PROJECT, firestoreInstance, ts)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	c := httputils.DefaultClientConfig().WithTokenSource(ts).Client()
	swarm, err := swarming.NewApiClient(c, swarming.SWARMING_SERVER)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	prpcClient := prpc.Client{
		C:       c,
		Host:    logdogHost,
		Options: prpc.DefaultOptions(),
	}
	coord := coordinator.NewClient(&prpcClient)
	tdLogs, err := logs.NewLogsManager(ctx, btProject, btInstance, ts)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	return &TaskDetailsClient{
		swarm:  swarm,
		td:     tdDB,
		tdLogs: tdLogs,
		ts:     tsDB,
		logdog: coord,
	}, nil
}

type getTaskStepsResult struct {
	TaskDriver *td_db.TaskDriverRun `json:"task_driver,omitempty"`

	Recipe                  *annopb.Step     `json:"recipe,omitempty"`
	RecipeStepStatusMapping map[int32]string `json:"recipe_step_status_mapping,omitempty"`

	SwarmingTaskID    string `json:"swarming_task_id,omitempty"`
	SwarmingTaskState string `json:"swarming_task_state,omitempty"`
	SwarmingTaskLogs  string `json:"swarming_task_logs,omitempty"`
}

func (c *TaskDetailsClient) GetTaskStepsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskID, err := req.RequireString(argTaskID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// The task might be run as a Task Driver, a Recipe, or just a plain task.
	// Try Task Driver first.
	var res getTaskStepsResult
	td, err := c.td.GetTaskDriver(ctx, taskID)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	if td != nil {
		res.TaskDriver = td
		return encodeJSONResponse(res)
	}

	// Fall back to Recipe steps via LogDog.
	task, err := c.ts.GetTaskById(ctx, taskID)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	step, err := c.fetchLogDogSteps(ctx, task.SwarmingTaskId)
	if err == nil {
		res.Recipe = step
		res.RecipeStepStatusMapping = annopb.Status_name
		// Populate SwarmingTaskID in case it's needed for log retrieval.
		res.SwarmingTaskID = task.SwarmingTaskId
		return encodeJSONResponse(res)
	} else if !strings.Contains(err.Error(), "coordinator: no access") {
		return nil, skerr.Wrap(err)
	}

	// If we couldn't find recipe steps, just return the Swarming task logs.
	swarmOutput, err := c.swarm.GetStdoutOfTask(ctx, task.SwarmingTaskId)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	res.SwarmingTaskState = swarmOutput.State
	res.SwarmingTaskLogs = swarmOutput.Output
	return encodeJSONResponse(res)
}

func (c *TaskDetailsClient) GetRecipeStepLogsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	swarmingTaskID, err := req.RequireString(argSwarmingTaskID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	logPath, err := req.RequireString(argLogPath)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	startIndex, err := req.RequireInt(argStartIndex)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	limit, err := req.RequireInt(argLimit)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	lines, err := c.fetchLogDogStepLogs(ctx, swarmingTaskID, logPath, startIndex, limit)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	return encodeJSONResponse(lines)
}

func (c *TaskDetailsClient) GetTaskDriverLogsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskID, err := req.RequireString(argTaskID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	stepID, err := req.RequireString(argStepID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	logID, err := req.RequireString(argLogID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	cursor := req.GetString(argCursor, "")
	limit := req.GetInt(argLimit, 0)

	logs, cursor, err := c.tdLogs.Search(ctx, taskID, stepID, logID, cursor, limit)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	response := struct {
		Logs   []string `json:"log_lines"`
		Cursor string   `json:"cursor"`
	}{
		Cursor: cursor,
	}
	for _, entry := range logs {
		if entry.TextPayload != "" {
			response.Logs = append(response.Logs, entry.TextPayload)
		} else if entry.JsonPayload != nil {
			b, err := json.Marshal(entry.JsonPayload)
			if err != nil {
				return nil, skerr.Wrap(err)
			} else {
				response.Logs = append(response.Logs, string(b))
			}
		} else {
			response.Logs = append(response.Logs, "")
		}
	}

	return encodeJSONResponse(response)
}

// fixupTaskID ensures that the given Swarming task ID is a *run* ID as opposed
// to a *request* ID. The request ID ends with a zero, while the first for a
// given request ends in a one.
func fixupTaskID(taskID string) string {
	if len(taskID) > 0 && taskID[len(taskID)-1] == '0' {
		return taskID[:len(taskID)-1] + "1"
	}
	return taskID
}

func (c *TaskDetailsClient) fetchLogDogSteps(ctx context.Context, taskID string) (*annopb.Step, error) {
	path := fmt.Sprintf(logdogPathTmplRun, fixupTaskID(taskID))
	stream := c.logdog.Stream(logdogProject, types.StreamPath(path))
	var state coordinator.LogStream
	le, err := stream.Tail(ctx, coordinator.WithState(&state), coordinator.Complete())
	if err != nil {
		return nil, skerr.Wrapf(err, "failed to tail stream")
	}
	if le == nil {
		return nil, skerr.Fmt("no annotation entries found in stream")
	}

	if state.Desc.ContentType != annopb.ContentTypeAnnotations {
		return nil, skerr.Fmt("expected annotations but found %s", state.Desc.ContentType)
	}
	dg := le.GetDatagram()
	if dg == nil {
		return nil, skerr.Fmt("no datagram found for step!")
	}
	var step annopb.Step
	if err := proto.Unmarshal(dg.Data, &step); err != nil {
		return nil, skerr.Wrapf(err, "failed to unmarshal datagram data")
	}
	return &step, nil
}

func (c *TaskDetailsClient) fetchLogDogStepLogs(ctx context.Context, taskID, logPath string, index, count int) ([]string, error) {
	path := fmt.Sprintf(logdogPathTmplStepLogs, fixupTaskID(taskID), logPath)
	f := c.logdog.Stream(logdogProject, types.StreamPath(path)).Fetcher(ctx, &fetcher.Options{
		Index: types.MessageIndex(index),
		Count: int64(count),
	})

	var logLines []string
	for {
		le, err := f.NextLogEntry()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, skerr.Wrapf(err, "failed to fetch log entry")
		}
		if text := le.GetText(); text != nil {
			for _, line := range text.Lines {
				logLines = append(logLines, string(line.Value))
			}
		}
	}

	return logLines, nil
}

func encodeJSONResponse(resp interface{}) (*mcp.CallToolResult, error) {
	b, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	return mcp.NewToolResultText(string(b)), nil
}
