package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
)

func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoding := r.Header.Get("Content-Encoding")

		if encoding == "" || encoding != "gzip" {
			next.ServeHTTP(w, r)
			return
		}

		zw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		defer zw.Close()

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(compressWriter{ResponseWriter: w, Writer: zw}, r)
	})

}

func DecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoding := r.Header.Get("Content-Encoding")

		if encoding == "" || encoding != "gzip" {
			next.ServeHTTP(w, r)
			return
		}

		cr, err := newCompressReader(r.Body)
		if err != nil {
			next.ServeHTTP(w, r)
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

func (c *compressReader) Read(p []byte) (n int, err error) {
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

func (c compressWriter) Write(p []byte) (int, error) {
	if c.Writer == nil {
		zw, err := gzip.NewWriterLevel(c.ResponseWriter, gzip.BestSpeed)
		if err != nil {
			return 0, err
		}

		c.Writer = zw
	}

	return c.Writer.Write(p)
}
