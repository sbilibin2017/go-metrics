package app

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"go-metrics/internal/logger"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	mockServer := http.NewServeMux()
	mockServer.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected POST method")
		assert.Contains(t, r.URL.Path, "/update/", "Expected URL path to contain '/update/'")
		if r.Header.Get("Content-Encoding") == "gzip" {
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				t.Fatalf("Error creating gzip reader: %v", err)
			}
			defer gzipReader.Close()
			var metricData map[string]interface{}
			decoder := json.NewDecoder(gzipReader)
			err = decoder.Decode(&metricData)
			assert.NoError(t, err, "Error decoding metric data")
			assert.Contains(t, metricData, "id", "Expected 'id' in metric data")
			assert.Contains(t, metricData, "type", "Expected 'type' in metric data")
			assert.True(t, metricData["type"] == "gauge" || metricData["type"] == "counter", "Invalid metric type")
			if metricData["type"] == "gauge" {
				assert.Contains(t, metricData, "value", "Expected 'value' for gauge metric")
			}
			if metricData["type"] == "counter" {
				assert.Contains(t, metricData, "delta", "Expected 'delta' for counter metric")
			}
		} else {
			t.Fatalf("Expected gzip content encoding, but got %s", r.Header.Get("Content-Encoding"))
		}
		w.WriteHeader(http.StatusOK)
	})
	logger.Init()
	server := &http.Server{Addr: ":8081", Handler: mockServer}
	serverErr := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()
	defer func() { _ = server.Close() }()
	config := &Config{
		PollInterval:   1,
		ReportInterval: 1,
		Address:        "http://localhost:8081",
	}
	ma := NewMetricAgent(config)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	agentErr := make(chan error, 1)
	go func() {
		err := ma.Start(ctx)
		if err != nil {
			agentErr <- err
		}
	}()
	select {
	case err := <-serverErr:
		t.Fatalf("Test server error: %v", err)
	case err := <-agentErr:
		t.Fatalf("Start failed: %v", err)
	case <-time.After(3 * time.Second):
	}
	assert.True(t, true, "Test finished successfully")
}
