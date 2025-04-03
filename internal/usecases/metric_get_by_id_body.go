package usecases

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/validation"
)

type MetricGetByIDBodyService interface {
	GetByID(ctx context.Context, id *domain.MetricID) (*domain.Metric, error)
}

type MetricGetByIDBodyUsecase struct {
	svc MetricGetByIDBodyService
}

func NewMetricGetByIDBodyUsecase(svc MetricGetByIDBodyService) *MetricGetByIDBodyUsecase {
	return &MetricGetByIDBodyUsecase{svc: svc}
}

func (uc *MetricGetByIDBodyUsecase) Execute(
	ctx context.Context,
	req *MetricGetByIDBodyRequest,
) (*MetricGetByIDBodyResponse, error) {
	err := ValidateMetricGetByIDBodyRequest(req)
	if err != nil {
		return nil, err
	}
	metricRequest := ConvertMetricGetByIDBodyRequestToDomain(req)
	metric, err := uc.svc.GetByID(ctx, metricRequest)
	if err != nil {
		return nil, err
	}
	return NewMetricGetByIDResponse(metric), nil
}

type MetricGetByIDBodyRequest struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func ValidateMetricGetByIDBodyRequest(req *MetricGetByIDBodyRequest) error {
	err := validation.ValidateMetricID(req.ID)
	if err != nil {
		return err
	}
	err = validation.ValidateMetricType(req.Type)
	if err != nil {
		return err
	}
	return nil
}

func ConvertMetricGetByIDBodyRequestToDomain(req *MetricGetByIDBodyRequest) *domain.MetricID {
	return &domain.MetricID{
		ID:   req.ID,
		Type: domain.MetricType(req.Type),
	}
}

type MetricGetByIDBodyResponse struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func NewMetricGetByIDResponse(metric *domain.Metric) *MetricGetByIDBodyResponse {
	return &MetricGetByIDBodyResponse{
		ID:    metric.ID,
		Type:  string(metric.Type),
		Delta: metric.Delta,
		Value: metric.Value,
	}
}
