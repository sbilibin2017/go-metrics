package converters

import (
	"errors"
	"strconv"
)

var ErrInvalidFloat64 = errors.New("invalid float64 value")

func ConvertToFloat64(value string) (*float64, error) {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, ErrInvalidFloat64
	}
	return &result, nil
}
