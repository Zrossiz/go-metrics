package send

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
)

func TestSendMetrics(t *testing.T) {
    metrics := []types.Metric{
        {Type: "gauge", Name: "metric1", Value: 1.23},
        {Type: "counter", Name: "metric2", Value: 42},
    }

    // Создаем тестовый HTTP-сервер, который всегда возвращает 200 OK
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))
    defer server.Close()

    // Вызываем функцию SendMetrics с тестовыми данными
    sendedMetrics := SendMetrics(metrics, server.Listener.Addr().String())

    // Проверяем, что все метрики были отправлены
    if len(sendedMetrics) != len(metrics) {
        t.Errorf("Expected %d metrics to be sent, but got %d", len(metrics), len(sendedMetrics))
    }

    for i, metric := range sendedMetrics {
        if metric != metrics[i] {
            t.Errorf("Expected metric %v, but got %v", metrics[i], metric)
        }
    }
}
