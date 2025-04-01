package app

import (
	"go-metrics/internal/domain"
	"go-metrics/internal/engines"
	"go-metrics/internal/repositories"
	"go-metrics/internal/services"
	"go-metrics/internal/unitofwork"
	"go-metrics/internal/usecases"
)

type Container struct {
	MemoryStorage            *engines.MemoryStorage[domain.MetricID, *domain.Metric]
	MemorySetter             *engines.MemorySetter[domain.MetricID, *domain.Metric]
	MemoryGetter             *engines.MemoryGetter[domain.MetricID, *domain.Metric]
	MemoryRanger             *engines.MemoryRanger[domain.MetricID, *domain.Metric]
	FileEngine               *engines.FileEngine
	FileWriterEngine         *engines.FileWriterEngine[*domain.Metric]
	FileGeneratorEngine      *engines.FileGeneratorEngine[*domain.Metric]
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
	memory := engines.NewMemoryStorage[domain.MetricID, *domain.Metric]()
	memorySetter := engines.NewMemorySetter(memory)
	memoryGetter := engines.NewMemoryGetter(memory)
	memoryRanger := engines.NewMemoryRanger(memory)

	var fileEngine *engines.FileEngine
	var fileWriterEngine *engines.FileWriterEngine[*domain.Metric]
	var fileGeneratorEngine *engines.FileGeneratorEngine[*domain.Metric]
	if config.GetFileStoragePath() != "" {
		fileEngine = engines.NewFileEngine()
		fileWriterEngine = engines.NewFileWriterEngine[*domain.Metric](fileEngine)
		fileGeneratorEngine = engines.NewFileGeneratorEngine[*domain.Metric](fileEngine)
	}

	var saveFileRepo *repositories.MetricFileSaveRepository
	var findFileRepo *repositories.MetricFileFindRepository
	if config.GetFileStoragePath() != "" {
		saveFileRepo = repositories.NewMetricFileSaveRepository(fileWriterEngine)
		findFileRepo = repositories.NewMetricFileFindRepository(fileGeneratorEngine)
	}

	saveMemoryRepo := repositories.NewMetricMemorySaveRepository(memorySetter)
	findMemoryRepo := repositories.NewMetricMemoryFindRepository(memoryGetter, memoryRanger)

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
		MemoryStorage:            memory,
		MemorySetter:             memorySetter,
		MemoryGetter:             memoryGetter,
		MemoryRanger:             memoryRanger,
		FileEngine:               fileEngine,
		FileWriterEngine:         fileWriterEngine,
		FileGeneratorEngine:      fileGeneratorEngine,
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
