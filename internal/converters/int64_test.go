package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToInt64(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected *int64
		err      error
	}{
		{
			name:     "valid int64",
			value:    "12345",
			expected: ptrInt64(12345),
			err:      nil,
		},
		{
			name:     "invalid int64",
			value:    "abc",
			expected: nil,
			err:      ErrInvalidInt64,
		},
		{
			name:     "valid zero",
			value:    "0",
			expected: ptrInt64(0),
			err:      nil,
		},
		{
			name:     "valid negative int64",
			value:    "-12345",
			expected: ptrInt64(-12345),
			err:      nil,
		},
		{
			name:     "overflow int64",
			value:    "9223372036854775808",
			expected: nil,
			err:      ErrInvalidInt64,
		},
		{
			name:     "underflow int64",
			value:    "-9223372036854775809",
			expected: nil,
			err:      ErrInvalidInt64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertToInt64(tt.value)
			assert.Equal(t, tt.err, err)
			if tt.expected != nil {
				assert.Equal(t, *tt.expected, *result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestFormatInt64(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected string
	}{
		{
			name:     "positive int64",
			value:    12345,
			expected: "12345",
		},
		{
			name:     "negative int64",
			value:    -12345,
			expected: "-12345",
		},
		{
			name:     "zero",
			value:    0,
			expected: "0",
		},
		{
			name:     "max int64",
			value:    9223372036854775807,
			expected: "9223372036854775807",
		},
		{
			name:     "min int64",
			value:    -9223372036854775808,
			expected: "-9223372036854775808",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatInt64(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func ptrInt64(i int64) *int64 {
	return &i
}
