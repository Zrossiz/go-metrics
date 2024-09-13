package dto

type GetMetricDto struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type PostMetricDto struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}
