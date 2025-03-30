package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToFloat64(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *float64
		wantErr  bool
	}{
		{
			name:     "ValidFloat",
			input:    "123.456",
			expected: func() *float64 { v := 123.456; return &v }(),
			wantErr:  false,
		},
		{
			name:     "ValidInteger",
			input:    "789",
			expected: func() *float64 { v := 789.0; return &v }(),
			wantErr:  false,
		},
		{
			name:     "InvalidString",
			input:    "abc",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "EmptyString",
			input:    "",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "WhitespaceString",
			input:    "  ",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ConvertToFloat64(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, *tc.expected, *result)
			}
		})
	}
}
