package app

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/pkg/log"
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
	log.Info("Server is starting, attempting to restore data...")
	w.restore(ctx)
	ticker := time.NewTicker(time.Duration(w.config.StoreInterval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		log.Info("Periodically saving data...")
		w.save(ctx)
	}
}

func (w *Worker) Stop(ctx context.Context) {
	log.Info("Server is stopping, attempting to save data...")
	w.save(ctx)
}

func (w *Worker) restore(ctx context.Context) {
	if w.config.Restore {
		var metrics []*domain.Metric
		var err error
		metricsMap, err := w.container.MetricFindFileRepo.Find(ctx, []*domain.MetricID{})
		if err != nil {
			log.Error("Error restoring data from file", "error", err)
		} else {
			for _, metric := range metricsMap {
				metrics = append(metrics, metric)
			}
			if w.config.GetDatabaseDSN() != "" {
				err = w.container.MetricSaveDBRepo.Save(ctx, metrics)
				if err != nil {
					log.Error("Error saving restored data to db", "error", err)
				} else {
					log.Info("Data successfully restored and saved to db")
				}
			} else {
				err = w.container.MetricSaveFileRepo.Save(ctx, metrics)
				if err != nil {
					log.Error("Error saving restored data to file", "error", err)
				} else {
					log.Info("Data successfully restored and saved to file")
				}
			}
		}
	}
}

func (w *Worker) save(ctx context.Context) {
	var metricsMap map[domain.MetricID]*domain.Metric
	var err error
	if w.config.GetDatabaseDSN() != "" {
		metricsMap, err = w.container.MetricFindDBRepo.Find(ctx, []*domain.MetricID{})
		if err != nil {
			log.Error("Error saving restored data to db", "error", err)
		} else {
			log.Info("Data successfully restored and saved to db")
		}
	} else {
		metricsMap, err = w.container.MetricFindFileRepo.Find(ctx, []*domain.MetricID{})
		if err != nil {
			log.Error("Error saving restored data to file", "error", err)
		} else {
			log.Info("Data successfully restored and saved to file")
		}
	}
	var metrics []*domain.Metric
	for _, metric := range metricsMap {
		metrics = append(metrics, metric)
	}
	err = w.container.MetricSaveFileRepo.Save(ctx, metrics)
	if err != nil {
		log.Error("Failed to save metrics to file", "error", err)
	} else {
		log.Info("Metrics successfully saved to file")
	}
}
