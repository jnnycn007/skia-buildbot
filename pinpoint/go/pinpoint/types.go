// TODO(b/500974820): Reuse types from `pinpoint/proto/v1/service.pb.go`.
package pinpoint

type TryJobCreateRequest struct {
	Name        string `json:"name"`
	BaseGitHash string `json:"base_git_hash"`
	// although "experiment" makes more sense in this context, the legacy Pinpoint API
	// explicitly defines the experiment commit as "end_git_hash" and defines
	// the experiment patch as "experiment_patch"
	EndGitHash      string `json:"end_git_hash"`
	BasePatch       string `json:"base_patch"`
	ExperimentPatch string `json:"experiment_patch"`
	Configuration   string `json:"configuration"`
	Benchmark       string `json:"benchmark"`
	Story           string `json:"story"`
	ExtraTestArgs   string `json:"extra_test_args"`
	Repository      string `json:"repository"`
	BugId           string `json:"bug_id"`
	User            string `json:"user"`
}

type BisectJobCreateRequest struct {
	ComparisonMode      string `json:"comparison_mode"`
	StartGitHash        string `json:"start_git_hash"`
	EndGitHash          string `json:"end_git_hash"`
	Configuration       string `json:"configuration"`
	Benchmark           string `json:"benchmark"`
	Story               string `json:"story"`
	Chart               string `json:"chart"`
	Statistic           string `json:"statistic"`
	ComparisonMagnitude string `json:"comparison_magnitude"`
	Pin                 string `json:"pin"`
	Project             string `json:"project"`
	BugId               string `json:"bug_id"`
	User                string `json:"user"`
	AlertIDs            string `json:"alert_ids"`
	TestPath            string `json:"test_path"`
}

type CreatePinpointResponse struct {
	JobID  string `json:"jobId"`
	JobURL string `json:"jobUrl"`
}
