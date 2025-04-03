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
	var updatedMetrics []*domain.Metric
	err := s.u.Do(ctx, func(tx *sql.Tx) error {
		metricMap := make(map[domain.MetricID]*domain.Metric)
		for _, metric := range metrics {
			metricID := domain.MetricID{
				ID:   metric.ID,
				Type: metric.Type,
			}
			if metric.Type == domain.Counter {
				if existingMetric, exists := metricMap[metricID]; exists {
					*existingMetric.Delta += *metric.Delta
				} else {
					metricMap[metricID] = metric
				}
			} else {
				metricMap[metricID] = metric
			}
		}
		metricIDs := make([]*domain.MetricID, 0, len(metricMap))
		for metricID := range metricMap {
			metricIDs = append(metricIDs, &metricID)
		}
		existingMetrics, err := s.f.Find(ctx, metricIDs)
		if err != nil {
			return err
		}
		if existingMetrics == nil {
			existingMetrics = make(map[domain.MetricID]*domain.Metric)
		}
		updatedMetrics = make([]*domain.Metric, 0, len(metricMap))
		for _, metric := range metricMap {
			if metric.Type == domain.Counter {
				if existingMetric, exists := existingMetrics[domain.MetricID{
					ID:   metric.ID,
					Type: metric.Type,
				}]; exists {
					*metric.Delta += *existingMetric.Delta
				}
				updatedMetrics = append(updatedMetrics, metric)
			} else {
				updatedMetrics = append(updatedMetrics, metric)
			}
		}
		if err := s.s.Save(ctx, updatedMetrics); err != nil {
			return errors.ErrMetricIsNotUpdated
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return updatedMetrics, nil
}
