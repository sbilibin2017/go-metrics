package repositories

import (
	"bufio"
	"context"
	"encoding/json"
	"go-metrics/internal/domain"
	"log"
	"os"
	"sync"
)

type MetricFileFindRepository struct {
	file    *os.File
	scanner *bufio.Scanner
	mu      sync.Mutex
}

func NewMetricFileFindRepository(file *os.File) *MetricFileFindRepository {
	return &MetricFileFindRepository{
		file:    file,
		scanner: bufio.NewScanner(file),
	}
}

func (repo *MetricFileFindRepository) Find(ctx context.Context, filters []*domain.MetricID) (map[domain.MetricID]*domain.Metric, error) {
	repo.mu.Lock()
	repo.updateScanner()
	defer repo.mu.Unlock()

	log.Printf("Filters: %+v\n", filters)

	// Создаем карту фильтров
	filterMap := make(map[domain.MetricID]struct{})
	for _, filter := range filters {
		filterMap[*filter] = struct{}{}
	}
	log.Printf("Filter map: %+v\n", filterMap)

	result := make(map[domain.MetricID]*domain.Metric)

	_, err := repo.file.Seek(0, 0)
	if err != nil {
		log.Printf("Error seeking file: %v\n", err)
		return nil, err
	}

	for repo.scanner.Scan() {
		line := repo.scanner.Text()
		log.Printf("Reading line: %s\n", line)

		var metric domain.Metric
		if err := json.Unmarshal([]byte(line), &metric); err != nil {
			log.Printf("Error unmarshalling JSON: %v\n", err)
			continue
		}

		metricID := domain.MetricID{ID: metric.ID, Type: metric.Type}

		if len(filters) != 0 {
			log.Printf("Checking metric: %+v\n", metric)
			if _, found := filterMap[metricID]; found {
				log.Printf("Metric found: %+v\n", metric)
				result[metricID] = &metric
			} else {
				log.Printf("Metric not found in filter map: %+v\n", metric)
			}
		} else {
			log.Printf("Metric no filter: %+v\n", metric)
			result[metricID] = &metric
		}
	}

	if err := repo.scanner.Err(); err != nil {
		log.Printf("Error reading file: %v\n", err)
		return nil, err
	}

	log.Printf("Find result: %+v\n", result)
	return result, nil
}

func (repo *MetricFileFindRepository) updateScanner() error {
	_, err := repo.file.Seek(0, 0)
	if err != nil {
		log.Printf("Error seeking file: %v\n", err)
		return err
	}
	repo.scanner = bufio.NewScanner(repo.file)
	return nil
}
