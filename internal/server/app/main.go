package app

import (
	"fmt"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/services"
	"github.com/go-chi/chi/v5"
)

func StartServer() {
	r := chi.NewRouter()

	r.Get("/", services.GetHTMLPageMetric)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", services.UpdateMetric)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", services.GetMetric)
	})

	fmt.Println("Starting server on port: 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println(err)
	}
}
