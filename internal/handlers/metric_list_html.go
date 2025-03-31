package handlers

import (
	"context"
	"go-metrics/internal/responses"
	"log"
	"net/http"
)

type MetricListHTMLUsecase interface {
	Execute(ctx context.Context) (*responses.MetricListHTMLResponse, error)
}

func MetricListHTMLHandler(uc MetricListHTMLUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		resp, err := uc.Execute(r.Context())
		if err != nil {
			log.Printf("Error processing request for %s: %v", r.URL.Path, err)
			handleMetricListHTMLError(w, err)
			return
		}
		log.Printf("Successfully processed request for %s", r.URL.Path)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(resp.ToResponse())
	}
}

func handleMetricListHTMLError(w http.ResponseWriter, err error) {
	// Error handling (Internal Server Error)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
