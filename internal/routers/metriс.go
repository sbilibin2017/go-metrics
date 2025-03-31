package routers

import (
	"net/http"

	"github.com/go-chi/chi"
)

func NewMetricRouter(
	h1 http.HandlerFunc,
	h2 http.HandlerFunc,
	h3 http.HandlerFunc,
) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", h1)
	r.Get("/value/{type}/{name}", h2)
	r.Get("/", h3)
	return r

}
