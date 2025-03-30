package validation

import (
	"errors"
	"regexp"
)

var (
	ErrEmptyName   = errors.New("invalid name: cannot be empty")
	ErrInvalidName = errors.New("invalid name: only english letters are allowed")
)

func ValidateName(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	if !regexp.MustCompile("^[A-Za-z]+$").MatchString(name) {
		return ErrInvalidName
	}
	return nil
}
