package repositories

import (
	"context"
	"go-metrics/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricMemorySaveRepository_Save(t *testing.T) {
	metric1 := &domain.Metric{MetricID: domain.MetricID{ID: "1", Type: domain.Counter}}
	metric2 := &domain.Metric{MetricID: domain.MetricID{ID: "2", Type: domain.Gauge}}
	repo := NewMetricMemorySaveRepository(make(map[domain.MetricID]*domain.Metric))
	err := repo.Save(context.Background(), []*domain.Metric{metric1, metric2})
	assert.NoError(t, err)
	assert.Contains(t, repo.data, domain.MetricID{ID: metric1.ID, Type: metric1.Type})
	assert.Contains(t, repo.data, domain.MetricID{ID: metric2.ID, Type: metric2.Type})
	assert.Equal(t, metric1, repo.data[domain.MetricID{ID: metric1.ID, Type: metric1.Type}])
	assert.Equal(t, metric2, repo.data[domain.MetricID{ID: metric2.ID, Type: metric2.Type}])
}
