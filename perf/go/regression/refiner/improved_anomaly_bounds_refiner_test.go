package refiner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.skia.org/infra/perf/go/alerts"
	"go.skia.org/infra/perf/go/clustering2"
	"go.skia.org/infra/perf/go/dataframe"
	"go.skia.org/infra/perf/go/git/provider"
	"go.skia.org/infra/perf/go/regression"
	regression_mocks "go.skia.org/infra/perf/go/regression/mocks"
	"go.skia.org/infra/perf/go/stepfit"
	tracestore_mocks "go.skia.org/infra/perf/go/tracestore/mocks"
	"go.skia.org/infra/perf/go/types"
	"go.skia.org/infra/perf/go/ui/frame"
)

func TestApplyImprovedLogic_NoPrevRegression_ReturnsOriginal(t *testing.T) {
	mockStore := &regression_mocks.Store{}
	r := &ImprovedAnomalyBoundsRefiner{store: mockStore}

	cr := &regression.ConfirmedRegression{
		DisplayCommitNumber: 100,
		Summary: &clustering2.ClusterSummaries{
			Clusters: []*clustering2.ClusterSummary{
				{Keys: []string{"trace1"}},
			},
		},
	}

	mockStore.On("GetRegressionsBefore", mock.Anything, "trace1", "", types.CommitNumber(100), 1).Return([]*regression.Regression{}, nil)

	res, err := r.applyImprovedLogic(context.Background(), cr, &alerts.Alert{}, nil)

	assert.NoError(t, err)
	assert.Equal(t, cr, res)
}

func TestApplyImprovedLogic_OverlapWithDB_FiltersOut(t *testing.T) {
	mockStore := &regression_mocks.Store{}
	r := &ImprovedAnomalyBoundsRefiner{store: mockStore}

	cr := &regression.ConfirmedRegression{
		DisplayCommitNumber: 100,
		PrevCommitNumber:    90,
		CommitNumber:        100,
		Summary: &clustering2.ClusterSummaries{
			Clusters: []*clustering2.ClusterSummary{
				{Keys: []string{"trace1"}},
			},
		},
	}

	dbReg := &regression.Regression{
		CommitNumber:     95,
		PrevCommitNumber: 85,
	}
	mockStore.On("GetRegressionsBefore", mock.Anything, "trace1", "", types.CommitNumber(100), 1).Return([]*regression.Regression{dbReg}, nil)

	res, err := r.applyImprovedLogic(context.Background(), cr, &alerts.Alert{}, nil)

	assert.NoError(t, err)
	assert.Nil(t, res) // Filtered out
}

func TestApplyImprovedLogic_OverlapWithInMemory_FiltersOut(t *testing.T) {
	mockStore := &regression_mocks.Store{}
	r := &ImprovedAnomalyBoundsRefiner{store: mockStore}

	cr := &regression.ConfirmedRegression{
		DisplayCommitNumber: 100,
		PrevCommitNumber:    90,
		CommitNumber:        100,
		Summary: &clustering2.ClusterSummaries{
			Clusters: []*clustering2.ClusterSummary{
				{Keys: []string{"trace1"}},
			},
		},
	}

	latestRefined := &regression.ConfirmedRegression{
		CommitNumber:     95,
		PrevCommitNumber: 85,
	}

	// Mock DB to return older regression
	dbReg := &regression.Regression{
		CommitNumber:     50,
		PrevCommitNumber: 40,
	}
	mockStore.On("GetRegressionsBefore", mock.Anything, "trace1", "", types.CommitNumber(100), 1).Return([]*regression.Regression{dbReg}, nil)

	res, err := r.applyImprovedLogic(context.Background(), cr, &alerts.Alert{}, latestRefined)

	assert.NoError(t, err)
	assert.Nil(t, res) // Filtered out
}

func TestApplyImprovedLogic_DryRun_SkipsDB(t *testing.T) {
	mockStore := &regression_mocks.Store{}
	r := &ImprovedAnomalyBoundsRefiner{store: mockStore}

	cr := &regression.ConfirmedRegression{
		DisplayCommitNumber: 100,
		Summary: &clustering2.ClusterSummaries{
			Clusters: []*clustering2.ClusterSummary{
				{Keys: []string{"trace1"}},
			},
		},
	}

	ctx := regression.WithDryRun(context.Background())

	res, err := r.applyImprovedLogic(ctx, cr, &alerts.Alert{}, nil)

	assert.NoError(t, err)
	assert.Equal(t, cr, res) // Fallback because no prev regression found in memory either
}

func TestApplyImprovedLogic_CohenThreshold_Capped(t *testing.T) {
	mockStore := &regression_mocks.Store{}
	mockTraceStore := &tracestore_mocks.TraceStore{}
	r := &ImprovedAnomalyBoundsRefiner{
		store:           mockStore,
		traceStore:      mockTraceStore,
		stdDevThreshold: 0.001,
	}

	cr := &regression.ConfirmedRegression{
		DisplayCommitNumber: 100,
		PrevCommitNumber:    90,
		CommitNumber:        100,
		Summary: &clustering2.ClusterSummaries{
			Clusters: []*clustering2.ClusterSummary{
				{
					Keys:     []string{"trace1"},
					StepFit:  &stepfit.StepFit{TurningPoint: 2},
					Centroid: []float32{10, 10, 20, 20},
				},
			},
		},
		RightMostSummary: &clustering2.ClusterSummaries{
			Clusters: []*clustering2.ClusterSummary{
				{
					Keys:     []string{"trace1"},
					StepFit:  &stepfit.StepFit{TurningPoint: 2},
					Centroid: []float32{10, 10, 20, 20},
				},
			},
		},
		RightMostFrame: &frame.FrameResponse{
			DataFrame: &dataframe.DataFrame{
				Header: []*dataframe.ColumnHeader{
					{Offset: 90}, {Offset: 91}, {Offset: 92}, {Offset: 93},
				},
			},
		},
	}

	dbReg := &regression.Regression{
		CommitNumber:     50,
		PrevCommitNumber: 40,
	}
	mockStore.On("GetRegressionsBefore", mock.Anything, "trace1", "", types.CommitNumber(100), 1).Return([]*regression.Regression{dbReg}, nil)

	mockTraceStore.On("ReadTracesForCommitRange", mock.Anything, []string{"trace1"}, types.CommitNumber(50), types.CommitNumber(90)).Return(
		types.TraceSet{"trace1": types.Trace{10, 10, 10, 10}},
		[]provider.Commit{{CommitNumber: 50}, {CommitNumber: 60}, {CommitNumber: 70}, {CommitNumber: 80}},
		nil,
		nil,
	)

	alert := &alerts.Alert{
		Step:        types.CohenStep,
		Interesting: 2.0, // Alert wants 2.0, but we cap at 1.2 in refiner!
		Radius:      2,
	}

	res, err := r.applyImprovedLogic(context.Background(), cr, alert, nil)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}
