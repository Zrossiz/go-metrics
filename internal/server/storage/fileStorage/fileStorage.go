package filestorage

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Zrossiz/go-metrics/internal/server/storage"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
	"go.uber.org/zap"
)

func CollectMetricsFromFile(relativePath string, store *memstorage.MemStorage) ([]storage.Metric, error) {
	var collectedMetrics []storage.Metric

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(workingDir, relativePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var metric storage.Metric
		line := scanner.Text()

		if err := json.Unmarshal([]byte(line), &metric); err != nil {
			continue
		}

		collectedMetrics = append(collectedMetrics, metric)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	for _, curMetric := range collectedMetrics {
		switch curMetric.Type {
		case storage.CounterType:
			if value, ok := curMetric.Value.(float64); ok {
				store.SetCounter(curMetric.Name, int64(value))
			}
		case storage.GaugeType:
			if value, ok := curMetric.Value.(float64); ok {
				store.SetGauge(curMetric.Name, value)
			}
		}
	}

	return collectedMetrics, nil
}

func UpdateMetrics(relativePath string, logger *zap.Logger, store *memstorage.MemStorage) error {
	workingDir, err := os.Getwd()
	if err != nil {
		logger.Error("error getting working directory:")
		return err
	}

	filePath := filepath.Join(workingDir, relativePath)

	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var newMetrics string

	for _, s := range store.Metrics {
		str, err := json.Marshal(s)
		if err != nil {
			logger.Sugar().Warnf("error to parse metric: %s", s.Name)
			continue
		}
		newMetrics += string(str)
		newMetrics += "\n"
	}

	_, err = file.WriteString(newMetrics)
	if err != nil {
		return err
	}

	return nil
}
