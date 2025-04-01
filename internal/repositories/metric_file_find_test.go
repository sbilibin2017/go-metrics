package repositories_test

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/engines"
	"go-metrics/internal/repositories"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricFileFindRepository_Find(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "metrics_test.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	fileEngine := engines.NewFileEngine()
	fsp := &MockFileStoragePathGetter{Path: tmpFile.Name()}
	require.NoError(t, fileEngine.Open(fsp))
	defer fileEngine.Close()
	generatorEngine := engines.NewFileGeneratorEngine[*domain.Metric](fileEngine)
	metrics := []*domain.Metric{
		{
			ID:    "1",
			Type:  "gauge",
			Value: float64Ptr(100),
		},
		{
			ID:    "2",
			Type:  "counter",
			Value: float64Ptr(200),
		},
	}
	var metricsValues []domain.Metric
	for _, metric := range metrics {
		metricsValues = append(metricsValues, *metric)
	}
	writerEngine := engines.NewFileWriterEngine[domain.Metric](fileEngine)
	err = writerEngine.Write(context.Background(), metricsValues)
	require.NoError(t, err)
	repo := repositories.NewMetricFileFindRepository(generatorEngine)
	filters := []*domain.MetricID{
		{ID: "1", Type: "gauge"},
	}
	result, err := repo.Find(context.Background(), filters)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Contains(t, result, domain.MetricID{ID: "1", Type: "gauge"})
	assert.NotNil(t, result[domain.MetricID{ID: "1", Type: "gauge"}])
}
