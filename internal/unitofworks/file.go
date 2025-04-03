package unitofworks

import (
	"context"
	"database/sql"
)

type FileUnitOfWork struct{}

func NewFileUnitOfWork() *FileUnitOfWork {
	return &FileUnitOfWork{}
}

func (f *FileUnitOfWork) Do(ctx context.Context, operation func(tx *sql.Tx) error) error {
	return operation(nil)
}
