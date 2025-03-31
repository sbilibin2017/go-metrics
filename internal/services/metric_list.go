package services

import (
	"context"
	"errors"
	"go-metrics/internal/domain"
	"sort"
)

type MetricListFindRepository interface {
	Find(ctx context.Context, filters []*domain.MetricID) (map[domain.MetricID]*domain.Metric, error)
}

type MetricListService struct {
	findRepo MetricListFindRepository
}

func NewMetricListService(
	findRepo MetricListFindRepository,
) *MetricListService {
	return &MetricListService{
		findRepo: findRepo,
	}
}

var (
	ErrMetricListInternal = errors.New("internal error")
)

func (s *MetricListService) List(
	ctx context.Context,
) ([]*domain.Metric, error) {
	metricsMap, err := s.findRepo.Find(ctx, []*domain.MetricID{})
	if err != nil {
		return nil, ErrMetricListInternal
	}
	var metrics []*domain.Metric
	for _, metric := range metricsMap {
		metrics = append(metrics, metric)
	}
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].ID < metrics[j].ID
	})
	return metrics, nil
}
