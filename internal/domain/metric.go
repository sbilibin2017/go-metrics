package domain

type Metric struct {
	ID    string
	MType string
	Delta *int64
	Value *float64
}
