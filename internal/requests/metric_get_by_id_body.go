package requests

import (
	"go-metrics/internal/domain"
	"go-metrics/internal/validation"
)

type MetricGetByIDBodyRequest struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (r *MetricGetByIDBodyRequest) Validate() error {
	err := validation.ValidateType(r.Type)
	if err != nil {
		return err
	}
	err = validation.ValidateName(r.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *MetricGetByIDBodyRequest) ToDomain() (*domain.MetricID, error) {
	metricID := &domain.MetricID{
		ID:   r.ID,
		Type: r.Type,
	}
	return metricID, nil
}
