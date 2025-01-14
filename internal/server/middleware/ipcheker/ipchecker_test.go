package ipchecker

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/server/config"
	"github.com/stretchr/testify/assert"
)

func TestIpCheckerMiddleware_Success(t *testing.T) {
	expectedIP := "1.1.1.1.2"
	config.AppConfig.TrustedSubnet = expectedIP

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("X-Real-IP", expectedIP)

	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	handler := TrustedIpCheckerMiddleware(nextHandler)

	handler.ServeHTTP(rr, req)

	assert.NotEqual(t, rr.Code, http.StatusForbidden)
}

func TestIpCheckerMiddleware_Fail(t *testing.T) {
	config.AppConfig.TrustedSubnet = "1.1.1.1.2"

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("X-Real-IP", "1.1.1.1.1")

	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	handler := TrustedIpCheckerMiddleware(nextHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusForbidden)
}
