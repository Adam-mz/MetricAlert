package middleware

import (
	"compress/gzip"
	"io"
	"mime"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	writer      *gzip.Writer
	wroteHeader bool
	compress    bool
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if w.compress {
		return w.writer.Write(b)
	}

	return w.ResponseWriter.Write(b)
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}

	w.wroteHeader = true
	w.compress = shouldCompressContentType(w.Header().Get("Content-Type"))

	if w.compress {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")
		w.writer = gzip.NewWriter(w.ResponseWriter)
	}

	w.ResponseWriter.WriteHeader(statusCode)
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if acceptsEncoding(r.Header.Get("Content-Encoding"), "gzip") {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "invalid gzip body", http.StatusBadRequest)
				return
			}
			defer reader.Close()

			r.Body = io.NopCloser(reader)
		}

		if !acceptsEncoding(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		wrapped := &gzipResponseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(wrapped, r)
		if wrapped.writer != nil {
			_ = wrapped.writer.Close()
		}
	})
}

func shouldCompressContentType(contentType string) bool {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	return mediaType == "application/json" || mediaType == "text/html"
}

func acceptsEncoding(header string, encoding string) bool {
	for _, value := range strings.Split(header, ",") {
		name, _, _ := strings.Cut(strings.TrimSpace(value), ";")
		if strings.EqualFold(name, encoding) {
			return true
		}
	}

	return false
}
