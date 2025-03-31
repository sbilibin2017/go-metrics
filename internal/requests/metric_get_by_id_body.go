package requests

import (
	"go-metrics/internal/domain"
	"go-metrics/internal/validation"
)

type MetricGetByIDBodyRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

func (r *MetricGetByIDBodyRequest) Validate() error {
	err := validation.ValidateType(r.MType)
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
		ID:    r.ID,
		MType: r.MType,
	}
	return metricID, nil
}
