package app

import (
	"fmt"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/config"
	"github.com/Zrossiz/go-metrics/internal/server/services/get"
	"github.com/Zrossiz/go-metrics/internal/server/services/update"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
	"github.com/go-chi/chi/v5"
)

func StartServer() {
	config.FlagParse()

	r := chi.NewRouter()

	store := memstorage.NewMemStorage()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		get.GetHTMLPageMetric(w, r, *store)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
			update.UpdateMetric(w, r, store)
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
			get.GetMetric(w, r, *store)
		})
	})

	fmt.Printf("Starting server on %v", config.RunAddr)
	if err := http.ListenAndServe(config.RunAddr, r); err != nil {
		fmt.Println(err)
	}
}
