package handlers

import (
	"cs-server-manager/constants"
	"cs-server-manager/gvalidator"
	"cs-server-manager/jfile"
	"cs-server-manager/server"
	"github.com/gofiber/fiber/v3"
)

type StartParameters struct {
	Hostname        string `json:"hostname" validate:"required,lte=128"`
	Password        string `json:"password" validate:"omitempty,alphanum,lte=32"`
	StartMap        string `json:"start_map" validate:"required,printascii,lte=32"`
	MaxPlayers      uint8  `json:"max_players" validate:"required,number,lte=64"`
	SteamLoginToken string `json:"steam_login_token" validate:"omitempty,alphanum,gt=31,lte=32"`
}

var loginTokenVisibleCount = 4

func GetSettingsHandler(c fiber.Ctx) error {
	startParametersJsonFile, err := GetFromLocals[*jfile.Instance[server.StartParameters]](c, constants.StartParametersJsonFileKey)
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

	resp := StartParameters{
		Hostname:        sp.Hostname,
		Password:        sp.Password,
		StartMap:        sp.StartMap,
		MaxPlayers:      sp.MaxPlayers,
		SteamLoginToken: redactedSteamLoginToken,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func UpdateSettingsHandler(c fiber.Ctx) error {
	startParametersJsonFile, err := GetFromLocals[*jfile.Instance[server.StartParameters]](c, constants.StartParametersJsonFileKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	sp := new(StartParameters)
	if err := c.Bind().JSON(sp); err != nil {
		return NewErrorWithInternal(c, fiber.StatusBadRequest, "request body is not valid", err)
	}

	if err := gvalidator.Instance().Struct(sp); err != nil {
		return NewErrorWithInternal(c, fiber.StatusBadRequest, "request body is not valid", err)
	}

	var localStartParameters server.StartParameters
	err = startParametersJsonFile.Update(func(currentData *server.StartParameters) {
		if sp.Hostname != currentData.Hostname {
			currentData.Hostname = sp.Hostname
		}

		if sp.Password != currentData.Password {
			currentData.Password = sp.Password
		}

		if sp.StartMap != currentData.StartMap {
			currentData.StartMap = sp.StartMap
		}

		if sp.MaxPlayers != currentData.MaxPlayers {
			currentData.MaxPlayers = sp.MaxPlayers
		}

		if len(sp.SteamLoginToken) == 0 {
			currentData.SteamLoginToken = ""
		} else if len(sp.SteamLoginToken) > loginTokenVisibleCount {
			if len(currentData.SteamLoginToken) == 0 {
				currentData.SteamLoginToken = sp.SteamLoginToken
			} else if sp.SteamLoginToken[:loginTokenVisibleCount] != currentData.SteamLoginToken[:loginTokenVisibleCount] {
				currentData.SteamLoginToken = sp.SteamLoginToken
			}
		}

		localStartParameters = *currentData
	})

	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(localStartParameters)
}
