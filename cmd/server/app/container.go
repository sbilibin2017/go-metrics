package app

import (
	"database/sql"
	"go-metrics/internal/domain"
	"go-metrics/internal/repositories"
	"go-metrics/internal/services"
	"go-metrics/internal/storages"
	"go-metrics/internal/unitofwork"
	"go-metrics/internal/usecases"
	"os"

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
	DBUOW                    *unitofwork.DBUnitOfWork
	FileUOW                  *unitofwork.FileUnitOfWork
	MemoryUOW                *unitofwork.MemoryUnitOfWork
	MetricUpdateService      *services.MetricUpdateService
	MetricGetByIDService     *services.MetricGetByIDService
	MetricListService        *services.MetricListService
	MetricUpdatePathUsecase  *usecases.MetricUpdatePathUsecase
	MetricGetByIDPathUsecase *usecases.MetricGetByIDPathUsecase
	MetricListHTMLUsecase    *usecases.MetricListHTMLUsecase
	MetricUpdateBodyUsecase  *usecases.MetricUpdateBodyUsecase
	MetricGetByIDBodyUsecase *usecases.MetricGetByIDBodyUsecase
}

func NewContainer(config *Config) (*Container, error) {
	container := &Container{
		Memory: make(map[domain.MetricID]*domain.Metric),
	}

	if dsn := config.GetDatabaseDSN(); dsn != "" {
		db, err := storages.NewDB(config)
		if err != nil {
			return nil, err
		}
		container.DB = db
		container.MetricSaveDBRepo = repositories.NewMetricDBSaveRepository(db)
		container.MetricFindDBRepo = repositories.NewMetricDBFindRepository(db)
		container.DBUOW = unitofwork.NewDBUnitOfWork(db)
	}

	if filePath := config.GetFileStoragePath(); filePath != "" {
		file, err := storages.NewFile(config)
		if err != nil {
			return nil, err
		}
		container.File = file
		container.MetricSaveFileRepo = repositories.NewMetricFileSaveRepository(file)
		container.MetricFindFileRepo = repositories.NewMetricFileFindRepository(file)
		container.FileUOW = unitofwork.NewFileUnitOfWork()
	}

	if container.DB == nil && container.File == nil {
		container.Memory = storages.NewMemory[domain.MetricID, *domain.Metric]()
		container.MetricSaveMemoryRepo = repositories.NewMetricMemorySaveRepository(container.Memory)
		container.MetricFindMemoryRepo = repositories.NewMetricMemoryFindRepository(container.Memory)
		container.MemoryUOW = unitofwork.NewMemoryUnitOfWork()
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

	return container, nil
}
