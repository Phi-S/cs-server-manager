package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now().UTC()

		lrw := NewLoggingResponseWriter(w)
		h.ServeHTTP(lrw, r)

		traceId := GetTraceId(r.Context())
		duration := time.Now().UTC().Sub(startTime)
		slog.Info("request",
			"traceId", traceId,
			"method", r.Method,
			"protocol", r.Proto,
			"url", r.URL,
			"remote", r.RemoteAddr,
			"duration", duration.String(),
			"duration_nanoseconds", duration.Nanoseconds(),
			"status", lrw.statusCode,
		)
	})
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &LoggingResponseWriter{w, http.StatusOK}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
