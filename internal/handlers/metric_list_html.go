package handlers

import (
	"context"
	"go-metrics/internal/errors"
	"go-metrics/internal/usecases"
	"net/http"
)

type MetricListHTMLUsecase interface {
	Execute(ctx context.Context) (*usecases.MetricListHTMLResponse, error)
}

func MetricListHTMLHandler(uc MetricListHTMLUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := uc.Execute(r.Context())
		if err != nil {
			errors.MakeMetricErrorResponse(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(resp.HTML))
		if err != nil {
			errors.MakeMetricErrorResponse(w, err)
		}
	}
}
