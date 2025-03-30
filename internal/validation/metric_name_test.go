package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "Empty name",
			input:    "",
			expected: ErrEmptyName,
		},
		{
			name:     "Valid name",
			input:    "John",
			expected: nil,
		},
		{
			name:     "Name with numbers",
			input:    "John123",
			expected: ErrInvalidName,
		},
		{
			name:     "Name with special characters",
			input:    "John_Doe",
			expected: ErrInvalidName,
		},
		{
			name:     "Name with mixed case",
			input:    "john",
			expected: nil,
		},
		{
			name:     "Name with spaces",
			input:    "John Doe",
			expected: ErrInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			assert.Equal(t, tt.expected, err)
		})
	}
}
