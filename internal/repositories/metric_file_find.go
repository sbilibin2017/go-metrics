package repositories

import (
	"bufio"
	"context"
	"encoding/json"
	"go-metrics/internal/domain"
	"os"
	"sync"
)

type MetricFileFindRepository struct {
	file *os.File
	mu   sync.Mutex
}

func NewMetricFileFindRepository(file *os.File) *MetricFileFindRepository {
	return &MetricFileFindRepository{
		file: file,
	}
}

func (repo *MetricFileFindRepository) Find(ctx context.Context, filters []*domain.MetricID) (map[domain.MetricID]*domain.Metric, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	var filterValues []domain.MetricID
	for _, filter := range filters {
		if filter != nil {
			filterValues = append(filterValues, *filter)
		}
	}
	filterMap := make(map[domain.MetricID]struct{})
	for _, filter := range filterValues {
		filterMap[filter] = struct{}{}
	}
	result := make(map[domain.MetricID]*domain.Metric)
	scanner := bufio.NewScanner(repo.file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		var metric domain.Metric
		err := json.Unmarshal([]byte(line), &metric)
		if err != nil {
			continue
		}
		metricID := domain.MetricID{
			ID:   metric.ID,
			Type: metric.Type,
		}
		if len(filters) == 0 {
			result[metricID] = &metric
		} else {
			if _, exists := filterMap[metricID]; exists {
				result[metricID] = &metric
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
