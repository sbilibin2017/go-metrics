package repositories_test

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricMemorySaveRepository_Save(t *testing.T) {
	delta1, delta2 := int64(10), int64(5)
	value1, value2 := 20.5, 30.5
	data := make(map[domain.MetricID]*domain.Metric)
	repo := repositories.NewMetricMemorySaveRepository(data)

	metrics := []*domain.Metric{
		{ID: "1", Type: "counter", Delta: &delta1},
		{ID: "2", Type: "gauge", Value: &value1},
		{ID: "3", Type: "counter", Delta: &delta2},
		{ID: "4", Type: "gauge", Value: &value2},
	}

	t.Run("Save new metrics", func(t *testing.T) {
		err := repo.Save(context.Background(), metrics)
		require.NoError(t, err)
		assert.Len(t, data, 4)

		assert.Equal(t, delta1, *data[domain.MetricID{ID: "1", Type: "counter"}].Delta)
		assert.Nil(t, data[domain.MetricID{ID: "1", Type: "counter"}].Value)

		assert.Equal(t, value1, *data[domain.MetricID{ID: "2", Type: "gauge"}].Value)
		assert.Nil(t, data[domain.MetricID{ID: "2", Type: "gauge"}].Delta)

		assert.Equal(t, delta2, *data[domain.MetricID{ID: "3", Type: "counter"}].Delta)
		assert.Nil(t, data[domain.MetricID{ID: "3", Type: "counter"}].Value)

		assert.Equal(t, value2, *data[domain.MetricID{ID: "4", Type: "gauge"}].Value)
		assert.Nil(t, data[domain.MetricID{ID: "4", Type: "gauge"}].Delta)
	})

	t.Run("Overwrite existing metrics", func(t *testing.T) {
		// Overwrite the first metric
		delta3 := int64(15)
		metricsToUpdate := []*domain.Metric{
			{ID: "1", Type: "counter", Delta: &delta3},
		}

		err := repo.Save(context.Background(), metricsToUpdate)
		require.NoError(t, err)

		// Verify that the metric is updated
		assert.Equal(t, delta3, *data[domain.MetricID{ID: "1", Type: "counter"}].Delta)
		assert.Nil(t, data[domain.MetricID{ID: "1", Type: "counter"}].Value)
	})

	t.Run("Save empty metrics", func(t *testing.T) {
		// Save no metrics
		err := repo.Save(context.Background(), []*domain.Metric{})
		require.NoError(t, err)
		assert.Len(t, data, 4) // The data should still be unchanged
	})
}
