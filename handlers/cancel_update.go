package handlers

import (
	"cs-server-controller/ctxex"
	"cs-server-controller/httpex/errorwrp"
	"net/http"
)

func CancelUpdateHandler(r *http.Request) (errorwrp.HttpResponse, *errorwrp.HttpError) {
	_, _, steamcmdInstance, err := ctxex.GetSteamcmdAndServerInstance(r.Context())
	if err != nil {
		return errorwrp.NewHttpError(http.StatusInternalServerError, "internal error", err)
	}

	if !steamcmdInstance.IsRunning() {
		return errorwrp.NewHttpErrorInternalServerError2("nothing to cancel. steamcmd is not running")
	}

	if err := steamcmdInstance.Cancel(); err != nil {
		return errorwrp.NewHttpError(http.StatusInternalServerError, "failed to cancel server update", err)
	}

	return errorwrp.NewOkHttpResponse()
}
