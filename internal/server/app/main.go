package app

import (
	"fmt"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/config"
	"github.com/Zrossiz/go-metrics/internal/server/services"
	"github.com/go-chi/chi/v5"
)

func StartServer() {
	config.FlagParse()

	r := chi.NewRouter()

	r.Get("/", services.GetHTMLPageMetric)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", services.UpdateMetric)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", services.GetMetric)
	})

	fmt.Printf("Starting server on %v", config.RunAddr)
	if err := http.ListenAndServe(config.RunAddr, r); err != nil {
		fmt.Println(err)
	}
}
