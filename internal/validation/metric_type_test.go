package validation

import (
	"go-metrics/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "Empty type",
			input:    "",
			expected: ErrInvalidType,
		},
		{
			name:     "Valid type 'gauge'",
			input:    domain.Gauge,
			expected: nil,
		},
		{
			name:     "Valid type 'counter'",
			input:    domain.Counter,
			expected: nil,
		},
		{
			name:     "Invalid type 'other'",
			input:    "other",
			expected: ErrInvalidType,
		},
		{
			name:     "Type with mixed case 'Gauge'",
			input:    "Gauge",
			expected: ErrInvalidType,
		},
		{
			name:     "Type with mixed case 'Counter'",
			input:    "Counter",
			expected: ErrInvalidType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateType(tt.input)
			assert.Equal(t, tt.expected, err)
		})
	}
}
