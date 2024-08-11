package app

import (
	"fmt"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/config"
	"github.com/Zrossiz/go-metrics/internal/server/services/get"
	"github.com/Zrossiz/go-metrics/internal/server/services/update"
	"github.com/go-chi/chi/v5"
)

func StartServer() {
	config.FlagParse()

	r := chi.NewRouter()

	r.Get("/", get.GetHTMLPageMetric)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", update.UpdateMetric)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", get.GetMetric)
	})

	fmt.Printf("Starting server on %v", config.RunAddr)
	if err := http.ListenAndServe(config.RunAddr, r); err != nil {
		fmt.Println(err)
	}
}
