package dto

type GetMetricDto struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type PostMetricDto struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
