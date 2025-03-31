package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid name - only letters",
			input:   "John",
			wantErr: nil,
		},
		{
			name:    "valid name - all lowercase letters",
			input:   "alice",
			wantErr: nil,
		},
		{
			name:    "valid name - all uppercase letters",
			input:   "BOB",
			wantErr: nil,
		},
		{
			name:    "invalid name - empty string",
			input:   "",
			wantErr: ErrEmptyName,
		},
		{
			name:    "valid name - contains letters and numbers",
			input:   "Alice123",
			wantErr: nil,
		},
		{
			name:    "invalid name - contains special characters",
			input:   "Bob@123",
			wantErr: ErrInvalidName,
		},
		{
			name:    "invalid name - contains spaces",
			input:   "John Doe",
			wantErr: ErrInvalidName,
		},
		{
			name:    "invalid name - contains hyphen",
			input:   "John-Doe",
			wantErr: ErrInvalidName,
		},
		{
			name:    "valid name - only one letter",
			input:   "A",
			wantErr: nil,
		},
		{
			name:    "valid name - contains letters and numbers",
			input:   "Alice123",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
