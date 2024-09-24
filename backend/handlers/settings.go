package handlers

import (
	"fmt"

	"github.com/Phi-S/cs-server-manager/constants"
	"github.com/Phi-S/cs-server-manager/gvalidator"
	"github.com/Phi-S/cs-server-manager/start_parameters_json"

	"github.com/gofiber/fiber/v3"
)

type SettingsModel struct {
	Hostname        string `json:"hostname" validate:"required,lte=128"`
	Password        string `json:"password" validate:"omitempty,alphanum,lte=32"`
	StartMap        string `json:"start_map" validate:"required,printascii,lte=32"`
	MaxPlayers      uint8  `json:"max_players" validate:"required,number,lte=128"`
	SteamLoginToken string `json:"steam_login_token" validate:"omitempty,alphanum,len=32"`
}

var loginTokenVisibleCount = 4

func RegisterSettings(r fiber.Router) {
	r.Get("/settings", getSettingsHandler)
	r.Post("/settings", updateSettingsHandler)
}

// @Summary				Get the current settings
// @Tags         		settings
// @Produce      		json
// @Success     		200  {object}  SettingsModel
// @Failure				400  {object}  handlers.ErrorResponse
// @Failure				500  {object}  handlers.ErrorResponse
// @Router       		/settings [get]
func getSettingsHandler(c fiber.Ctx) error {
	startParametersJsonFile, err := GetFromLocals[*start_parameters_json.Instance](c, constants.StartParametersJsonFileKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	sp, err := startParametersJsonFile.Read()
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	redactedSteamLoginToken := ""
	for i, v := range sp.SteamLoginToken {
		if i <= loginTokenVisibleCount {
			redactedSteamLoginToken += string(v)
		} else {
			redactedSteamLoginToken += "X"
		}
	}

	resp := SettingsModel{
		Hostname:        sp.Hostname,
		Password:        sp.Password,
		StartMap:        sp.StartMap,
		MaxPlayers:      sp.MaxPlayers,
		SteamLoginToken: redactedSteamLoginToken,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// @Summary					Update settings
// @Tags         			settings
// @Accept       			json
// @Produce      			json
// @Param		 			settings body SettingsModel true "The updated settings"
// @Success     			200  {object}  SettingsModel
// @Failure					400  {object}  handlers.ErrorResponse
// @Failure					500  {object}  handlers.ErrorResponse
// @Router       			/settings [post]
func updateSettingsHandler(c fiber.Ctx) error {
	startParametersJsonFile, err := GetFromLocals[*start_parameters_json.Instance](c, constants.StartParametersJsonFileKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	sp := new(SettingsModel)
	if err := c.Bind().JSON(sp); err != nil {
		return NewErrorWithInternal(c, fiber.StatusBadRequest, "request body is not valid", err)
	}

	if err := gvalidator.Instance().Struct(sp); err != nil {
		return NewErrorValidation(c, err)
	}

	startParameters, err := startParametersJsonFile.Read()
	if err != nil {
		return NewInternalServerErrorWithInternal(c, fmt.Errorf("startParametersJsonFile.Read(): %w", err))
	}

	if sp.Hostname != startParameters.Hostname {
		startParameters.Hostname = sp.Hostname
	}

	if sp.Password != startParameters.Password {
		startParameters.Password = sp.Password
	}

	if sp.StartMap != startParameters.StartMap {
		startParameters.StartMap = sp.StartMap
	}

	if sp.MaxPlayers != startParameters.MaxPlayers {
		startParameters.MaxPlayers = sp.MaxPlayers
	}

	if len(sp.SteamLoginToken) == 0 {
		startParameters.SteamLoginToken = ""
	} else if len(sp.SteamLoginToken) > loginTokenVisibleCount {
		if len(startParameters.SteamLoginToken) == 0 {
			startParameters.SteamLoginToken = sp.SteamLoginToken
		} else if sp.SteamLoginToken[:loginTokenVisibleCount] != startParameters.SteamLoginToken[:loginTokenVisibleCount] {
			startParameters.SteamLoginToken = sp.SteamLoginToken
		}
	}

	if err := startParametersJsonFile.Write(startParameters); err != nil {
		return NewInternalServerErrorWithInternal(c, fmt.Errorf("startParametersJsonFile.Write: %w", err))
	}

	return c.Status(fiber.StatusOK).JSON(startParameters)
}
