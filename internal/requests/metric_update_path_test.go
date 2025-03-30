package requests

import (
	"go-metrics/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricUpdatePathRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		request   MetricUpdatePathRequest
		expectErr bool
	}{
		{
			name: "invalid type (invalid value for Type)",
			request: MetricUpdatePathRequest{
				Type:  "invalid",
				Name:  "metric",
				Value: "100",
			},
			expectErr: true,
		},
		{
			name: "valid Counter type",
			request: MetricUpdatePathRequest{
				Type:  domain.Counter,
				Name:  "metric",
				Value: "100",
			},
			expectErr: false,
		},
		{
			name: "invalid name (empty name)",
			request: MetricUpdatePathRequest{
				Type:  domain.Counter,
				Name:  "",
				Value: "100",
			},
			expectErr: true,
		},
		{
			name: "valid name with Counter type",
			request: MetricUpdatePathRequest{
				Type:  domain.Counter,
				Name:  "metric",
				Value: "100",
			},
			expectErr: false,
		},
		{
			name: "invalid value (non-numeric value for Counter)",
			request: MetricUpdatePathRequest{
				Type:  domain.Counter,
				Name:  "metric",
				Value: "not-a-num",
			},
			expectErr: true,
		},
		{
			name: "invalid value (non-numeric value for Gauge)",
			request: MetricUpdatePathRequest{
				Type:  domain.Gauge,
				Name:  "metric",
				Value: "not-a-num",
			},
			expectErr: true,
		},
		{
			name: "empty value for Counter",
			request: MetricUpdatePathRequest{
				Type:  domain.Counter,
				Name:  "metric",
				Value: "",
			},
			expectErr: true,
		},
		{
			name: "empty value for Gauge",
			request: MetricUpdatePathRequest{
				Type:  domain.Gauge,
				Name:  "metric",
				Value: "",
			},
			expectErr: true,
		},
		{
			name: "valid value for Counter (numeric)",
			request: MetricUpdatePathRequest{
				Type:  domain.Counter,
				Name:  "metric",
				Value: "100",
			},
			expectErr: false,
		},
		{
			name: "valid value for Gauge (float)",
			request: MetricUpdatePathRequest{
				Type:  domain.Gauge,
				Name:  "metric",
				Value: "10.5",
			},
			expectErr: false,
		},
		{
			name: "invalid value for Counter (non-integer)",
			request: MetricUpdatePathRequest{
				Type:  domain.Counter,
				Name:  "metric",
				Value: "10.5",
			},
			expectErr: true,
		},
		{
			name: "valid integer value for Counter",
			request: MetricUpdatePathRequest{
				Type:  domain.Counter,
				Name:  "metric",
				Value: "10",
			},
			expectErr: false,
		},
		{
			name: "invalid value for Gauge (non-float)",
			request: MetricUpdatePathRequest{
				Type:  domain.Gauge,
				Name:  "metric",
				Value: "not-a-num",
			},
			expectErr: true,
		},
		{
			name: "valid float value for Gauge",
			request: MetricUpdatePathRequest{
				Type:  domain.Gauge,
				Name:  "metric",
				Value: "10.5",
			},
			expectErr: false,
		},
		{
			name: "invalid value for second ValidateValue call (non-numeric)",
			request: MetricUpdatePathRequest{
				Type:  domain.Counter,
				Name:  "metric",
				Value: "non-numeric",
			},
			expectErr: true,
		},
		{
			name: "valid value for second ValidateValue call",
			request: MetricUpdatePathRequest{
				Type:  domain.Counter,
				Name:  "metric",
				Value: "100",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMetricUpdatePathRequest_ToDomain(t *testing.T) {
	delta100 := int64(100)
	value10_5 := float64(10.5)
	tests := []struct {
		name      string
		request   MetricUpdatePathRequest
		expected  *domain.Metric
		expectErr bool
	}{
		{
			name: "valid Counter type with integer value",
			request: MetricUpdatePathRequest{
				Type:  domain.Counter,
				Name:  "metric1",
				Value: "100",
			},
			expected: &domain.Metric{
				ID:    "metric1",
				Type:  domain.Counter,
				Delta: &delta100,
			},
			expectErr: false,
		},
		{
			name: "valid Gauge type with float value",
			request: MetricUpdatePathRequest{
				Type:  domain.Gauge,
				Name:  "metric2",
				Value: "10.5",
			},
			expected: &domain.Metric{
				ID:    "metric2",
				Type:  domain.Gauge,
				Value: &value10_5,
			},
			expectErr: false,
		},
		{
			name: "valid Gauge type with non-numeric value (should be converted)",
			request: MetricUpdatePathRequest{
				Type:  domain.Gauge,
				Name:  "metric3",
				Value: "10.5",
			},
			expected: &domain.Metric{
				ID:    "metric3",
				Type:  domain.Gauge,
				Value: &value10_5,
			},
			expectErr: false,
		},
		{
			name: "valid Counter type with invalid value (non-integer)",
			request: MetricUpdatePathRequest{
				Type:  domain.Counter,
				Name:  "metric4",
				Value: "invalid_value",
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "valid Counter type with non-numeric value (should return error)",
			request: MetricUpdatePathRequest{
				Type:  domain.Counter,
				Name:  "metric5",
				Value: "invalid_value",
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "valid Gauge type with invalid float value",
			request: MetricUpdatePathRequest{
				Type:  domain.Gauge,
				Name:  "metric6",
				Value: "not_a_float",
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric, err := tt.request.ToDomain()

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, metric)
			}
		})
	}
}
