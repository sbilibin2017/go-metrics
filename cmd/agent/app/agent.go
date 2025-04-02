package app

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"go-metrics/internal/domain"
	"go-metrics/internal/logger"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type MetricAgent struct {
	config *Config
	client *resty.Client
}

func NewMetricAgent(config *Config) *MetricAgent {
	return &MetricAgent{
		config: config,
		client: resty.New(),
	}
}

func (ma *MetricAgent) Start(ctx context.Context) error {
	tickerPoll := time.NewTicker(time.Duration(ma.config.PollInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(ma.config.ReportInterval) * time.Second)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()

	var metrics []domain.Metric
	for {
		select {
		case <-tickerPoll.C:
			metrics = ma.collectMetrics(metrics)
			logger.Logger.Infow("Metrics collected", "metrics_count", len(metrics))
		case <-tickerReport.C:
			if len(metrics) > 0 {
				// Send the metrics in batches
				err := ma.sendMetrics(ctx, metrics)
				if err != nil {
					logger.Logger.Errorw("Failed to send metrics", "error", err)
				} else {
					logger.Logger.Infow("Metrics sent successfully", "metrics_count", len(metrics))
				}
				// Clear the metrics after sending
				metrics = nil
			}
		case <-ctx.Done():
			logger.Logger.Info("Shutting down metric agent")
			return nil
		}
	}
}

func (ma *MetricAgent) collectMetrics(metrics []domain.Metric) []domain.Metric {
	metrics = ma.collectCounterMetrics(metrics)
	metrics = ma.collectGaugeMetrics(metrics)
	return metrics
}

func (ma *MetricAgent) collectGaugeMetrics(metrics []domain.Metric) []domain.Metric {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	float64ptr := func(value float64) *float64 { return &value }
	metrics = append(metrics, []domain.Metric{
		{ID: "Alloc", Type: domain.Gauge, Value: float64ptr(float64(memStats.Alloc))},
		{ID: "BuckHashSys", Type: domain.Gauge, Value: float64ptr(float64(memStats.BuckHashSys))},
		{ID: "Frees", Type: domain.Gauge, Value: float64ptr(float64(memStats.Frees))},
		{ID: "GCCPUFraction", Type: domain.Gauge, Value: &memStats.GCCPUFraction},
		{ID: "GCSys", Type: domain.Gauge, Value: float64ptr(float64(memStats.GCSys))},
		{ID: "HeapAlloc", Type: domain.Gauge, Value: float64ptr(float64(memStats.HeapAlloc))},
		{ID: "HeapIdle", Type: domain.Gauge, Value: float64ptr(float64(memStats.HeapIdle))},
		{ID: "HeapInuse", Type: domain.Gauge, Value: float64ptr(float64(memStats.HeapInuse))},
		{ID: "HeapObjects", Type: domain.Gauge, Value: float64ptr(float64(memStats.HeapObjects))},
		{ID: "HeapReleased", Type: domain.Gauge, Value: float64ptr(float64(memStats.HeapReleased))},
		{ID: "HeapSys", Type: domain.Gauge, Value: float64ptr(float64(memStats.HeapSys))},
		{ID: "LastGC", Type: domain.Gauge, Value: float64ptr(float64(memStats.LastGC))},
		{ID: "Lookups", Type: domain.Gauge, Value: float64ptr(float64(memStats.Lookups))},
		{ID: "MCacheInuse", Type: domain.Gauge, Value: float64ptr(float64(memStats.MCacheInuse))},
		{ID: "MCacheSys", Type: domain.Gauge, Value: float64ptr(float64(memStats.MCacheSys))},
		{ID: "MSpanInuse", Type: domain.Gauge, Value: float64ptr(float64(memStats.MSpanInuse))},
		{ID: "MSpanSys", Type: domain.Gauge, Value: float64ptr(float64(memStats.MSpanSys))},
		{ID: "Mallocs", Type: domain.Gauge, Value: float64ptr(float64(memStats.Mallocs))},
		{ID: "NextGC", Type: domain.Gauge, Value: float64ptr(float64(memStats.NextGC))},
		{ID: "NumForcedGC", Type: domain.Gauge, Value: float64ptr(float64(memStats.NumForcedGC))},
		{ID: "NumGC", Type: domain.Gauge, Value: float64ptr(float64(memStats.NumGC))},
		{ID: "OtherSys", Type: domain.Gauge, Value: float64ptr(float64(memStats.OtherSys))},
		{ID: "PauseTotalNs", Type: domain.Gauge, Value: float64ptr(float64(memStats.PauseTotalNs))},
		{ID: "StackInuse", Type: domain.Gauge, Value: float64ptr(float64(memStats.StackInuse))},
		{ID: "StackSys", Type: domain.Gauge, Value: float64ptr(float64(memStats.StackSys))},
		{ID: "Sys", Type: domain.Gauge, Value: float64ptr(float64(memStats.Sys))},
		{ID: "TotalAlloc", Type: domain.Gauge, Value: float64ptr(float64(memStats.TotalAlloc))},
		{ID: "RandomValue", Type: domain.Gauge, Value: float64ptr(rand.Float64())},
	}...)
	return metrics
}

func (ma *MetricAgent) collectCounterMetrics(metrics []domain.Metric) []domain.Metric {
	int64ptr := func(value int64) *int64 { return &value }
	metrics = append(metrics, domain.Metric{
		ID:    "PollCount",
		Type:  domain.Counter,
		Delta: int64ptr(1),
	})
	return metrics
}

func (ma *MetricAgent) sendMetrics(ctx context.Context, metrics []domain.Metric) error {
	normalizeAddress := func(address string) string {
		if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
			return "http://" + address
		}
		return address
	}

	address := normalizeAddress(ma.config.Address)

	// Create batches of metrics (for example, batches of 100 metrics each)
	batchSize := 100
	for i := 0; i < len(metrics); i += batchSize {
		end := i + batchSize
		if end > len(metrics) {
			end = len(metrics)
		}

		batch := metrics[i:end]
		url := address + "/updates/"

		// Marshal the batch of metrics
		var buf bytes.Buffer
		gzipWriter := gzip.NewWriter(&buf)
		body, err := json.Marshal(batch)
		if err != nil {
			return err
		}

		_, err = gzipWriter.Write(body)
		if err != nil {
			return err
		}

		err = gzipWriter.Close()
		if err != nil {
			return err
		}

		// Send the batch
		resp, err := ma.client.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(buf.Bytes()).
			Post(url)

		if err != nil {
			return err
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}
	}

	return nil
}
