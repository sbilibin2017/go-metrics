package usecases

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/validation"
)

type MetricUpdateBodyService interface {
	Update(ctx context.Context, metrics []*domain.Metric) ([]*domain.Metric, error)
}

type MetricUpdateBodyUsecase struct {
	svc MetricUpdateBodyService
}

func NewMetricUpdateBodyUsecase(svc MetricUpdateBodyService) *MetricUpdateBodyUsecase {
	return &MetricUpdateBodyUsecase{svc: svc}
}

func (uc *MetricUpdateBodyUsecase) Execute(
	ctx context.Context,
	req *MetricUpdateBodyRequest,
) (*MetricUpdateBodyResponse, error) {
	err := ValidateMetricUpdateBodyRequest(req)
	if err != nil {
		return nil, err
	}
	m := ConvertMetricUpdateBodyRequestToDomain(req)
	metric, err := uc.svc.Update(ctx, []*domain.Metric{m})
	if err != nil {
		return nil, err
	}
	return NewMetricUpdateBodyResponse(metric), nil
}

type MetricUpdateBodyRequest struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func ValidateMetricUpdateBodyRequest(req *MetricUpdateBodyRequest) error {
	err := validation.ValidateMetricID(req.ID)
	if err != nil {
		return err
	}
	err = validation.ValidateMetricType(req.Type)
	if err != nil {
		return err
	}
	if req.Type == string(domain.Counter) {
		err = validation.ValidateCounterPtrValue(req.Delta)
		if err != nil {
			return err
		}
	} else {
		err = validation.ValidateGaugePtrValue(req.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func ConvertMetricUpdateBodyRequestToDomain(req *MetricUpdateBodyRequest) *domain.Metric {
	m := domain.Metric{
		MetricID: domain.MetricID{
			ID:   req.ID,
			Type: domain.MetricType(req.Type),
		},
	}
	if req.Type == string(domain.Counter) {
		m.Delta = req.Delta
	} else {
		m.Value = req.Value
	}
	return &m
}

type MetricUpdateBodyResponse MetricUpdateBodyRequest

func NewMetricUpdateBodyResponse(metrics []*domain.Metric) *MetricUpdateBodyResponse {
	return &MetricUpdateBodyResponse{
		ID:    metrics[0].ID,
		Type:  string(metrics[0].Type),
		Delta: metrics[0].Delta,
		Value: metrics[0].Value,
	}
}
