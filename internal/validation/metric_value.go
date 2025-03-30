package validation

import (
	"errors"
	"go-metrics/internal/converters"
)

var (
	ErrEmptyValue          = errors.New("invalid value: cannot be empty")
	ErrInvalidCounterValue = errors.New("invalid 'counter' value: must be int64")
	ErrInvalidGaugeValue   = errors.New("invalid 'gauge' value: must be float64")
)

func ValidateValue(value string) error {
	if value == "" {
		return ErrEmptyValue
	}
	return nil
}

func ValidateCounterValueString(value string) error {
	_, err := converters.ConvertToInt64(value)
	if err != nil {
		return ErrInvalidCounterValue
	}
	return nil
}

func ValidateGaugeValueString(value string) error {
	_, err := converters.ConvertToFloat64(value)
	if err != nil {
		return ErrInvalidGaugeValue
	}
	return nil
}
