package repositories

import (
	"context"
	"go-metrics/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricMemoryFindRepository_Find_AllMetrics(t *testing.T) {
	metric1 := &domain.Metric{MetricID: domain.MetricID{ID: "1", Type: domain.Counter}}
	metric2 := &domain.Metric{MetricID: domain.MetricID{ID: "2", Type: domain.Gauge}}
	repo := NewMetricMemoryFindRepository(map[domain.MetricID]*domain.Metric{
		metric1.MetricID: metric1,
		metric2.MetricID: metric2,
	})

	result, err := repo.Find(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, metric1.MetricID)
	assert.Contains(t, result, metric2.MetricID)
}

func TestMetricMemoryFindRepository_Find_ByFilters(t *testing.T) {
	metric1 := &domain.Metric{MetricID: domain.MetricID{ID: "1", Type: domain.Counter}}
	metric2 := &domain.Metric{MetricID: domain.MetricID{ID: "2", Type: domain.Gauge}}
	repo := NewMetricMemoryFindRepository(map[domain.MetricID]*domain.Metric{
		metric1.MetricID: metric1,
		metric2.MetricID: metric2,
	})

	filters := []*domain.MetricID{
		&metric1.MetricID,
		&domain.MetricID{ID: "nonexistent", Type: domain.Counter},
	}

	result, err := repo.Find(context.Background(), filters)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Contains(t, result, metric1.MetricID)
	assert.NotContains(t, result, metric2.MetricID)
}

func TestMetricMemoryFindRepository_Find_EmptyFilter(t *testing.T) {
	metric1 := &domain.Metric{MetricID: domain.MetricID{ID: "1", Type: domain.Counter}}
	metric2 := &domain.Metric{MetricID: domain.MetricID{ID: "2", Type: domain.Gauge}}
	repo := NewMetricMemoryFindRepository(map[domain.MetricID]*domain.Metric{
		metric1.MetricID: metric1,
		metric2.MetricID: metric2,
	})

	filters := []*domain.MetricID{}
	result, err := repo.Find(context.Background(), filters)
	assert.NoError(t, err)
	assert.Len(t, result, 2) // Now expects 2 metrics because no filters are provided
	assert.Contains(t, result, metric1.MetricID)
	assert.Contains(t, result, metric2.MetricID)
}
