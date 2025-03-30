package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToInt64(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *int64
		wantErr  bool
	}{
		{
			name:     "ValidInteger",
			input:    "123456",
			expected: func() *int64 { v := int64(123456); return &v }(),
			wantErr:  false,
		},
		{
			name:     "ZeroValue",
			input:    "0",
			expected: func() *int64 { v := int64(0); return &v }(),
			wantErr:  false,
		},
		{
			name:     "NegativeInteger",
			input:    "-98765",
			expected: func() *int64 { v := int64(-98765); return &v }(),
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
			result, err := ConvertToInt64(tc.input)
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
