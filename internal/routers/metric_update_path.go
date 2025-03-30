package routers

import (
	"go-metrics/internal/router"

	"github.com/julienschmidt/httprouter"
)

func RegisterMetricUpdatePathRouter(
	r *router.Router,
	h httprouter.Handle,
) {
	r.AddHandler(
		"POST",
		"/update/:type/:name/:value",
		h,
	)
}
