package server

import (
	"fmt"
	"net/http"

	serverServices "github.com/Zrossiz/go-metrics/internal/services/server"
	"github.com/go-chi/chi/v5"
)

func StartServer() {
	r := chi.NewRouter()

	r.Get("/", serverServices.GetHtmlPageMetric)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", serverServices.UpdateMetric)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", serverServices.GetMetric)
	})

	fmt.Println("Starting server on port: 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println(err)
	}
}
