package regression

import (
	"context"

	"github.com/jackc/pgx/v4"
	"go.skia.org/infra/perf/go/alerts"
	"go.skia.org/infra/perf/go/clustering2"
	"go.skia.org/infra/perf/go/progress"
	pb "go.skia.org/infra/perf/go/subscription/proto/v1"
	"go.skia.org/infra/perf/go/types"
	"go.skia.org/infra/perf/go/ui/frame"
)

// Store persists Regressions.
type Store interface {
	// Range returns a map from types.CommitNumber to *Regressions that exist in the
	// given range of commits. Note that if begin==end that results
	// will be returned for begin.
	Range(ctx context.Context, begin, end types.CommitNumber) (map[types.CommitNumber]*AllRegressionsForCommit, error)

	// RangeFiltered gets all regressions in the given commit range and trace names.
	RangeFiltered(ctx context.Context, begin, end types.CommitNumber, traceNames []string) ([]*Regression, error)

	// SetHigh sets the ClusterSummary for a high regression at the given commit and alertID.
	SetHigh(ctx context.Context, commitNumber types.CommitNumber, alertID string, df *frame.FrameResponse, high *clustering2.ClusterSummary) (bool, string, error)

	// SetLow sets the ClusterSummary for a low regression at the given commit and alertID.
	SetLow(ctx context.Context, commitNumber types.CommitNumber, alertID string, df *frame.FrameResponse, low *clustering2.ClusterSummary) (bool, string, error)

	// TriageLow sets the triage status for the low cluster at the given commit and alertID.
	TriageLow(ctx context.Context, commitNumber types.CommitNumber, alertID string, tr TriageStatus) error

	// TriageHigh sets the triage status for the high cluster at the given commit and alertID.
	TriageHigh(ctx context.Context, commitNumber types.CommitNumber, alertID string, tr TriageStatus) error

	// Write the Regressions to the store. The provided 'regressions' maps from
	// types.CommitNumber to all the regressions for that commit.
	Write(ctx context.Context, regressions map[types.CommitNumber]*AllRegressionsForCommit) error

	// Given the subscription name GetRegressionsBySubName gets all the regressions against
	// the specified subscription. The response will be paginated according to the provided
	// limit and offset.
	GetRegressionsBySubName(ctx context.Context, req GetAnomalyListRequest, limit int) ([]*Regression, error)

	// Given a list of regression IDs (only in the regression2store),
	// return a list of regressions.
	GetByIDs(ctx context.Context, ids []string) ([]*Regression, error)

	// GetIdsByManualTriageBugID returns a list of distinct regression ids with given manual triage bug id.
	GetIdsByManualTriageBugID(ctx context.Context, bugID int) ([]string, error)

	// Return a list of regressions satisfying: previous_commit < rev <= commit.
	GetByRevision(ctx context.Context, rev string) ([]*Regression, error)

	// GetOldestCommit returns the commit with the lowest commit number
	GetOldestCommit(ctx context.Context) (*types.CommitNumber, error)

	// GetRegression returns the regression info at the given commit for specific alert.
	GetRegression(ctx context.Context, commitNumber types.CommitNumber, alertID string) (*Regression, error)

	// DeleteByCommit deletes a regression from the Regression table via the CommitNumber.
	// Use with caution.
	DeleteByCommit(ctx context.Context, commitNumber types.CommitNumber, tx pgx.Tx) error

	// SetBugID associates a set of regressions, identified by their IDs, with a bug ID.
	SetBugID(ctx context.Context, regressionIDs []string, bugID int) error

	// IgnoreAnomalies sets the triage status to Ignored and message to IgnoredMessage for the given regressions.
	IgnoreAnomalies(ctx context.Context, regressionIDs []string) error

	// ResetAnomalies sets the triage status to Untriaged, message to ResetMessage, and bugID to 0 for the given regressions.
	ResetAnomalies(ctx context.Context, regressionIDs []string) error

	// NudgeAndResetAnomalies updates the commit number and previous commit number for the given regressions,
	// and also sets the triage status to Untriaged, message to NudgedMessage, and bugID to 0.
	NudgeAndResetAnomalies(ctx context.Context, regressionIDs []string, commitNumber, prevCommitNumber types.CommitNumber) error

	// GetBugIdsForRegressions queries all bugs from regressions2, culprits and anomalygroups for given regressions.
	GetBugIdsForRegressions(ctx context.Context, regressions []*Regression) ([]*Regression, error)

	// GetSubscriptionsForRegressions returns a subset of subscription fields for given regressions, together with regression and alert ids.
	GetSubscriptionsForRegressions(ctx context.Context, regressionIDs []string) ([]string, []int64, []*pb.Subscription, error)
}

