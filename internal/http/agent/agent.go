package agent

import (
	"fmt"
	"net/http"

	agentServices "github.com/Zrossiz/go-metrics/internal/services/agent"
)

const serverAddress = "localhost:8080"

func SendMetrics(metrics []agentServices.Metric) {
	for i := 0; i < len(metrics); i++ {
		reqURL := fmt.Sprintf("http://%s/update/%s/%s/%d", serverAddress, metrics[i].Type, metrics[i].Name, metrics[i].Value)
		resp, err := http.Post(reqURL, "text/plain", nil)
		if err != nil {
			fmt.Println("Request:", reqURL, "failed, err:", err)
			continue
		}
		defer resp.Body.Close()
	}
}
