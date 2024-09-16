package filestorage

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"golang.org/x/exp/rand"
)

type FileStorage struct {
	data []models.Metric
	mu   sync.Mutex
	path string
}

func New(filePath string) *FileStorage {
	return &FileStorage{
		data: make([]models.Metric, 0),
		path: filePath,
	}
}

func (f *FileStorage) SetGauge(metric dto.PostMetricDto) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i := 0; i < len(f.data); i++ {
		if metric.ID == f.data[i].Name {
			f.data[i].Value = metric.Value
			return nil
		}
	}

	f.data = append(f.data, models.Metric{
		ID:    uint(rand.Int63()),
		Name:  metric.ID,
		Type:  models.GaugeType,
		Value: metric.Value,
	})

	return nil
}

func (f *FileStorage) SetCounter(metric dto.PostMetricDto) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i := 0; i < len(f.data); i++ {
		if metric.ID == f.data[i].Name {
			newValue := *f.data[i].Delta + int64(*metric.Delta)
			f.data[i].Delta = &newValue
			return nil
		}
	}

	f.data = append(f.data, models.Metric{
		ID:    uint(rand.Int63()),
		Name:  metric.ID,
		Type:  models.CounterType,
		Delta: metric.Delta,
	})

	return nil
}

func (f *FileStorage) Get(name string) (*models.Metric, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i := 0; i < len(f.data); i++ {
		if f.data[i].Name == name {
			return &f.data[i], nil
		}
	}

	return nil, nil
}

func (f *FileStorage) GetAll() (*[]models.Metric, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return &f.data, nil
}

func (f *FileStorage) Load(filePath string) error {
	var collectedMetrics []models.Metric

	file, err := os.Open(f.path)
	if os.IsNotExist(err) {
		file, err = os.Create(f.path)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var metric models.Metric
		line := scanner.Text()

		if err := json.Unmarshal([]byte(line), &metric); err != nil {
			continue
		}

		collectedMetrics = append(collectedMetrics, metric)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for _, curMetric := range collectedMetrics {
		metricDTO := dto.PostMetricDto{
			ID:    curMetric.Name,
			MType: curMetric.Type,
		}

		switch curMetric.Type {
		case models.CounterType:
			metricDTO.Delta = curMetric.Delta
			f.SetCounter(metricDTO)
		case models.GaugeType:
			metricDTO.Value = curMetric.Value
			f.SetGauge(metricDTO)
		}
	}

	return nil
}

func (f *FileStorage) Save(filePath string) error {
	file, err := os.OpenFile(f.path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var newMetrics string

	for _, s := range f.data {
		str, err := json.Marshal(s)
		if err != nil {
			return err
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

func (f *FileStorage) SetBatch(body []dto.PostMetricDto) error {
	for i := 0; i < len(body); i++ {
		if body[i].MType == models.CounterType {
			_ = f.SetCounter(body[i])
			continue
		}

		_ = f.SetGauge(body[i])
	}

	return nil
}

func (f *FileStorage) Close(filePath string) error {
	err := f.Save(filePath)
	if err != nil {
		return err
	}

	return nil
}
