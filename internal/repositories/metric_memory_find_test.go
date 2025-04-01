package repositories_test

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricMemoryFindRepository_Find(t *testing.T) {
	delta1, delta2 := int64(10), int64(5)
	value1, value2 := 20.5, 30.5
	data := map[domain.MetricID]*domain.Metric{
		{ID: "1", Type: "counter"}: {ID: "1", Type: "counter", Delta: &delta1},
		{ID: "2", Type: "gauge"}:   {ID: "2", Type: "gauge", Value: &value1},
		{ID: "3", Type: "counter"}: {ID: "3", Type: "counter", Delta: &delta2},
		{ID: "4", Type: "gauge"}:   {ID: "4", Type: "gauge", Value: &value2},
	}
	repo := repositories.NewMetricMemoryFindRepository(data)

	t.Run("Find all metrics", func(t *testing.T) {
		result, err := repo.Find(context.Background(), nil)
		require.NoError(t, err)
		assert.Len(t, result, 4)
		assert.Equal(t, data, result)
	})

	t.Run("Find specific metrics", func(t *testing.T) {
		filters := []*domain.MetricID{
			{ID: "1", Type: "counter"},
			{ID: "2", Type: "gauge"},
		}
		result, err := repo.Find(context.Background(), filters)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, data[domain.MetricID{ID: "1", Type: "counter"}], result[domain.MetricID{ID: "1", Type: "counter"}])
		assert.Equal(t, data[domain.MetricID{ID: "2", Type: "gauge"}], result[domain.MetricID{ID: "2", Type: "gauge"}])
	})

	t.Run("Find with non-existing metrics", func(t *testing.T) {
		filters := []*domain.MetricID{
			{ID: "100", Type: "counter"},
			{ID: "200", Type: "gauge"},
		}
		result, err := repo.Find(context.Background(), filters)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Find with nil filters", func(t *testing.T) {
		filters := []*domain.MetricID{
			nil,
			{ID: "2", Type: "gauge"},
		}
		result, err := repo.Find(context.Background(), filters)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, data[domain.MetricID{ID: "2", Type: "gauge"}], result[domain.MetricID{ID: "2", Type: "gauge"}])
	})
}
