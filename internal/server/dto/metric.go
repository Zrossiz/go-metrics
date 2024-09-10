package dto

type GetMetricDto struct {
	Name  string `json:"name"`
	MType string `json:"type"`
}

type PostMetricDto struct {
	Name  string  `json:"name"`
	MType string  `json:"type"`
	Value float64 `json:"value"`
}
