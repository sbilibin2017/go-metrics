package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type Addresser interface {
	GetAddress() string
}

type Server struct {
	server HTTPServer
}

func NewServer(a Addresser, r *chi.Mux) *Server {
	return &Server{
		server: &http.Server{
			Addr:    a.GetAddress(),
			Handler: r,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() error {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.server.Shutdown(shutdownCtx)
	return nil
}
