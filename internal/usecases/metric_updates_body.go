package usecases

import (
	"context"
	"go-metrics/internal/domain"
)

type MetricUpdatesBodyService interface {
	Update(ctx context.Context, metrics []*domain.Metric) ([]*domain.Metric, error)
}

type MetricUpdatesBodyUsecase struct {
	svc MetricUpdatesBodyService
}

func NewMetricUpdatesBodyUsecase(svc MetricUpdatesBodyService) *MetricUpdatesBodyUsecase {
	return &MetricUpdatesBodyUsecase{svc: svc}
}

func (uc *MetricUpdatesBodyUsecase) Execute(
	ctx context.Context,
	req []*MetricUpdateBodyRequest,
) ([]*MetricUpdateBodyResponse, error) {
	for _, r := range req {
		err := ValidateMetricUpdateBodyRequest(r)
		if err != nil {
			return nil, err
		}
	}
	var metricsDomain []*domain.Metric
	for _, r := range req {
		m := ConvertMetricUpdateBodyRequestToDomain(r)
		metricsDomain = append(metricsDomain, m)
	}
	metrics, err := uc.svc.Update(ctx, metricsDomain)
	if err != nil {
		return nil, err
	}
	var metricsResponse []*MetricUpdateBodyResponse
	for _, m := range metrics {
		metricsResponse = append(metricsResponse, NewMetricUpdateBodyResponse([]*domain.Metric{m}))
	}
	return metricsResponse, nil
}
