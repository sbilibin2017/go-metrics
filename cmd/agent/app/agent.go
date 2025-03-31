package app

import (
	"context"
	"go-metrics/internal/configs"
	"go-metrics/internal/domain"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type MetricAgent struct {
	config *configs.AgentConfig
	client *http.Client
}

func NewMetricAgent(config *configs.AgentConfig) *MetricAgent {
	return &MetricAgent{
		config: config,
		client: &http.Client{},
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
		case <-tickerReport.C:
			ma.sendMetrics(ctx, metrics)
			metrics = nil
		case <-ctx.Done():
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
		{ID: "Alloc", MType: domain.Gauge, Value: float64ptr(float64(memStats.Alloc))},
		{ID: "BuckHashSys", MType: domain.Gauge, Value: float64ptr(float64(memStats.BuckHashSys))},
		{ID: "Frees", MType: domain.Gauge, Value: float64ptr(float64(memStats.Frees))},
		{ID: "GCCPUFraction", MType: domain.Gauge, Value: &memStats.GCCPUFraction},
		{ID: "GCSys", MType: domain.Gauge, Value: float64ptr(float64(memStats.GCSys))},
		{ID: "HeapAlloc", MType: domain.Gauge, Value: float64ptr(float64(memStats.HeapAlloc))},
		{ID: "HeapIdle", MType: domain.Gauge, Value: float64ptr(float64(memStats.HeapIdle))},
		{ID: "HeapInuse", MType: domain.Gauge, Value: float64ptr(float64(memStats.HeapInuse))},
		{ID: "HeapObjects", MType: domain.Gauge, Value: float64ptr(float64(memStats.HeapObjects))},
		{ID: "HeapReleased", MType: domain.Gauge, Value: float64ptr(float64(memStats.HeapReleased))},
		{ID: "HeapSys", MType: domain.Gauge, Value: float64ptr(float64(memStats.HeapSys))},
		{ID: "LastGC", MType: domain.Gauge, Value: float64ptr(float64(memStats.LastGC))},
		{ID: "Lookups", MType: domain.Gauge, Value: float64ptr(float64(memStats.Lookups))},
		{ID: "MCacheInuse", MType: domain.Gauge, Value: float64ptr(float64(memStats.MCacheInuse))},
		{ID: "MCacheSys", MType: domain.Gauge, Value: float64ptr(float64(memStats.MCacheSys))},
		{ID: "MSpanInuse", MType: domain.Gauge, Value: float64ptr(float64(memStats.MSpanInuse))},
		{ID: "MSpanSys", MType: domain.Gauge, Value: float64ptr(float64(memStats.MSpanSys))},
		{ID: "Mallocs", MType: domain.Gauge, Value: float64ptr(float64(memStats.Mallocs))},
		{ID: "NextGC", MType: domain.Gauge, Value: float64ptr(float64(memStats.NextGC))},
		{ID: "NumForcedGC", MType: domain.Gauge, Value: float64ptr(float64(memStats.NumForcedGC))},
		{ID: "NumGC", MType: domain.Gauge, Value: float64ptr(float64(memStats.NumGC))},
		{ID: "OtherSys", MType: domain.Gauge, Value: float64ptr(float64(memStats.OtherSys))},
		{ID: "PauseTotalNs", MType: domain.Gauge, Value: float64ptr(float64(memStats.PauseTotalNs))},
		{ID: "StackInuse", MType: domain.Gauge, Value: float64ptr(float64(memStats.StackInuse))},
		{ID: "StackSys", MType: domain.Gauge, Value: float64ptr(float64(memStats.StackSys))},
		{ID: "Sys", MType: domain.Gauge, Value: float64ptr(float64(memStats.Sys))},
		{ID: "TotalAlloc", MType: domain.Gauge, Value: float64ptr(float64(memStats.TotalAlloc))},
		{ID: "RandomValue", MType: domain.Gauge, Value: float64ptr(rand.Float64())},
	}...)
	return metrics
}

func (ma *MetricAgent) collectCounterMetrics(metrics []domain.Metric) []domain.Metric {
	int64ptr := func(value int64) *int64 { return &value }
	metrics = append(metrics, domain.Metric{
		ID:    "PollCount",
		MType: domain.Counter,
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
	for _, metric := range metrics {
		url := address + "/update/" + metric.MType + "/" + metric.ID
		req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "text/plain")
		resp, err := ma.client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
	}
	return nil
}
