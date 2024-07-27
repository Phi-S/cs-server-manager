package handlers

import (
	"cs-server-controller/context_values"
	"net/http"
)

func CancelUpdateHandler(w http.ResponseWriter, r *http.Request) {
	_, _, steamcmdInstance, err := context_values.GetSteamcmdAndServerInstance(r.Context())
	if err != nil {
		WriteProblemDetail2(w, http.StatusInternalServerError, "internal error", "")
		return
	}

	if err := steamcmdInstance.Cancel(); err != nil {
		WriteProblemDetail2(w, http.StatusInternalServerError, "failed to cancel server update", "")
		return
	}
}
