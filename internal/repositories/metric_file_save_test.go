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

type MockFileStoragePathGetter struct {
	Path string
}

func (m *MockFileStoragePathGetter) GetFileStoragePath() string {
	return m.Path
}

func TestMetricFileSaveRepository_Save(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "metrics_test.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	fileEngine := engines.NewFileEngine()
	fsp := &MockFileStoragePathGetter{Path: tmpFile.Name()}
	require.NoError(t, fileEngine.Open(fsp))
	defer fileEngine.File.Close()
	writerEngine := engines.NewFileWriterEngine[*domain.Metric](fileEngine)
	repo := repositories.NewMetricFileSaveRepository(writerEngine)
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
	err = repo.Save(context.Background(), metrics)
	require.NoError(t, err)
	content, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)
	expectedContent := `[
		{
			"id": "1",
			"type": "gauge",
			"value": 100
		},
		{
			"id": "2",
			"type": "counter",
			"value": 200
		}
	]`
	assert.JSONEq(t, expectedContent, string(content))
}

func float64Ptr(value float64) *float64 {
	return &value
}
