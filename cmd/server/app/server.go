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
	metricRouter.Get("/ping", PingDBHandler(container.DB))

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
	defer func() {
		if err := s.container.File.Sync(); err != nil {
			logger.Logger.Errorw("failed to sync file", "error", err)
		}
		if err := s.container.File.Close(); err != nil {
			logger.Logger.Errorw("failed to close file", "error", err)
		}
	}()
	go func() {
		s.worker.Start(ctx)
	}()

	if s.config.GetDatabaseDSN() != "" {
		logger.Logger.Infow("Opening database connection", "dsn", s.config.GetDatabaseDSN())
		defer func() {
			if err := s.container.DB.Close(); err != nil {
				logger.Logger.Errorw("failed to close db connection", "error", err)
			} else {
				logger.Logger.Infow("Database connection closed successfully")
			}
		}()
		CreateMetricTable(s.container.DB)
	}

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

	s.worker.Save(shutdownCtx)

	s.server.Shutdown(shutdownCtx)

	return nil
}

func PingDBHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		if err := db.Ping(); err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Database connection successful"))
	}
}

func CreateMetricTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS metrics (
		id VARCHAR(255) NOT NULL,             
		type VARCHAR(255) NOT NULL,  
		delta BIGINT,                
		value DOUBLE PRECISION,      
		PRIMARY KEY (id, type)     
	);`)
	if err != nil {
		return err
	}
	return nil
}
