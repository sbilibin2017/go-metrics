package repositories

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/engines"
	"sync"
)

type MetricMemoryFindRepository struct {
	g  *engines.MemoryGetter[domain.MetricID, *domain.Metric]
	r  *engines.MemoryRanger[domain.MetricID, *domain.Metric]
	mu sync.Mutex
}

func NewMetricMemoryFindRepository(
	g *engines.MemoryGetter[domain.MetricID, *domain.Metric],
	r *engines.MemoryRanger[domain.MetricID, *domain.Metric],
) *MetricMemoryFindRepository {
	return &MetricMemoryFindRepository{g: g, r: r}
}

func (repo *MetricMemoryFindRepository) Find(
	ctx context.Context, filters []*domain.MetricID,
) (map[domain.MetricID]*domain.Metric, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	filterMap := make(map[domain.MetricID]struct{})
	for _, filter := range filters {
		if filter != nil {
			filterMap[domain.MetricID{ID: filter.ID, Type: filter.Type}] = struct{}{}
		}
	}
	result := make(map[domain.MetricID]*domain.Metric)
	if len(filters) == 0 {
		repo.r.Range(func(key domain.MetricID, value *domain.Metric) bool {
			result[key] = value
			return true
		})
	} else {
		for filterKey := range filterMap {
			if metric, found := repo.g.Get(filterKey); found {
				result[filterKey] = metric
			}
		}
	}
	return result, nil
}
