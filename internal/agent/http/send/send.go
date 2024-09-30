package send

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Zrossiz/go-metrics/internal/agent/config"
	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
	"github.com/Zrossiz/go-metrics/internal/agent/dto"
)

const maxRetries = 3
const retryDelay = 1 * time.Second

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

func GzipMetrics(metrics []types.Metric, addr string, cfg config.Config) []types.Metric {
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

		hash := computeHash(gzippedData, cfg.Key)

		request, err := getRequest("POST", reqURL, gzippedData, hash)
		if err != nil {
			log.Println("error get request:", err)
		}

		err = sendWithRetry(request)
		if err != nil {
			log.Println("error send metric", err)
		}

		sendedMetrics = append(sendedMetrics, metrics[i])
	}
	return sendedMetrics
}

func BatchGzipMetrics(metrics []types.Metric, addr string, cfg *config.Config) {
	reqURL := fmt.Sprintf("http://%s/updates/", addr)

	postDtoMetrics := []dto.MetricDTO{}

	for i := 0; i < len(metrics); i++ {
		postMetricDto := dto.MetricDTO{
			ID:    metrics[i].Name,
			MType: metrics[i].Type,
		}

		float64Value, okFloat := metrics[i].Value.(float64)
		if okFloat {
			postMetricDto.Value = &float64Value
		}

		int64Value, okInt := metrics[i].Value.(int64)
		if okInt {
			postMetricDto.Delta = &int64Value
		}

		postDtoMetrics = append(postDtoMetrics, postMetricDto)
	}

	bytesData, err := json.Marshal(postDtoMetrics)
	if err != nil {
		log.Println("failed to marshal metrics to JSON:", err)
		return
	}

	var gzippedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&gzippedData)

	_, err = gzipWriter.Write(bytesData)
	if err != nil {
		log.Println("failed to write JSON to gzip:", err)
		gzipWriter.Close()
		return
	}

	gzipWriter.Close()

	hash := computeHash(gzippedData, cfg.Key)

	request, err := getRequest("POST", reqURL, gzippedData, hash)
	if err != nil {
		log.Println("error get request:", err)
	}

	err = sendWithRetry(request)
	if err != nil {
		log.Println("error send metric", err)
	}
}

func getBytesMetricDTO(metric types.Metric) ([]byte, error) {
	jsonBody := dto.MetricDTO{
		ID:    metric.Name,
		MType: metric.Type,
	}

	float64Value, okFloat := metric.Value.(float64)
	if okFloat {
		jsonBody.Value = &float64Value
	}

	int64Value, okInt := metric.Value.(int64)
	if okInt {
		jsonBody.Delta = &int64Value
	}

	jsonData, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func getRequest(method string, reqURL string, data bytes.Buffer, hash string) (*http.Request, error) {
	req, err := http.NewRequest(method, reqURL, &data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HashSHA256", hash)

	return req, nil
}

func sendWithRetry(request *http.Request) error {
	delay := retryDelay
	for i := 0; i < maxRetries; i++ {
		client := &http.Client{}
		resp, err := client.Do(request)
		resp.Body.Close()

		if err != nil {
			log.Printf("Failed to send request: %v\n", err)
		} else if resp.StatusCode == 200 {
			return nil
		} else {
			log.Printf("Failed to send request: status code %d\n", resp.StatusCode)
		}

		time.Sleep(delay)
		delay += 2 * time.Second
	}
	return fmt.Errorf("failed to send request after %d attempts", maxRetries)
}

func computeHash(data bytes.Buffer, cfgKey string) string {
	if cfgKey == "" {
		return ""
	}
	h := sha256.New()
	h.Write(data.Bytes())
	h.Write([]byte(cfgKey))
	return hex.EncodeToString(h.Sum(nil))
}
