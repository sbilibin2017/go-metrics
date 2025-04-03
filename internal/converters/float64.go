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

func FormatFloat64(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
