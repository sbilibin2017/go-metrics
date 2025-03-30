package responses

type MetricUpdatePathResponse struct{}

func (r MetricUpdatePathResponse) ToResponse() []byte {
	return []byte("Metric updated successfully")
}
