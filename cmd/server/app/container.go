package app

import (
	"bufio"
	"database/sql"
	"go-metrics/internal/domain"
	"go-metrics/internal/repositories"
	"go-metrics/internal/services"
	"go-metrics/internal/unitofworks"
	"go-metrics/internal/usecases"
	"go-metrics/pkg/log"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Container struct {
	DB                       *sql.DB
	File                     *os.File
	Memory                   map[domain.MetricID]*domain.Metric
	MetricSaveDBRepo         *repositories.MetricDBSaveRepository
	MetricFindDBRepo         *repositories.MetricDBFindRepository
	MetricSaveFileRepo       *repositories.MetricFileSaveRepository
	MetricFindFileRepo       *repositories.MetricFileFindRepository
	MetricSaveMemoryRepo     *repositories.MetricMemorySaveRepository
	MetricFindMemoryRepo     *repositories.MetricMemoryFindRepository
	DBUOW                    *unitofworks.DBUnitOfWork
	FileUOW                  *unitofworks.FileUnitOfWork
	MemoryUOW                *unitofworks.MemoryUnitOfWork
	MetricUpdateService      *services.MetricUpdateService
	MetricGetByIDService     *services.MetricGetByIDService
	MetricListService        *services.MetricListService
	MetricUpdatePathUsecase  *usecases.MetricUpdatePathUsecase
	MetricGetByIDPathUsecase *usecases.MetricGetByIDPathUsecase
	MetricListHTMLUsecase    *usecases.MetricListHTMLUsecase
	MetricUpdateBodyUsecase  *usecases.MetricUpdateBodyUsecase
	MetricGetByIDBodyUsecase *usecases.MetricGetByIDBodyUsecase
	MetricUpdatesBodyUsecase *usecases.MetricUpdatesBodyUsecase
}

func NewContainer(config *Config) (*Container, error) {
	container := &Container{
		Memory: make(map[domain.MetricID]*domain.Metric),
	}
	if dsn := config.GetDatabaseDSN(); dsn != "" {
		log.Info("Connecting to database", "dsn", config.GetDatabaseDSN())
		db, err := sql.Open("pgx", config.GetDatabaseDSN())
		if err != nil {
			log.Error("Failed to open database connection", "error", err)
			return nil, err
		}
		if err := db.Ping(); err != nil {
			log.Error("Failed to ping database", "error", err)
			db.Close()
			return nil, err
		}
		log.Info("Database connection established successfully")
		container.DB = db
		container.MetricSaveDBRepo = repositories.NewMetricDBSaveRepository(db)
		container.MetricFindDBRepo = repositories.NewMetricDBFindRepository(db)
		container.DBUOW = unitofworks.NewDBUnitOfWork(db)
	}
	if filePath := config.GetFileStoragePath(); filePath != "" {
		dir := filepath.Dir(config.GetFileStoragePath())
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Error("Failed to create directories", "error", err)
			return nil, err
		}
		file, err := os.OpenFile(config.GetFileStoragePath(), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Error("Failed to open file", "error", err)
			return nil, err
		}
		scanner := bufio.NewScanner(file)
		container.File = file
		container.MetricSaveFileRepo = repositories.NewMetricFileSaveRepository(file)
		container.MetricFindFileRepo = repositories.NewMetricFileFindRepository(file, scanner)
		container.FileUOW = unitofworks.NewFileUnitOfWork()
	}
	if container.DB == nil && container.File == nil {
		container.Memory = make(map[domain.MetricID]*domain.Metric)
		container.MetricSaveMemoryRepo = repositories.NewMetricMemorySaveRepository(container.Memory)
		container.MetricFindMemoryRepo = repositories.NewMetricMemoryFindRepository(container.Memory)
		container.MemoryUOW = unitofworks.NewMemoryUnitOfWork()
	}
	if container.MetricSaveDBRepo != nil {
		container.MetricUpdateService = services.NewMetricUpdateService(
			container.MetricSaveDBRepo,
			container.MetricFindDBRepo,
			container.DBUOW,
		)
		container.MetricGetByIDService = services.NewMetricGetByIDService(
			container.MetricFindDBRepo,
		)
		container.MetricListService = services.NewMetricListService(
			container.MetricFindDBRepo,
		)
	} else if container.MetricSaveFileRepo != nil {
		container.MetricUpdateService = services.NewMetricUpdateService(
			container.MetricSaveFileRepo,
			container.MetricFindFileRepo,
			container.FileUOW,
		)
		container.MetricGetByIDService = services.NewMetricGetByIDService(
			container.MetricFindFileRepo,
		)
		container.MetricListService = services.NewMetricListService(
			container.MetricFindFileRepo,
		)
	} else {
		container.MetricUpdateService = services.NewMetricUpdateService(
			container.MetricSaveMemoryRepo,
			container.MetricFindMemoryRepo,
			container.MemoryUOW,
		)
		container.MetricGetByIDService = services.NewMetricGetByIDService(
			container.MetricFindMemoryRepo,
		)
		container.MetricListService = services.NewMetricListService(
			container.MetricFindMemoryRepo,
		)
	}
	container.MetricUpdatePathUsecase = usecases.NewMetricUpdatePathUsecase(container.MetricUpdateService)
	container.MetricUpdateBodyUsecase = usecases.NewMetricUpdateBodyUsecase(container.MetricUpdateService)
	container.MetricGetByIDPathUsecase = usecases.NewMetricGetByIDPathUsecase(container.MetricGetByIDService)
	container.MetricGetByIDBodyUsecase = usecases.NewMetricGetByIDBodyUsecase(container.MetricGetByIDService)
	container.MetricListHTMLUsecase = usecases.NewMetricListHTMLUsecase(container.MetricListService)
	container.MetricUpdatesBodyUsecase = usecases.NewMetricUpdatesBodyUsecase(container.MetricUpdateService)
	return container, nil
}
