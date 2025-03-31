package validation

import (
	"errors"
	"go-metrics/internal/domain"
	"strings"
)

var (
	ErrInvalidType = errors.New("invalid mtype: must be 'gauge' or 'counter'")
)

func ValidateType(mType string) error {
	mType = strings.ToLower(mType)
	if mType != domain.Gauge && mType != domain.Counter {
		return ErrInvalidType
	}
	return nil
}
