package validation

import (
	"errors"
	"go-metrics/internal/domain"
)

var (
	ErrInvalidType = errors.New("invalid mtype: must be 'gauge' or 'counter'")
)

func ValidateType(mType string) error {
	if mType == "" || (mType != domain.Gauge && mType != domain.Counter) {
		return ErrInvalidType
	}
	return nil
}
