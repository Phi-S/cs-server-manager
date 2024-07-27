package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ProblemDetail struct {
	Status  int    `json:"status"`
	Message string `json:"title"`
	TraceId string `json:"trace-id"`
}

func WriteProblemDetail2(w http.ResponseWriter, code int, msg string, traceId string) {
	WriteProblemDetail(w, ProblemDetail{
		Status:  code,
		Message: msg,
		TraceId: traceId,
	})
}

func WriteProblemDetail(w http.ResponseWriter, pd ProblemDetail) {
	w.WriteHeader(pd.Status)
	w.Header().Set("Content-Type", "application/problem+json")
	if err := json.NewEncoder(w).Encode(pd); err != nil {
		slog.Error("failed to encode problem detail", "ProblemDetail", pd)
	}
}

func WriteJson(w http.ResponseWriter, p any, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(p); err != nil {
		slog.Error("failed to encode problem detail", "ProblemDetail", p)
	}
}
