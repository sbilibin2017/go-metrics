package repositories

import (
	"context"
	"go-metrics/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricMemorySaveRepository_Save(t *testing.T) {
	data := make(map[domain.MetricID]*domain.Metric)
	repo := NewMetricMemorySaveRepository(data)

	value1 := 42.5
	delta2 := int64(10)

	metrics := []*domain.Metric{
		{ID: "metric1", Type: "gauge", Value: &value1},
		{ID: "metric2", Type: domain.Counter, Delta: &delta2},
	}

	err := repo.Save(context.Background(), metrics)
	assert.NoError(t, err)

	for _, metric := range metrics {
		key := domain.MetricID{ID: metric.ID, Type: metric.Type}
		storedMetric, exists := data[key]
		assert.True(t, exists)
		assert.Equal(t, metric, storedMetric)
	}
}
