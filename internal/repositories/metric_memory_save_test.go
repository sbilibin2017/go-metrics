package repositories_test

import (
	"context"
	"testing"

	"go-metrics/internal/domain"
	"go-metrics/internal/engines"
	"go-metrics/internal/repositories"

	"github.com/stretchr/testify/assert"
)

func TestSaveAndRetrieveMetrics(t *testing.T) {
	// Initialize the storage and the repository
	storage := engines.NewMemoryStorage[domain.MetricID, *domain.Metric]()
	storageSetter := &engines.MemorySetter[domain.MetricID, *domain.Metric]{MemoryStorage: storage}
	repo := repositories.NewMetricMemorySaveRepository(storageSetter)

	// Create sample metrics with *float64 and *int64 values
	value1 := 10.0
	value2 := 20.0
	metrics := []*domain.Metric{
		{ID: "metric1", Type: "counter", Value: &value1}, // Value is a pointer to float64
		{ID: "metric2", Type: "gauge", Value: &value2},   // Value is a pointer to float64
	}

	// Save the metrics
	err := repo.Save(context.Background(), metrics)
	assert.NoError(t, err, "Save should not return an error")

}
