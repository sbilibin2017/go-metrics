package usecases

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/requests"
	"go-metrics/internal/responses"
)

type MetricUpdateBodyService interface {
	Update(ctx context.Context, metrics []*domain.Metric) ([]*domain.Metric, error)
}

type MetricUpdateBodyUsecase struct {
	svc MetricUpdateBodyService
}

func NewMetricUpdateBodyUsecase(
	svc MetricUpdateBodyService,
) *MetricUpdateBodyUsecase {
	return &MetricUpdateBodyUsecase{
		svc: svc,
	}
}

func (uc *MetricUpdateBodyUsecase) Execute(
	ctx context.Context,
	req *requests.MetricUpdateBodyRequest,
) (*responses.MetricUpdateBodyResponse, error) {

	err := req.Validate()
	if err != nil {
		return nil, err
	}
	metrics, err := req.ToDomain()
	if err != nil {
		return nil, err
	}
	metrics, err = uc.svc.Update(ctx, metrics)
	if err != nil {
		return nil, err
	}
	return responses.NewMetricUpdateBodyResponse(metrics), nil
}
