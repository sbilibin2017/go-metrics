package repositories

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/engines"
)

type MetricFileSaveRepository struct {
	e *engines.FileWriterEngine[*domain.Metric]
}

func NewMetricFileSaveRepository(e *engines.FileWriterEngine[*domain.Metric]) *MetricFileSaveRepository {
	return &MetricFileSaveRepository{
		e: e,
	}
}

func (repo *MetricFileSaveRepository) Save(ctx context.Context, metrics []*domain.Metric) error {
	return repo.e.Write(ctx, metrics)
}
