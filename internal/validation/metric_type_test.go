package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateType(t *testing.T) {
	tests := []struct {
		name    string
		Type    string
		wantErr error
	}{
		{
			name:    "valid type - counter",
			Type:    "counter",
			wantErr: nil,
		},
		{
			name:    "valid type - gauge",
			Type:    "gauge",
			wantErr: nil,
		},
		{
			name:    "invalid type - empty string",
			Type:    "",
			wantErr: ErrInvalidType,
		},
		{
			name:    "invalid type - unknown value",
			Type:    "unknown",
			wantErr: ErrInvalidType,
		},
		{
			name:    "valid type with mixed case - Counter",
			Type:    "Counter",
			wantErr: nil,
		},
		{
			name:    "valid type with mixed case - Gauge",
			Type:    "Gauge",
			wantErr: nil,
		},
		{
			name:    "invalid type with mixed case - SomeOtherType",
			Type:    "SomeOtherType",
			wantErr: ErrInvalidType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateType(tt.Type)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
