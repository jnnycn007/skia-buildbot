package refiner

import (
	"math"

	"go.skia.org/infra/go/vec32"
	"go.skia.org/infra/perf/go/alerts"
	"go.skia.org/infra/perf/go/stepfit"
	"go.skia.org/infra/perf/go/types"
)

const (
	defaultCohenDThreshold = 2.0
)

// isAnomaly determines if the given value is an anomaly compared to the baseline.
//
// Note: It is safe to use this only when you want to improve the accuracy of an anomaly range (refinement).
// Otherwise, it is typically not recommended to use this logic, as a single data point is usually not enough
// to make a decision about whether an anomaly happened or not.
func isAnomaly(val float32, baseline []float32, cfg *alerts.Alert, stdDevThreshold float32) bool {

	// 1. Prepare data
	baseline_mean := vec32.Mean(baseline)
	treatment_mean := val

	baseline_len := len(baseline)

	// Since we are analyzing a single treatment data point, we lack sufficient
	// data to calculate its true distribution statistics. Because some anomaly
	// detection algorithms (like Cohen's d) require the sample size and standard deviation,
	// we simulate them by making the assumption that the median/mean may change,
	// but the other statistical properties of the dataset (size and stddev) remain
	// identical to the baseline data.
	treatment_len := baseline_len

	baseline_stddev := vec32.StdDev(baseline, baseline_mean)
	treatment_stddev := baseline_stddev

	algo := cfg.Step
	interesting := cfg.Interesting

	var regression float32

	switch algo {
	case types.AbsoluteStep:
		_, regression = stepfit.CalcAbsoluteStep(baseline_mean, treatment_mean)
	case types.Const:
		_, regression = stepfit.CalcConstStep(val, interesting)
	case types.PercentStep:
		_, regression = stepfit.CalcPercentStep(baseline_mean, treatment_mean)
	case types.CohenStep:
		_, regression = stepfit.CalcCohenStep(baseline_mean, treatment_mean, baseline_stddev, treatment_stddev, baseline_len, treatment_len, stdDevThreshold)
	default:
		// Other algorithms
		_, regression = stepfit.CalcCohenStep(baseline_mean, treatment_mean, baseline_stddev, treatment_stddev, baseline_len, treatment_len, stdDevThreshold)
		interesting = defaultCohenDThreshold
	}

	// Check against interesting threshold
	if math.Abs(float64(regression)) >= float64(interesting) {
		return true
	}
	return false
}
