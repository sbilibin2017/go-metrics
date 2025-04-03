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

func NewMetricMemoryFindRepository(data map[domain.MetricID]*domain.Metric) *MetricMemoryFindRepository {
	return &MetricMemoryFindRepository{data: data}
}

func (repo *MetricMemoryFindRepository) Find(ctx context.Context, filters []*domain.MetricID) (map[domain.MetricID]*domain.Metric, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	result := make(map[domain.MetricID]*domain.Metric)
	if len(filters) == 0 {
		for key, value := range repo.data {
			result[key] = value
		}
		return result, nil
	}
	for _, filter := range filters {
		if filter != nil {
			if metric, found := repo.data[*filter]; found {
				result[*filter] = metric
			}
		}
	}
	return result, nil
}
