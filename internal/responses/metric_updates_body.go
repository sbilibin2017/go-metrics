package responses

import "go-metrics/internal/domain"

func NewMetricUpdatesBodyResponse(metrics []*domain.Metric) []*MetricUpdateBodyResponse {
	var responses []*MetricUpdateBodyResponse
	for _, metric := range metrics {
		responses = append(responses, &MetricUpdateBodyResponse{
			ID:    metric.ID,
			Type:  metric.Type,
			Delta: metric.Delta,
			Value: metric.Value,
		})
	}
	return responses
}
