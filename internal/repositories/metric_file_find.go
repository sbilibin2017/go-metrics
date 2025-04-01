package repositories

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/engines"
	"sync"
)

type MetricFileFindRepository struct {
	engine *engines.FileGeneratorEngine[*domain.Metric]
	mu     sync.Mutex
}

func NewMetricFileFindRepository(engine *engines.FileGeneratorEngine[*domain.Metric]) *MetricFileFindRepository {
	return &MetricFileFindRepository{
		engine: engine,
	}
}

func (repo *MetricFileFindRepository) Find(ctx context.Context, filters []*domain.MetricID) (map[domain.MetricID]*domain.Metric, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	result := make(map[domain.MetricID]*domain.Metric)
	filterMap := make(map[domain.MetricID]bool)
	for _, filter := range filters {
		if filter != nil {
			filterMap[domain.MetricID{ID: filter.ID, Type: filter.Type}] = true
		}
	}
	for metric := range repo.engine.Generate(ctx) {
		metricID := domain.MetricID{ID: metric.ID, Type: metric.Type}
		if len(filters) == 0 || filterMap[metricID] {
			result[metricID] = metric
		}
	}
	return result, nil
}
