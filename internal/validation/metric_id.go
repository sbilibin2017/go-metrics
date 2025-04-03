package validation

import (
	"go-metrics/internal/errors"
	"regexp"
)

func ValidateMetricID(id string) error {
	if id == "" {
		return errors.ErrInvalidMetricID
	}
	if !regexp.MustCompile("^[A-Za-z0-9_]+$").MatchString(id) {
		return errors.ErrInvalidMetricID
	}
	return nil
}
