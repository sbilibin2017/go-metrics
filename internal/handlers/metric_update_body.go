package handlers

import (
	"context"
	"encoding/json"
	"go-metrics/internal/errors"
	"go-metrics/internal/usecases"
	"net/http"
)

type MetricUpdateBodyUsecase interface {
	Execute(ctx context.Context, req *usecases.MetricUpdateBodyRequest) (*usecases.MetricUpdateBodyResponse, error)
}

func MetricUpdateBodyHandler(uc MetricUpdateBodyUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req usecases.MetricUpdateBodyRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		resp, err := uc.Execute(r.Context(), &req)
		if err != nil {
			errors.MakeMetricErrorResponse(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			errors.MakeMetricErrorResponse(w, err)
			return
		}
	}
}
