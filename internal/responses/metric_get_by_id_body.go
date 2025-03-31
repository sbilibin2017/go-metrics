package responses

import "go-metrics/internal/domain"

type MetricGetByIDBodyResponse struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func NewMetricGetByIDBodyResponse(metric *domain.Metric) *MetricGetByIDBodyResponse {
	response := &MetricGetByIDBodyResponse{
		ID:    metric.ID,
		MType: metric.MType,
		Delta: metric.Delta,
		Value: metric.Value,
	}
	return response
}
