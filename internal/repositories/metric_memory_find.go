package repositories

import (
	"context"
	"go-metrics/internal/domain"
	"sync"
)

type MetricMemoryFindRepository struct {
	data map[domain.MetricID]*domain.Metric
	mu   sync.Mutex
}

func NewMetricMemoryFindRepository(
	data map[domain.MetricID]*domain.Metric,
) *MetricMemoryFindRepository {
	return &MetricMemoryFindRepository{data: data}
}

func (repo *MetricMemoryFindRepository) Find(
	ctx context.Context, filters []domain.MetricID,
) (map[domain.MetricID]*domain.Metric, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	filterMap := make(map[domain.MetricID]struct{})
	for _, filter := range filters {
		filterMap[filter] = struct{}{}
	}
	result := make(map[domain.MetricID]*domain.Metric)
	for metricID, metric := range repo.data {
		if _, exists := filterMap[metricID]; exists {
			result[metricID] = metric
		}
	}
	return result, nil
}
