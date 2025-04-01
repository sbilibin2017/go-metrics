package repositories

import (
	"context"
	"fmt"
	"go-metrics/internal/domain"
	"go-metrics/internal/engines"
)

type MetricDBFindRepository struct {
	db *engines.DBEngine
}

func NewMetricDBFindRepository(db *engines.DBEngine) *MetricDBFindRepository {
	return &MetricDBFindRepository{db: db}
}

var baseMetricFindQuery = "SELECT id, type, delta, value FROM metrics"

func buildMetricFindQuery(filters []*domain.MetricID) (string, []any) {
	query := baseMetricFindQuery
	args := []interface{}{}
	if len(filters) > 0 {
		query += " WHERE "
		for i, filter := range filters {
			if i > 0 {
				query += " OR "
			}
			query += fmt.Sprintf("(id = $%d AND type = $%d)", i*2+1, i*2+2)
			args = append(args, filter.ID, filter.Type)
		}
	}
	return query, args
}

func (repo *MetricDBFindRepository) Find(ctx context.Context, filters []*domain.MetricID) (map[domain.MetricID]*domain.Metric, error) {
	result := make(map[domain.MetricID]*domain.Metric)
	query, args := buildMetricFindQuery(filters)
	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var metric domain.Metric
		if err := rows.Scan(&metric.ID, &metric.Type, &metric.Delta, &metric.Value); err != nil {
			return nil, err
		}
		result[domain.MetricID{ID: metric.ID, Type: metric.Type}] = &metric
	}
	return result, nil
}
