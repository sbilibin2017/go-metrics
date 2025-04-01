package repositories_test

import (
	"bytes"
	"context"
	"encoding/json"
	"go-metrics/internal/domain"
	"go-metrics/internal/repositories"

	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestFileWithMetrics(t *testing.T, metrics []domain.Metric) *os.File {
	var lines [][]byte
	for _, metric := range metrics {
		data, err := json.Marshal(metric)
		require.NoError(t, err)
		lines = append(lines, data)
	}
	fileContent := bytes.Join(lines, []byte("\n"))
	file, err := os.CreateTemp("", "metrics_test_")
	require.NoError(t, err)
	_, err = file.Write(fileContent)
	require.NoError(t, err)
	err = file.Sync()
	require.NoError(t, err)
	_, err = file.Seek(0, 0)
	require.NoError(t, err)
	return file
}

func TestFindWithoutFilters(t *testing.T) {
	metrics := []domain.Metric{
		{ID: "1", Type: "counter", Delta: new(int64), Value: nil},
		{ID: "2", Type: "gauge", Delta: nil, Value: new(float64)},
	}
	file := createTestFileWithMetrics(t, metrics)
	defer os.Remove(file.Name())
	repo := repositories.NewMetricFileFindRepository(file)
	result, err := repo.Find(context.Background(), nil)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, domain.MetricID{ID: "1", Type: "counter"})
	assert.Contains(t, result, domain.MetricID{ID: "2", Type: "gauge"})
}

func TestFindWithFilters(t *testing.T) {
	metrics := []domain.Metric{
		{ID: "1", Type: "counter", Delta: new(int64), Value: nil},
		{ID: "2", Type: "gauge", Delta: nil, Value: new(float64)},
		{ID: "3", Type: "counter", Delta: new(int64), Value: nil},
	}
	file := createTestFileWithMetrics(t, metrics)
	defer os.Remove(file.Name())
	repo := repositories.NewMetricFileFindRepository(file)
	filters := []*domain.MetricID{
		{ID: "1", Type: "counter"},
		{ID: "2", Type: "gauge"},
	}
	result, err := repo.Find(context.Background(), filters)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, domain.MetricID{ID: "1", Type: "counter"})
	assert.Contains(t, result, domain.MetricID{ID: "2", Type: "gauge"})
	assert.NotContains(t, result, domain.MetricID{ID: "3", Type: "counter"})
}
