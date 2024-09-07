package send

import (
	"bytes"
	"compress/gzip"
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

func GzipMetrics(metrics []types.Metric, addr string) []types.Metric {
	var sendedMetrics []types.Metric

	for i := 0; i < len(metrics); i++ {
		reqURL := fmt.Sprintf("http://%s/update/", addr)

		var gzippedData bytes.Buffer
		gzipWriter := gzip.NewWriter(&gzippedData)

		bytesData, err := getBytesMetricDTO(metrics[i])
		if err != nil {
			log.Println("failed get bytes from metric: ", err)
			continue
		}

		_, err = gzipWriter.Write(bytesData)
		if err != nil {
			log.Println("failed to write json to gzip:", err)
			gzipWriter.Close()
			continue
		}

		gzipWriter.Close()

		_, err = sendMetric("POST", reqURL, gzippedData)
		if err != nil {
			log.Println("error send metric")
			continue
		}

		sendedMetrics = append(sendedMetrics, metrics[i])
	}
	return sendedMetrics
}

func getBytesMetricDTO(metric types.Metric) ([]byte, error) {
	jsonBody := dto.MetricDTO{
		ID:    metric.Name,
		MType: metric.Type,
	}

	switch v := metric.Value.(type) {
	case int64:
		jsonBody.Delta = &v
	case float64:
		jsonBody.Value = &v
	default:
		log.Println("unsupported metric type for metric:", metric.Name)
	}

	jsonData, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func sendMetric(method string, reqURL string, data bytes.Buffer) (bool, error) {
	req, err := http.NewRequest(method, reqURL, &data)
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	resp.Body.Close()

	return true, nil
}
