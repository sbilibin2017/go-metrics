package unitofworks

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryUnitOfWork_Do(t *testing.T) {
	uow := NewMemoryUnitOfWork()
	operation := func(tx *sql.Tx) error {
		assert.Nil(t, tx)
		return nil
	}
	err := uow.Do(context.Background(), operation)
	assert.NoError(t, err)
}
