package converters

import (
	"strconv"
)

func ConvertToInt64(value string) (*int64, error) {
	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}
	return &parsedValue, nil
}
