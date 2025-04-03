package repositories

import (
	"bufio"
	"context"
	"encoding/json"
	"go-metrics/internal/domain"
	"sync"
)

type File interface {
	Seek(offset int64, whence int) (int64, error)
	Read(p []byte) (n int, err error)
}

type Scanner interface {
	Scan() bool
	Text() string
	Err() error
}

type MetricFileFindRepository struct {
	file    File
	scanner Scanner
	mu      sync.Mutex
}

func NewMetricFileFindRepository(file File, scanner Scanner) *MetricFileFindRepository {
	return &MetricFileFindRepository{
		file:    file,
		scanner: scanner,
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
	filterMap := make(map[domain.MetricID]bool)
	for _, filter := range filters {
		filterMap[*filter] = true
	}
	result := make(map[domain.MetricID]*domain.Metric)
	for repo.scanner.Scan() {
		line := repo.scanner.Text()
		var metric domain.Metric
		if err := json.Unmarshal([]byte(line), &metric); err != nil {
			continue
		}
		metricID := domain.MetricID{ID: metric.ID, Type: metric.Type}
		if len(filters) == 0 || filterMap[metricID] {
			result[metricID] = &metric
		}
	}
	if err := repo.scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
