package engines

import (
	"context"
	"database/sql"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DatabaseDSNGetter interface {
	GetDatabaseDSN() string
}

type DBEngine struct {
	*sql.DB
	once sync.Once
}

func NewDBEngine() *DBEngine {
	return &DBEngine{}
}

func (e *DBEngine) Open(ctx context.Context, g DatabaseDSNGetter) error {
	var err error
	e.once.Do(func() {
		dsn := g.GetDatabaseDSN()
		e.DB, err = sql.Open("pgx", dsn)
		if err != nil {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = e.DB.PingContext(ctx); err != nil {
			e.DB.Close()
		}
	})
	return err
}
