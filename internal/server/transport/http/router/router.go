package router

import (
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/middleware/gzip"
	"github.com/Zrossiz/go-metrics/internal/server/middleware/logger"
	"github.com/Zrossiz/go-metrics/internal/server/service"
	"github.com/Zrossiz/go-metrics/internal/server/transport/http/handler"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func New(s *service.MetricService, log *zap.Logger) http.Handler {
	handl := handler.New(s, log)

	r := chi.NewRouter()

	r.Use(logger.WithLogs(log))

	r.Use(gzip.DecompressMiddleware)

	r.Use(gzip.CompressMiddleware)

	return r
}
