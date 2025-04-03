package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected *float64
		err      error
	}{
		{
			name:     "valid float64",
			value:    "123.45",
			expected: ptrFloat64(123.45),
			err:      nil,
		},
		{
			name:     "invalid float64",
			value:    "abc",
			expected: nil,
			err:      ErrInvalidFloat64,
		},
		{
			name:     "valid zero",
			value:    "0",
			expected: ptrFloat64(0),
			err:      nil,
		},
		{
			name:     "valid negative float64",
			value:    "-123.45",
			expected: ptrFloat64(-123.45),
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertToFloat64(tt.value)
			assert.Equal(t, tt.err, err)
			if tt.expected != nil {
				assert.Equal(t, *tt.expected, *result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestFormatFloat64(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected string
	}{
		{
			name:     "positive float64",
			value:    123.45,
			expected: "123.45",
		},
		{
			name:     "negative float64",
			value:    -123.45,
			expected: "-123.45",
		},
		{
			name:     "zero",
			value:    0,
			expected: "0",
		},
		{
			name:     "large float64",
			value:    1.2345678901234567e+10,
			expected: "12345678901.234568",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFloat64(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func ptrFloat64(f float64) *float64 {
	return &f
}
