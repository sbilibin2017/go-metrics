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
	Find(ctx context.Context, filters []*domain.MetricID) (map[domain.MetricID]*domain.Metric, error)
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

func (s *MetricUpdateService) Update(
	ctx context.Context, metrics []*domain.Metric,
) ([]*domain.Metric, error) {
	var updatedMetrics []*domain.Metric
	err := s.uow.Do(ctx, func() error {
		existingMetrics, err := s.findMetrics(ctx, metrics)
		if err != nil {
			return ErrMetricIsNotUpdated
		}
		if err := s.updateMetrics(ctx, metrics, existingMetrics); err != nil {
			return err
		}
		updatedMetrics = metrics
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updatedMetrics, nil
}

var ErrMetricIsNotUpdated = errors.New("metric is not updated")

func (s *MetricUpdateService) findMetrics(
	ctx context.Context, metrics []*domain.Metric,
) (map[domain.MetricID]*domain.Metric, error) {
	metricIDs := make([]*domain.MetricID, len(metrics))
	for i, metric := range metrics {
		metricIDs[i] = &domain.MetricID{ID: metric.ID, MType: metric.MType}
	}
	return s.findRepo.Find(ctx, metricIDs)
}

func (s *MetricUpdateService) updateMetrics(
	ctx context.Context, metrics []*domain.Metric, existingMetrics map[domain.MetricID]*domain.Metric,
) error {
	for i, metric := range metrics {
		switch metric.MType {
		case domain.Counter:
			if existingMetric, exists := existingMetrics[domain.MetricID{
				ID:    metric.ID,
				MType: metric.MType,
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
}
