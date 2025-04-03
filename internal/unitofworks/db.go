package unitofworks

import (
	"context"
	"database/sql"
	"fmt"
	"go-metrics/internal/errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBUnitOfWork struct {
	db *sql.DB
}

func NewDBUnitOfWork(db *sql.DB) *DBUnitOfWork {
	return &DBUnitOfWork{db: db}
}

func (uow *DBUnitOfWork) Do(ctx context.Context, operation func(tx *sql.Tx) error) error {
	var err error
	var tx *sql.Tx
	retryIntervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	attempts := 0
	for attempts < len(retryIntervals)+1 {
		tx, err = uow.db.BeginTx(ctx, nil)
		if err != nil {
			if errors.IsRetriableError(err) && attempts < len(retryIntervals) {
				attempts++
				time.Sleep(retryIntervals[attempts-1])
				continue
			}
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		err = operation(tx)
		if err != nil {
			if errors.IsRetriableError(err) && attempts < len(retryIntervals) {
				tx.Rollback()
				attempts++
				time.Sleep(retryIntervals[attempts-1])
				continue
			}
			tx.Rollback()
			return fmt.Errorf("operation failed: %w", err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		return nil
	}
	return fmt.Errorf("failed to perform operation after multiple attempts")
}
