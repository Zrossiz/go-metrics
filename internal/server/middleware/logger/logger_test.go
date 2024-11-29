package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestWithLogs_SuccessfulRequest(t *testing.T) {

	core, observedLogs := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, Logger!"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler := WithLogs(logger)(nextHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Hello, Logger!", rr.Body.String())

	assert.Equal(t, 1, observedLogs.Len())
	logEntry := observedLogs.All()[0]

	assert.Equal(t, "request", logEntry.Message)
	assert.Equal(t, "/test", logEntry.ContextMap()["uri"])
	assert.Equal(t, "GET", logEntry.ContextMap()["method"])
	assert.Equal(t, int64(http.StatusOK), logEntry.ContextMap()["status"].(int64))
	assert.Equal(t, int64(len("Hello, Logger!")), logEntry.ContextMap()["size"].(int64))
	assert.NotZero(t, logEntry.ContextMap()["duration"].(time.Duration))
}

func TestWithLogs_NotFoundRequest(t *testing.T) {

	core, observedLogs := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
	rr := httptest.NewRecorder()

	handler := WithLogs(logger)(nextHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	assert.Equal(t, 1, observedLogs.Len())
	logEntry := observedLogs.All()[0]

	assert.Equal(t, "request", logEntry.Message)
	assert.Equal(t, "/notfound", logEntry.ContextMap()["uri"])
	assert.Equal(t, "GET", logEntry.ContextMap()["method"])
	assert.Equal(t, int64(http.StatusNotFound), logEntry.ContextMap()["status"].(int64))

	assert.Greater(t, logEntry.ContextMap()["size"].(int64), int64(0))
	assert.NotZero(t, logEntry.ContextMap()["duration"].(time.Duration))
}

func TestWithLogs_InternalServerError(t *testing.T) {

	core, observedLogs := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	})

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	rr := httptest.NewRecorder()

	handler := WithLogs(logger)(nextHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	assert.Equal(t, 1, observedLogs.Len())
	logEntry := observedLogs.All()[0]

	assert.Equal(t, "request", logEntry.Message)
	assert.Equal(t, "/error", logEntry.ContextMap()["uri"])
	assert.Equal(t, "GET", logEntry.ContextMap()["method"])
	assert.Equal(t, int64(http.StatusInternalServerError), logEntry.ContextMap()["status"].(int64))
	assert.Equal(t, int64(len("internal server error\n")), logEntry.ContextMap()["size"].(int64))
	assert.NotZero(t, logEntry.ContextMap()["duration"].(time.Duration))
}
