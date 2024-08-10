package send

import (
	"fmt"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
)

func SendMetrics(metrics []types.Metric, addr string) []types.Metric {
	var sendedMetrics []types.Metric

	for i := 0; i < len(metrics); i++ {
		reqURL := fmt.Sprintf("http://%s/update/%s/%s/%v", addr, metrics[i].Type, metrics[i].Name, metrics[i].Value)
		resp, err := http.Post(reqURL, "text/plain", nil)
		if err != nil {
			fmt.Println("Request:", reqURL, "failed, err:", err)
			continue
		}

		sendedMetrics = append(sendedMetrics, metrics[i])
		resp.Body.Close()
	}
	return sendedMetrics
}
