package main

import (
	"fmt"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/http/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handlers.UpdateMetricHandler)
	fmt.Println("Starting server on port 8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
