package validation

import (
	"go-metrics/internal/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMetricType(t *testing.T) {
	tests := []struct {
		metricType string
		expectErr  error
	}{
		{"gauge", nil},
		{"counter", nil},
		{"gauge123", errors.ErrInvalidMetricType},
		{"invalidType", errors.ErrInvalidMetricType},
		{"", errors.ErrInvalidMetricType},
		{"counter123", errors.ErrInvalidMetricType},
		{"Counter", errors.ErrInvalidMetricType},
		{"Gauge", errors.ErrInvalidMetricType},
	}

	for _, tt := range tests {
		t.Run(tt.metricType, func(t *testing.T) {
			err := ValidateMetricType(tt.metricType)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}
