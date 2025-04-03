package repositories

import (
	"context"
	"encoding/json"
	"go-metrics/internal/domain"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricFileSaveRepository_Save_Success(t *testing.T) {
	metric1 := &domain.Metric{MetricID: domain.MetricID{ID: "1", Type: domain.Counter}}
	metric2 := &domain.Metric{MetricID: domain.MetricID{ID: "2", Type: domain.Gauge}}
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	repo := NewMetricFileSaveRepository(tmpFile)
	err = repo.Save(context.Background(), []*domain.Metric{metric1, metric2})
	assert.NoError(t, err)
	tmpFile.Seek(0, 0)
	var savedMetrics []*domain.Metric
	decoder := json.NewDecoder(tmpFile)
	for {
		var m domain.Metric
		if err := decoder.Decode(&m); err != nil {
			break
		}
		savedMetrics = append(savedMetrics, &m)
	}
	assert.Len(t, savedMetrics, 2)
	assert.Equal(t, metric1, savedMetrics[0])
	assert.Equal(t, metric2, savedMetrics[1])
}
