package repositories

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"go-metrics/internal/domain"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricFileFindRepository_Find_WithFilters(t *testing.T) {
	metric1 := &domain.Metric{MetricID: domain.MetricID{ID: "1", Type: domain.Counter}}
	metric2 := &domain.Metric{MetricID: domain.MetricID{ID: "2", Type: domain.Gauge}}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.Encode(metric1)
	encoder.Encode(metric2)
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Write(buf.Bytes())
	scanner := bufio.NewScanner(tmpFile)
	repo := NewMetricFileFindRepository(tmpFile, scanner)
	filters := []*domain.MetricID{
		&metric1.MetricID,
	}
	result, err := repo.Find(context.Background(), filters)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Contains(t, result, metric1.MetricID)
	assert.NotContains(t, result, metric2.MetricID)
}

func TestMetricFileFindRepository_Find_WithoutFilters(t *testing.T) {
	metric1 := &domain.Metric{MetricID: domain.MetricID{ID: "1", Type: domain.Counter}}
	metric2 := &domain.Metric{MetricID: domain.MetricID{ID: "2", Type: domain.Gauge}}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.Encode(metric1)
	encoder.Encode(metric2)
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Write(buf.Bytes())
	scanner := bufio.NewScanner(tmpFile)
	repo := NewMetricFileFindRepository(tmpFile, scanner)
	result, err := repo.Find(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, metric1.MetricID)
	assert.Contains(t, result, metric2.MetricID)
}

func TestMetricFileFindRepository_Find_EmptyFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	scanner := bufio.NewScanner(tmpFile)
	repo := NewMetricFileFindRepository(tmpFile, scanner)
	result, err := repo.Find(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, result, 0)
}
