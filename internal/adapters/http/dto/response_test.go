package dto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccess(t *testing.T) {
	resp := Success(200, map[string]string{"id": "1"})

	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, "application/json", resp.Headers["Content-Type"])

	var body StandardResponse
	require.NoError(t, json.Unmarshal([]byte(resp.Body), &body))
	require.True(t, body.Success)
	require.Nil(t, body.Error)
}

func TestFail(t *testing.T) {
	resp := Fail(400, "bad_request", "invalid input")

	require.Equal(t, 400, resp.StatusCode)
	require.Equal(t, "application/json", resp.Headers["Content-Type"])

	var body StandardResponse
	require.NoError(t, json.Unmarshal([]byte(resp.Body), &body))
	require.False(t, body.Success)
	require.NotNil(t, body.Error)
	require.Equal(t, "bad_request", body.Error.Code)
	require.Equal(t, "invalid input", body.Error.Message)
}

func TestJSONFallbackOnMarshalError(t *testing.T) {
	// A channel cannot be JSON-serialized, triggering the fallback branch.
	resp := JSON(200, StandardResponse{Data: make(chan int)})

	require.Equal(t, 500, resp.StatusCode)
	require.Equal(t, "application/json", resp.Headers["Content-Type"])
	require.Contains(t, resp.Body, "internal_error")
}
