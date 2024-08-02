package main

import (
	"fmt"
	"net/http"
)

func defaultRoute(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Write([]byte(`
		{
			"success": true,
			"message": "hello world"
		}
		`))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultRoute)
	fmt.Println("Starting server on port 8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
