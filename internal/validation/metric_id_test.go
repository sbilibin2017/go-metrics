package validation

import (
	"go-metrics/internal/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateID(t *testing.T) {
	tests := []struct {
		id        string
		expectErr error
	}{
		{"", errors.ErrInvalidMetricID},
		{"valid123", nil},
		{"invalid_id!", errors.ErrInvalidMetricID},
		{"123456", nil},
		{"validID", nil},
		{"@invalidID", errors.ErrInvalidMetricID},
		{"another-valid123", errors.ErrInvalidMetricID},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			err := ValidateMetricID(tt.id)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}
