package repositories

import (
	"context"
	"database/sql"
	"go-metrics/internal/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type MetricDBSaveRepository struct {
	db *sql.DB
}

func NewMetricDBSaveRepository(db *sql.DB) *MetricDBSaveRepository {
	return &MetricDBSaveRepository{db: db}
}

var metricSaveQuery = `
	INSERT INTO metrics (id, type, delta, value) 
	VALUES ($1, $2, $3, $4) 
	ON CONFLICT (id, type) DO UPDATE 
	SET delta = EXCLUDED.delta, value = EXCLUDED.value;
`

func (repo *MetricDBSaveRepository) Save(ctx context.Context, metrics []*domain.Metric) error {
	stmt, err := repo.db.PrepareContext(ctx, metricSaveQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, metric := range metrics {
		_, err := stmt.ExecContext(ctx, metric.ID, metric.Type, metric.Delta, metric.Value)
		if err != nil {
			return err
		}
	}
	return nil
}
