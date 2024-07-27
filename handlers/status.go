package handlers

import (
	"cs-server-controller/context_values"
	"net/http"
)

type StatusResponse struct {
	ServerRunning   bool `json:"server-running"`
	SteamCmdRunning bool `json:"steamcmd-running"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	_, server, steamcmd, err := context_values.GetSteamcmdAndServerInstance(r.Context())
	if err != nil {
		WriteProblemDetail2(w, http.StatusInternalServerError, "internal error", "")
		return
	}

	response := StatusResponse{
		ServerRunning:   server.IsRunning(),
		SteamCmdRunning: steamcmd.IsRunning(),
	}
	WriteJson(w, response, http.StatusOK)
}
