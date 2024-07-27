package handlers

import (
	"cs-server-controller/context_values"
	json_file "cs-server-controller/jsonfile"
	"cs-server-controller/server"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func StartHandler(w http.ResponseWriter, r *http.Request) {
	lock, s, steamcmd, err := context_values.GetSteamcmdAndServerInstance(r.Context())
	if err != nil {
		WriteProblemDetail2(w, http.StatusInternalServerError, "internal error", "")
		return
	}

	if steamcmd.IsRunning() {
		http.Error(w, "steamcmd is running", http.StatusInternalServerError)
		return
	}

	if s.IsRunning() {
		http.Error(w, "server is already running", http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()

	//startParametersJsonPath := "/home/desk/programming/code/go/data/cs-server-controller/start-parameters.json"
	startParametersJsonPath := "..\\start-parameters.json"
	startParameters, err := getStarParameters(startParametersJsonPath, q)
	if err != nil {
		slog.Debug("failed to read json file. Using default start parameters", "path", startParametersJsonPath)
		startParameters = server.DefaultStartParameters()
	}

	lock.Lock()
	defer lock.Unlock()

	if err := s.Start(*startParameters); err != nil {
		http.Error(w, "failed to start server", http.StatusInternalServerError)
		return
	}

	startParametersJsonFileInstance, err := json_file.Get[server.StartParameters](startParametersJsonPath)
	if err == nil {
		if err := startParametersJsonFileInstance.Write(*startParameters); err != nil {
			slog.Warn("server started but failed to save valid start parameters to file. " + err.Error())
		}
	} else {
		slog.Warn("failed to get instance of start parameters json file")
	}
}

func getStarParameters(startParametersJsonPath string, query url.Values) (*server.StartParameters, error) {
	startParametersJsonFileInstance, err := json_file.Get[server.StartParameters](startParametersJsonPath)
	if err != nil {
		return nil, err
	}

	startParameters, err := startParametersJsonFileInstance.Read()
	if err != nil {
		slog.Debug("failed to read json file. Using default start parameters", "path", startParametersJsonFileInstance.GetPath())
		startParameters = server.DefaultStartParameters()
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
			return nil, fmt.Errorf("maxPlayers parameter must be a valid number between 1 and 255", http.StatusBadRequest)
		}

		startParameters.MaxPlayers = uint8(maxPlayersUint64)
	}

	if loginToken := strings.TrimSpace(query.Get("loginToken")); loginToken != "" {
		startParameters.SteamLoginToken = loginToken
	}

	return startParameters, nil
}
