package dto

type GetMetricDto struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type PostMetricDto struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type ResponseMerticDto struct {
	ID    string      `json:"id"`
	MType string      `json:"type"`
	Value interface{} `json:"value,omitempty"`
}
