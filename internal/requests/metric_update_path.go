package requests

import (
	"go-metrics/internal/converters"
	"go-metrics/internal/domain"
	"go-metrics/internal/validation"
)

type MetricUpdatePathRequest struct {
	Type  string
	Name  string
	Value string
}

func (r *MetricUpdatePathRequest) Validate() error {
	err := validation.ValidateType(r.Type)
	if err != nil {
		return err
	}
	err = validation.ValidateName(r.Name)
	if err != nil {
		return err
	}
	err = validation.ValidateValue(r.Value)
	if err != nil {
		return err
	}
	if r.Type == domain.Counter {
		err = validation.ValidateCounterValueString(r.Value)
		if err != nil {
			return err
		}
	} else if r.Type == domain.Gauge {
		err = validation.ValidateGaugeValueString(r.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MetricUpdatePathRequest) ToDomain() (*domain.Metric, error) {
	metric := &domain.Metric{
		ID:   r.Name,
		Type: r.Type,
	}
	if r.Type == domain.Gauge {
		value, err := converters.ConvertToFloat64(r.Value)
		if err != nil {
			return nil, err
		}
		metric.Value = value
	} else if r.Type == domain.Counter {
		delta, err := converters.ConvertToInt64(r.Value)
		if err != nil {
			return nil, err
		}
		metric.Delta = delta
	}
	return metric, nil
}
