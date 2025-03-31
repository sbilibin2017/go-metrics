package app

import (
	"context"
	"go-metrics/internal/configs"
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
		w.WriteHeader(http.StatusOK)
	})
	server := &http.Server{Addr: ":8081", Handler: mockServer}
	serverErr := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()
	defer func() { _ = server.Close() }()
	config := &configs.AgentConfig{
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
	assert.True(t, true)
}
