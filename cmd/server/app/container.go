package app

import (
	"database/sql"
	"go-metrics/internal/domain"
	"go-metrics/internal/logger"
	"go-metrics/internal/repositories"
	"go-metrics/internal/services"
	"go-metrics/internal/unitofwork"
	"go-metrics/internal/usecases"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Container struct {
	File                     *os.File
	DB                       *sql.DB
	SaveFileRepo             *repositories.MetricFileSaveRepository
	FindFileRepo             *repositories.MetricFileFindRepository
	SaveMemoryRepo           *repositories.MetricMemorySaveRepository
	FindMemoryRepo           *repositories.MetricMemoryFindRepository
	UOW                      *unitofwork.MemoryUnitOfWork
	MetricUpdateService      *services.MetricUpdateService
	MetricGetByIDService     *services.MetricGetByIDService
	MetricListService        *services.MetricListService
	MetricUpdatePathUsecase  *usecases.MetricUpdatePathUsecase
	MetricGetByIDPathUsecase *usecases.MetricGetByIDPathUsecase
	MetricListHTMLUsecase    *usecases.MetricListHTMLUsecase
	MetricUpdateBodyUsecase  *usecases.MetricUpdateBodyUsecase
	MetricGetByIDBodyUsecase *usecases.MetricGetByIDBodyUsecase
}

func NewContainer(config *Config) *Container {
	var file *os.File
	if config.GetFileStoragePath() != "" {
		logger.Logger.Infow("Opening file storage", "path", config.GetFileStoragePath())
		dir := filepath.Dir(config.GetFileStoragePath())
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Logger.Errorw("Failed to create directories", "error", err)
			return nil
		}
		var err error
		file, err = os.OpenFile(config.GetFileStoragePath(), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			logger.Logger.Errorw("Failed to open file", "error", err)
			return nil
		}
	}

	var db *sql.DB
	if config.GetDatabaseDSN() != "" {
		logger.Logger.Infow("Connecting to database", "dsn", config.GetDatabaseDSN())
		var err error
		db, err = sql.Open("pgx", config.GetDatabaseDSN())
		if err != nil {
			logger.Logger.Errorw("Failed to open database connection", "error", err)
			return nil
		}
		if err := db.Ping(); err != nil {
			logger.Logger.Errorw("Failed to ping database", "error", err)
			db.Close()
			return nil
		}
		logger.Logger.Infow("Database connection established successfully")
	}

	var saveFileRepo *repositories.MetricFileSaveRepository
	var findFileRepo *repositories.MetricFileFindRepository
	if file != nil {
		saveFileRepo = repositories.NewMetricFileSaveRepository(file)
		findFileRepo = repositories.NewMetricFileFindRepository(file)
	}

	data := make(map[domain.MetricID]*domain.Metric)
	saveMemoryRepo := repositories.NewMetricMemorySaveRepository(data)
	findMemoryRepo := repositories.NewMetricMemoryFindRepository(data)

	uow := unitofwork.NewMemoryUnitOfWork()

	metricUpdateService := services.NewMetricUpdateService(saveMemoryRepo, findMemoryRepo, uow)
	metricGetByIDService := services.NewMetricGetByIDService(findMemoryRepo)
	metricListService := services.NewMetricListService(findMemoryRepo)

	metricUpdatePathUsecase := usecases.NewMetricUpdatePathUsecase(metricUpdateService)
	metricGetByIDPathUsecase := usecases.NewMetricGetByIDPathUsecase(metricGetByIDService)
	metricListHTMLUsecase := usecases.NewMetricListHTMLUsecase(metricListService)
	metricUpdateBodyUsecase := usecases.NewMetricUpdateBodyUsecase(metricUpdateService)
	metricGetByIDBodyUsecase := usecases.NewMetricGetByIDBodyUsecase(metricGetByIDService)

	return &Container{
		File:                     file,
		DB:                       db,
		SaveFileRepo:             saveFileRepo,
		FindFileRepo:             findFileRepo,
		SaveMemoryRepo:           saveMemoryRepo,
		FindMemoryRepo:           findMemoryRepo,
		UOW:                      uow,
		MetricUpdateService:      metricUpdateService,
		MetricGetByIDService:     metricGetByIDService,
		MetricListService:        metricListService,
		MetricUpdatePathUsecase:  metricUpdatePathUsecase,
		MetricGetByIDPathUsecase: metricGetByIDPathUsecase,
		MetricListHTMLUsecase:    metricListHTMLUsecase,
		MetricUpdateBodyUsecase:  metricUpdateBodyUsecase,
		MetricGetByIDBodyUsecase: metricGetByIDBodyUsecase,
	}
}
