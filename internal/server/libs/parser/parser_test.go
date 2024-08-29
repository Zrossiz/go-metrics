package parser

import (
	"fmt"
	"os"
	"testing"

	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
)

func TestCollectMetricsFromFile(t *testing.T) {
	testString := `
		{"Name":"Alloc","Type":"gauge","Value":2942224}
		{"Name":"BuckHashSys","Type":"gauge","Value":7544}
		{"Name":"Frees","Type":"gauge","Value":25251}
	`

	filePath := "storage/storage.txt"

	err := os.MkdirAll("storage", 0755)
	if err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	store := memstorage.NewMemStorage()

	file, err := os.Create(filePath)
	if err != nil {
		t.Errorf("create storage file error")
	}
	defer file.Close()

	_, err = file.WriteString(testString)
	if err != nil {
		t.Errorf("insert data in file error")
	}

	CollectMetricsFromFile(filePath, store)

	fmt.Print(store.Metrics)

	if len(store.Metrics) != 3 {
		t.Errorf("epxpected 3 metrics, got %v", len(store.Metrics))
	}
}
