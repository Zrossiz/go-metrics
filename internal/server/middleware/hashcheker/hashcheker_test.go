package hashcheker

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/server/config"
	"github.com/Zrossiz/go-metrics/internal/server/libs/hashgenerator"
	"github.com/stretchr/testify/assert"
)

func TestHashCheker_ValidHash(t *testing.T) {
	config.AppConfig = config.Config{Key: "test_key"}
	body := []byte("test_body")

	validHash := hashgenerator.Generate(body, config.AppConfig.Key)

	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(body))
	req.Header.Set("HashSHA256", validHash)

	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := HashCheker(nextHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestHashCheker_InvalidHash(t *testing.T) {
	config.AppConfig = config.Config{Key: "test_key"}
	body := []byte("test_body")

	invalidHash := "invalid_hash"

	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(body))
	req.Header.Set("HashSHA256", invalidHash)

	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := HashCheker(nextHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "invalid hash")
}

func TestHashCheker_EmptyHashHeader(t *testing.T) {
	config.AppConfig = config.Config{Key: "test_key"}
	body := []byte("test_body")

	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(body))

	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := HashCheker(nextHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestHashCheker_EmptyConfigKey(t *testing.T) {
	config.AppConfig = config.Config{Key: ""}
	body := []byte("test_body")

	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(body))
	req.Header.Set("HashSHA256", "any_hash")

	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := HashCheker(nextHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
