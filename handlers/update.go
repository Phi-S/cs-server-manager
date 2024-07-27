package handlers

import (
	"cs-server-controller/context_values"
	"net/http"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	lock, serverInstance, steamcmdInstance, err := context_values.GetSteamcmdAndServerInstance(r.Context())
	if err != nil {
		WriteProblemDetail2(w, http.StatusInternalServerError, "internal error", "")
		return
	}

	lock.Lock()
	defer lock.Unlock()

	if steamcmdInstance.IsRunning() {
		WriteProblemDetail2(w, http.StatusInternalServerError, "another update is already running", "")
		return
	}

	if serverInstance.IsRunning() {
		WriteProblemDetail2(w, http.StatusInternalServerError, "cant update server. Server is running", "")
		return
	}

	if err := steamcmdInstance.Update(false); err != nil {
		WriteProblemDetail2(w, http.StatusInternalServerError, "failed update server", "")
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
