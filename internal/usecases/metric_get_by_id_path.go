package usecases

import (
	"context"
	"go-metrics/internal/converters"
	"go-metrics/internal/domain"
	"go-metrics/internal/validation"
)

type MetricGetByIDPathService interface {
	GetByID(ctx context.Context, id *domain.MetricID) (*domain.Metric, error)
}

type MetricGetByIDPathUsecase struct {
	svc MetricGetByIDPathService
}

func NewMetricGetByIDPathUsecase(svc MetricGetByIDPathService) *MetricGetByIDPathUsecase {
	return &MetricGetByIDPathUsecase{svc: svc}
}

func (uc *MetricGetByIDPathUsecase) Execute(
	ctx context.Context,
	req *MetricGetByIDPathRequest,
) (*MetricGetByIDPathResponse, error) {
	err := ValidateMetricGetByIDPathRequest(req)
	if err != nil {
		return nil, err
	}
	metricRequest := ConvertMetricGetByIDPathRequestToDomain(req)
	metric, err := uc.svc.GetByID(ctx, metricRequest)
	if err != nil {
		return nil, err
	}
	return NewMetricGetByIDPathResponse(metric), nil
}

type MetricGetByIDPathRequest struct {
	Type string
	Name string
}

func ValidateMetricGetByIDPathRequest(req *MetricGetByIDPathRequest) error {
	err := validation.ValidateMetricID(req.Name)
	if err != nil {
		return err
	}
	err = validation.ValidateMetricType(req.Type)
	if err != nil {
		return err
	}
	return nil
}

func ConvertMetricGetByIDPathRequestToDomain(req *MetricGetByIDPathRequest) *domain.MetricID {
	return &domain.MetricID{
		ID:   req.Name,
		Type: domain.MetricType(req.Type),
	}
}

type MetricGetByIDPathResponse string

func NewMetricGetByIDPathResponse(metrics *domain.Metric) *MetricGetByIDPathResponse {
	if metrics.Type == domain.Counter {
		v := converters.FormatInt64(*metrics.Delta)
		s := MetricGetByIDPathResponse(v)
		return &s
	} else {
		v := converters.FormatFloat64(*metrics.Value)
		s := MetricGetByIDPathResponse(v)
		return &s
	}
}
