package services

import (
	"context"
	"errors"
	"go-metrics/internal/domain"
)

type MetricGetByIDFindBatchRepository interface {
	Find(ctx context.Context, filters []domain.MetricID) (map[domain.MetricID]*domain.Metric, error)
}

type MetricGetByIDService struct {
	findRepo MetricUpdateFindBatchRepository
}

func NewMetricGetByIDService(
	findRepo MetricGetByIDFindBatchRepository,
) *MetricGetByIDService {
	return &MetricGetByIDService{
		findRepo: findRepo,
	}
}

var (
	ErrMetricNotFound        = errors.New("metric not found")
	ErrMetricGetByIDInternal = errors.New("internal error")
)

func (s *MetricGetByIDService) GetByID(
	ctx context.Context, id *domain.MetricID,
) (*domain.Metric, error) {
	metrics, err := s.findRepo.Find(ctx, []domain.MetricID{*id})
	if err != nil {
		return nil, ErrMetricGetByIDInternal
	}
	metric, found := metrics[*id]
	if !found {
		return nil, ErrMetricNotFound
	}
	return metric, nil
}
