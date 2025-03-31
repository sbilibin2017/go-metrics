package responses

type MetricUpdatePathResponse struct {
	Message string
}

func NewMetricUpdatePathResponse() *MetricUpdatePathResponse {
	return &MetricUpdatePathResponse{Message: "Metric updated successfully"}
}

func (r *MetricUpdatePathResponse) ToResponse() []byte {
	return []byte(r.Message)
}
