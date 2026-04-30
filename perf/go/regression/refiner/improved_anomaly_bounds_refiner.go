package refiner

import (
	"context"
	"fmt"
	"math"

	"go.skia.org/infra/go/sklog"
	"go.skia.org/infra/go/vec32"
	"go.skia.org/infra/perf/go/alerts"
	"go.skia.org/infra/perf/go/anomalies"
	"go.skia.org/infra/perf/go/git"
	"go.skia.org/infra/perf/go/regression"
	"go.skia.org/infra/perf/go/stepfit"
	"go.skia.org/infra/perf/go/tracestore"
	"go.skia.org/infra/perf/go/types"
)

// ImprovedAnomalyBoundsRefiner implements regression.RegressionRefiner.
// It runs the standard AnomalyBoundsRefiner first, and then applies
// additional refinement logic based on previous regressions found in the database
// and loading raw data from the database.
type ImprovedAnomalyBoundsRefiner struct {
	base            *AnomalyBoundsRefiner
	anomalyStore    anomalies.Store
	store           regression.Store
	traceStore      tracestore.TraceStore
	perfGit         git.Git
	stdDevThreshold float32
}

// NewImprovedAnomalyBoundsRefiner returns a new instance of ImprovedAnomalyBoundsRefiner.
func NewImprovedAnomalyBoundsRefiner(anomalyStore anomalies.Store, store regression.Store, traceStore tracestore.TraceStore, perfGit git.Git, stdDevThreshold float32) *ImprovedAnomalyBoundsRefiner {
	return &ImprovedAnomalyBoundsRefiner{
		base:            &AnomalyBoundsRefiner{stdDevThreshold: stdDevThreshold},
		anomalyStore:    anomalyStore,
		store:           store,
		traceStore:      traceStore,
		perfGit:         perfGit,
		stdDevThreshold: stdDevThreshold,
	}
}

// Process implements the regression.RegressionRefiner interface.
func (r *ImprovedAnomalyBoundsRefiner) Process(ctx context.Context, cfg *alerts.Alert, responses []*regression.RegressionDetectionResponse) ([]*regression.ConfirmedRegression, error) {
	// 1. Run the base AnomalyBoundsRefiner logic.
	confirmed, err := r.base.Process(ctx, cfg, responses)
	if err != nil {
		return nil, err
	}

	var refined []*regression.ConfirmedRegression

	// 2. Apply additional action on confirmed regressions.
	for _, cr := range confirmed {
		newCr := r.applyImprovedLogic(ctx, cr, cfg)
		if newCr != nil {
			refined = append(refined, newCr)
		}
	}

	return refined, nil
}

