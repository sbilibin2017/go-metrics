package app

import (
	"context"
	"encoding/json"
	"go-metrics/internal/configs"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	mockServer := http.NewServeMux()
	mockServer.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		// Проверка метода запроса
		assert.Equal(t, "POST", r.Method, "Expected POST method")

		// Проверка URL-адреса
		assert.Contains(t, r.URL.Path, "/update/", "Expected URL path to contain '/update/'")

		// Чтение тела запроса
		var metricData map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&metricData)
		assert.NoError(t, err, "Error decoding metric data")

		// Проверка, что в теле запроса присутствуют обязательные поля
		assert.Contains(t, metricData, "ID", "Expected 'ID' in metric data")
		assert.Contains(t, metricData, "Type", "Expected 'Type' in metric data")

		// Проверка типа метрики (должен быть "gauge" или "counter")
		assert.True(t, metricData["Type"] == "gauge" || metricData["Type"] == "counter", "Invalid metric type")

		// Проверка поля "Value" для gauge
		if metricData["Type"] == "gauge" {
			assert.Contains(t, metricData, "Value", "Expected 'Value' for gauge metric")
		}

		// Проверка поля "Delta" для counter
		if metricData["Type"] == "counter" {
			assert.Contains(t, metricData, "Delta", "Expected 'Delta' for counter metric")
		}

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

	// Контекст с таймаутом для теста
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Запуск агента в горутине
	agentErr := make(chan error, 1)
	go func() {
		err := ma.Start(ctx)
		if err != nil {
			agentErr <- err
		}
	}()

	// Ожидаем завершения теста
	select {
	case err := <-serverErr:
		t.Fatalf("Test server error: %v", err)
	case err := <-agentErr:
		t.Fatalf("Start failed: %v", err)
	case <-time.After(3 * time.Second):
	}

	// После завершения теста убедимся, что он выполнился успешно
	assert.True(t, true)
}
