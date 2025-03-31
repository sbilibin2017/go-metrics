package handlers

import (
	"context"
	"go-metrics/internal/handlers/utils"
	"go-metrics/internal/requests"
	"go-metrics/internal/responses"
	"net/http"
)

type MetricUpdateBodyUsecase interface {
	Execute(
		ctx context.Context,
		req *requests.MetricUpdateBodyRequest,
	) (*responses.MetricUpdateBodyResponse, error)
}

func MetricUpdateBodyHandler(uc MetricUpdateBodyUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.MetricUpdateBodyRequest
		err := utils.ParseJSONRequest(r, &req)
		if err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := uc.Execute(r.Context(), &req)
		if err != nil {
			handleMetricError(w, err)
			return
		}
		utils.SendJSONResponse(w, http.StatusOK, resp)
	}
}
