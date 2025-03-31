package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateType(t *testing.T) {
	tests := []struct {
		name    string
		mType   string
		wantErr error
	}{
		{
			name:    "valid type - counter",
			mType:   "counter",
			wantErr: nil,
		},
		{
			name:    "valid type - gauge",
			mType:   "gauge",
			wantErr: nil,
		},
		{
			name:    "invalid type - empty string",
			mType:   "",
			wantErr: ErrInvalidType,
		},
		{
			name:    "invalid type - unknown value",
			mType:   "unknown",
			wantErr: ErrInvalidType,
		},
		{
			name:    "valid type with mixed case - Counter",
			mType:   "Counter",
			wantErr: nil,
		},
		{
			name:    "valid type with mixed case - Gauge",
			mType:   "Gauge",
			wantErr: nil,
		},
		{
			name:    "invalid type with mixed case - SomeOtherType",
			mType:   "SomeOtherType",
			wantErr: ErrInvalidType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateType(tt.mType)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
