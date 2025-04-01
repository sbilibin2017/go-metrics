package app

import (
	"database/sql"
	"fmt"
	"go-metrics/internal/domain"
	"go-metrics/internal/handlers"
	"go-metrics/internal/logger"
	"go-metrics/internal/repositories"
	"go-metrics/internal/routers"
	"go-metrics/internal/services"
	"go-metrics/internal/unitofwork"
	"go-metrics/internal/usecases"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Container struct {
	File                 *os.File
	DB                   *sql.DB
	SaveFileRepo         *repositories.MetricFileSaveRepository
	FindFileRepo         *repositories.MetricFileFindRepository
	UOWFile              *unitofwork.FileUnitOfWork
	SaveDBRepo           *repositories.MetricDBSaveRepository
	FindDBRepo           *repositories.MetricDBFindRepository
	UOWDB                *unitofwork.DBUnitOfWork
	MetricUpdateService  *services.MetricUpdateService
	MetricGetByIDService *services.MetricGetByIDService
	MetricListService    *services.MetricListService
	MetricRouter         *chi.Mux
}

func NewContainer(config *Config) (*Container, error) {
	var file *os.File
	if config.GetFileStoragePath() != "" {
		logger.Logger.Infow("Opening file storage", "path", config.GetFileStoragePath())
		dir := filepath.Dir(config.GetFileStoragePath())
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Logger.Errorw("Failed to create directories", "error", err)
			return nil, err
		}
		var err error
		file, err = os.OpenFile(config.GetFileStoragePath(), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			logger.Logger.Errorw("Failed to open file", "error", err)
			return nil, err
		}
	}

	var db *sql.DB
	if config.GetDatabaseDSN() != "" {
		logger.Logger.Infow("Connecting to database", "dsn", config.GetDatabaseDSN())
		var err error
		db, err = sql.Open("pgx", config.GetDatabaseDSN())
		if err != nil {
			logger.Logger.Errorw("Failed to open database connection", "error", err)
			return nil, err
		}
		if err := db.Ping(); err != nil {
			logger.Logger.Errorw("Failed to ping database", "error", err)
			db.Close()
			return nil, err
		}
		logger.Logger.Infow("Database connection established successfully")
	}

	var saveFileRepo *repositories.MetricFileSaveRepository
	var findFileRepo *repositories.MetricFileFindRepository
	if file != nil {
		saveFileRepo = repositories.NewMetricFileSaveRepository(file)
		findFileRepo = repositories.NewMetricFileFindRepository(file)
	}

	var saveDBRepo *repositories.MetricDBSaveRepository
	var findDBRepo *repositories.MetricDBFindRepository
	if db != nil {
		saveDBRepo = repositories.NewMetricDBSaveRepository(db)
		findDBRepo = repositories.NewMetricDBFindRepository(db)
	}

	uowFile := unitofwork.NewFileUnitOfWork()
	var uowDB *unitofwork.DBUnitOfWork
	if db != nil {
		uowDB = unitofwork.NewDBUnitOfWork(db)
	}

	var metricUpdateService *services.MetricUpdateService
	var metricGetByIDService *services.MetricGetByIDService
	var metricListService *services.MetricListService

	if saveDBRepo != nil {
		logger.Logger.Infow("Using database storage")
		metricUpdateService = services.NewMetricUpdateService(saveDBRepo, findDBRepo, uowDB)
		metricGetByIDService = services.NewMetricGetByIDService(findDBRepo)
		metricListService = services.NewMetricListService(findDBRepo)
	} else if saveFileRepo != nil {
		logger.Logger.Infow("Using file storage")
		metricUpdateService = services.NewMetricUpdateService(saveFileRepo, findFileRepo, uowFile)
		metricGetByIDService = services.NewMetricGetByIDService(findFileRepo)
		metricListService = services.NewMetricListService(findFileRepo)
	} else {
		logger.Logger.Infow("Using memory storage")
		data := make(map[domain.MetricID]*domain.Metric)
		saveMemoryRepo := repositories.NewMetricMemorySaveRepository(data)
		findMemoryRepo := repositories.NewMetricMemoryFindRepository(data)
		uowMemory := unitofwork.NewMemoryUnitOfWork()
		metricUpdateService = services.NewMetricUpdateService(saveMemoryRepo, findMemoryRepo, uowMemory)
		metricGetByIDService = services.NewMetricGetByIDService(findMemoryRepo)
		metricListService = services.NewMetricListService(findMemoryRepo)
	}

	metricUpdatePathUsecase := usecases.NewMetricUpdatePathUsecase(metricUpdateService)
	metricGetByIDPathUsecase := usecases.NewMetricGetByIDPathUsecase(metricGetByIDService)
	metricListHTMLUsecase := usecases.NewMetricListHTMLUsecase(metricListService)
	metricUpdateBodyUsecase := usecases.NewMetricUpdateBodyUsecase(metricUpdateService)
	metricGetByIDBodyUsecase := usecases.NewMetricGetByIDBodyUsecase(metricGetByIDService)

	metricUpdateHandler := handlers.MetricUpdatePathHandler(metricUpdatePathUsecase)
	metricGetByIDHandler := handlers.MetricGetByIDPathHandler(metricGetByIDPathUsecase)
	metricListHTMLHandler := handlers.MetricListHTMLHandler(metricListHTMLUsecase)
	metricUpdateBodyHandler := handlers.MetricUpdateBodyHandler(metricUpdateBodyUsecase)
	metricGetByIDBodyHandler := handlers.MetricGetByIDBodyHandler(metricGetByIDBodyUsecase)

	metricRouter := routers.NewMetricRouter(
		metricUpdateHandler,
		metricGetByIDHandler,
		metricListHTMLHandler,
		metricUpdateBodyHandler,
		metricGetByIDBodyHandler,
	)
	metricRouter.Get("/ping", PingDBHandler(db))

	return &Container{
		File:                 file,
		DB:                   db,
		SaveFileRepo:         saveFileRepo,
		FindFileRepo:         findFileRepo,
		UOWFile:              uowFile,
		SaveDBRepo:           saveDBRepo,
		FindDBRepo:           findDBRepo,
		UOWDB:                uowDB,
		MetricUpdateService:  metricUpdateService,
		MetricGetByIDService: metricGetByIDService,
		MetricListService:    metricListService,
		MetricRouter:         metricRouter,
	}, nil
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
	}
}

func CreateMetricTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS metrics (
			id VARCHAR(255),
			type VARCHAR(50) NOT NULL,
			delta BIGINT,
			value DOUBLE PRECISION,
			PRIMARY KEY (id, type)
		);
	`
	_, err := db.Exec(query)
	if err != nil {
		logger.Logger.Errorw("failed to create table", "error", err)
		return fmt.Errorf("failed to create table: %w", err)
	}
	logger.Logger.Infow("metrics table created or already exists")
	return nil
}
