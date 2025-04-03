package app

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"go-metrics/internal/domain"
	"go-metrics/pkg/log"
	"math/rand"
	"net/http"
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
			log.Info("Metrics collected", "metrics_count", len(metrics))
		case <-tickerReport.C:
			err := ma.sendMetrics(ctx, metrics)
			if err != nil {
				log.Error("Failed to send metrics", "error", err)
			} else {
				log.Info("Metrics sent successfully", "metrics_count", len(metrics))
			}
			metrics = nil
		case <-ctx.Done():
			log.Info("Shutting down metric agent")
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
		{MetricID: domain.MetricID{ID: "Alloc", Type: domain.Gauge}, Value: float64ptr(float64(memStats.Alloc))},
		{MetricID: domain.MetricID{ID: "BuckHashSys", Type: domain.Gauge}, Value: float64ptr(float64(memStats.BuckHashSys))},
		{MetricID: domain.MetricID{ID: "Frees", Type: domain.Gauge}, Value: float64ptr(float64(memStats.Frees))},
		{MetricID: domain.MetricID{ID: "GCCPUFraction", Type: domain.Gauge}, Value: &memStats.GCCPUFraction},
		{MetricID: domain.MetricID{ID: "GCSys", Type: domain.Gauge}, Value: float64ptr(float64(memStats.GCSys))},
		{MetricID: domain.MetricID{ID: "HeapAlloc", Type: domain.Gauge}, Value: float64ptr(float64(memStats.HeapAlloc))},
		{MetricID: domain.MetricID{ID: "HeapIdle", Type: domain.Gauge}, Value: float64ptr(float64(memStats.HeapIdle))},
		{MetricID: domain.MetricID{ID: "HeapInuse", Type: domain.Gauge}, Value: float64ptr(float64(memStats.HeapInuse))},
		{MetricID: domain.MetricID{ID: "HeapObjects", Type: domain.Gauge}, Value: float64ptr(float64(memStats.HeapObjects))},
		{MetricID: domain.MetricID{ID: "HeapReleased", Type: domain.Gauge}, Value: float64ptr(float64(memStats.HeapReleased))},
		{MetricID: domain.MetricID{ID: "HeapSys", Type: domain.Gauge}, Value: float64ptr(float64(memStats.HeapSys))},
		{MetricID: domain.MetricID{ID: "LastGC", Type: domain.Gauge}, Value: float64ptr(float64(memStats.LastGC))},
		{MetricID: domain.MetricID{ID: "Lookups", Type: domain.Gauge}, Value: float64ptr(float64(memStats.Lookups))},
		{MetricID: domain.MetricID{ID: "MCacheInuse", Type: domain.Gauge}, Value: float64ptr(float64(memStats.MCacheInuse))},
		{MetricID: domain.MetricID{ID: "MCacheSys", Type: domain.Gauge}, Value: float64ptr(float64(memStats.MCacheSys))},
		{MetricID: domain.MetricID{ID: "MSpanInuse", Type: domain.Gauge}, Value: float64ptr(float64(memStats.MSpanInuse))},
		{MetricID: domain.MetricID{ID: "MSpanSys", Type: domain.Gauge}, Value: float64ptr(float64(memStats.MSpanSys))},
		{MetricID: domain.MetricID{ID: "Mallocs", Type: domain.Gauge}, Value: float64ptr(float64(memStats.Mallocs))},
		{MetricID: domain.MetricID{ID: "NextGC", Type: domain.Gauge}, Value: float64ptr(float64(memStats.NextGC))},
		{MetricID: domain.MetricID{ID: "NumForcedGC", Type: domain.Gauge}, Value: float64ptr(float64(memStats.NumForcedGC))},
		{MetricID: domain.MetricID{ID: "NumGC", Type: domain.Gauge}, Value: float64ptr(float64(memStats.NumGC))},
		{MetricID: domain.MetricID{ID: "OtherSys", Type: domain.Gauge}, Value: float64ptr(float64(memStats.OtherSys))},
		{MetricID: domain.MetricID{ID: "PauseTotalNs", Type: domain.Gauge}, Value: float64ptr(float64(memStats.PauseTotalNs))},
		{MetricID: domain.MetricID{ID: "StackInuse", Type: domain.Gauge}, Value: float64ptr(float64(memStats.StackInuse))},
		{MetricID: domain.MetricID{ID: "StackSys", Type: domain.Gauge}, Value: float64ptr(float64(memStats.StackSys))},
		{MetricID: domain.MetricID{ID: "Sys", Type: domain.Gauge}, Value: float64ptr(float64(memStats.Sys))},
		{MetricID: domain.MetricID{ID: "TotalAlloc", Type: domain.Gauge}, Value: float64ptr(float64(memStats.TotalAlloc))},
		{MetricID: domain.MetricID{ID: "RandomValue", Type: domain.Gauge}, Value: float64ptr(rand.Float64())},
	}...)
	return metrics
}

func (ma *MetricAgent) collectCounterMetrics(metrics []domain.Metric) []domain.Metric {
	int64ptr := func(value int64) *int64 { return &value }
	metrics = append(metrics, domain.Metric{
		MetricID: domain.MetricID{ID: "PollCount", Type: domain.Counter},
		Delta:    int64ptr(1),
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
	url := address + "/updates/"
	body, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}
	compressedBody, err := compress(body)
	if err != nil {
		return fmt.Errorf("failed to compress metrics: %w", err)
	}
	resp, err := ma.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(compressedBody).
		Post(url)

	if err != nil {
		return fmt.Errorf("failed to send metrics: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to send metrics, status code: %d", resp.StatusCode())
	}
	log.Info("Metrics sent successfully", "metrics_count", len(metrics))
	return nil
}

func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err := gzipWriter.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write to gzip: %w", err)
	}
	err = gzipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}
	return buf.Bytes(), nil
}
