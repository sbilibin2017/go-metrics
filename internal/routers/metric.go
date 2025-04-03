package routers

import (
	"go-metrics/internal/middlewares"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewMetricRouter(
	h1 http.HandlerFunc,
	h2 http.HandlerFunc,
	h3 http.HandlerFunc,
	h4 http.HandlerFunc,
	h5 http.HandlerFunc,
	h6 http.HandlerFunc,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)
	r.Post("/update/{type}/{name}/{value}", h1)
	r.Post("/update/", h2)
	r.Post("/updates/", h3)
	r.Get("/value/{type}/{name}", h4)
	r.Post("/value/", h5)
	r.Get("/", h6)
	return r

}
