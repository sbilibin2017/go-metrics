package handlers

import (
	"context"
	"go-metrics/internal/errors"
	"go-metrics/internal/usecases"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MetricGetByIDPathUsecase interface {
	Execute(ctx context.Context, req *usecases.MetricGetByIDPathRequest) (*usecases.MetricGetByIDPathResponse, error)
}

func MetricGetByIDPathHandler(uc MetricGetByIDPathUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req usecases.MetricGetByIDPathRequest
		req.Type = chi.URLParam(r, "type")
		req.Name = chi.URLParam(r, "name")
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
