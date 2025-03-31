package validation

import (
	"errors"
	"go-metrics/internal/domain"
	"strings"
)

var (
	ErrInvalidType = errors.New("invalid Type: must be 'gauge' or 'counter'")
)

func ValidateType(Type string) error {
	Type = strings.ToLower(Type)
	if Type != domain.Gauge && Type != domain.Counter {
		return ErrInvalidType
	}
	return nil
}
