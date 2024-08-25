package dto

type MetricDTO struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type GetMetricDto struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}
