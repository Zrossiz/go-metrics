package send

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
	"github.com/Zrossiz/go-metrics/internal/agent/dto"
)

func Metrics(metrics []types.Metric, addr string) []types.Metric {
	var sendedMetrics []types.Metric

	for i := 0; i < len(metrics); i++ {
		reqURL := fmt.Sprintf("http://%s/update/", addr)
		jsonBody := dto.MetricDTO{
			ID:    metrics[i].Name,
			MType: metrics[i].Type,
		}

		switch v := metrics[i].Value.(type) {
		case int64:
			jsonBody.Delta = &v
		case float64:
			jsonBody.Value = &v
		default:
			log.Println("Unsupported metric type for metric:", metrics[i].Name)
			continue
		}

		jsonData, err := json.Marshal(jsonBody)
		if err != nil {
			log.Println("Failed to marshal jsonBody:", err)
			continue
		}

		resp, err := http.Post(reqURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Request:", reqURL, "failed, err:", err)
			continue
		}

		sendedMetrics = append(sendedMetrics, metrics[i])
		resp.Body.Close()
	}
	return sendedMetrics
}
