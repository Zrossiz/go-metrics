package parser

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Zrossiz/go-metrics/internal/server/storage"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
	"go.uber.org/zap"
)

func CollectMetricsFromFile(filePath string, logger *zap.Logger, store *memstorage.MemStorage) []storage.Metric {
	var collectedMetrics []storage.Metric

	file, err := os.Open(filePath)
	if err != nil {
		logger.Fatal("file not exist")
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Fatal("failed to read file")
	}

	err = json.Unmarshal(byteValue, &collectedMetrics)
	if err != nil {
		logger.Fatal("failed to unmarshal json")
	}

	for i := 0; i < len(collectedMetrics); i++ {
		curMetric := collectedMetrics[i]
		if curMetric.Type == storage.CounterType {
			store.SetCounter(curMetric.Name, curMetric.Value.(int64))
		}
		if curMetric.Type == storage.GaugeType {
			store.SetGauge(curMetric.Name, curMetric.Value.(float64))
		}
	}

	return collectedMetrics
}

// TODO: сделать обновление метрики
// func UpdateMetrics(filePath string) error {

// }
