package handlers

import (
	"errors"
	"go-metrics/internal/services"
	"go-metrics/internal/validation"
	"net/http"
)

func handleMetricError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, validation.ErrEmptyName):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, validation.ErrInvalidType):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, validation.ErrEmptyValue):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, validation.ErrInvalidCounterValue):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, validation.ErrInvalidGaugeValue):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, services.ErrMetricNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
