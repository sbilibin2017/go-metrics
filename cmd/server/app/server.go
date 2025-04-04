package app

import (
	"context"
	"database/sql"
	"go-metrics/internal/handlers"
	"go-metrics/internal/routers"
	"go-metrics/pkg/log"
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
	log.Init(log.LevelInfo)
	defer log.Sync()

	log.Info("Initializing server")

	metricUpdateHandler := handlers.MetricUpdatePathHandler(container.MetricUpdatePathUsecase)
	metricGetByIDHandler := handlers.MetricGetByIDPathHandler(container.MetricGetByIDPathUsecase)
	metricListHTMLHandler := handlers.MetricListHTMLHandler(container.MetricListHTMLUsecase)
	metricUpdateBodyHandler := handlers.MetricUpdateBodyHandler(container.MetricUpdateBodyUsecase)
	metricGetByIDBodyHandler := handlers.MetricGetByIDBodyHandler(container.MetricGetByIDBodyUsecase)
	metricUpdatesHandler := handlers.MetricUpdatesBodyHandler(container.MetricUpdatesBodyUsecase)

	metricRouter := routers.NewMetricRouter(
		config,
		metricUpdateHandler,
		metricUpdateBodyHandler,
		metricUpdatesHandler,
		metricGetByIDHandler,
		metricGetByIDBodyHandler,
		metricListHTMLHandler,
	)
	metricRouter.Get("/ping", PingDBHandler(container.DB))

	server := &http.Server{
		Addr:    config.GetAddress(),
		Handler: metricRouter,
	}

	log.Info("Server initialized", "address", config.GetAddress())

	return &Server{
		config:    config,
		container: container,
		server:    server,
		worker:    worker,
	}
}

func (s *Server) Start(ctx context.Context) error {
	log.Info("Starting server")
	defer func() {
		if err := s.container.File.Sync(); err != nil {
			log.Error("Failed to sync file", "error", err)
		}
		if err := s.container.File.Close(); err != nil {
			log.Error("Failed to close file", "error", err)
		}
	}()

	if s.config.GetDatabaseDSN() != "" {
		log.Info("Opening database connection", "dsn", s.config.GetDatabaseDSN())
		defer func() {
			if err := s.container.DB.Close(); err != nil {
				log.Error("Failed to close DB connection", "error", err)
			} else {
				log.Info("Database connection closed successfully")
			}
		}()
		if err := CreateMetricTable(s.container.DB); err != nil {
			log.Error("Failed to create metrics table", "error", err)
		}
	}

	go func() {
		log.Info("Starting HTTP server", "address", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server error", "error", err)
		}
	}()

	go func() {
		log.Info("Starting worker")
		s.worker.Start(ctx)
	}()

	<-ctx.Done()

	log.Info("Shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.worker.Stop(shutdownCtx)
	s.server.Shutdown(shutdownCtx)

	log.Info("Server shutdown complete")
	return nil
}

func PingDBHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Ping request received")
		if db == nil {
			log.Error("Database connection is nil")
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		if err := db.Ping(); err != nil {
			log.Error("Database ping failed", "error", err)
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		log.Info("Database connection successful")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Database connection successful"))
	}
}

func CreateMetricTable(db *sql.DB) error {
	log.Info("Creating metrics table if not exists")
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS metrics (
		id VARCHAR(255) NOT NULL,
		type VARCHAR(255) NOT NULL,
		delta BIGINT,
		value DOUBLE PRECISION,
		PRIMARY KEY (id, type)
	);`)
	if err != nil {
		log.Error("Failed to create metrics table", "error", err)
		return err
	}
	log.Info("Metrics table ready")
	return nil
}
