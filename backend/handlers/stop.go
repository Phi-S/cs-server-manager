package handlers

import (
	"cs-server-controller/ctxex"
	"cs-server-controller/httpex/errorwrp"
	"net/http"
)

func StopHandler(r *http.Request) (errorwrp.HttpResponse, *errorwrp.HttpError) {
	_, server, _, err := ctxex.GetSteamcmdAndServerInstance(r.Context())
	if err != nil {
		return errorwrp.NewHttpErrorInternalServerError("internal error", err)
	}

	if !server.IsRunning() {
		return errorwrp.NewHttpErrorInternalServerError2("nothing to stop. server is not running")
	}

	if err := server.Stop(); err != nil {
		return errorwrp.NewHttpErrorInternalServerError("failed to stop server", err)
	}

	return errorwrp.NewOkHttpResponse()
}
