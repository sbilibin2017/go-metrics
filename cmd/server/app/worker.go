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
	w.restore(ctx)
	ticker := time.NewTicker(time.Duration(w.config.StoreInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Logger.Infow("Server is shutting down, saving data...")
			w.save(ctx)
			return
		case <-ticker.C:
			logger.Logger.Infow("Periodically saving data...")
			w.save(ctx)
		}
	}
}

func (w *Worker) restore(ctx context.Context) {
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
				err = w.container.SaveMemoryRepo.Save(ctx, metrics)
				if err != nil {
					logger.Logger.Errorw("Error saving restored data to memory", "error", err)
				} else {
					logger.Logger.Infow("Data successfully restored and saved to memory")
				}
			}
		}
	}
}

func (w *Worker) save(ctx context.Context) {
	metricsMap, err := w.container.FindMemoryRepo.Find(ctx, []*domain.MetricID{})
	if err != nil {
		logger.Logger.Errorw("Failed to retrieve metrics from memory", "error", err)
		return
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
