package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.skia.org/infra/go/httputils"
	"go.skia.org/infra/go/metrics2"
	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/go/sklog"

	"go.skia.org/infra/go/auth"
	"golang.org/x/oauth2/google"
)

const (
	pinpointLegacyURL         = "https://pinpoint-dot-chromeperf.appspot.com/api/new"
	contentType               = "application/json"
	tryJobComparisonMode      = "try"
	chromeperfLegacyBisectURL = "https://chromeperf.appspot.com/pinpoint/new/bisect"
)

type LegacyClient struct {
	httpClient          *http.Client
	createBisectCalled  metrics2.Counter
	createBisectFailed  metrics2.Counter
	createTryJobCalled  metrics2.Counter
	createTryJobFailed  metrics2.Counter
	fetchJobStateCalled metrics2.Counter
	fetchJobStateFailed metrics2.Counter
}

// New returns a new LegacyClient instance.
func NewLegacyClient(ctx context.Context) (*LegacyClient, error) {
	tokenSource, err := google.DefaultTokenSource(ctx, auth.ScopeUserinfoEmail)
	if err != nil {
		return nil, skerr.Wrapf(err, "Failed to create pinpoint client.")
	}

	client := httputils.DefaultClientConfig().WithTokenSource(tokenSource).Client()
	return &LegacyClient{
		httpClient:          client,
		createBisectCalled:  metrics2.GetCounter("pinpoint_create_bisect_called"),
		createBisectFailed:  metrics2.GetCounter("pinpoint_create_bisect_failed"),
		createTryJobCalled:  metrics2.GetCounter("pinpoint_create_try_job_called"),
		createTryJobFailed:  metrics2.GetCounter("pinpoint_create_try_job_failed"),
		fetchJobStateCalled: metrics2.GetCounter("pinpoint_fetch_job_state_called"),
		fetchJobStateFailed: metrics2.GetCounter("pinpoint_fetch_job_state_failed"),
	}, nil
}

// CreateTryJob calls the legacy pinpoint API to create a try job.
func (pc *LegacyClient) CreateTryJob(
	ctx context.Context,
	req TryJobCreateRequest,
) (resp *CreatePinpointResponse, err error) {
	pc.createTryJobCalled.Inc(1)
	defer func() { trackError(pc.createTryJobFailed, err) }()

	requestURL, err := buildTryJobRequestURL(req)
	if err != nil {
		return nil, skerr.Wrapf(err, "Failed to generate Pinpoint request URL.")
	}

	httpResp, err := pc.doPostRequest(ctx, requestURL)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	body, err := pc.readResponseBody(httpResp)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, skerr.Wrapf(err, "Failed to parse pinpoint response body.")
	}
	return resp, nil
}

// CreateBisect calls pinpoint API to create bisect job.
func (pc *LegacyClient) CreateBisect(
	ctx context.Context,
	req BisectJobCreateRequest,
	isNewAnomaly bool,
) (resp *CreatePinpointResponse, err error) {
	pc.createBisectCalled.Inc(1)
	defer func() { trackError(pc.createBisectFailed, err) }()

	requestURL := buildBisectJobRequestURL(req, isNewAnomaly)
	httpResp, err := pc.doPostRequest(ctx, requestURL)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	body, err := pc.readResponseBody(httpResp)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, skerr.Wrapf(err, "Failed to parse pinpoint response body.")
	}
	return resp, nil
}

// FetchJobState queries the legacy pinpoint API to retrieve job details.
func (pc *LegacyClient) FetchJobState(
	ctx context.Context,
	req FetchJobStateRequest,
) (resp *FetchJobStateResponse, err error) {
	pc.fetchJobStateCalled.Inc(1)
	defer func() { trackError(pc.fetchJobStateFailed, err) }()

	requestURL := fmt.Sprintf(
		"https://pinpoint-dot-chromeperf.appspot.com/api/job/%s?o=STATE",
		url.PathEscape(req.JobID),
	)
	httpResp, err := pc.doGetRequest(ctx, requestURL)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	body, err := pc.readResponseBody(httpResp)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, skerr.Wrapf(err, "Failed to parse pinpoint response body.")
	}
	return resp, err
}

func trackError(counter metrics2.Counter, err error) {
	if err != nil {
		counter.Inc(1)
	}
}

