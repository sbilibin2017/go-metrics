package usecases

import (
	"context"
	"go-metrics/internal/domain"
	responses "go-metrics/internal/responses"
)

type MetricListHTMLService interface {
	List(ctx context.Context) ([]*domain.Metric, error)
}

type MetricListHTMLUsecase struct {
	svc MetricListHTMLService
}

func NewMetricListHTMLUsecase(
	svc MetricListHTMLService,
) *MetricListHTMLUsecase {
	return &MetricListHTMLUsecase{
		svc: svc,
	}
}

func (uc *MetricListHTMLUsecase) Execute(
	ctx context.Context,
) (*responses.MetricListHTMLResponse, error) {
	metrics, err := uc.svc.List(ctx)
	if err != nil {
		return nil, err
	}
	return responses.NewMetricListHTMLResponse(metrics), nil
}
