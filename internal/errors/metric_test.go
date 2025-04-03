package errors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeMetricErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		statusCode int
		expected   string
	}{
		{
			name:       "ErrInvalidMetricID",
			err:        ErrInvalidMetricID,
			statusCode: http.StatusBadRequest,
			expected:   "invalid id: only letters and numbers are allowed",
		},
		{
			name:       "ErrInvalidMetricType",
			err:        ErrInvalidMetricType,
			statusCode: http.StatusBadRequest,
			expected:   "invalid Type: must be 'gauge' or 'counter'",
		},
		{
			name:       "ErrEmptyMetricValue",
			err:        ErrEmptyMetricValue,
			statusCode: http.StatusBadRequest,
			expected:   "invalid metric value: cannot be empty",
		},
		{
			name:       "ErrInvalidCounterMetricValue",
			err:        ErrInvalidCounterMetricValue,
			statusCode: http.StatusBadRequest,
			expected:   "invalid 'counter' metric value: must be int64",
		},
		{
			name:       "ErrInvalidGaugeMetricValue",
			err:        ErrInvalidGaugeMetricValue,
			statusCode: http.StatusBadRequest,
			expected:   "invalid 'gauge' metric value: must be float64",
		},
		{
			name:       "ErrMetricNotFound",
			err:        ErrMetricNotFound,
			statusCode: http.StatusNotFound,
			expected:   "metric not found",
		},
		{
			name:       "ErrMetricGetByIDInternal",
			err:        ErrMetricGetByIDInternal,
			statusCode: http.StatusInternalServerError,
			expected:   "internal error",
		},
		{
			name:       "Unknown error",
			err:        errors.New("some unknown error"),
			statusCode: http.StatusInternalServerError,
			expected:   "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			MakeMetricErrorResponse(rr, tt.err)
			assert.Equal(t, tt.statusCode, rr.Code)
			assert.Equal(t, tt.expected, strings.TrimRight(rr.Body.String(), "\n"))
		})
	}
}
