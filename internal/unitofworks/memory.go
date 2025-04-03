package unitofworks

import (
	"context"
	"database/sql"
)

type MemoryUnitOfWork struct{}

func NewMemoryUnitOfWork() *MemoryUnitOfWork {
	return &MemoryUnitOfWork{}
}

func (m *MemoryUnitOfWork) Do(ctx context.Context, operation func(tx *sql.Tx) error) error {
	return operation(nil)
}
