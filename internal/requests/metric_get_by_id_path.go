package requests

import (
	"go-metrics/internal/domain"
	"go-metrics/internal/validation"
)

type MetricGetByIDPathRequest struct {
	Type string
	Name string
}

func NewMetricGetByIDPathRequest(
	mtype string,
	name string,
) *MetricGetByIDPathRequest {
	return &MetricGetByIDPathRequest{
		Type: mtype,
		Name: name,
	}
}

func (r *MetricGetByIDPathRequest) Validate() error {
	err := validation.ValidateType(r.Type)
	if err != nil {
		return err
	}
	err = validation.ValidateName(r.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *MetricGetByIDPathRequest) ToDomain() (*domain.MetricID, error) {
	return &domain.MetricID{
		ID:   r.Name,
		Type: r.Type,
	}, nil
}
