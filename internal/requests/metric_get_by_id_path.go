package requests

import (
	"go-metrics/internal/domain"
	"go-metrics/internal/validation"
)

type MetricGetByIDPathRequest struct {
	Type string
	name string
}

func NewMetricGetByIDPathRequest(
	Type string,
	name string,
) *MetricGetByIDPathRequest {
	return &MetricGetByIDPathRequest{
		Type: Type,
		name: name,
	}
}

func (r *MetricGetByIDPathRequest) Validate() error {
	err := validation.ValidateType(r.Type)
	if err != nil {
		return err
	}
	err = validation.ValidateName(r.name)
	if err != nil {
		return err
	}
	return nil
}

func (r *MetricGetByIDPathRequest) ToDomain() (*domain.MetricID, error) {
	return &domain.MetricID{
		ID:   r.name,
		Type: r.Type,
	}, nil
}
