package server

import (
	"context"
	"go-metrics/internal/router"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Addresser interface {
	GetAddress() string
}

type Router interface {
	AddHandler(method, path string, handler httprouter.Handle)
	GetRoutes() []router.Route
	ServerHTTP() http.Handler
}

type HTTPServer struct {
	server *http.Server
	router Router
}

func NewHTTPServer(a Addresser) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:    a.GetAddress(),
			Handler: http.NewServeMux(),
		},
		router: router.NewRouter(),
	}
}

func (s *HTTPServer) AddRouter(rtr Router) {
	if s.router == nil {
		s.router = rtr
		s.server.Handler = rtr.ServerHTTP()
		return
	}
	for _, r := range rtr.GetRoutes() {
		s.router.AddHandler(r.Method, r.Path, r.Handler)
	}
	s.server.Handler = s.router.ServerHTTP()
}

func (s *HTTPServer) Start(ctx context.Context) error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	return nil
}
