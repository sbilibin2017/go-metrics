package domain

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type MetricID struct {
	ID   string     `json:"id"`
	Type MetricType `json:"type"`
}

type Metric struct {
	MetricID
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
