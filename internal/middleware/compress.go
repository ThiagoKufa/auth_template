package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
	headerWritten bool
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	if !w.headerWritten {
		w.Header().Del("Content-Length")
		w.ResponseWriter.WriteHeader(status)
		w.headerWritten = true
	}
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.headerWritten {
		w.WriteHeader(http.StatusOK)
	}
	return w.Writer.Write(b)
}

// Compress retorna um middleware que comprime as respostas usando gzip
func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzw := &gzipResponseWriter{
			Writer:         gz,
			ResponseWriter: w,
		}

		gzw.Header().Set("Content-Encoding", "gzip")
		gzw.Header().Add("Vary", "Accept-Encoding")

		next.ServeHTTP(gzw, r)
	})
}
