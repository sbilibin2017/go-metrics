package repositories_test

import (
	"context"
	"encoding/json"
	"go-metrics/internal/domain"
	"go-metrics/internal/repositories"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricFileSaveRepository_Save(t *testing.T) {
	tempFile, err := os.CreateTemp("", "metrics_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	repo := repositories.NewMetricFileSaveRepository(tempFile)
	metrics := []*domain.Metric{
		{ID: "metric1", Type: "gauge", Value: float64Ptr(12.34)},
		{ID: "metric2", Type: "counter", Delta: int64Ptr(42)},
	}

	err = repo.Save(context.Background(), metrics)
	assert.NoError(t, err)

	content, err := os.ReadFile(tempFile.Name())
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	assert.Len(t, lines, len(metrics))

	for i, line := range lines {
		var metric domain.Metric
		err := json.Unmarshal([]byte(line), &metric)
		assert.NoError(t, err)
		assert.Equal(t, metrics[i], &metric)
	}
}

func float64Ptr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}
