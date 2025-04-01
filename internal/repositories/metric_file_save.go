package repositories

import (
	"context"
	"encoding/json"
	"go-metrics/internal/domain"
	"os"
	"sync"
)

type MetricFileSaveRepository struct {
	file *os.File
	mu   sync.Mutex
}

func NewMetricFileSaveRepository(file *os.File) *MetricFileSaveRepository {
	return &MetricFileSaveRepository{
		file: file,
	}
}

func (repo *MetricFileSaveRepository) Save(ctx context.Context, metrics []*domain.Metric) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, metric := range metrics {
		data, err := json.Marshal(metric)
		if err != nil {
			return err
		}
		_, err = repo.file.Write(append(data, '\n'))
		if err != nil {
			return err
		}
	}
	err := repo.file.Sync()
	if err != nil {
		return err
	}

	return nil
}
