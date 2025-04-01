package app

import (
	"context"
	"go-metrics/internal/handlers"
	"go-metrics/internal/logger"
	"go-metrics/internal/routers"
	"net/http"
	"time"
)

type Server struct {
	config    *Config
	container *Container
	server    *http.Server
	worker    *Worker
}

func NewServer(config *Config, container *Container, worker *Worker) *Server {
	logger.Init()
	defer logger.Logger.Sync()

	metricUpdateHandler := handlers.MetricUpdatePathHandler(container.MetricUpdatePathUsecase)
	metricGetByIDHandler := handlers.MetricGetByIDPathHandler(container.MetricGetByIDPathUsecase)
	metricListHTMLHandler := handlers.MetricListHTMLHandler(container.MetricListHTMLUsecase)
	metricUpdateBodyHandler := handlers.MetricUpdateBodyHandler(container.MetricUpdateBodyUsecase)
	metricGetByIDBodyHandler := handlers.MetricGetByIDBodyHandler(container.MetricGetByIDBodyUsecase)

	metricRouter := routers.NewMetricRouter(
		metricUpdateHandler,
		metricGetByIDHandler,
		metricListHTMLHandler,
		metricUpdateBodyHandler,
		metricGetByIDBodyHandler,
	)

	server := &http.Server{
		Addr:    config.GetAddress(),
		Handler: metricRouter,
	}

	return &Server{
		config:    config,
		container: container,
		server:    server,
		worker:    worker,
	}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.container.FileEngine.Open(s.config); err != nil {
		logger.Logger.Errorw("failed to open file", "error", err)
		return err
	}

	defer func() {

		if err := s.container.FileEngine.Sync(); err != nil {
			logger.Logger.Errorw("failed to sync file", "error", err)
		}

		if err := s.container.FileEngine.Close(); err != nil {
			logger.Logger.Errorw("failed to close file", "error", err)
		}
	}()

	go func() {
		s.worker.Start(ctx)
	}()

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