func buildTryJobRequestURL(req TryJobCreateRequest) (string, error) {
	if req.Benchmark == "" {
		return "", skerr.Fmt("Benchmark must be specified but is empty.")
	}
	if req.Configuration == "" {
		return "", skerr.Fmt("Configuration must be specified but is empty.")
	}

	params := url.Values{}
	// Pinpoint try jobs always use comparison mode try
	params.Set("comparison_mode", tryJobComparisonMode)
	setIfNotEmpty(params, "name", req.Name)
	setIfNotEmpty(params, "base_git_hash", req.BaseGitHash)
	setIfNotEmpty(params, "end_git_hash", req.EndGitHash)
	setIfNotEmpty(params, "base_patch", req.BasePatch)
	setIfNotEmpty(params, "experiment_patch", req.ExperimentPatch)
	setIfNotEmpty(params, "configuration", req.Configuration)
	setIfNotEmpty(params, "benchmark", req.Benchmark)
	setIfNotEmpty(params, "story", req.Story)
	setIfNotEmpty(params, "extra_test_args", req.ExtraTestArgs)
	setIfNotEmpty(params, "repository", req.Repository)
	setIfNotEmpty(params, "bug_id", req.BugId)
	setIfNotEmpty(params, "user", req.User)
	params.Set("tags", "{\"origin\":\"skia_perf\"}")

	return fmt.Sprintf("%s?%s", pinpointLegacyURL, params.Encode()), nil
}

func buildBisectJobRequestURL(req BisectJobCreateRequest, isNewAnomaly bool) string {
	params := url.Values{}
	setIfNotEmpty(params, "comparison_mode", req.ComparisonMode)
	setIfNotEmpty(params, "start_git_hash", req.StartGitHash)
	setIfNotEmpty(params, "end_git_hash", req.EndGitHash)
	setIfNotEmpty(params, "configuration", req.Configuration)
	setIfNotEmpty(params, "benchmark", req.Benchmark)
	setIfNotEmpty(params, "story", req.Story)
	setIfNotEmpty(params, "chart", req.Chart)
	setIfNotEmpty(params, "statistic", req.Statistic)
	setIfNotEmpty(params, "comparison_magnitude", req.ComparisonMagnitude)
	setIfNotEmpty(params, "pin", req.Pin)
	setIfNotEmpty(params, "project", req.Project)
	setIfNotEmpty(params, "user", req.User)
	if !isNewAnomaly {
		setIfNotEmpty(params, "alert_ids", req.AlertIDs)
	}
	// Bug ID must present otherwise chromeperf returns an error.
	params.Set("bug_id", req.BugId)
	params.Set("test_path", req.TestPath)
	return fmt.Sprintf("%s?%s", chromeperfLegacyBisectURL, params.Encode())
}

func extractErrorMessage(responseBody []byte) string {
	var errorResponse struct {
		Error string `json:"error"`
	}
	err := json.Unmarshal(responseBody, &errorResponse)
	if err == nil && errorResponse.Error != "" {
		return errorResponse.Error
	}
	return string(responseBody)
}

func setIfNotEmpty(params url.Values, key, value string) {
	if value != "" {
		params.Set(key, value)
	}
}

func (pc *LegacyClient) doPostRequest(
	ctx context.Context,
	requestURL string,
) (*http.Response, error) {
	sklog.Debugf("Preparing to send a Pinpoint POST request to: %s", requestURL)
	resp, err := httputils.PostWithContext(ctx, pc.httpClient, requestURL, contentType, nil)
	if err != nil {
		return nil, skerr.Wrapf(err, "Failed to get pinpoint response.")
	}
	sklog.Debugf("Got response from Pinpoint service: %+v", resp)
	return resp, nil
}

func (pc *LegacyClient) doGetRequest(
	ctx context.Context,
	requestURL string,
) (*http.Response, error) {
	sklog.Debugf("Preparing to send a Pinpoint GET request to: %s", requestURL)
	resp, err := httputils.GetWithContext(ctx, pc.httpClient, requestURL)
	if err != nil {
		return nil, skerr.Wrapf(err, "Failed to get pinpoint response.")
	}
	sklog.Debugf("Got response from Pinpoint service: %+v", resp)
	return resp, nil
}

func (pc *LegacyClient) readResponseBody(
	resp *http.Response,
) (body []byte, err error) {
	defer resp.Body.Close()
	if body, err = io.ReadAll(resp.Body); err != nil {
		return nil, skerr.Wrapf(err, "Failed to read body from pinpoint response.")
	}
	if resp.StatusCode != http.StatusOK {
		requestErrorMessage := extractErrorMessage(body)

		// A response must contain a request with a URL. Condition is just to make
		// sure we never panic here.
		url := "Unknown URL"
		if resp.Request != nil && resp.Request.URL != nil {
			url = resp.Request.URL.String()
		}
		errMsg := fmt.Sprintf(
			"Request to %s failed with status code %d and error: %s",
			url,
			resp.StatusCode,
			requestErrorMessage,
		)
		return nil, skerr.Wrap(errors.New(errMsg))
	}

	return body, err
}
