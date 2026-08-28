package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSON_Success(t *testing.T) {
	w := httptest.NewRecorder()

	JSON(w, http.StatusOK, HTTPResponse{
		Success: true,
		Message: "test message",
		Data:    map[string]string{"key": "value"},
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp HTTPResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "test message", resp.Message)
}

func TestJSON_Error(t *testing.T) {
	w := httptest.NewRecorder()

	JSON(w, http.StatusBadRequest, HTTPResponse{
		Success: false,
		Error:   "bad request",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp HTTPResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "bad request", resp.Error)
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()

	Error(w, http.StatusUnauthorized, "unauthorized")

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp HTTPResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "unauthorized", resp.Error)
}

func TestJSON_EmptyData(t *testing.T) {
	w := httptest.NewRecorder()

	JSON(w, http.StatusOK, HTTPResponse{
		Success: true,
	})

	assert.Equal(t, http.StatusOK, w.Code)

	var resp HTTPResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Empty(t, resp.Message)
}