// FullSummary describes a single regression.
type FullSummary struct {
	Summary clustering2.ClusterSummary `json:"summary"`
	Triage  TriageStatus               `json:"triage"`
	Frame   frame.FrameResponse        `json:"frame"`
}

// Request object for the request from the anomaly table UI.
type GetAnomalyListRequest struct {
	SubName             string `json:"sheriff"`
	IncludeTriaged      bool   `json:"triaged"`
	IncludeImprovements bool   `json:"improvements"`
	QueryCursor         string `json:"anomaly_cursor"`
	Host                string `json:"host"`
	PaginationOffset    int    `json:"pagination_offset,omitempty"`
}

// RegressionDetectionRequest is all the info needed to start a clustering run,
// an Alert and the Domain over which to run that Alert.
type RegressionDetectionRequest struct {
	Alert  *alerts.Alert `json:"alert"`
	Domain types.Domain  `json:"domain"`

	// query is the exact query being run. It may be more specific than the one
	// in the Alert if the Alert has a non-empty GroupBy.
	query string

	// Step/TotalQueries is the current percent of all the queries that have been processed.
	Step int `json:"step"`

	// TotalQueries is the number of sub-queries to be processed based on the
	// GroupBy setting in the Alert.
	TotalQueries int `json:"total_queries"`

	// Progress of the detection request.
	Progress progress.Progress `json:"-"`
}

// Query returns the query that the RegressionDetectionRequest process is
// running.
//
// Note that it may be more specific than the Alert.Query if the Alert has a
// non-empty GroupBy value.
func (r *RegressionDetectionRequest) Query() string {
	if r.query != "" {
		return r.query
	}
	if r.Alert != nil {
		return r.Alert.Query
	}
	return ""
}

// SetQuery sets a more refined query for the RegressionDetectionRequest.
func (r *RegressionDetectionRequest) SetQuery(q string) {
	r.query = q
}

// NewRegressionDetectionRequest returns a new RegressionDetectionRequest.
func NewRegressionDetectionRequest() *RegressionDetectionRequest {
	return &RegressionDetectionRequest{
		Progress: progress.New(),
	}
}

// RegressionDetectionResponse is the response from running a RegressionDetectionRequest.
type RegressionDetectionResponse struct {
	Summary *clustering2.ClusterSummaries `json:"summary"`
	Frame   *frame.FrameResponse          `json:"frame"`

	// Message contains context about the detection for this specific response,
	// such as trace filtering statistics.
	Message string `json:"-"` // Using json:"-" prevents it from being serialized by default.
}

// ConfirmedRegression is an alias for RegressionDetectionResponse used by the RegressionRefiner
// and ConfirmedRegressionHandler to represent regressions that have been validated and approved
// for saving or alerting.
type ConfirmedRegression RegressionDetectionResponse

// RegressionRefiner defines an interface for modules that process a complete
// set of regression detection results before they are sent for storage.
type RegressionRefiner interface {
	// Process takes a slice of RegressionDetectionResponse (the raw results of
	// regression detection). It returns a processed slice of ConfirmedRegression,
	// which contains the anomalies (e.g. high or low status when exceeding a threshold)
	// that we want to save into the database, send notifications and etc.
	// ConfirmedRegression is an alias for RegressionDetectionResponse used by the RegressionRefiner
	// and ConfirmedRegressionHandler to represent regressions that have been validated and approved
	// for saving or alerting.
	Process(ctx context.Context, cfg *alerts.Alert, responses []*RegressionDetectionResponse) ([]*ConfirmedRegression, error)
}
