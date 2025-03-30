package app

import (
	"context"
	"fmt"
	"go-metrics/internal/configs"
	"go-metrics/internal/domain"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type MetricAgent struct {
	config    *configs.AgentConfig
	client    *http.Client
	pollCount int64
	metrics   []domain.Metric
}

func NewMetricAgent(config *configs.AgentConfig) *MetricAgent {
	return &MetricAgent{
		config: config,
		client: &http.Client{},
	}
}

func (ma *MetricAgent) Start(ctx context.Context) error {
	log.Println("Starting MetricAgent...")
	tickerPoll := time.NewTicker(time.Duration(ma.config.PollInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(ma.config.ReportInterval) * time.Second)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()

	for {
		select {
		case <-tickerPoll.C:
			ma.collectCounterMetrics()
			ma.collectGaugeMetrics()
		case <-tickerReport.C:
			ma.sendMetrics()
			ma.metrics = nil
		case <-ctx.Done():
			log.Println("MetricAgent received shutdown signal")
			return nil
		}
	}
}

func (ma *MetricAgent) collectGaugeMetrics() {
	log.Println("Collecting gauge metrics...")
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	float64ptr := func(value float64) *float64 { return &value }
	ma.metrics = append(ma.metrics, []domain.Metric{
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

}

func (ma *MetricAgent) collectCounterMetrics() {
	ma.pollCount++
	int64ptr := func(value int64) *int64 { return &value }
	ma.metrics = append(ma.metrics, domain.Metric{
		ID:    "PollCount",
		Type:  domain.Counter,
		Delta: int64ptr(ma.pollCount),
	})
	log.Println("Collecting counter metrics...")
}

func (ma *MetricAgent) sendMetrics() error {
	log.Println("Sending metrics...")
	getMetricValueAsString := func(metric domain.Metric) string {
		if metric.Type == domain.Gauge {
			return fmt.Sprintf("%f", *metric.Value)
		} else if metric.Type == domain.Counter {
			return fmt.Sprintf("%d", *metric.Delta)
		}
		return ""
	}
	normalizeAddress := func(address string) string {
		if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
			return "http://" + address
		}
		return address
	}
	address := normalizeAddress(ma.config.Address)
	for _, metric := range ma.metrics {
		url := fmt.Sprintf("%s/update/%s/%s/%s", address, metric.Type, metric.ID, getMetricValueAsString(metric))
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			log.Println("Error creating request:", err)
			return err
		}
		req.Header.Set("Content-Type", "text/plain")
		resp, err := ma.client.Do(req)
		if err != nil {
			log.Println("Error sending request:", err)
			return err
		}
		defer resp.Body.Close()
		log.Printf("Successfully sent metric: %s/%s/%s", metric.Type, metric.ID, getMetricValueAsString(metric))
	}
	log.Println("All metrics sent successfully")
	return nil
}
