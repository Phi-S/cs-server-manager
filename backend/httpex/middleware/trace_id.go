package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type traceIdKeyType uint

const TraceIdKey traceIdKeyType = 0

func TraceId(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, TraceIdKey, uuid.NewString())
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetTraceId Return the trace id or an empty string If trace id is not present, an error
// will be printed
func GetTraceId(ctx context.Context) string {
	traceId, ok := ctx.Value(TraceIdKey).(string)
	if traceId == "" || !ok {
		slog.Warn("failed to get traceId from context")
	}

	return traceId
}
