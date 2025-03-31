package responses

import (
	"go-metrics/internal/converters"
	"go-metrics/internal/domain"
)

type MetricGetByIDPathResponse struct {
	Value *string
}

func NewMetricGetByIDPathResponse(metric *domain.Metric) *MetricGetByIDPathResponse {
	var result string
	if metric.MType == domain.Counter && metric.Delta != nil {
		result = converters.FormatInt64(*metric.Delta)
	} else if metric.MType == domain.Gauge && metric.Value != nil {
		result = converters.FormatFloat64(*metric.Value)
	} else {
		return nil
	}
	return &MetricGetByIDPathResponse{Value: &result}
}

func (r *MetricGetByIDPathResponse) ToResponse() []byte {
	return []byte(*r.Value)
}
