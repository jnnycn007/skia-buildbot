package pinpoint

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPinpointClient_NoTargetNewPinpointArg_LegacyClient(t *testing.T) {
	args := map[string]any{}
	c := NewPinpointClient(args)

	assert.True(t, c.Url == LegacyPinpointUrl)
}

func TestPinpointClient_FalseTargetNewPinpointArg_LegacyClient(t *testing.T) {
	args := map[string]any{
		"target_new_pinpoint": false,
	}
	c := NewPinpointClient(args)

	assert.True(t, c.Url == LegacyPinpointUrl)
}

func TestPinpointClient_TargetNewPinpointArg_LegacyClient(t *testing.T) {
	args := map[string]any{
		"target_new_pinpoint": true,
	}
	c := NewPinpointClient(args)

	assert.True(t, c.Url == PinpointUrl)
}

func TestTryJob_TargetingNewPinpoint_Nil(t *testing.T) {
	args := map[string]any{
		"target_new_pinpoint": true,
	}
	ctx := context.Background()
	c := NewPinpointClient(args)
	resp, err := c.TryJob(ctx, nil)

	assert.Nil(t, err)
	assert.Nil(t, resp)
}

func TestTryJob_LegacyPinpoint_NotOK_Err(t *testing.T) {
	args := map[string]any{}
	ctx := context.Background()
	c := NewPinpointClient(args)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	c.Url = ts.URL
	defer ts.Close()

	resp, err := c.TryJob(ctx, ts.Client())

	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "failed with request 403")
}

func TestTryJob_LegacyPinpoint_OK_NoErr(t *testing.T) {
	args := map[string]any{}
	ctx := context.Background()
	c := NewPinpointClient(args)

	expectedJobId := "12345"
	expectedJobUrl := "http://some.url/job/12345"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := json.Marshal(map[string]string{
			"jobId":  expectedJobId,
			"jobUrl": expectedJobUrl,
		})
		require.NoError(t, err)
		w.Header().Add("Content-Type", "application/json")
		_, err = w.Write(content)
		require.NoError(t, err)
	}))
	c.Url = ts.URL
	defer ts.Close()

	resp, err := c.TryJob(ctx, ts.Client())

	require.NoError(t, err)
	assert.Equal(t, resp.JobID, expectedJobId)
	assert.Equal(t, resp.JobURL, expectedJobUrl)
}
