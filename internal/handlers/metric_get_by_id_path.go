package handlers

import (
	"context"
	"errors"
	"go-metrics/internal/requests"
	"go-metrics/internal/responses"
	"go-metrics/internal/services"
	"go-metrics/internal/validation"
	"net/http"

	"github.com/go-chi/chi"
)

type MetricGetByIDPathUsecase interface {
	Execute(ctx context.Context, req *requests.MetricGetByIDPathRequest) (*responses.MetricGetByIDPathResponse, error)
}

func MetricGetByIDPathHandler(uc MetricGetByIDPathUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := requests.NewMetricGetByIDPathRequest(
			chi.URLParam(r, "type"),
			chi.URLParam(r, "name"),
		)
		resp, err := uc.Execute(r.Context(), req)
		if err != nil {
			handleMetricGetByIDPathError(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write(resp.ToResponse())
	}
}

func handleMetricGetByIDPathError(w http.ResponseWriter, err error) {
	if errors.Is(err, validation.ErrEmptyName) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if errors.Is(err, validation.ErrInvalidType) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if errors.Is(err, services.ErrMetricNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
