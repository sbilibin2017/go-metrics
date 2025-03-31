package handlers

import (
	"context"
	"go-metrics/internal/handlers/utils"
	"go-metrics/internal/requests"
	"go-metrics/internal/responses"
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
			handleMetricError(w, err)
			return
		}
		utils.SendTextResponse(w, http.StatusOK, resp.ToResponse())

	}
}
