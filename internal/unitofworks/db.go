package unitofworks

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBUnitOfWork struct {
	db *sql.DB
}

func NewDBUnitOfWork(db *sql.DB) *DBUnitOfWork {
	return &DBUnitOfWork{db: db}
}

func (uow *DBUnitOfWork) Do(ctx context.Context, operation func(tx *sql.Tx) error) error {
	tx, err := uow.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = operation(tx)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