func (r *ImprovedAnomalyBoundsRefiner) applyImprovedLogic(ctx context.Context, cr *regression.ConfirmedRegression, cfg *alerts.Alert) *regression.ConfirmedRegression {
	if len(cr.Summary.Clusters) == 0 || len(cr.Summary.Clusters[0].Keys) == 0 {
		sklog.Infof("[ImprovedAnomalyBoundsRefiner] Skipping improved logic for regression at %d because it has no clusters or keys", cr.DisplayCommitNumber)
		return cr
	}
	traceName := cr.Summary.Clusters[0].Keys[0]
	pickOffset := cr.DisplayCommitNumber // Use DisplayCommitNumber as pick point

	// 1. Load last regression <= pickOffset from DB.
	regressions, err := r.store.GetRegressionsBefore(ctx, traceName, pickOffset, 1)
	if err != nil {
		sklog.Errorf("[ImprovedAnomalyBoundsRefiner] Failed to get regressions before %d: %s", pickOffset, err)
		return cr // Return original if we can't check DB
	}

	var lastRegression *regression.Regression
	if len(regressions) > 0 {
		lastRegression = regressions[0]
	}

	if lastRegression == nil {
		sklog.Infof("[ImprovedAnomalyBoundsRefiner] No previous regression found for trace %s before %d. Keeping original regression.", traceName, pickOffset)
		return cr
	}

	// Check for overlap.
	if lastRegression.CommitNumber >= cr.PrevCommitNumber && lastRegression.PrevCommitNumber <= cr.CommitNumber {
		sklog.Infof("[ImprovedAnomalyBoundsRefiner] Filtering out regression at %d due to overlap with existing regression at %d", pickOffset, lastRegression.CommitNumber)
		return nil // Filter out
	}

	// 2. Load data points directly using TraceStore instead of DataFrameBuilder.
	// "can we crqte method jsut by reading trace values limit 200 where rcommint number <= ... order by commit number somerthign like this"
	// We use ReadTracesForCommitRange which is more direct than DataFrameBuilder.

	startCommit := lastRegression.CommitNumber

	// Read traces for the range [startCommit, PrevCommitNumber].
	traceSet, commits, _, err := r.traceStore.ReadTracesForCommitRange(ctx, []string{traceName}, startCommit, cr.PrevCommitNumber)
	if err != nil {
		sklog.Errorf("[ImprovedAnomalyBoundsRefiner] Failed to read traces for range [%d, %d]: %s", startCommit, pickOffset, err)
		return cr
	}

	traceData, ok := traceSet[traceName]
	if !ok || len(traceData) < cfg.Radius {
		sklog.Errorf("[ImprovedAnomalyBoundsRefiner] Not enough data found for trace %s in range [%d, %d]", traceName, startCommit, cr.PrevCommitNumber)
		return cr
	}

	// 3. Extract left side data.
	// We take the points from the loaded trace.
	// We filter out missing data first, and then take the last 200 points.
	var leftData []float32
	var leftCommits []types.CommitNumber
	for i, v := range traceData {
		if v != vec32.MissingDataSentinel {
			leftData = append(leftData, v)
			leftCommits = append(leftCommits, commits[i].CommitNumber)
		}
	}

	// Limit to 200 points.
	if len(leftData) > 200 {
		leftData = leftData[len(leftData)-200:]
		leftCommits = leftCommits[len(leftCommits)-200:]
	}

	if len(leftData) < 3 {
		return cr
	}

	// 4. Extract right side data.
	tpIndex := cr.Summary.Clusters[0].StepFit.TurningPoint
	rightData := cr.Summary.Clusters[0].Centroid[tpIndex:]

	var cleanRightData []float32
	var rightCommits []types.CommitNumber
	for i := tpIndex; i < len(cr.Summary.Clusters[0].Centroid); i++ {
		v := cr.Summary.Clusters[0].Centroid[i]
		if v != vec32.MissingDataSentinel {
			cleanRightData = append(cleanRightData, v)
			rightCommits = append(rightCommits, cr.Frame.DataFrame.Header[i].Offset)
		}
	}
	rightData = cleanRightData

	if len(rightData) < 3 {
		return cr
	}

	// 5. Run anomaly detection (math check).
	y0 := vec32.Mean(leftData)
	y1 := vec32.Mean(rightData)
	s1 := vec32.StdDev(leftData, y0)
	s2 := vec32.StdDev(rightData, y1)
	n1 := len(leftData)
	n2 := len(rightData)

	var regressionVal float32
	var stepSize float32

	switch cfg.Step {
	case types.AbsoluteStep:
		stepSize, regressionVal = stepfit.CalcAbsoluteStep(y0, y1)
	case types.PercentStep:
		stepSize, regressionVal = stepfit.CalcPercentStep(y0, y1)
	case types.CohenStep:
		stepSize, regressionVal = stepfit.CalcValidCohenStep(y0, y1, s1, s2, n1, n2, r.stdDevThreshold)
	default:
		stepSize, regressionVal = stepfit.CalcValidCohenStep(y0, y1, s1, s2, n1, n2, r.stdDevThreshold)
	}

	interesting := cfg.Interesting
	if interesting == 0 {
		interesting = 2.0
	}

	isConfirmed := false
	if math.Abs(float64(regressionVal)) >= float64(interesting) {
		isConfirmed = true
	}

	if !isConfirmed {
		leftStart := leftCommits[0]
		leftEnd := leftCommits[len(leftCommits)-1]
		rightStart := rightCommits[0]
		rightEnd := rightCommits[len(rightCommits)-1]
		sklog.Infof("[ImprovedAnomalyBoundsRefiner] Filtering out regression for trace %s at offset %d. Failed strict check. RegressionVal: %f, Threshold: %f, Left(mean=%f, stddev=%f, n=%d, range=[%d, %d]), Right(mean=%f, stddev=%f, n=%d, range=[%d, %d]), Pick Range: [%d, %d]",
			traceName, pickOffset, regressionVal, interesting, y0, s1, n1, leftStart, leftEnd, y1, s2, n2, rightStart, rightEnd, cr.PrevCommitNumber, cr.CommitNumber)
		return nil
	}

	cr.Message = fmt.Sprintf("%s | Confirmed by ImprovedAnomalyBoundsRefiner with regression value: %f, step size: %f", cr.Message, regressionVal, stepSize)

	return cr
}
