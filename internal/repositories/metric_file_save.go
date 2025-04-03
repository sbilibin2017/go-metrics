package repositories

import (
	"context"
	"encoding/json"
	"go-metrics/internal/domain"
	"os"
	"sync"
)

type MetricFileSaveRepository struct {
	file    *os.File
	encoder *json.Encoder
	mu      sync.Mutex
}

func NewMetricFileSaveRepository(file *os.File) *MetricFileSaveRepository {
	return &MetricFileSaveRepository{
		file:    file,
		encoder: json.NewEncoder(file),
	}
}

func (repo *MetricFileSaveRepository) Save(ctx context.Context, metrics []*domain.Metric) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, metric := range metrics {
		repo.encoder.Encode(metric)
	}
	return nil
}
