package validation

import (
	"errors"
	"regexp"
)

var (
	ErrEmptyName   = errors.New("invalid name: cannot be empty")
	ErrInvalidName = errors.New("invalid name: only letters and numbers are allowed")
)

func ValidateName(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	if !regexp.MustCompile("^[A-Za-z0-9]+$").MatchString(name) {
		return ErrInvalidName
	}
	return nil
}
