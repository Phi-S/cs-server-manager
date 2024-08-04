package handlers

import (
	"cs-server-manager/constants"
	"cs-server-manager/status"

	"github.com/gofiber/fiber/v3"
)

type StatusResponse struct {
	ServerRunning   bool `json:"server-running"`
	SteamCmdRunning bool `json:"steamcmd-running"`
}

func StatusHandler(c fiber.Ctx) error {
	status, err := GetFromLocals[*status.Status](c, constants.StatusKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(status.Status())
}
