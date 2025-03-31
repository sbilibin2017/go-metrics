package usecases

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/requests"
	"go-metrics/internal/responses"
)

type MetricGetByIDPathService interface {
	GetByID(ctx context.Context, id *domain.MetricID) (*domain.Metric, error)
}

type MetricGetByIDPathUsecase struct {
	svc MetricGetByIDPathService
}

func NewMetricGetByIDPathUsecase(
	svc MetricGetByIDPathService,
) *MetricGetByIDPathUsecase {
	return &MetricGetByIDPathUsecase{
		svc: svc,
	}
}

func (uc *MetricGetByIDPathUsecase) Execute(
	ctx context.Context, req *requests.MetricGetByIDPathRequest,
) (*responses.MetricGetByIDPathResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	metricIDs, err := req.ToDomain()
	if err != nil {
		return nil, err
	}
	metric, err := uc.svc.GetByID(ctx, metricIDs)
	if err != nil {
		return nil, err
	}
	return responses.NewMetricGetByIDPathResponse(metric), nil
}
