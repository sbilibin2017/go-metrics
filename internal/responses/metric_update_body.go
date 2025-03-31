package responses

import "go-metrics/internal/domain"

type MetricUpdateBodyResponse struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func NewMetricUpdateBodyResponse(metrics []*domain.Metric) *MetricUpdateBodyResponse {
	if len(metrics) == 0 {
		return nil
	}
	metric := metrics[0]
	response := &MetricUpdateBodyResponse{
		ID:    metric.ID,
		MType: metric.MType,
		Delta: metric.Delta,
		Value: metric.Value,
	}
	return response
}
