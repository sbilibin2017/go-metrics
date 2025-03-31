package handlers

import (
	"context"
	"errors"
	"go-metrics/internal/requests"
	"go-metrics/internal/responses"
	"go-metrics/internal/validation"
	"net/http"

	"github.com/go-chi/chi"
)

type MetricUpdatePathUsecase interface {
	Execute(
		ctx context.Context,
		req *requests.MetricUpdatePathRequest,
	) (*responses.MetricUpdatePathResponse, error)
}

func MetricUpdatePathHandler(uc MetricUpdatePathUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parsing URL parameters using chi
		req := parseMetricUpdateURLParams(r)

		// Executing usecase
		resp, err := uc.Execute(r.Context(), req)
		if err != nil {
			handleMetricUpdatePathError(w, err)
			return
		}

		// Writing the response
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write(resp.ToResponse())
	}
}

func parseMetricUpdateURLParams(r *http.Request) *requests.MetricUpdatePathRequest {
	// Using chi to parse URL parameters
	return requests.NewMetricUpdatePathRequest(
		chi.URLParam(r, "type"),
		chi.URLParam(r, "name"),
		chi.URLParam(r, "value"),
	)
}

func handleMetricUpdatePathError(w http.ResponseWriter, err error) {
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
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
