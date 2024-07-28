package handlers

import (
	"cs-server-controller/httpex/errorwrp"
	"cs-server-controller/middleware"
	"net/http"
	"strings"
)

func SendCommandHandler(r *http.Request) (errorwrp.HttpResponse, *errorwrp.HttpError) {
	_, serverInstance, _, err := middleware.GetSteamcmdAndServerInstance(r.Context())
	if err != nil {
		return errorwrp.NewHttpErrorInternalServerError("internal error", err)
	}

	if !serverInstance.IsRunning() {
		return errorwrp.NewHttpErrorInternalServerError2("server is not running")
	}

	q := r.URL.Query()

	command := strings.TrimSpace(q.Get("command"))
	if command == "" {
		return errorwrp.NewHttpErrorInternalServerError2("cant execute empty command")
	}

	out, err := serverInstance.SendCommand(command)
	if err != nil {
		return errorwrp.NewHttpErrorInternalServerError("internal error", err)
	}

	resp := struct {
		Output []string `json:"output"`
	}{
		Output: strings.Split(out, "\n"),
	}

	return errorwrp.NewJsonHttpResponse(http.StatusOK, resp)
}
