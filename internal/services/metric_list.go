package services

import (
	"context"
	"errors"
	"go-metrics/internal/domain"
)

type MetricListFindBatchRepository interface {
	Find(ctx context.Context, filters []domain.MetricID) (map[domain.MetricID]*domain.Metric, error)
}

type MetricListService struct {
	findRepo MetricListFindBatchRepository
}

// Конструктор для создания сервиса
func NewMetricListService(
	findRepo MetricListFindBatchRepository,
) *MetricListService {
	return &MetricListService{
		findRepo: findRepo,
	}
}

var (
	ErrMetricListInternal = errors.New("internal error")
)

// Метод List для получения списка всех метрик
func (s *MetricListService) List(
	ctx context.Context,
) ([]*domain.Metric, error) {
	metricsMap, err := s.findRepo.Find(ctx, []domain.MetricID{})
	if err != nil {
		return nil, ErrMetricListInternal
	}
	var metrics []*domain.Metric
	for _, metric := range metricsMap {
		metrics = append(metrics, metric)
	}
	return metrics, nil
}
