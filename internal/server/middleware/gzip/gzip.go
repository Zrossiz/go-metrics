package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoding := r.Header.Get("Accept-Encoding")
		if !strings.Contains(encoding, "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")
		zw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		defer zw.Close()

		next.ServeHTTP(&compressWriter{ResponseWriter: w, Writer: zw}, r)
	})
}

func DecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") != "gzip" {
			next.ServeHTTP(w, r)
			return
		}

		cr, err := newCompressReader(r.Body)
		if err != nil {
			http.Error(w, "failed to decompress request", http.StatusBadRequest)
			return
		}
		defer cr.Close()

		r.Body = cr
		next.ServeHTTP(w, r)
	})
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c *compressReader) Read(p []byte) (int, error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

type compressWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.Writer.Write(p)
}

func (c *compressWriter) Close() error {
	if writer, ok := c.Writer.(*gzip.Writer); ok {
		return writer.Close()
	}
	return nil
}
