package handlers

import (
	"context"
	"go-metrics/internal/errors"
	"go-metrics/internal/usecases"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MetricUpdatePathUsecase interface {
	Execute(ctx context.Context, req *usecases.MetricUpdatePathRequest) (*usecases.MetricUpdatePathResponse, error)
}

func MetricUpdatePathHandler(uc MetricUpdatePathUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req usecases.MetricUpdatePathRequest
		req.Type = chi.URLParam(r, "type")
		req.Name = chi.URLParam(r, "name")
		req.Value = chi.URLParam(r, "value")
		resp, err := uc.Execute(r.Context(), &req)
		if err != nil {
			errors.MakeMetricErrorResponse(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(*resp)))
	}
}
