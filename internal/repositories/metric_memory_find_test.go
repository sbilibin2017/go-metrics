package repositories

import (
	"context"
	"go-metrics/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricMemoryFindRepository_Find(t *testing.T) {
	data := map[domain.MetricID]*domain.Metric{
		{ID: "metric1", Type: "gauge"}:        {ID: "metric1", Type: "gauge", Value: func() *float64 { v := 42.5; return &v }()},
		{ID: "metric2", Type: domain.Counter}: {ID: "metric2", Type: domain.Counter, Delta: func() *int64 { v := int64(10); return &v }()},
		{ID: "metric3", Type: "gauge"}:        {ID: "metric3", Type: "gauge", Value: func() *float64 { v := 15.7; return &v }()},
	}

	repo := NewMetricMemoryFindRepository(data)

	testCases := []struct {
		name     string
		filters  []*domain.MetricID // Используем указатели на MetricID
		expected map[domain.MetricID]*domain.Metric
	}{
		{
			name: "FindExistingMetrics",
			filters: []*domain.MetricID{
				&domain.MetricID{ID: "metric1", Type: "gauge"},
				&domain.MetricID{ID: "metric2", Type: domain.Counter},
			},
			expected: map[domain.MetricID]*domain.Metric{
				{ID: "metric1", Type: "gauge"}:        {ID: "metric1", Type: "gauge", Value: func() *float64 { v := 42.5; return &v }()},
				{ID: "metric2", Type: domain.Counter}: {ID: "metric2", Type: domain.Counter, Delta: func() *int64 { v := int64(10); return &v }()},
			},
		},
		{
			name:     "FindNonExistingMetric",
			filters:  []*domain.MetricID{{ID: "metric4", Type: "gauge"}},
			expected: map[domain.MetricID]*domain.Metric{},
		},
		{
			name:    "FindAllMetricsWhenNoFilters",
			filters: []*domain.MetricID{}, // Empty filters
			expected: map[domain.MetricID]*domain.Metric{
				{ID: "metric1", Type: "gauge"}:        {ID: "metric1", Type: "gauge", Value: func() *float64 { v := 42.5; return &v }()},
				{ID: "metric2", Type: domain.Counter}: {ID: "metric2", Type: domain.Counter, Delta: func() *int64 { v := int64(10); return &v }()},
				{ID: "metric3", Type: "gauge"}:        {ID: "metric3", Type: "gauge", Value: func() *float64 { v := 15.7; return &v }()},
			},
		},
		{
			name:    "FindAllMetricsWhenFilterMapEmpty",
			filters: []*domain.MetricID{}, // Empty filters, which results in filterMap being empty
			expected: map[domain.MetricID]*domain.Metric{
				{ID: "metric1", Type: "gauge"}:        {ID: "metric1", Type: "gauge", Value: func() *float64 { v := 42.5; return &v }()},
				{ID: "metric2", Type: domain.Counter}: {ID: "metric2", Type: domain.Counter, Delta: func() *int64 { v := int64(10); return &v }()},
				{ID: "metric3", Type: "gauge"}:        {ID: "metric3", Type: "gauge", Value: func() *float64 { v := 15.7; return &v }()},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := repo.Find(context.Background(), tc.filters)
			assert.NoError(t, err)

			// Here we ignore comparing by pointers, comparing only the content
			assert.Len(t, result, len(tc.expected))
			for key, expectedValue := range tc.expected {
				actualValue, exists := result[key]
				assert.True(t, exists)
				assert.Equal(t, expectedValue, actualValue)
			}
		})
	}
}
