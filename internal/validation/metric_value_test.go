package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "Empty value",
			input:    "",
			expected: ErrEmptyValue,
		},
		{
			name:     "Valid value",
			input:    "12345",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateValue(tt.input)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestValidateCounterValueString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "Empty value",
			input:    "",
			expected: ErrInvalidCounterValue,
		},
		{
			name:     "Valid counter value '123'",
			input:    "123",
			expected: nil,
		},
		{
			name:     "Invalid counter value '12.34'",
			input:    "12.34",
			expected: ErrInvalidCounterValue,
		},
		{
			name:     "Invalid counter value 'abc'",
			input:    "abc",
			expected: ErrInvalidCounterValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCounterValueString(tt.input)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestValidateGaugeValueString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "Empty value",
			input:    "",
			expected: ErrInvalidGaugeValue,
		},
		{
			name:     "Valid gauge value '12.34'",
			input:    "12.34",
			expected: nil,
		},
		{
			name:     "Valid gauge value '123'",
			input:    "123",
			expected: nil,
		},
		{
			name:     "Invalid gauge value 'abc'",
			input:    "abc",
			expected: ErrInvalidGaugeValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGaugeValueString(tt.input)
			assert.Equal(t, tt.expected, err)
		})
	}
}
