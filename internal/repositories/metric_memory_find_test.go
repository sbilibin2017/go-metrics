package repositories_test

import (
	"context"
	"testing"

	"go-metrics/internal/domain"
	"go-metrics/internal/engines"
	"go-metrics/internal/repositories"

	"github.com/stretchr/testify/assert"
)

func TestFindAllMetrics(t *testing.T) {
	// Initialize the storage and the repository
	storage := engines.NewMemoryStorage[domain.MetricID, *domain.Metric]()
	storageGetter := &engines.MemoryGetter[domain.MetricID, *domain.Metric]{MemoryStorage: storage}
	storageRanger := &engines.MemoryRanger[domain.MetricID, *domain.Metric]{MemoryStorage: storage}

	// Create sample metrics and add them to the storage
	metric1 := &domain.Metric{ID: "metric1", Type: "counter", Value: nil}
	metric2 := &domain.Metric{ID: "metric2", Type: "gauge", Value: nil}
	storageSetter := &engines.MemorySetter[domain.MetricID, *domain.Metric]{MemoryStorage: storage}
	storageSetter.Set(domain.MetricID{ID: "metric1", Type: "counter"}, metric1)
	storageSetter.Set(domain.MetricID{ID: "metric2", Type: "gauge"}, metric2)

	// Creating the repository
	repo := repositories.NewMetricMemoryFindRepository(storageGetter, storageRanger)

	// Call Find with no filters (fetch all metrics)
	result, err := repo.Find(context.Background(), nil)

	// Assertions
	assert.NoError(t, err, "Find should not return an error")
	assert.Len(t, result, 2, "There should be 2 metrics in the result")
	assert.Equal(t, metric1, result[domain.MetricID{ID: "metric1", Type: "counter"}], "Metric 1 should be in the result")
	assert.Equal(t, metric2, result[domain.MetricID{ID: "metric2", Type: "gauge"}], "Metric 2 should be in the result")
}

func TestFindWithFilters(t *testing.T) {
	// Initialize the storage and the repository
	storage := engines.NewMemoryStorage[domain.MetricID, *domain.Metric]()
	storageGetter := &engines.MemoryGetter[domain.MetricID, *domain.Metric]{MemoryStorage: storage}
	storageRanger := &engines.MemoryRanger[domain.MetricID, *domain.Metric]{MemoryStorage: storage}

	// Create sample metrics and add them to the storage
	metric1 := &domain.Metric{ID: "metric1", Type: "counter", Value: nil}
	metric2 := &domain.Metric{ID: "metric2", Type: "gauge", Value: nil}
	storageSetter := &engines.MemorySetter[domain.MetricID, *domain.Metric]{MemoryStorage: storage}
	storageSetter.Set(domain.MetricID{ID: "metric1", Type: "counter"}, metric1)
	storageSetter.Set(domain.MetricID{ID: "metric2", Type: "gauge"}, metric2)

	// Creating the repository
	repo := repositories.NewMetricMemoryFindRepository(storageGetter, storageRanger)

	// Filters to find specific metrics
	filters := []*domain.MetricID{
		{ID: "metric1", Type: "counter"},
	}

	// Call Find with filters
	result, err := repo.Find(context.Background(), filters)

	// Assertions
	assert.NoError(t, err, "Find should not return an error")
	assert.Len(t, result, 1, "There should be 1 metric in the result")

	// Check that the expected metric is returned
	expectedMetric := &domain.Metric{ID: "metric1", Type: "counter", Value: nil}
	assert.Equal(t, expectedMetric, result[domain.MetricID{ID: "metric1", Type: "counter"}], "Metric 1 should be in the result")
}
