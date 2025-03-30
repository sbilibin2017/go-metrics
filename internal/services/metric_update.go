package services

import (
	"context"
	"errors"
	"go-metrics/internal/domain"
)

type MetricUpdateSaveBatchRepository interface {
	Save(ctx context.Context, metrics []*domain.Metric) error
}

type MetricUpdateFindBatchRepository interface {
	Find(ctx context.Context, filters []domain.MetricID) (map[domain.MetricID]*domain.Metric, error)
}

type UnitOfWork interface {
	Do(ctx context.Context, operation func() error) error
}

type MetricUpdateService struct {
	saveRepo MetricUpdateSaveBatchRepository
	findRepo MetricUpdateFindBatchRepository
	uow      UnitOfWork
}

func NewMetricUpdateService(
	saveRepo MetricUpdateSaveBatchRepository,
	findRepo MetricUpdateFindBatchRepository,
	uow UnitOfWork,
) *MetricUpdateService {
	return &MetricUpdateService{
		saveRepo: saveRepo,
		findRepo: findRepo,
		uow:      uow,
	}
}

var ErrMetricIsNotUpdated = errors.New("metric is not updated")

func (s *MetricUpdateService) Update(
	ctx context.Context, metrics []*domain.Metric,
) ([]*domain.Metric, error) {
	err := s.uow.Do(ctx, func() error {
		metricIDs := make([]domain.MetricID, len(metrics))
		for i, metric := range metrics {
			metricIDs[i] = domain.MetricID{ID: metric.ID, Type: metric.Type}
		}
		existingMetrics, err := s.findRepo.Find(ctx, metricIDs)
		if err != nil {
			return ErrMetricIsNotUpdated
		}
		for i, metric := range metrics {
			switch metric.Type {
			case domain.Counter:
				if existingMetric, exists := existingMetrics[domain.MetricID{
					ID:   metric.ID,
					Type: metric.Type,
				}]; exists {
					*metric.Delta += *existingMetric.Delta
				}
			}
			metrics[i] = metric
		}
		if err := s.saveRepo.Save(ctx, metrics); err != nil {
			return ErrMetricIsNotUpdated
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return metrics, nil
}
