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
	ctx context.Context, filters []*domain.MetricID,
) (map[domain.MetricID]*domain.Metric, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	var result map[domain.MetricID]*domain.Metric
	if len(filters) == 0 {
		result = make(map[domain.MetricID]*domain.Metric)
		for metricID, metric := range repo.data {
			result[metricID] = metric
		}
	} else {
		filterMap := make(map[string]map[string]struct{})
		for _, filter := range filters {
			if filter != nil {
				if _, ok := filterMap[filter.ID]; !ok {
					filterMap[filter.ID] = make(map[string]struct{})
				}
				filterMap[filter.ID][filter.Type] = struct{}{}
			}
		}
		result = make(map[domain.MetricID]*domain.Metric)
		for metricID, metric := range repo.data {
			if types, exists := filterMap[metricID.ID]; exists {
				if _, exists := types[metricID.Type]; exists {
					result[metricID] = metric
				}
			}
		}
	}
	return result, nil
}
