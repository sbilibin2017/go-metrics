package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-metrics/internal/domain"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type MetricDBFindRepository struct {
	db *sql.DB
}

func NewMetricDBFindRepository(db *sql.DB) *MetricDBFindRepository {
	return &MetricDBFindRepository{db: db}
}

var baseMetricFindQuery = "SELECT id, type, delta, value FROM metrics"

func buildMetricFindQuery(filters []*domain.MetricID) (string, []any) {
	var sb strings.Builder
	sb.WriteString(baseMetricFindQuery)
	args := make([]any, 0, len(filters)*2)
	if len(filters) > 0 {
		sb.WriteString(" WHERE ")
		for i, filter := range filters {
			if i > 0 {
				sb.WriteString(" OR ")
			}
			sb.WriteString(fmt.Sprintf("(id = $%d AND type = $%d)", i*2+1, i*2+2))
			args = append(args, filter.ID, filter.Type)
		}
	}
	return sb.String(), args
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
	if rows.Err() != nil {
		return nil, err
	}
	return result, nil
}
