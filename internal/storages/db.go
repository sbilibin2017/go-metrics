package storages

import (
	"database/sql"
	"go-metrics/internal/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DatabaseConfig interface {
	GetDatabaseDSN() string
}

func NewDB(config DatabaseConfig) (*sql.DB, error) {
	logger.Logger.Infow("Connecting to database", "dsn", config.GetDatabaseDSN())
	db, err := sql.Open("pgx", config.GetDatabaseDSN())
	if err != nil {
		logger.Logger.Errorw("Failed to open database connection", "error", err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		logger.Logger.Errorw("Failed to ping database", "error", err)
		db.Close()
		return nil, err
	}
	logger.Logger.Infow("Database connection established successfully")
	return db, nil
}
