package routers

import (
	"github.com/julienschmidt/httprouter"
)

type Router interface {
	AddHandler(method, path string, handler httprouter.Handle)
}

func RegisterMetricUpdatePathRouter(
	r Router,
	h httprouter.Handle,
) {
	r.AddHandler(
		"POST",
		"/update/:type/:name/:value",
		h,
	)
}
