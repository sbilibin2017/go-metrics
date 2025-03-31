package requests

import (
	"go-metrics/internal/converters"
	"go-metrics/internal/domain"
	"go-metrics/internal/validation"
)

type MetricUpdatePathRequest struct {
	mtype string
	name  string
	value string
}

func NewMetricUpdatePathRequest(
	mtype string,
	name string,
	value string,
) *MetricUpdatePathRequest {
	return &MetricUpdatePathRequest{
		mtype: mtype,
		name:  name,
		value: value,
	}
}

func (r *MetricUpdatePathRequest) Validate() error {
	err := validation.ValidateType(r.mtype)
	if err != nil {
		return err
	}
	err = validation.ValidateName(r.name)
	if err != nil {
		return err
	}
	if r.mtype == domain.Counter {
		err = validation.ValidateCounterValue(r.value)
		if err != nil {
			return err
		}
	} else if r.mtype == domain.Gauge {
		err = validation.ValidateGaugeValue(r.value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MetricUpdatePathRequest) ToDomain() ([]*domain.Metric, error) {
	var metrics []*domain.Metric
	metric := &domain.Metric{
		ID:    r.name,
		MType: r.mtype,
	}
	if r.mtype == domain.Gauge {
		value, err := converters.ConvertToFloat64(r.value)
		if err != nil {
			return nil, err
		}
		metric.Value = value
	} else if r.mtype == domain.Counter {
		delta, err := converters.ConvertToInt64(r.value)
		if err != nil {
			return nil, err
		}
		metric.Delta = delta
	}
	metrics = append(metrics, metric)
	return metrics, nil
}
