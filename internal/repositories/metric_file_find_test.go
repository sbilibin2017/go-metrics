package repositories

import (
	"context"
	"encoding/json"
	"go-metrics/internal/domain"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind_Success(t *testing.T) {
	// Создаем временный файл для теста
	tmpFile, err := os.CreateTemp("", "metrics_*.json")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Удалим файл после теста

	// Записываем тестовые данные в файл
	metrics := []domain.Metric{
		{ID: "metric1", Type: "counter", Delta: new(int64)},
		{ID: "metric2", Type: "gauge", Value: new(float64)},
		{ID: "metric3", Type: "counter", Delta: new(int64)},
	}

	// Для каждой метрики сериализуем ее в JSON и записываем в файл
	for _, metric := range metrics {
		data, err := json.Marshal(metric)
		if err != nil {
			t.Fatalf("Error marshalling metric: %v", err)
		}
		_, err = tmpFile.Write(append(data, '\n'))
		if err != nil {
			t.Fatalf("Error writing to temporary file: %v", err)
		}
	}

	// Перемещаем указатель файла в начало для чтения
	tmpFile.Seek(0, 0)

	// Создаем репозиторий с файлом
	repo := NewMetricFileFindRepository(tmpFile)

	// Создаем фильтры
	filters := []*domain.MetricID{
		{ID: "metric1", Type: "counter"},
		{ID: "metric3", Type: "counter"},
	}

	// Выполняем поиск
	result, err := repo.Find(context.Background(), filters)

	// Проверяем ошибки
	assert.NoError(t, err, "Error while finding metrics")

	// Проверяем, что в результате есть нужные метрики
	assert.Len(t, result, 2, "Expected 2 metrics to be found")

	// Проверяем, что найденные метрики совпадают с ожидаемыми
	assert.Contains(t, result, domain.MetricID{ID: "metric1", Type: "counter"}, "Metric1 should be in the result")
	assert.Contains(t, result, domain.MetricID{ID: "metric3", Type: "counter"}, "Metric3 should be in the result")

	// Дополнительные проверки на метрики
	assert.Equal(t, result[domain.MetricID{ID: "metric1", Type: "counter"}].ID, "metric1")
	assert.Equal(t, result[domain.MetricID{ID: "metric3", Type: "counter"}].ID, "metric3")
}

func TestFind_NoResults(t *testing.T) {
	// Создаем временный файл для теста
	tmpFile, err := os.CreateTemp("", "metrics_*.json")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Удалим файл после теста

	// Записываем тестовые данные в файл
	metrics := []domain.Metric{
		{ID: "metric1", Type: "counter", Delta: new(int64)},
		{ID: "metric2", Type: "gauge", Value: new(float64)},
		{ID: "metric3", Type: "counter", Delta: new(int64)},
	}

	// Для каждой метрики сериализуем ее в JSON и записываем в файл
	for _, metric := range metrics {
		data, err := json.Marshal(metric)
		if err != nil {
			t.Fatalf("Error marshalling metric: %v", err)
		}
		_, err = tmpFile.Write(append(data, '\n'))
		if err != nil {
			t.Fatalf("Error writing to temporary file: %v", err)
		}
	}

	// Перемещаем указатель файла в начало для чтения
	tmpFile.Seek(0, 0)

	// Создаем репозиторий с файлом
	repo := NewMetricFileFindRepository(tmpFile)

	// Создаем фильтры, которые не совпадают с данными в файле
	filters := []*domain.MetricID{
		{ID: "metric4", Type: "counter"},
		{ID: "metric5", Type: "counter"},
	}

	// Выполняем поиск
	result, err := repo.Find(context.Background(), filters)

	// Проверяем ошибки
	assert.NoError(t, err, "Error while finding metrics")

	// Проверяем, что результат пустой
	assert.Len(t, result, 0, "Expected no metrics to be found")
}

func TestFind_FileError(t *testing.T) {
	// Попробуем открыть несуществующий файл
	_, err := os.Open("non_existing_file.json")
	if err == nil {
		t.Fatal("Expected an error when opening a non-existing file")
	}
}

func TestFind_NoFilters(t *testing.T) {
	// Создаем временный файл для теста
	tmpFile, err := os.CreateTemp("", "metrics_*.json")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Удалим файл после теста

	// Записываем тестовые данные в файл
	metrics := []domain.Metric{
		{ID: "metric1", Type: "counter", Delta: new(int64)},
		{ID: "metric2", Type: "gauge", Value: new(float64)},
		{ID: "metric3", Type: "counter", Delta: new(int64)},
	}

	// Для каждой метрики сериализуем ее в JSON и записываем в файл
	for _, metric := range metrics {
		data, err := json.Marshal(metric)
		if err != nil {
			t.Fatalf("Error marshalling metric: %v", err)
		}
		_, err = tmpFile.Write(append(data, '\n'))
		if err != nil {
			t.Fatalf("Error writing to temporary file: %v", err)
		}
	}

	// Перемещаем указатель файла в начало для чтения
	tmpFile.Seek(0, 0)

	// Создаем репозиторий с файлом
	repo := NewMetricFileFindRepository(tmpFile)

	// Создаем пустой массив фильтров (это имитирует отсутствие фильтров)
	filters := []*domain.MetricID{}

	// Выполняем поиск
	result, err := repo.Find(context.Background(), filters)

	// Проверяем ошибки
	assert.NoError(t, err, "Error while finding metrics")

	// Проверяем, что все метрики из файла вернулись
	assert.Len(t, result, 3, "Expected all metrics to be found")

	// Проверяем, что все записи присутствуют в результате
	assert.Contains(t, result, domain.MetricID{ID: "metric1", Type: "counter"}, "Metric1 should be in the result")
	assert.Contains(t, result, domain.MetricID{ID: "metric2", Type: "gauge"}, "Metric2 should be in the result")
	assert.Contains(t, result, domain.MetricID{ID: "metric3", Type: "counter"}, "Metric3 should be in the result")
}
