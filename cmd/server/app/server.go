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
	metricUpdateHandler := handlers.MetricUpdatePathHandler(container.MetricUpdatePathUsecase)
	metricGetByIDHandler := handlers.MetricGetByIDPathHandler(container.MetricGetByIDPathUsecase)
	metricListHTMLHandler := handlers.MetricListHTMLHandler(container.MetricListHTMLUsecase)
	metricUpdateBodyHandler := handlers.MetricUpdateBodyHandler(container.MetricUpdateBodyUsecase)
	metricGetByIDBodyHandler := handlers.MetricGetByIDBodyHandler(container.MetricGetByIDBodyUsecase)
	pingDBHandler := handlers.PingDBHandler(container.DBEngine)

	metricRouter := routers.NewMetricRouter(
		metricUpdateHandler,
		metricGetByIDHandler,
		metricListHTMLHandler,
		metricUpdateBodyHandler,
		metricGetByIDBodyHandler,
	)
	metricRouter.Get("/ping", pingDBHandler)

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
	logger.Init()
	defer logger.Logger.Sync()

	logger.Logger.Infow("Server starting", "address", s.config.GetAddress())

	if s.config.GetFileStoragePath() != "" {
		logger.Logger.Infow("Opening file storage", "path", s.config.GetFileStoragePath())
		if err := s.container.FileEngine.Open(ctx, s.config); err != nil {
			logger.Logger.Errorw("failed to open file", "error", err)
			return err
		}
		defer func() {
			if err := s.container.FileEngine.Sync(); err != nil {
				logger.Logger.Errorw("failed to sync file", "error", err)
			} else {
				logger.Logger.Infow("File synced successfully")
			}
			if err := s.container.FileEngine.Close(); err != nil {
				logger.Logger.Errorw("failed to close file", "error", err)
			} else {
				logger.Logger.Infow("File closed successfully")
			}
		}()
	}

	if s.config.GetDatabaseDSN() != "" {
		logger.Logger.Infow("Opening database connection", "dsn", s.config.GetDatabaseDSN())
		err := s.container.DBEngine.Open(ctx, s.config)
		if err != nil {
			logger.Logger.Errorw("failed to connect to db", "error", err)
			return err
		} else {
			logger.Logger.Infow("Database connection established successfully")
		}
		defer func() {
			if err := s.container.DBEngine.Close(); err != nil {
				logger.Logger.Errorw("failed to close db connection", "error", err)
			} else {
				logger.Logger.Infow("Database connection closed successfully")
			}
		}()
		go func() {
			logger.Logger.Infow("Starting worker")
			s.worker.Start(ctx)
		}()
	}

	go func() error {
		logger.Logger.Infow("Starting HTTP server", "address", s.config.GetAddress())
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Logger.Errorw("HTTP server error", "error", err)
			return err
		}
		return nil
	}()

	<-ctx.Done()

	logger.Logger.Infow("Server shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Errorw("Server shutdown failed", "error", err)
	} else {
		logger.Logger.Infow("Server shutdown successfully")
	}

	return nil
}
