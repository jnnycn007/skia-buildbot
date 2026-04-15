package internal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.skia.org/infra/go/httputils"
	"go.skia.org/infra/go/metrics2"
)

func TestExtractErrorMessageReturnsErrorMessage(t *testing.T) {
	input := []byte(`{"error": "something went wrong"}`)
	expected := "something went wrong"
	actual := extractErrorMessage(input)
	assert.Equal(t, expected, actual)
}

func TestExtractErrorMessageReturnsRawString(t *testing.T) {
	input := []byte(`Internal Server Error`)
	expected := "Internal Server Error"
	actual := extractErrorMessage(input)
	assert.Equal(t, expected, actual)
}

func TestExtractErrorMessageReturnsRawStringWhenEmptyError(t *testing.T) {
	input := []byte(`{"error": ""}`)
	expected := `{"error": ""}`
	actual := extractErrorMessage(input)
	assert.Equal(t, expected, actual)
}

func TestExtractErrorMessageReturnsRawStringWhenNoError(t *testing.T) {
	input := []byte(`{"message": "some error"}`)
	expected := `{"message": "some error"}`
	actual := extractErrorMessage(input)
	assert.Equal(t, expected, actual)
}

func TestBuildBisectJobRequestUrlPopulatesAllFieldsForOldAnomaly(t *testing.T) {
	req := BisectJobCreateRequest{
		ComparisonMode:      "performance",
		StartGitHash:        "start_hash",
		EndGitHash:          "end_hash",
		Configuration:       "config",
		Benchmark:           "benchmark",
		Story:               "story",
		Chart:               "chart",
		Statistic:           "statistic",
		ComparisonMagnitude: "magnitude",
		Pin:                 "pin",
		Project:             "project",
		BugId:               "123",
		User:                "user",
		AlertIDs:            "456",
		TestPath:            "test_path",
	}

	builtURL := buildBisectJobRequestURL(req, false)
	assert.Contains(t, builtURL, chromeperfLegacyBisectURL)

	parsedURL, err := url.Parse(builtURL)
	assert.NoError(t, err)

	expected := url.Values{
		"comparison_mode":      []string{"performance"},
		"start_git_hash":       []string{"start_hash"},
		"end_git_hash":         []string{"end_hash"},
		"configuration":        []string{"config"},
		"benchmark":            []string{"benchmark"},
		"story":                []string{"story"},
		"chart":                []string{"chart"},
		"statistic":            []string{"statistic"},
		"comparison_magnitude": []string{"magnitude"},
		"pin":                  []string{"pin"},
		"project":              []string{"project"},
		"bug_id":               []string{"123"},
		"user":                 []string{"user"},
		"alert_ids":            []string{"456"},
		"test_path":            []string{"test_path"},
	}
	assert.Equal(t, expected, parsedURL.Query())
}

func TestBuildBisectJobRequestUrlPopulatesAllFieldsForNewAnomaly(t *testing.T) {
	req := BisectJobCreateRequest{
		ComparisonMode:      "performance",
		StartGitHash:        "start_hash",
		EndGitHash:          "end_hash",
		Configuration:       "config",
		Benchmark:           "benchmark",
		Story:               "story",
		Chart:               "chart",
		Statistic:           "statistic",
		ComparisonMagnitude: "magnitude",
		Pin:                 "pin",
		Project:             "project",
		BugId:               "123",
		User:                "user",
		AlertIDs:            "456",
		TestPath:            "test_path",
	}

	builtURL := buildBisectJobRequestURL(req, true)
	assert.Contains(t, builtURL, chromeperfLegacyBisectURL)

	parsedURL, err := url.Parse(builtURL)
	assert.NoError(t, err)

	// Alert IDs should not be present.
	expected := url.Values{
		"comparison_mode":      []string{"performance"},
		"start_git_hash":       []string{"start_hash"},
		"end_git_hash":         []string{"end_hash"},
		"configuration":        []string{"config"},
		"benchmark":            []string{"benchmark"},
		"story":                []string{"story"},
		"chart":                []string{"chart"},
		"statistic":            []string{"statistic"},
		"comparison_magnitude": []string{"magnitude"},
		"pin":                  []string{"pin"},
		"project":              []string{"project"},
		"bug_id":               []string{"123"},
		"user":                 []string{"user"},
		"test_path":            []string{"test_path"},
	}
	assert.Equal(t, expected, parsedURL.Query())
}

func TestBuildBisectJobRequestUrlPopulatesRequiredFields(t *testing.T) {
	req := BisectJobCreateRequest{}
	builtURL := buildBisectJobRequestURL(req, false)
	parsedURL, err := url.Parse(builtURL)
	assert.NoError(t, err)

	expected := url.Values{
		"bug_id":    []string{""},
		"test_path": []string{""},
	}
	assert.Equal(t, expected, parsedURL.Query())
}

type mockTransport struct {
	handler http.HandlerFunc
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	recorder := httptest.NewRecorder()
	m.handler(recorder, req)
	return recorder.Result(), nil
}

