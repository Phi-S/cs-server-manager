package middleware

import (
    "log/slog"
    "net/http"
    "runtime/debug"
)

func Recover(next http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
                traceId := GetTraceId(r.Context())
                debugStack := debug.Stack()
                slog.Error("handler panics", "value", rvr, "trace-id", traceId, "stack", debugStack)
                w.WriteHeader(http.StatusInternalServerError)
            }
        }()

        next.ServeHTTP(w, r)
    }

    return http.HandlerFunc(fn)
}
