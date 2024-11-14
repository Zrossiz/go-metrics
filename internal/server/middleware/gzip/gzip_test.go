package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressMiddleware_NoGzipHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	compressHandler := CompressMiddleware(handler)
	compressHandler.ServeHTTP(rr, req)

	assert.NotContains(t, rr.Header().Get("Content-Encoding"), "gzip")
	assert.Equal(t, "Hello, World!", rr.Body.String())
}

func TestCompressMiddleware_WithGzipHeader(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Gzip!"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()

	compressHandler := CompressMiddleware(handler)
	compressHandler.ServeHTTP(rr, req)

	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

	gzipReader, err := gzip.NewReader(rr.Body)
	assert.NoError(t, err)
	defer gzipReader.Close()

	uncompressedBody, err := io.ReadAll(gzipReader)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, Gzip!", string(uncompressedBody))
}

func TestDecompressMiddleware_WithGzipEncodedRequest(t *testing.T) {

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err := gzipWriter.Write([]byte("Hello, Decompress!"))
	assert.NoError(t, err)
	gzipWriter.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Encoding", "gzip")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		w.Write(body)
	})

	rr := httptest.NewRecorder()

	decompressHandler := DecompressMiddleware(handler)
	decompressHandler.ServeHTTP(rr, req)

	assert.Equal(t, "Hello, Decompress!", rr.Body.String())
}

func TestDecompressMiddleware_NoGzipEncodingHeader(t *testing.T) {

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("Hello, No Gzip!"))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		w.Write(body)
	})

	rr := httptest.NewRecorder()

	decompressHandler := DecompressMiddleware(handler)
	decompressHandler.ServeHTTP(rr, req)

	assert.Equal(t, "Hello, No Gzip!", rr.Body.String())
}
