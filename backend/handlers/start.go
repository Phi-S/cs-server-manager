package handlers

import (
	"cs-server-manager/constants"
	"cs-server-manager/gvalidator"
	"cs-server-manager/start_parameters_json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type StartBody struct {
	Hostname        string `json:"hostname" validate:"omitempty,lte=128"`
	Password        string `json:"password" validate:"omitempty,alphanum,lte=32"`
	StartMap        string `json:"start_map" validate:"omitempty,printascii,lte=32"`
	MaxPlayers      uint8  `json:"max_players" validate:"omitempty,number,lte=128"`
	SteamLoginToken string `json:"steam_login_token" validate:"omitempty,alphanum,eq=32"`
}

// StartHandler
// @Summary      Start the server
// @Description	 Starts the server with the given start parameters
// @Tags         server
// @Accept       json
// @Param 		 startParameters body StartBody false "You can provide no, all or only a few start parameters. The provided start parameters will overwrite the saved start parameters in the start-parameters.json file if the server started successfully."
// @Success      200
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /start [post]
func StartHandler(c fiber.Ctx) error {
	lock, serverInstance, steamcmd, err := GetServerSteamcmdInstances(c)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	lock.Lock()
	defer lock.Unlock()

	if steamcmd.IsRunning() {
		return fiber.NewError(fiber.StatusInternalServerError, "can not start server while server is updating")
	}

	if serverInstance.IsRunning() {
		return fiber.NewError(fiber.StatusInternalServerError, "server is already running")
	}

	startParameterJsonFile, err := GetFromLocals[*start_parameters_json.Instance](c, constants.StartParametersJsonFileKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	startParameters, err := startParameterJsonFile.Read()
	if err != nil {
		return NewInternalServerErrorWithInternal(c, fmt.Errorf("startParameterJsonFile.Read(): %w", err))
	}

	if len(c.Body()) > 0 {
		startBody := new(StartBody)
		if err := c.Bind().JSON(startBody); err != nil {
			return NewErrorWithInternal(c, fiber.StatusBadRequest, "start parameters are not not valid", fmt.Errorf("c.Bind().JSON(startBody): %w", err))
		}

		if err := gvalidator.Instance().Struct(startBody); err != nil {
			return NewErrorWithInternal(c, fiber.StatusBadRequest, "start parameters are not not valid", fmt.Errorf("gvalidator.Instance().Struct(startBody): %w", err))
		}

		if hostname := strings.TrimSpace(startBody.Hostname); hostname != "" {
			startParameters.Hostname = hostname
		}

		if password := strings.TrimSpace(startBody.Password); password != "" {
			startParameters.Password = password
		}

		if startMap := strings.TrimSpace(startBody.StartMap); startMap != "" {
			startParameters.StartMap = startMap
		}

		if maxPlayers := startBody.MaxPlayers; maxPlayers != 0 {
			startParameters.MaxPlayers = maxPlayers
		}

		if loginToken := strings.TrimSpace(startBody.SteamLoginToken); loginToken != "" {
			startParameters.SteamLoginToken = loginToken
		}
	}

	if err := serverInstance.Start(startParameters); err != nil {
		return NewErrorWithInternal(c, fiber.StatusInternalServerError, "failed to start server", err)
	}

	if err := startParameterJsonFile.Write(startParameters); err != nil {
		slog.Warn("server started but failed to save valid start parameters to file. " + err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
