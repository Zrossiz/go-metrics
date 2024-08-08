package app

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/services"
	"github.com/go-chi/chi/v5"
)

func StartServer() {
	PORT := flag.String("a", "8080", "Port to run the server on")
	flag.Parse()

	addr := fmt.Sprintf(":%s", *PORT)

	r := chi.NewRouter()

	r.Get("/", services.GetHTMLPageMetric)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", services.UpdateMetric)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", services.GetMetric)
	})

	fmt.Printf("Starting server on port: %s", *PORT)
	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Println(err)
	}
}
