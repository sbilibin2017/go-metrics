package repositories

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/engines"
	"sync"
)

type MetricMemorySaveRepository struct {
	e  *engines.MemorySetter[domain.MetricID, *domain.Metric]
	mu sync.Mutex
}

func NewMetricMemorySaveRepository(
	e *engines.MemorySetter[domain.MetricID, *domain.Metric],
) *MetricMemorySaveRepository {
	return &MetricMemorySaveRepository{
		e: e,
	}
}

func (repo *MetricMemorySaveRepository) Save(
	ctx context.Context, metrics []*domain.Metric,
) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, metric := range metrics {
		repo.e.Set(domain.MetricID{ID: metric.ID, Type: metric.Type}, metric)
	}
	return nil
}
