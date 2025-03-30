package app

import (
	"context"
	"go-metrics/internal/configs"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestMetricAgent_Start(t *testing.T) {
	config := &configs.AgentConfig{
		Address:        "http://localhost:8080",
		PollInterval:   1,
		ReportInterval: 2,
	}
	ma := NewMetricAgent(config)
	go func() {
		http.HandleFunc("/update_8080", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Expected POST request, but got %s", r.Method)
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	time.Sleep(500 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		err := ma.Start(ctx)
		if err != nil && err.Error() != "MetricAgent received shutdown signal" {
			t.Errorf("MetricAgent failed with error: %v", err)
		}
	}()
	time.Sleep(3 * time.Second)
}

func TestMetricAgent_Start_ErrorSendingMetrics(t *testing.T) {
	// Создаем конфигурацию для агента
	config := &configs.AgentConfig{
		Address:        "http://localhost:8081",
		PollInterval:   1,
		ReportInterval: 2,
	}
	ma := NewMetricAgent(config)
	go func() {
		http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Expected POST request, but got %s", r.Method)
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
		})
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()
	time.Sleep(500 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		err := ma.Start(ctx)
		if err != nil && err.Error() != "MetricAgent received shutdown signal" {
			t.Errorf("MetricAgent failed with error: %v", err)
		}
	}()
	time.Sleep(3 * time.Second)
	if len(ma.metrics) > 0 {
		t.Errorf("Expected no successful metrics to be sent, but some metrics were sent")
	}
}
