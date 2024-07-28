package handlers

import (
	"cs-server-controller/httpex/errorwrp"
	"cs-server-controller/middleware"
	"net/http"
)

type StatusResponse struct {
	ServerRunning   bool `json:"server-running"`
	SteamCmdRunning bool `json:"steamcmd-running"`
}

func StatusHandler(r *http.Request) (errorwrp.HttpResponse, *errorwrp.HttpError) {
	_, server, steamcmd, err := middleware.GetSteamcmdAndServerInstance(r.Context())
	if err != nil {
		return errorwrp.NewHttpErrorInternalServerError("internal error", err)
	}

	resp := StatusResponse{
		ServerRunning:   server.IsRunning(),
		SteamCmdRunning: steamcmd.IsRunning(),
	}
	return errorwrp.NewJsonHttpResponse(http.StatusOK, resp)
}
