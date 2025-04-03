package errors

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidMetricID           = errors.New("invalid id: only letters and numbers are allowed")
	ErrInvalidMetricType         = errors.New("invalid Type: must be 'gauge' or 'counter'")
	ErrEmptyMetricValue          = errors.New("invalid metric value: cannot be empty")
	ErrInvalidCounterMetricValue = errors.New("invalid 'counter' metric value: must be int64")
	ErrInvalidGaugeMetricValue   = errors.New("invalid 'gauge' metric value: must be float64")
	ErrMetricNotFound            = errors.New("metric not found")
	ErrMetricGetByIDInternal     = errors.New("internal error")
	ErrMetricListInternal        = errors.New("internal error")
	ErrMetricIsNotUpdated        = errors.New("metric is not updated")
)

func MakeMetricErrorResponse(w http.ResponseWriter, err error) {
	switch err {
	case ErrInvalidMetricID:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case ErrInvalidMetricType:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case ErrEmptyMetricValue:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case ErrInvalidCounterMetricValue:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case ErrInvalidGaugeMetricValue:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case ErrMetricNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	case ErrMetricGetByIDInternal, ErrMetricListInternal, ErrMetricIsNotUpdated:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
