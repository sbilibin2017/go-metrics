package converters

import (
	"errors"
	"strconv"
)

var ErrInvalidInt64 = errors.New("invalid int64 value")

func ConvertToInt64(value string) (*int64, error) {
	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, ErrInvalidInt64
	}
	return &parsedValue, nil
}

func FormatInt64(value int64) string {
	return strconv.FormatInt(value, 10)
}
