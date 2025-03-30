package domain

const (
	Gauge   string = "gauge"
	Counter string = "counter"
)

type MetricID struct {
	ID   string
	Type string
}

type Metric struct {
	ID    string
	Type  string
	Delta *int64
	Value *float64
}
