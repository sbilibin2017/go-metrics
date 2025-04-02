package app

import (
	"context"
	"database/sql"
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

	logger.Logger.Infow("Initializing server")

	metricUpdateHandler := handlers.MetricUpdatePathHandler(container.MetricUpdatePathUsecase)
	metricGetByIDHandler := handlers.MetricGetByIDPathHandler(container.MetricGetByIDPathUsecase)
	metricListHTMLHandler := handlers.MetricListHTMLHandler(container.MetricListHTMLUsecase)
	metricUpdateBodyHandler := handlers.MetricUpdateBodyHandler(container.MetricUpdateBodyUsecase)
	metricGetByIDBodyHandler := handlers.MetricGetByIDBodyHandler(container.MetricGetByIDBodyUsecase)
	metricUpdatesBodyHandler := handlers.MetricUpdatesBodyHandler(container.MetricUpdatesBodyUsecase)

	metricRouter := routers.NewMetricRouter(
		metricUpdateHandler,
		metricGetByIDHandler,
		metricListHTMLHandler,
		metricUpdateBodyHandler,
		metricGetByIDBodyHandler,
		metricUpdatesBodyHandler,
	)
	metricRouter.Get("/ping", PingDBHandler(container.DB))

	server := &http.Server{
		Addr:    config.GetAddress(),
		Handler: metricRouter,
	}

	logger.Logger.Infow("Server initialized", "address", config.GetAddress())

	return &Server{
		config:    config,
		container: container,
		server:    server,
		worker:    worker,
	}
}

func (s *Server) Start(ctx context.Context) error {
	logger.Logger.Infow("Starting server")
	defer func() {
		if err := s.container.File.Sync(); err != nil {
			logger.Logger.Errorw("Failed to sync file", "error", err)
		}
		if err := s.container.File.Close(); err != nil {
			logger.Logger.Errorw("Failed to close file", "error", err)
		}
	}()

	if s.config.GetDatabaseDSN() != "" {
		logger.Logger.Infow("Opening database connection", "dsn", s.config.GetDatabaseDSN())
		defer func() {
			if err := s.container.DB.Close(); err != nil {
				logger.Logger.Errorw("Failed to close DB connection", "error", err)
			} else {
				logger.Logger.Infow("Database connection closed successfully")
			}
		}()
		if err := CreateMetricTable(s.container.DB); err != nil {
			logger.Logger.Errorw("Failed to create metrics table", "error", err)
		}
	}

	go func() {
		logger.Logger.Infow("Starting HTTP server", "address", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Errorw("HTTP server error", "error", err)
		}
	}()

	go func() {
		logger.Logger.Infow("Starting worker")
		s.worker.Start(ctx)
	}()

	<-ctx.Done()

	logger.Logger.Infow("Shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.worker.Stop(shutdownCtx)
	s.server.Shutdown(shutdownCtx)

	logger.Logger.Infow("Server shutdown complete")
	return nil
}

func PingDBHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Infow("Ping request received")
		if db == nil {
			logger.Logger.Error("Database connection is nil")
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		if err := db.Ping(); err != nil {
			logger.Logger.Errorw("Database ping failed", "error", err)
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		logger.Logger.Infow("Database connection successful")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Database connection successful"))
	}
}

func CreateMetricTable(db *sql.DB) error {
	logger.Logger.Infow("Creating metrics table if not exists")
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS metrics (
		id VARCHAR(255) NOT NULL,
		type VARCHAR(255) NOT NULL,
		delta BIGINT,
		value DOUBLE PRECISION,
		PRIMARY KEY (id, type)
	);`)
	if err != nil {
		logger.Logger.Errorw("Failed to create metrics table", "error", err)
		return err
	}
	logger.Logger.Infow("Metrics table ready")
	return nil
}
