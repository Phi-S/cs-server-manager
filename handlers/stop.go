package handlers

import (
	"cs-server-controller/context_values"
	"net/http"
)

func StopHandler(w http.ResponseWriter, r *http.Request) {
	_, server, _, err := context_values.GetSteamcmdAndServerInstance(r.Context())
	if err != nil {
		WriteProblemDetail2(w, http.StatusInternalServerError, "internal error", "")
		return
	}

	if err := server.Stop(); err != nil {
		WriteProblemDetail2(w, 551, "failed to stop server", "")
		http.Error(w, "failed to stop server", http.StatusInternalServerError)
		return
	}
}
