package handlers

import (
	"cs-server-controller/context_values"
	"net/http"
	"strings"
)

func SendCommandHandler(w http.ResponseWriter, r *http.Request) {
	_, serverInstance, _, err := context_values.GetSteamcmdAndServerInstance(r.Context())
	if err != nil {
		WriteProblemDetail2(w, http.StatusInternalServerError, "internal error", "")
		return
	}

	if !serverInstance.IsRunning() {
		WriteProblemDetail2(w, http.StatusInternalServerError, "server is not running", "")
		return
	}

	q := r.URL.Query()

	command := strings.TrimSpace(q.Get("command"))
	if command == "" {
		WriteProblemDetail2(w, http.StatusInternalServerError, "cant execute empty command", "")
		return
	}

	out, err := serverInstance.SendCommand(command)
	if err != nil {
		WriteProblemDetail2(w, http.StatusInternalServerError, "failed to execute command", "")
		return
	}

	WriteJson(w, out, http.StatusOK)
}
