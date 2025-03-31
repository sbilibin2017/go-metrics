package handlers

import (
	"context"
	"go-metrics/internal/handlers/utils"
	"go-metrics/internal/requests"
	"go-metrics/internal/responses"
	"net/http"
)

type MetricGetByIDBodyUsecase interface {
	Execute(ctx context.Context, req *requests.MetricGetByIDBodyRequest) (*responses.MetricGetByIDBodyResponse, error)
}

func MetricGetByIDBodyHandler(uc MetricGetByIDBodyUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.MetricGetByIDBodyRequest
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
