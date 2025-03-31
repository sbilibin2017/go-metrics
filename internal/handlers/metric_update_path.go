package handlers

import (
	"context"
	"go-metrics/internal/handlers/utils"
	"go-metrics/internal/requests"
	"go-metrics/internal/responses"
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
		req := parseMetricUpdateURLParams(r)
		resp, err := uc.Execute(r.Context(), req)
		if err != nil {
			handleMetricError(w, err)
			return
		}
		utils.SendTextResponse(w, http.StatusOK, resp.ToResponse())
	}
}

func parseMetricUpdateURLParams(r *http.Request) *requests.MetricUpdatePathRequest {
	return requests.NewMetricUpdatePathRequest(
		chi.URLParam(r, "type"),
		chi.URLParam(r, "name"),
		chi.URLParam(r, "value"),
	)
}
