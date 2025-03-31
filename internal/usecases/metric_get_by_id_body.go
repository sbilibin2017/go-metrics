package usecases

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/requests"
	"go-metrics/internal/responses"
)

type MetricGetByIDBodyService interface {
	GetByID(ctx context.Context, id *domain.MetricID) (*domain.Metric, error)
}

type MetricGetByIDBodyUsecase struct {
	svc MetricGetByIDPathService
}

func NewMetricGetByIDBodyUsecase(
	svc MetricGetByIDBodyService,
) *MetricGetByIDBodyUsecase {
	return &MetricGetByIDBodyUsecase{
		svc: svc,
	}
}

func (uc *MetricGetByIDBodyUsecase) Execute(
	ctx context.Context, req *requests.MetricGetByIDBodyRequest,
) (*responses.MetricGetByIDBodyResponse, error) {
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
	return responses.NewMetricGetByIDBodyResponse(metric), nil
}
