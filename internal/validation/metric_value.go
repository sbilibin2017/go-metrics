package validation

import (
	"go-metrics/internal/converters"
	"go-metrics/internal/errors"
)

func ValidateCounterValue(value string) error {
	_, err := converters.ConvertToInt64(value)
	if err != nil {
		return errors.ErrInvalidCounterMetricValue
	}
	return nil
}

func ValidateGaugeValue(value string) error {
	_, err := converters.ConvertToFloat64(value)
	if err != nil {
		return errors.ErrInvalidGaugeMetricValue
	}
	return nil
}

func ValidateCounterPtrValue(value *int64) error {
	if value == nil {
		return errors.ErrInvalidCounterMetricValue
	}
	return nil
}

func ValidateGaugePtrValue(value *float64) error {
	if value == nil {
		return errors.ErrInvalidGaugeMetricValue
	}
	return nil
}
