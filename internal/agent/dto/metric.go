package dto

type MetricDTO struct {
	Name  string      `json:"name"`
	MType string      `json:"type"`
	Value interface{} `json:"value,omitempty"`
}

type GetMetricDto struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}
