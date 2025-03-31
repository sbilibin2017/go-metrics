package requests

import (
	"go-metrics/internal/domain"
	"go-metrics/internal/validation"
)

type MetricUpdateBodyRequest struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (r *MetricUpdateBodyRequest) Validate() error {
	err := validation.ValidateType(r.Type)
	if err != nil {
		return err
	}
	err = validation.ValidateName(r.ID)
	if err != nil {
		return err
	}
	if r.Type == domain.Counter {
		err = validation.ValidateCounterPtrValue(r.Delta)
		if err != nil {
			return err
		}
	}
	if r.Type == domain.Gauge {
		err = validation.ValidateGaugePtrValue(r.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MetricUpdateBodyRequest) ToDomain() ([]*domain.Metric, error) {
	metrics := []*domain.Metric{
		{
			ID:    r.ID,
			Type:  r.Type,
			Delta: r.Delta,
			Value: r.Value,
		},
	}
	return metrics, nil
}
