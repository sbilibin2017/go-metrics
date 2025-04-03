package services

import (
	"context"
	"database/sql"
	"go-metrics/internal/domain"
	"go-metrics/internal/errors"
)

type MetricUpdateSaveRepository interface {
	Save(ctx context.Context, metrics []*domain.Metric) error
}

type MetricUpdateFindRepository interface {
	Find(ctx context.Context, filters []*domain.MetricID) (map[domain.MetricID]*domain.Metric, error)
}

type UnitOfWork interface {
	Do(ctx context.Context, operation func(tx *sql.Tx) error) error
}

type MetricUpdateService struct {
	s MetricUpdateSaveRepository
	f MetricUpdateFindRepository
	u UnitOfWork
}

func NewMetricUpdateService(
	s MetricUpdateSaveRepository,
	f MetricUpdateFindRepository,
	u UnitOfWork,
) *MetricUpdateService {
	return &MetricUpdateService{
		s: s,
		f: f,
		u: u,
	}
}

func (s *MetricUpdateService) Update(
	ctx context.Context, metrics []*domain.Metric,
) ([]*domain.Metric, error) {
	err := s.u.Do(ctx, func(tx *sql.Tx) error {
		metricIDs := make([]*domain.MetricID, len(metrics))
		for i, metric := range metrics {
			metricIDs[i] = &domain.MetricID{ID: metric.ID, Type: metric.Type}
		}
		existingMetrics, err := s.f.Find(ctx, metricIDs)
		if err != nil {
			return err
		}
		if existingMetrics == nil {
			existingMetrics = make(map[domain.MetricID]*domain.Metric)
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
		if err := s.s.Save(ctx, metrics); err != nil {
			return errors.ErrMetricIsNotUpdated
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
