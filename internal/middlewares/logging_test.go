package middlewares

import (
	"go-metrics/pkg/log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggingMiddleware(t *testing.T) {
	err := log.Init(log.LevelInfo)
	require.NoError(t, err)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})
	middleware := LoggingMiddleware(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "test response")

}
