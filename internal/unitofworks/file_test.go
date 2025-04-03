package unitofworks

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileUnitOfWork_Do(t *testing.T) {
	uow := NewFileUnitOfWork()
	operation := func(tx *sql.Tx) error {
		assert.Nil(t, tx)
		return nil
	}
	err := uow.Do(context.Background(), operation)
	assert.NoError(t, err)
}
