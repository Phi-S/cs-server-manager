package handlers

import (
	"cs-server-controller/ctxex"
	"cs-server-controller/httpex/errorwrp"
	json_file "cs-server-controller/jsonfile"
	"cs-server-controller/server"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func StartHandler(r *http.Request) (errorwrp.HttpResponse, *errorwrp.HttpError) {
	lock, s, steamcmd, err := ctxex.GetSteamcmdAndServerInstance(r.Context())
	if err != nil {
		return errorwrp.NewHttpErrorInternalServerError("internal error", err)
	}

	if steamcmd.IsRunning() {
		return errorwrp.NewHttpErrorInternalServerError2("steamcmd is running")
	}

	if s.IsRunning() {
		return errorwrp.NewHttpErrorInternalServerError2("server is already running")
	}

	startParameterJsonFile, err := ctxex.Get[*json_file.JsonFile[server.StartParameters]](r.Context(), ctxex.StartParametersJsonFileKey)
	if err != nil {
		return errorwrp.HttpResponse{}, nil
	}

	q := r.URL.Query()

	startParameters, err := getStarParameters(startParameterJsonFile, q)
	if err != nil {
		return errorwrp.NewHttpErrorInternalServerError("internal error", err)
	}

	lock.Lock()
	defer lock.Unlock()

	if err := s.Start(*startParameters); err != nil {
		return errorwrp.NewHttpErrorInternalServerError("failed to start server", err)
	}

	if err := startParameterJsonFile.Write(*startParameters); err != nil {
		slog.Warn("server started but failed to save valid start parameters to file. " + err.Error())
	}

	return errorwrp.NewOkHttpResponse()
}

func getStarParameters(startParametersJsonFileInstance *json_file.JsonFile[server.StartParameters], query url.Values) (*server.StartParameters, error) {
	startParameters, err := startParametersJsonFileInstance.Read()
	if err != nil {
		return nil, err
	}

	if hostname := strings.TrimSpace(query.Get("name")); hostname != "" {
		startParameters.Hostname = hostname
	}

	if password := strings.TrimSpace(query.Get("pw")); password != "" {
		startParameters.Password = password
	}

	if startMap := strings.TrimSpace(query.Get("map")); startMap != "" {
		startParameters.StartMap = startMap
	}

	if maxPlayersString := strings.TrimSpace(query.Get("maxPlayers")); maxPlayersString != "" {
		maxPlayersUint64, err := strconv.ParseUint(maxPlayersString, 10, 8)
		if err != nil {
			return nil, err
		}

		if maxPlayersUint64 > 255 {
			return nil, fmt.Errorf("maxPlayers parameter must be a valid number between 1 and 255")
		}

		startParameters.MaxPlayers = uint8(maxPlayersUint64)
	}

	if loginToken := strings.TrimSpace(query.Get("loginToken")); loginToken != "" {
		startParameters.SteamLoginToken = loginToken
	}

	return startParameters, nil
}
