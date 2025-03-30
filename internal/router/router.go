package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	handler     *httprouter.Router
	routes      []Route
	middlewares []func(http.Handler) http.Handler
}

type Route struct {
	Method  string
	Path    string
	Handler httprouter.Handle
}

func NewRouter() *Router {
	return &Router{
		handler:     httprouter.New(),
		routes:      []Route{},
		middlewares: []func(http.Handler) http.Handler{},
	}
}

func (r *Router) AddHandler(method, path string, handler httprouter.Handle) {
	r.routes = append(r.routes, Route{Method: method, Path: path, Handler: handler})
	r.handler.Handle(method, path, handler)
}

func (r *Router) AddMiddleware(middleware func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) applyMiddlewares(h http.Handler) http.Handler {
	for _, middleware := range r.middlewares {
		h = middleware(h)
	}
	return h
}

func (r *Router) ServerHTTP() http.Handler {
	h := r.applyMiddlewares(r.handler)
	return h
}

func (r *Router) GetRoutes() []Route {
	return r.routes
}
