package unitofwork

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBUnitOfWork struct {
	db *sql.DB
}

func NewDBUnitOfWork(db *sql.DB) *DBUnitOfWork {
	return &DBUnitOfWork{db: db}
}

func (uow *DBUnitOfWork) Do(ctx context.Context, operation func() error) error {
	tx, err := uow.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()
	err = operation()
	if err != nil {
		return fmt.Errorf("operation failed: %v", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}
