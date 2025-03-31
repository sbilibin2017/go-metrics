package responses

import (
	"fmt"
	"go-metrics/internal/domain"
)

type MetricListHTMLResponse struct {
	metrics []*domain.Metric
}

func NewMetricListHTMLResponse(metrics []*domain.Metric) *MetricListHTMLResponse {
	return &MetricListHTMLResponse{metrics: metrics}
}

func (r *MetricListHTMLResponse) ToResponse() []byte {
	html := "<html><head><title>Metrics</title></head><body>"
	html += "<h1>Metrics List</h1><table border='1'><tr><th>ID</th><th>Value</th></tr>"
	for _, metric := range r.metrics {
		var value string
		if metric.Type == domain.Counter {
			if metric.Delta != nil {
				value = fmt.Sprintf("%d", *metric.Delta)
			} else {
				value = "N/A"
			}
		} else if metric.Type == domain.Gauge {
			if metric.Value != nil {
				value = fmt.Sprintf("%f", *metric.Value)
			} else {
				value = "N/A"
			}
		}
		html += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", metric.ID, value)
	}
	html += "</table></body></html>"
	return []byte(html)
}
