package requests

import (
	"go-metrics/internal/domain"
	"go-metrics/internal/validation"
)

func Validate(metrics []*MetricUpdateBodyRequest) error {
	for _, m := range metrics {
		err := validation.ValidateType(m.Type)
		if err != nil {
			return err
		}
		err = validation.ValidateName(m.ID)
		if err != nil {
			return err
		}
		if m.Type == domain.Counter {
			err = validation.ValidateCounterPtrValue(m.Delta)
			if err != nil {
				return err
			}
		}
		if m.Type == domain.Gauge {
			err = validation.ValidateGaugePtrValue(m.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ToDomain(metrics []*MetricUpdateBodyRequest) []*domain.Metric {
	var domainMetrics []*domain.Metric
	for _, metric := range metrics {
		domainMetrics = append(domainMetrics, &domain.Metric{
			ID:    metric.ID,
			Type:  metric.Type,
			Delta: metric.Delta,
			Value: metric.Value,
		})
	}
	return domainMetrics
}
