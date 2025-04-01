package unitofwork

import "context"

type FileUnitOfWork struct{}

func NewFileUnitOfWork() *FileUnitOfWork {
	return &FileUnitOfWork{}
}

func (m *FileUnitOfWork) Do(ctx context.Context, operation func() error) error {
	return operation()
}
