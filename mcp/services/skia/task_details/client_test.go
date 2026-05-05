package task_details

import (
	"testing"

	"github.com/stretchr/testify/assert"
	annopb "go.chromium.org/luci/luciexe/legacy/annotee/proto"
	"go.skia.org/infra/task_driver/go/display"
	"go.skia.org/infra/task_driver/go/td"
)

func TestGetTaskStepsResult_String_TaskDriver(t *testing.T) {
	res := GetTaskStepsResult{
		TaskDriver: &display.TaskDriverRunDisplay{
			StepDisplay: &display.StepDisplay{
				StepProperties: &td.StepProperties{Name: "Root Step"},
				Result:         td.StepResultSuccess,
				Steps: []*display.StepDisplay{
					{
						StepProperties: &td.StepProperties{Name: "Sub Step 1"},
						Result:         td.StepResultSuccess,
						Steps: []*display.StepDisplay{
							{
								StepProperties: &td.StepProperties{Name: "Sub-sub Step 1"},
								Result:         td.StepResultSuccess,
							},
						},
					},
					{
						StepProperties: &td.StepProperties{Name: "Sub Step 2"},
						Result:         td.StepResultFailure,
					},
				},
			},
		},
	}

	expected := `# Task Driver

- Root Step (SUCCESS)
  - Sub Step 1 (SUCCESS)
    - Sub-sub Step 1 (SUCCESS)
  - Sub Step 2 (FAILURE)
`
	assert.Equal(t, expected, res.String())
}

func TestGetTaskStepsResult_String_Recipe(t *testing.T) {
	res := GetTaskStepsResult{
		Recipe: &annopb.Step{
			Name:   "Root Step",
			Status: annopb.Status_SUCCESS,
			Substep: []*annopb.Step_Substep{
				{
					Substep: &annopb.Step_Substep_Step{
						Step: &annopb.Step{
							Name:   "Sub Step 1",
							Status: annopb.Status_SUCCESS,
							Substep: []*annopb.Step_Substep{
								{
									Substep: &annopb.Step_Substep_Step{
										Step: &annopb.Step{
											Name:   "Sub-sub Step 1",
											Status: annopb.Status_SUCCESS,
										},
									},
								},
							},
						},
					},
				},
				{
					Substep: &annopb.Step_Substep_Step{
						Step: &annopb.Step{
							Name:   "Sub Step 2",
							Status: annopb.Status_FAILURE,
						},
					},
				},
			},
		},
		SwarmingTaskID: "abc123",
	}

	expected := `# Recipe

**Swarming Task ID:** abc123
**Steps:**
- Root Step (SUCCESS)
  - Sub Step 1 (SUCCESS)
    - Sub-sub Step 1 (SUCCESS)
  - Sub Step 2 (FAILURE)
`
	assert.Equal(t, expected, res.String())
}

func TestGetTaskStepsResult_String_Swarming(t *testing.T) {
	res := GetTaskStepsResult{
		SwarmingTaskID:    "abc123",
		SwarmingTaskState: "SUCCESS",
		SwarmingTaskLogs:  "Log line 1\nLog line 2",
	}

	expected := `# Swarming Task

**ID:**    abc123
**State:** SUCCESS
**Logs:**
` + "```" + `
Log line 1
Log line 2
` + "```\n"

	assert.Equal(t, expected, res.String())
}
