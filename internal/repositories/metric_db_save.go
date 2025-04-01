package repositories

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/engines"
)

type MetricDBSaveRepository struct {
	db *engines.DBEngine
}

func NewMetricDBSaveRepository(db *engines.DBEngine) *MetricDBSaveRepository {
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
