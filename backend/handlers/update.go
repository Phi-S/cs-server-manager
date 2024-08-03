package handlers

import (
	"cs-server-controller/ctxex"
	"cs-server-controller/httpex/errorwrp"
	"net/http"
)

func UpdateHandler(r *http.Request) (errorwrp.HttpResponse, *errorwrp.HttpError) {
	lock, serverInstance, steamcmdInstance, err := ctxex.GetSteamcmdAndServerInstance(r.Context())
	if err != nil {
		return errorwrp.NewHttpErrorInternalServerError("internal error", err)
	}

	lock.Lock()
	defer lock.Unlock()

	if steamcmdInstance.IsRunning() {
		return errorwrp.NewHttpErrorInternalServerError2("update is already running")
	}

	if serverInstance.IsRunning() {
		return errorwrp.NewHttpErrorInternalServerError2("can not update server. server is still running")
	}

	if err := steamcmdInstance.Update(false); err != nil {
		return errorwrp.NewHttpErrorInternalServerError("failed to update server", err)
	}

	return errorwrp.NewHttpResponse(http.StatusAccepted)
}
