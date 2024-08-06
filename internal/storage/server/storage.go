package server

const (
	CounterType = "counter"
	GaugeType   = "gauge"
)

type Metric struct {
	Name  string
	Type  string
	Value interface{}
}
