package converters

import (
	"strconv"
)

func ConvertToFloat64(value string) (*float64, error) {
	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, err
	}
	return &parsedValue, nil
}
