package pinpoint

import (
	"context"

	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/pinpoint/go/pinpoint/internal"
)

type Client struct {
	legacyClient *internal.LegacyClient
}

// New returns a new PinpointClient instance.
func New(ctx context.Context) (*Client, error) {
	legacyClient, err := internal.NewLegacyClient(ctx)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	return &Client{legacyClient: legacyClient}, nil
}

// CreateTryJob calls the legacy pinpoint API to create a try job.
func (c *Client) CreateTryJob(ctx context.Context, req TryJobCreateRequest) (*CreatePinpointResponse, error) {
	return c.legacyClient.CreateTryJob(ctx, req)
}

// CreateBisect calls pinpoint API to create bisect job.
func (c *Client) CreateBisect(ctx context.Context, req BisectJobCreateRequest, isNewAnomaly bool) (*CreatePinpointResponse, error) {
	return c.legacyClient.CreateBisect(ctx, req, isNewAnomaly)
}
