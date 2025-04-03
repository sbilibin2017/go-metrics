package usecases

import (
	"context"
	"go-metrics/internal/converters"
	"go-metrics/internal/domain"
	"go-metrics/internal/validation"
)

type MetricUpdatePathService interface {
	Update(ctx context.Context, metrics []*domain.Metric) ([]*domain.Metric, error)
}

type MetricUpdatePathUsecase struct {
	svc MetricUpdatePathService
}

func NewMetricUpdatePathUsecase(svc MetricUpdatePathService) *MetricUpdatePathUsecase {
	return &MetricUpdatePathUsecase{svc: svc}
}

func (uc *MetricUpdatePathUsecase) Execute(
	ctx context.Context,
	req *MetricUpdatePathRequest,
) (*MetricUpdatePathResponse, error) {
	err := ValidateMetricUpdatePathRequest(req)
	if err != nil {
		return nil, err
	}
	metric := ConvertMetricUpdatePathRequestToDomain(req)
	_, err = uc.svc.Update(ctx, []*domain.Metric{metric})
	if err != nil {
		return nil, err
	}
	return NewMetricUpdatePathResponse(), nil
}

type MetricUpdatePathRequest struct {
	Type  string
	Name  string
	Value string
}

func ValidateMetricUpdatePathRequest(req *MetricUpdatePathRequest) error {
	err := validation.ValidateMetricID(req.Name)
	if err != nil {
		return err
	}
	err = validation.ValidateMetricType(req.Type)
	if err != nil {
		return err
	}
	if req.Type == string(domain.Counter) {
		err = validation.ValidateCounterValue(req.Value)
		if err != nil {
			return err
		}
	} else {
		err = validation.ValidateGaugeValue(req.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func ConvertMetricUpdatePathRequestToDomain(req *MetricUpdatePathRequest) *domain.Metric {
	m := domain.Metric{
		MetricID: domain.MetricID{
			ID:   req.Name,
			Type: domain.MetricType(req.Type),
		},
	}
	if req.Type == string(domain.Counter) {
		v, _ := converters.ConvertToInt64(req.Value)
		m.Delta = v
	} else {
		v, _ := converters.ConvertToFloat64(req.Value)
		m.Value = v
	}
	return &m
}

type MetricUpdatePathResponse string

func NewMetricUpdatePathResponse() *MetricUpdatePathResponse {
	s := MetricUpdatePathResponse("Metric updated successfully")
	return &s
}
