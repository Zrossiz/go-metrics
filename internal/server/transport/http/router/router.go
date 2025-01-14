// Package router provides HTTP routes for metrics management,
// connecting HTTP requests to the respective handlers
package router

import (
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/middleware/cryptochecker"
	"github.com/Zrossiz/go-metrics/internal/server/middleware/gzip"
	ipchecker "github.com/Zrossiz/go-metrics/internal/server/middleware/ipcheker"
	"github.com/Zrossiz/go-metrics/internal/server/middleware/logger"
	"github.com/Zrossiz/go-metrics/internal/server/security"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// MetricRouter defines an interface for handling metric-related routes.
// Each method represents a handler for a specific HTTP route.
type MetricRouter interface {
	// GetHTML handles requests to render metrics in HTML format.
	GetHTML(rw http.ResponseWriter, r *http.Request)
	// CreateParamMetric handles requests to create a metric using URL parameters.
	CreateParamMetric(rw http.ResponseWriter, r *http.Request)
	// CreateJSONMetric handles requests to create a metric from a JSON payload
	CreateJSONMetric(rw http.ResponseWriter, r *http.Request)
	// GetStringMetric retrieves the value of a metric in string format.
	GetStringMetric(rw http.ResponseWriter, r *http.Request)
	// GetJSONMetric retrieves a metric and returns it in JSON format.
	GetJSONMetric(rw http.ResponseWriter, r *http.Request)
	// PingDB checks the database connection.
	PingDB(rw http.ResponseWriter, _ *http.Request)
	//CreateBatchJSONMetrics processes batch creation of multiple metrics from a JSON array/
	CreateBatchJSONMetrics(rw http.ResponseWriter, _ *http.Request)
}

// New initializes and returns a new HTTP router.
func New(handl MetricRouter, log *zap.Logger) http.Handler {
	// Create a new router
	r := chi.NewRouter()

	// Attach logging middleware
	r.Use(logger.WithLogs(log))
	// Attach GZIP decompression and compression middlewares
	r.Use(gzip.DecompressMiddleware)
	r.Use(gzip.CompressMiddleware)

	// Route for rendering metrics in HTML format
	r.Get("/", handl.GetHTML)

	// Route for checking database connect
	r.Get("/ping", handl.PingDB)

	// Routes for updating metrics
	r.Route("/update", func(r chi.Router) {
		// Update a metric using URL parameters
		r.Post("/{type}/{name}/{value}", handl.CreateParamMetric)
		// Update a metric using a JSON payload
		r.With(
			func(next http.Handler) http.Handler {
				return cryptochecker.DecryptMiddleware(next, security.CheckCryptoBody)
			},
			func(next http.Handler) http.Handler {
				return ipchecker.TrustedIPCheckerMiddleware(next)
			},
		).Post("/", handl.CreateJSONMetric)
	})

	// Route for batch updates of metrics
	r.Post("/updates/", handl.CreateBatchJSONMetrics)

	// Routes for retrieving metric values
	r.Route("/value", func(r chi.Router) {
		// Get the value of a metric by its type and name
		r.Get("/{type}/{name}", handl.GetStringMetric)
		// Get a metric useing a JSON payload
		r.Post("/", handl.GetJSONMetric)
	})

	return r
}
