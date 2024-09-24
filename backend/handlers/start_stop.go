package handlers

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/Phi-S/cs-server-manager/constants"
	"github.com/Phi-S/cs-server-manager/gvalidator"
	"github.com/Phi-S/cs-server-manager/start_parameters_json"
	"github.com/Phi-S/cs-server-manager/status"

	"github.com/gofiber/fiber/v3"
)

type StartBody struct {
	Hostname        string `json:"hostname" validate:"omitempty,lte=128"`
	Password        string `json:"password" validate:"omitempty,alphanum,lte=32"`
	StartMap        string `json:"start_map" validate:"omitempty,printascii,lte=32"`
	MaxPlayers      uint8  `json:"max_players" validate:"omitempty,number,lte=128"`
	SteamLoginToken string `json:"steam_login_token" validate:"omitempty,alphanum,len=32"`
}

func RegisterStartStop(r fiber.Router) {
	r.Post("/start", startHandler)
	r.Post("/stop", stopHandler)
}

// @Summary      Start the server
// @Description	 Starts the server with the given start parameters
// @Tags         server
// @Accept       json
// @Param 		 startParameters body StartBody false "You can provide no, all or only a few start parameters. The provided start parameters will overwrite the saved start parameters in the start-parameters.json file if the server started successfully."
// @Success      200
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /start [post]
func startHandler(c fiber.Ctx) error {
	lock, serverInstance, steamcmd, err := GetServerSteamcmdInstances(c)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	status, err := GetFromLocals[*status.Status](c, constants.StatusKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	if !status.Status().IsGameServerInstalled {
		return NewErrorWithMessage(c, fiber.StatusInternalServerError, "Game server is not yet installed")
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
			return NewErrorValidation(c, err)
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

// @Summary 	Stop the server
// @Description Stops the server of if the server is not running, returns 200 OK
// @Tags        server
// @Success     200
// @Failure     400  {object}  handlers.ErrorResponse
// @Failure     500  {object}  handlers.ErrorResponse
// @Router      /stop [post]
func stopHandler(c fiber.Ctx) error {
	lock, server, _, err := GetServerSteamcmdInstances(c)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	lock.Lock()
	defer lock.Unlock()

	if !server.IsRunning() {
		return c.SendStatus(fiber.StatusOK)
	}

	if err := server.Stop(); err != nil {
		return NewErrorWithInternal(c, fiber.StatusInternalServerError, "failed to stop server", err)
	}

	return c.SendStatus(fiber.StatusOK)
}
