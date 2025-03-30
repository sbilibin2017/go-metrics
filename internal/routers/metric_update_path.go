package routers

import (
	"net/http"

	"github.com/go-chi/chi"
)

func RegisterMetricUpdatePathRouter(
	r chi.Router,
	h http.HandlerFunc,
) {
	r.Post("/update/{type}/{name}/{value}", h)
}
