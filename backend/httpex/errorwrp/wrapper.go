package errorwrp

import (
	"cs-server-manager/httpex/middleware"
	"encoding/json"
	"log/slog"
	"net/http"
)

func GET(mux *http.ServeMux, path string, handler ErrorHandlerFunc) {
	mux.HandleFunc("GET "+path, WrapErrorHandler(handler))
}

func POST(mux *http.ServeMux, path string, handler ErrorHandlerFunc) {
	mux.HandleFunc("POST "+path, WrapErrorHandler(handler))
}

type ErrorHandlerFunc func(r *http.Request) (HttpResponse, *HttpError)

func WrapErrorHandler(errorHandler ErrorHandlerFunc) http.HandlerFunc {
	hand := func(w http.ResponseWriter, r *http.Request) {
		resp, err := errorHandler(r)

		if err != nil {
			traceId := middleware.GetTraceId(r.Context())
			WriteErrorResponse(w, err.Status, err.ResponseMessage, traceId)
			slog.Error("request failed", "trace-id", traceId, "error", err.InternalError)
		} else {
			if resp.Response != nil {
				WriteJsonResponse(w, resp.Status, resp.Response)
			} else {
				WriteResponse(w, resp.Status)
			}
		}
	}

	return hand
}

func WriteErrorResponse(w http.ResponseWriter, status int, msg string, traceId string) {
	resp := struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		TraceId string `json:"trace-id"`
	}{
		Status:  status,
		Message: msg,
		TraceId: traceId,
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode error response", "response", resp)
	}
}

func WriteJsonResponse(w http.ResponseWriter, status int, resp any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json response", "response", resp)
	}
}

func WriteResponse(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}
