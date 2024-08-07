package http

import (
	"fmt"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/agent/constants"
	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
)

func SendMetrics(metrics []types.Metric) {
	for i := 0; i < len(metrics); i++ {
		reqURL := fmt.Sprintf("http://%s/update/%s/%s/%v", constants.ServerAddress, metrics[i].Type, metrics[i].Name, metrics[i].Value)
		resp, err := http.Post(reqURL, "text/plain", nil)
		if err != nil {
			fmt.Println("Request:", reqURL, "failed, err:", err)
			continue
		}
		defer resp.Body.Close()
	}
}
