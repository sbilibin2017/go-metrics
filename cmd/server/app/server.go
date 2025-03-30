package app

import (
	"go-metrics/internal/configs"
	"go-metrics/internal/domain"
	"go-metrics/internal/handlers"
	"go-metrics/internal/repositories"
	"go-metrics/internal/routers"
	"go-metrics/internal/server"
	"go-metrics/internal/services"
	"go-metrics/internal/unitofwork"
	"go-metrics/internal/usecases"

	"github.com/go-chi/chi"
)

func NewServer(config *configs.ServerConfig) *server.Server {
	data := make(map[domain.MetricID]*domain.Metric)
	saveRepo := repositories.NewMetricMemorySaveRepository(data)
	findRepo := repositories.NewMetricMemoryFindRepository(data)
	uow := unitofwork.NewMemoryUnitOfWork()
	metricUpdateService := services.NewMetricUpdateService(saveRepo, findRepo, uow)
	metricUpdatePathUsecase := usecases.NewMetricUpdatePathUsecase(metricUpdateService)
	metricUpdateHandler := handlers.MetricUpdatePathHandler(metricUpdatePathUsecase)
	r := chi.NewRouter()
	routers.RegisterMetricUpdatePathRouter(r, metricUpdateHandler)
	server := server.NewServer(config)
	server.AddRouter(r)
	return server
}
