package send

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMetrics(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	// metrics := []types.Metric{
	// 	{Name: "metricCounter", Type: constants.Counter, Value: 100},
	// 	{Name: "metricGauge", Type: constants.Gauge, Value: 0.2},
	// }

	// expectedMetrics := metrics

	// sendedMetrics := SendMetrics(metrics, "localhost:8080")

	// if len(sendedMetrics) != len(expectedMetrics) {
	// 	t.Errorf("Expectd %d metrics to be sent, but got %d", len(expectedMetrics), len(sendedMetrics))
	// }

	// for i, metric := range sendedMetrics {
	// 	if metric != expectedMetrics[i] {
	// 		t.Errorf("Expected metric %v, but got %v", expectedMetrics[i], metric)
	// 	}
	// }
}
