package app

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/logger"
	"time"
)

type Worker struct {
	config    *Config
	container *Container
}

func NewWorker(config *Config, container *Container) *Worker {
	return &Worker{
		config:    config,
		container: container,
	}
}

func (w *Worker) Start(ctx context.Context) {
	logger.Logger.Infow("Server is starting, attempting to restore data...")
	w.Restore(ctx)
	ticker := time.NewTicker(time.Duration(w.config.StoreInterval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		logger.Logger.Infow("Periodically saving data...")
		w.Save(ctx)
	}
}

func (w *Worker) Restore(ctx context.Context) {
	if w.config.Restore {
		var metrics []*domain.Metric
		var err error
		metricsMap, err := w.container.FindFileRepo.Find(ctx, []*domain.MetricID{})
		if err != nil {
			logger.Logger.Errorw("Error restoring data from file", "error", err)
		} else {
			for _, metric := range metricsMap {
				metrics = append(metrics, metric)
			}
			if len(metrics) > 0 {
				if w.config.GetDatabaseDSN() != "" {
					err = w.container.SaveDBRepo.Save(ctx, metrics)
					if err != nil {
						logger.Logger.Errorw("Error saving restored data to memory", "error", err)
					} else {
						logger.Logger.Infow("Data successfully restored and saved to memory")
					}
				} else {
					err = w.container.SaveFileRepo.Save(ctx, metrics)
					if err != nil {
						logger.Logger.Errorw("Error saving restored data to memory", "error", err)
					} else {
						logger.Logger.Infow("Data successfully restored and saved to memory")
					}
				}
			}
		}
	}
}

func (w *Worker) Save(ctx context.Context) {
	var metricsMap map[domain.MetricID]*domain.Metric
	var err error
	if w.config.GetDatabaseDSN() != "" {
		metricsMap, err = w.container.FindDBRepo.Find(ctx, []*domain.MetricID{})
		if err != nil {
			logger.Logger.Errorw("Failed to retrieve metrics from memory", "error", err)
			return
		}
	} else {
		metricsMap, err = w.container.FindFileRepo.Find(ctx, []*domain.MetricID{})
		if err != nil {
			logger.Logger.Errorw("Failed to retrieve metrics from memory", "error", err)
			return
		}

	}
	var metrics []*domain.Metric
	for _, metric := range metricsMap {
		metrics = append(metrics, metric)
	}
	if len(metrics) > 0 {
		err = w.container.SaveFileRepo.Save(ctx, metrics)
		if err != nil {
			logger.Logger.Errorw("Failed to save metrics to file", "error", err)
		} else {
			logger.Logger.Infow("Metrics successfully saved to file")
		}
	} else {
		logger.Logger.Infow("No metrics to save")
	}
}
