package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateValue(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{
			name:    "valid value",
			value:   "some value",
			wantErr: nil,
		},
		{
			name:    "empty value",
			value:   "",
			wantErr: ErrEmptyValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateValue(tt.value)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestValidateCounterValue(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{
			name:    "valid counter value",
			value:   "123",
			wantErr: nil,
		},
		{
			name:    "invalid counter value",
			value:   "not_a_number",
			wantErr: ErrInvalidCounterValue,
		},
		{
			name:    "empty counter value",
			value:   "",
			wantErr: ErrInvalidCounterValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCounterValue(tt.value)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestValidateGaugeValue(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{
			name:    "valid gauge value",
			value:   "12.34",
			wantErr: nil,
		},
		{
			name:    "invalid gauge value",
			value:   "not_a_float",
			wantErr: ErrInvalidGaugeValue,
		},
		{
			name:    "empty gauge value",
			value:   "",
			wantErr: ErrInvalidGaugeValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGaugeValue(tt.value)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
