package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
	"github.com/stretchr/testify/require"
)

func MakeRequest(
	t *testing.T,
	method string,
	url string,
	body interface{},
	token string,
) *httptest.ResponseRecorder {
	var bodyReader io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(jsonBody)
	}

	req := httptest.NewRequest(method, url, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := TestApp.Test(req)
	require.NoError(t, err)

	// Read and rewind the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	recorder := httptest.NewRecorder()
	recorder.Body = bytes.NewBuffer(bodyBytes)
	recorder.Code = resp.StatusCode

	return recorder
}

func ParseResponse(t *testing.T, recorder *httptest.ResponseRecorder) helpers.BaseResponse {
	// First check if response is JSON
	contentType := recorder.Header().Get("Content-Type")
	require.Contains(t, contentType, "application/json",
		"Expected JSON response, got: %s", contentType)

	var response helpers.BaseResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err,
		"Failed to parse JSON response: %s", recorder.Body.String())

	return response
}
