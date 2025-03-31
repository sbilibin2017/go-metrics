package app

import (
	"go-metrics/internal/configs"
	"go-metrics/internal/domain"
	"go-metrics/internal/handlers"
	"go-metrics/internal/logger"
	"go-metrics/internal/repositories"
	"go-metrics/internal/routers"
	"go-metrics/internal/server"
	"go-metrics/internal/services"
	"go-metrics/internal/unitofwork"
	"go-metrics/internal/usecases"
)

func NewServer(config *configs.ServerConfig) *server.Server {
	logger.Init()
	defer logger.Logger.Sync()

	data := make(map[domain.MetricID]*domain.Metric)

	saveRepo := repositories.NewMetricMemorySaveRepository(data)
	findRepo := repositories.NewMetricMemoryFindRepository(data)

	uow := unitofwork.NewMemoryUnitOfWork()

	metricUpdateService := services.NewMetricUpdateService(saveRepo, findRepo, uow)
	metricGetByIDService := services.NewMetricGetByIDService(findRepo)
	metricListService := services.NewMetricListService(findRepo)

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

	server := server.NewServer(config, metricRouter)

	return server
}
