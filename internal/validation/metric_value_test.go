package validation

import (
	"go-metrics/internal/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCounterValue(t *testing.T) {
	tests := []struct {
		value     string
		expectErr error
	}{
		{"123", nil},
		{"-123", nil},
		{"abc", errors.ErrInvalidCounterMetricValue},
		{"123.45", errors.ErrInvalidCounterMetricValue},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			err := ValidateCounterValue(tt.value)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}

func TestValidateGaugeValue(t *testing.T) {
	tests := []struct {
		value     string
		expectErr error
	}{
		{"123.45", nil},
		{"-123.45", nil},
		{"123", nil},
		{"abc", errors.ErrInvalidGaugeMetricValue},
		{"123.45.67", errors.ErrInvalidGaugeMetricValue},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			err := ValidateGaugeValue(tt.value)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}

func TestValidateCounterPtrValue(t *testing.T) {
	tests := []struct {
		value     *int64
		expectErr error
	}{
		{nil, errors.ErrInvalidCounterMetricValue},
		{new(int64), nil},
	}

	for _, tt := range tests {
		t.Run("ValidateCounterPtrValue", func(t *testing.T) {
			err := ValidateCounterPtrValue(tt.value)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}

func TestValidateGaugePtrValue(t *testing.T) {
	tests := []struct {
		value     *float64
		expectErr error
	}{
		{nil, errors.ErrInvalidGaugeMetricValue},
		{new(float64), nil},
	}

	for _, tt := range tests {
		t.Run("ValidateGaugePtrValue", func(t *testing.T) {
			err := ValidateGaugePtrValue(tt.value)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}
