package unitofwork

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryUnitOfWork_Do(t *testing.T) {
	uow := NewMemoryUnitOfWork()
	ctx := context.Background()

	tests := []struct {
		name          string
		operation     func() error
		expectedError error
	}{
		{
			name: "successful operation",
			operation: func() error {
				return nil
			},
			expectedError: nil,
		},
		{
			name: "failed operation",
			operation: func() error {
				return errors.New("operation failed")
			},
			expectedError: errors.New("operation failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uow.Do(ctx, tt.operation)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			}
		})
	}
}
