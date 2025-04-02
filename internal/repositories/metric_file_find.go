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
	file    *os.File
	scanner *bufio.Scanner
	mu      sync.Mutex
}

func NewMetricFileFindRepository(file *os.File) *MetricFileFindRepository {
	return &MetricFileFindRepository{
		file:    file,
		scanner: bufio.NewScanner(file),
	}
}

func (repo *MetricFileFindRepository) Find(ctx context.Context, filters []*domain.MetricID) (map[domain.MetricID]*domain.Metric, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	_, err := repo.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	repo.scanner = bufio.NewScanner(repo.file)
	filterMap := make(map[domain.MetricID]struct{})
	for _, filter := range filters {
		filterMap[*filter] = struct{}{}
	}
	result := make(map[domain.MetricID]*domain.Metric)
	for repo.scanner.Scan() {
		line := repo.scanner.Text()
		var metric domain.Metric
		if err := json.Unmarshal([]byte(line), &metric); err != nil {
			continue
		}
		metricID := domain.MetricID{ID: metric.ID, Type: metric.Type}
		if len(filters) != 0 {
			if _, found := filterMap[metricID]; found {
				result[metricID] = &metric
			}
		} else {
			result[metricID] = &metric
		}
	}
	if err := repo.scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
