package routers

import (
	"go-metrics/internal/middlewares"
	"net/http"

	"github.com/go-chi/chi"
)

func NewMetricRouter(
	h1 http.HandlerFunc,
	h2 http.HandlerFunc,
	h3 http.HandlerFunc,
	h4 http.HandlerFunc,
	h5 http.HandlerFunc,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)
	r.Post("/update/{type}/{name}/{value}", h1)
	r.Get("/value/{type}/{name}", h2)
	r.Get("/", h3)
	r.Post("/update/", h4)
	r.Post("/value/", h5)
	return r

}