func setupTestMocks(t *testing.T, handler http.HandlerFunc) *LegacyClient {
	client := httputils.DefaultClientConfig().WithoutRetries().Client()
	client.Transport = &mockTransport{
		handler: handler,
	}

	pc := &LegacyClient{
		httpClient:          client,
		createBisectCalled:  metrics2.GetCounter("pinpoint_create_bisect_called"),
		createBisectFailed:  metrics2.GetCounter("pinpoint_create_bisect_failed"),
		createTryJobCalled:  metrics2.GetCounter("pinpoint_create_try_job_called"),
		createTryJobFailed:  metrics2.GetCounter("pinpoint_create_try_job_failed"),
		fetchJobStateCalled: metrics2.GetCounter("pinpoint_fetch_job_state_called"),
		fetchJobStateFailed: metrics2.GetCounter("pinpoint_fetch_job_state_failed"),
	}

	return pc
}

func TestDoPostRequest(t *testing.T) {
	t.Run("Returns parsed response on success", func(t *testing.T) {
		expectedResponseBody := []byte(
			`{"jobId": "12345", "jobUrl": "https://example.com/job/12345"}`,
		)
		pc := setupTestMocks(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/new", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(expectedResponseBody)
		})

		resp, err := pc.doPostRequest(context.Background(), pinpointLegacyURL)
		assert.NoError(t, err)
		body, err := pc.readResponseBody(resp)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponseBody, body)
	})

	t.Run("Returns error if non-200 status code", func(t *testing.T) {
		pc := setupTestMocks(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal Server Error"}`))
		})

		resp, err := pc.doPostRequest(context.Background(), pinpointLegacyURL)
		assert.NoError(t, err)
		body, err := pc.readResponseBody(resp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Internal Server Error")
		assert.Nil(t, body)
	})
}

func TestDoGetRequest(t *testing.T) {
	t.Run("Returns parsed response on success", func(t *testing.T) {
		expectedResponseBody := []byte(`{"job_id": "12345", "status": "completed"}`)
		testURL := "https://example.com/api/job/12345?o=STATE"
		pc := setupTestMocks(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/api/job/12345", r.URL.Path)
			assert.Equal(t, "STATE", r.URL.Query().Get("o"))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(expectedResponseBody)
		})

		resp, err := pc.doGetRequest(context.Background(), testURL)
		assert.NoError(t, err)
		body, err := pc.readResponseBody(resp)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponseBody, body)
	})

	t.Run("Returns error if non-200 status code", func(t *testing.T) {
		testURL := "https://example.com/api/job/12345?o=STATE"
		pc := setupTestMocks(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal Server Error"}`))
		})

		resp, err := pc.doGetRequest(context.Background(), testURL)
		assert.NoError(t, err)
		body, err := pc.readResponseBody(resp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Internal Server Error")
		assert.Nil(t, body)
	})
}

func TestBuildTryJobRequestUrlPopulatesRequiredFields(t *testing.T) {
	req := TryJobCreateRequest{
		Name:            "test-job",
		BaseGitHash:     "base_hash",
		EndGitHash:      "end_hash",
		BasePatch:       "base_patch",
		ExperimentPatch: "experiment_patch",
		Configuration:   "config",
		Benchmark:       "benchmark",
		Story:           "story",
		ExtraTestArgs:   "args",
		Repository:      "repo",
		BugId:           "123",
		User:            "user",
	}

	urlStr, err := buildTryJobRequestURL(req)
	assert.NoError(t, err)

	parsedURL, err := url.Parse(urlStr)
	assert.NoError(t, err)

	expected := url.Values{
		"comparison_mode":  []string{tryJobComparisonMode},
		"name":             []string{"test-job"},
		"base_git_hash":    []string{"base_hash"},
		"end_git_hash":     []string{"end_hash"},
		"base_patch":       []string{"base_patch"},
		"experiment_patch": []string{"experiment_patch"},
		"configuration":    []string{"config"},
		"benchmark":        []string{"benchmark"},
		"story":            []string{"story"},
		"extra_test_args":  []string{"args"},
		"repository":       []string{"repo"},
		"bug_id":           []string{"123"},
		"user":             []string{"user"},
		"tags":             []string{"{\"origin\":\"skia_perf\"}"},
	}
	assert.Equal(t, expected, parsedURL.Query())
}

func TestBuildTryJobRequestUrlVerifiesMissingBenchmark(t *testing.T) {
	req := TryJobCreateRequest{
		Configuration: "config",
	}
	_, err := buildTryJobRequestURL(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Benchmark must be specified")
}

func TestBuildTryJobRequestUrlVerifiesMissingConfiguration(t *testing.T) {
	req := TryJobCreateRequest{
		Benchmark: "benchmark",
	}
	_, err := buildTryJobRequestURL(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Configuration must be specified")
}

func TestFetchJobState(t *testing.T) {
	t.Run("Returns parsed response on success", func(t *testing.T) {
		pc := setupTestMocks(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/job/12345", r.URL.Path)
			assert.Equal(t, "STATE", r.URL.Query().Get("o"))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"job_id": "12345", "status": "completed"}`))
		})

		resp, err := pc.FetchJobState(context.Background(), FetchJobStateRequest{JobID: "12345"})
		assert.NoError(t, err)
		assert.Equal(
			t,
			resp,
			&FetchJobStateResponse{
				JobID:  "12345",
				Status: "completed",
			},
		)
	})

	t.Run("Returns error if non-200 status code", func(t *testing.T) {
		pc := setupTestMocks(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal Server Error"}`))
		})

		resp, err := pc.FetchJobState(context.Background(), FetchJobStateRequest{JobID: "12345"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Internal Server Error")
		assert.Nil(t, resp)
	})
}
