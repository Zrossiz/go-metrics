package handlers

import (
	"fmt"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/agent/lib/collector"
)

const serverAddress = "localhost:8080"

func SendMetrics(metrics []collector.Metric) {
	for i := 0; i < len(metrics); i++ {
		reqURL := fmt.Sprintf("http://%s/update/%s/%s/%f", serverAddress, metrics[i].Type, metrics[i].Name, metrics[i].Value)
		resp, err := http.Post(reqURL, "text/plain", nil)
		if err != nil {
			fmt.Println("Request:", reqURL, "failed, err:", err)
			continue
		}
		defer resp.Body.Close()
	}
}
