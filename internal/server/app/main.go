package app

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/services"
	"github.com/go-chi/chi/v5"
)

func StartServer() {
	ADDRESS := flag.String("a", "localhost:8080", "Port to run the server on")
	flag.Parse()

	r := chi.NewRouter()

	r.Get("/", services.GetHTMLPageMetric)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", services.UpdateMetric)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", services.GetMetric)
	})

	fmt.Printf("Starting server on %s", *ADDRESS)
	if err := http.ListenAndServe(*ADDRESS, r); err != nil {
		fmt.Println(err)
	}
}
