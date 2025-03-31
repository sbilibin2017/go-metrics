package usecases

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/requests"
	"go-metrics/internal/responses"
)

type MetricUpdatePathService interface {
	Update(ctx context.Context, metrics []*domain.Metric) ([]*domain.Metric, error)
}

type MetricUpdatePathUsecase struct {
	svc MetricUpdatePathService
}

func NewMetricUpdatePathUsecase(
	svc MetricUpdatePathService,
) *MetricUpdatePathUsecase {
	return &MetricUpdatePathUsecase{
		svc: svc,
	}
}

func (uc *MetricUpdatePathUsecase) Execute(
	ctx context.Context,
	req *requests.MetricUpdatePathRequest,
) (*responses.MetricUpdatePathResponse, error) {

	err := req.Validate()
	if err != nil {
		return nil, err
	}
	metrics, err := req.ToDomain()
	if err != nil {
		return nil, err
	}
	_, err = uc.svc.Update(ctx, metrics)
	if err != nil {
		return nil, err
	}
	return responses.NewMetricUpdatePathResponse(), nil
}
