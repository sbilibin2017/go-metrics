package services

import (
	"context"

	"go-metrics/internal/domain"
	"go-metrics/internal/errors"
)

type MetricGetByIDFindRepository interface {
	Find(ctx context.Context, filters []*domain.MetricID) (map[domain.MetricID]*domain.Metric, error)
}

type MetricGetByIDService struct {
	f MetricGetByIDFindRepository
}

func NewMetricGetByIDService(
	f MetricGetByIDFindRepository,
) *MetricGetByIDService {
	return &MetricGetByIDService{
		f: f,
	}
}

func (s *MetricGetByIDService) GetByID(
	ctx context.Context, id *domain.MetricID,
) (*domain.Metric, error) {
	metrics, err := s.f.Find(ctx, []*domain.MetricID{id})
	if err != nil {
		return nil, errors.ErrMetricGetByIDInternal
	}
	metric, found := metrics[*id]
	if !found {
		return nil, errors.ErrMetricNotFound
	}
	return metric, nil
}
