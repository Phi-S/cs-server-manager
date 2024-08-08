package handlers

import (
	"cs-server-manager/constants"
	"cs-server-manager/jfile"
	"cs-server-manager/server"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func StartHandler(c fiber.Ctx) error {
	lock, serverInstance, steamcmd, err := GetServerSteamcmdInstances(c)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	lock.Lock()
	defer lock.Unlock()

	if steamcmd.IsRunning() {
		return fiber.NewError(fiber.StatusInternalServerError, "can not start server while steamcmd is running")
	}

	if serverInstance.IsRunning() {
		return fiber.NewError(fiber.StatusInternalServerError, "server is already running")
	}

	startParameterJsonFile, err := GetFromLocals[*jfile.Instance[server.StartParameters]](c, constants.StartParametersJsonFileKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	startParameters, err := getStarParameters(startParameterJsonFile, c.Queries())
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	if err := serverInstance.Start(*startParameters); err != nil {
		return NewErrorWithInternal(c, fiber.StatusInternalServerError, "failed to start server", err)
	}

	if err := startParameterJsonFile.Write(*startParameters); err != nil {
		slog.Warn("server started but failed to save valid start parameters to file. " + err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func getStarParameters(startParametersJsonFileInstance *jfile.Instance[server.StartParameters], query map[string]string) (*server.StartParameters, error) {
	startParameters, err := startParametersJsonFileInstance.Read()
	if err != nil {
		return nil, err
	}

	if hostname := strings.TrimSpace(query["name"]); hostname != "" {
		startParameters.Hostname = hostname
	}

	if password := strings.TrimSpace(query["pw"]); password != "" {
		startParameters.Password = password
	}

	if startMap := strings.TrimSpace(query["map"]); startMap != "" {
		startParameters.StartMap = startMap
	}

	if maxPlayersString := strings.TrimSpace(query["max_player_count"]); maxPlayersString != "" {
		maxPlayersUint64, err := strconv.ParseUint(maxPlayersString, 10, 8)
		if err != nil {
			return nil, err
		}

		if maxPlayersUint64 > 255 {
			return nil, fmt.Errorf("maxPlayers parameter must be a valid number between 1 and 255")
		}

		startParameters.MaxPlayers = uint8(maxPlayersUint64)
	}

	if loginToken := strings.TrimSpace(query["loginToken"]); loginToken != "" {
		startParameters.SteamLoginToken = loginToken
	}

	return startParameters, nil
}
