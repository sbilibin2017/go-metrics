package engines

import "context"

type MemoryUnitOfWork struct{}

func NewMemoryUnitOfWork() *MemoryUnitOfWork {
	return &MemoryUnitOfWork{}
}

func (m *MemoryUnitOfWork) Do(ctx context.Context, operation func() error) error {
	return nil
}
