package repositories

import (
	"context"
	"go-metrics/internal/domain"
	"sync"
)

type MetricMemorySaveRepository struct {
	data map[domain.MetricID]*domain.Metric
	mu   sync.Mutex
}

func NewMetricMemorySaveRepository(
	data map[domain.MetricID]*domain.Metric,
) *MetricMemorySaveRepository {
	return &MetricMemorySaveRepository{
		data: data,
	}
}

func (repo *MetricMemorySaveRepository) Save(
	ctx context.Context, metrics []*domain.Metric,
) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, metric := range metrics {
		repo.data[domain.MetricID{ID: metric.ID, MType: metric.MType}] = metric
	}
	return nil
}
