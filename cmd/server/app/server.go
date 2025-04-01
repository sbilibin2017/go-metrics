package app

import (
	"context"
	"go-metrics/internal/logger"
	"net/http"
	"time"
)

type Server struct {
	config    *Config
	container *Container
	worker    *Worker
	server    *http.Server
}

func NewServer(
	config *Config,
	container *Container,
	worker *Worker,
) *Server {
	server := &http.Server{
		Addr:    config.GetAddress(),
		Handler: container.MetricRouter,
	}
	return &Server{
		config:    config,
		container: container,
		worker:    worker,
		server:    server,
	}
}

func (s *Server) Start(ctx context.Context) {
	if s.config.GetFileStoragePath() != "" {
		go func() {
			logger.Logger.Infow("Starting worker")
			s.worker.Start(ctx)
		}()
	}

	if s.config.GetDatabaseDSN() != "" {
		CreateMetricTable(s.container.DB)
	}

	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Logger.Errorw("Server error", "error", err)
		}
	}()

	logger.Logger.Infow("Server is running, waiting for shutdown signal...")

	<-ctx.Done()

	logger.Logger.Infow("Shutdown signal received, starting server shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.worker.Save(shutdownCtx)

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Errorw("Error during server shutdown", "error", err)
	} else {
		logger.Logger.Infow("Server shutdown completed successfully")
	}
}
