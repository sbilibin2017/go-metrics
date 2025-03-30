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
	Execute(ctx context.Context, req *requests.MetricUpdatePathRequest) (*responses.MetricUpdatePathResponse, error)
}

func MetricUpdatePathHandler(uc MetricUpdatePathUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Type := chi.URLParam(r, "type")
		Name := chi.URLParam(r, "name")
		Value := chi.URLParam(r, "value")
		req := &requests.MetricUpdatePathRequest{
			Type:  Type,
			Name:  Name,
			Value: Value,
		}
		resp, err := uc.Execute(r.Context(), req)
		if err != nil {
			handleMetricUpdatePathError(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write(resp.ToResponse())
	}
}

func handleMetricUpdatePathError(w http.ResponseWriter, err error) {
	if errors.Is(err, validation.ErrEmptyName) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if errors.Is(err, validation.ErrInvalidType) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if errors.Is(err, validation.ErrEmptyValue) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if errors.Is(err, validation.ErrInvalidCounterValue) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if errors.Is(err, validation.ErrInvalidGaugeValue) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
