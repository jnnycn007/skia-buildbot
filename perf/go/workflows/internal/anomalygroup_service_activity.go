package internal

import (
	"context"
	"time"

	"go.skia.org/infra/go/skerr"
	pb "go.skia.org/infra/perf/go/anomalygroup/proto/v1"
	backend "go.skia.org/infra/perf/go/backend/client"

	// TODO(b/500974820): Replace `legacyPinpoint` with `pinpoint`.
	legacyPinpoint "go.skia.org/infra/pinpoint/go/pinpoint"
	"golang.org/x/time/rate"
)

type AnomalyGroupServiceActivity struct {
	insecure_conn             bool
	legacyPinpointRateLimiter *rate.Limiter
	useLegacyPinpoint         bool
	legacyPinpointClient      *legacyPinpoint.Client
}

func NewAnomalyGroupServiceActivity(client *legacyPinpoint.Client) *AnomalyGroupServiceActivity {
	return &AnomalyGroupServiceActivity{
		insecure_conn: false,
		// Protects legacy Pinpoint from overloading with bisection job requests.
		// Set to 1 hour for testing purposes.
		legacyPinpointRateLimiter: rate.NewLimiter(rate.Every(time.Hour), 1),
		useLegacyPinpoint:         true,
		legacyPinpointClient:      client,
	}
}

func (agsa *AnomalyGroupServiceActivity) CheckBisectionAllowed(ctx context.Context) (bool, error) {
	if agsa.legacyPinpointRateLimiter == nil {
		return false, skerr.Fmt("Legacy Pinpoint rate limiter is not initialized")
	}
	return agsa.legacyPinpointRateLimiter.Allow(), nil
}

func (agsa *AnomalyGroupServiceActivity) ShouldUseLegacyPinpoint(ctx context.Context) (bool, error) {
	return agsa.useLegacyPinpoint, nil
}

func (agsa *AnomalyGroupServiceActivity) LoadAnomalyGroupByID(ctx context.Context, anomalygroupServiceUrl string, req *pb.LoadAnomalyGroupByIDRequest) (*pb.LoadAnomalyGroupByIDResponse, error) {
	client, err := backend.NewAnomalyGroupServiceClient(anomalygroupServiceUrl, agsa.insecure_conn)
	if err != nil {
		return nil, err
	}
	resp, err := client.LoadAnomalyGroupByID(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (agsa *AnomalyGroupServiceActivity) FindTopAnomalies(ctx context.Context, anomalygroupServiceUrl string, req *pb.FindTopAnomaliesRequest) (*pb.FindTopAnomaliesResponse, error) {
	client, err := backend.NewAnomalyGroupServiceClient(anomalygroupServiceUrl, agsa.insecure_conn)
	if err != nil {
		return nil, err
	}
	resp, err := client.FindTopAnomalies(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (agsa *AnomalyGroupServiceActivity) UpdateAnomalyGroup(ctx context.Context, anomalygroupServiceUrl string, req *pb.UpdateAnomalyGroupRequest) (*pb.UpdateAnomalyGroupResponse, error) {
	client, err := backend.NewAnomalyGroupServiceClient(anomalygroupServiceUrl, agsa.insecure_conn)
	if err != nil {
		return nil, err
	}
	resp, err := client.UpdateAnomalyGroup(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (agsa *AnomalyGroupServiceActivity) CreateLegacyBisectJob(
	ctx context.Context,
	req *legacyPinpoint.BisectJobCreateRequest,
) (*legacyPinpoint.CreatePinpointResponse, error) {
	if !agsa.useLegacyPinpoint {
		return nil, skerr.Fmt("Legacy Pinpoint should not be used")
	}
	if agsa.legacyPinpointClient == nil {
		return nil, skerr.Fmt("legacyPinpointClient is not initialized")
	}
	isNewAnomaly := true
	return agsa.legacyPinpointClient.CreateBisect(ctx, *req, isNewAnomaly)
}
