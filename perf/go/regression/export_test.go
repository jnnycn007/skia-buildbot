package regression

import (
	"context"

	"go.skia.org/infra/perf/go/alerts"
	"go.skia.org/infra/perf/go/dataframe"
	"go.skia.org/infra/perf/go/dfiter"
	"go.skia.org/infra/perf/go/progress"
)

// ExportedRegressionDetectionProcess allows black-box tests to access the unexported regressionDetectionProcess.
type ExportedRegressionDetectionProcess struct {
	p *regressionDetectionProcess
}

func NewExportedRegressionDetectionProcess(prog progress.Progress, iter dfiter.DataFrameIterator, alertConfig *alerts.Alert, refiner RegressionRefiner) *ExportedRegressionDetectionProcess {
	return &ExportedRegressionDetectionProcess{
		p: &regressionDetectionProcess{
			iter: iter,
			request: &RegressionDetectionRequest{
				Progress: prog,
				Alert:    alertConfig,
			},
			regressionRefiner: refiner,
		},
	}
}

func (e *ExportedRegressionDetectionProcess) DetectRegressionsOnDataFrame(ctx context.Context, df *dataframe.DataFrame) (*RegressionDetectionResponse, error) {
	return e.p.detectRegressionsOnDataFrame(ctx, df)
}
