package domain

type Metric struct {
	ID    string
	Type  string
	Delta *int64
	Value *float64
}
