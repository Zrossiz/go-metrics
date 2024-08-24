package storage

const (
	CounterType = "counter"
	GaugeType   = "gauge"
)

type Metric struct {
	Name  string
	Type  string
	Value interface{}
}

type MetricsStorage interface {
	SetGauge(name string, value float64) bool
	SetCounter(name string, value int64) bool
	GetMetric(name string) Metric
}
