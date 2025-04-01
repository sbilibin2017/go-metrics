package repositories_test

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/repositories"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	repo := repositories.NewMetricFileSaveRepository(tmpFile)
	delta := int64(10)
	value := 20.5
	metrics := []*domain.Metric{
		{ID: "1", Type: "counter", Delta: &delta},
		{ID: "2", Type: "gauge", Value: &value},
	}
	err = repo.Save(context.Background(), metrics)
	require.NoError(t, err)
	fileData, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)
	expectedData := `{"id":"1","type":"counter","delta":10}
{"id":"2","type":"gauge","value":20.5}
`
	assert.Equal(t, expectedData, string(fileData))
	lines := strings.Split(strings.TrimSpace(string(fileData)), "\n")
	for _, line := range lines {
		assert.JSONEq(t, line, line)
	}
}
