package usecases

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/requests"
	"go-metrics/internal/responses"
)

type MetricUpdatesBodyService interface {
	Update(ctx context.Context, metrics []*domain.Metric) ([]*domain.Metric, error)
}

type MetricUpdatesBodyUsecase struct {
	svc MetricUpdatesBodyService
}

func NewMetricUpdatesBodyUsecase(
	svc MetricUpdatesBodyService,
) *MetricUpdatesBodyUsecase {
	return &MetricUpdatesBodyUsecase{
		svc: svc,
	}
}

func (uc *MetricUpdatesBodyUsecase) Execute(
	ctx context.Context,
	req []*requests.MetricUpdateBodyRequest,
) ([]*responses.MetricUpdateBodyResponse, error) {
	err := requests.Validate(req)
	if err != nil {
		return nil, err
	}
	metrics, err := uc.svc.Update(ctx, requests.ToDomain(req))
	if err != nil {
		return nil, err
	}
	return responses.NewMetricUpdatesBodyResponse(metrics), nil
}
