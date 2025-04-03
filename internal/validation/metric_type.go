package validation

import (
	"go-metrics/internal/domain"
	"go-metrics/internal/errors"
)

func ValidateMetricType(metricType string) error {
	mtype := domain.MetricType(metricType)
	if mtype != domain.Gauge && mtype != domain.Counter {
		return errors.ErrInvalidMetricType
	}
	return nil
}
