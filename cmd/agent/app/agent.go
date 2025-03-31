package app

import (
	"context"
	"go-metrics/internal/configs"
	"go-metrics/internal/domain"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type MetricAgent struct {
	config *configs.AgentConfig
	client *resty.Client
}

func NewMetricAgent(config *configs.AgentConfig) *MetricAgent {
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
		case <-tickerReport.C:
			ma.sendMetrics(ctx, metrics)
			metrics = nil
		case <-ctx.Done():
			log.Println("Metric agent stopping...")
			return nil
		}
	}
}

func (ma *MetricAgent) collectMetrics(metrics []domain.Metric) []domain.Metric {
	log.Println("Collecting metrics...")
	metrics = ma.collectCounterMetrics(metrics)
	metrics = ma.collectGaugeMetrics(metrics)
	log.Printf("Collected %d metrics\n", len(metrics))
	return metrics
}

func (ma *MetricAgent) collectGaugeMetrics(metrics []domain.Metric) []domain.Metric {
	log.Println("Collecting gauge metrics...")
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
	log.Println("Finished collecting gauge metrics.")
	return metrics
}

func (ma *MetricAgent) collectCounterMetrics(metrics []domain.Metric) []domain.Metric {
	log.Println("Collecting counter metrics...")
	int64ptr := func(value int64) *int64 { return &value }
	metrics = append(metrics, domain.Metric{
		ID:    "PollCount",
		Type:  domain.Counter,
		Delta: int64ptr(1),
	})
	log.Println("Finished collecting counter metrics.")
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
		url := address + "/update/"
		log.Printf("Sending metric %s with data: %+v", metric.ID, metric)
		resp, err := ma.client.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(metric).
			Post(url)

		if err != nil {
			log.Printf("Error sending metric %s: %v", metric.ID, err)
			continue
		}
		log.Printf("Successfully sent metric %s, status code: %d", metric.ID, resp.StatusCode())
		if resp.StatusCode() != 200 {
			log.Printf("Response body: %s", resp.String())
		}
	}
	return nil
}
