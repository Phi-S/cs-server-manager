package handlers

import (
	"cs-server-manager/constants"
	"cs-server-manager/server"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func SendCommandHandler(c fiber.Ctx) error {

	serverInstance, err := GetFromLocals[*server.Instance](c, constants.ServerInstanceKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	if !serverInstance.IsRunning() {
		return fiber.NewError(fiber.StatusInternalServerError, "server is not running")
	}

	command := strings.TrimSpace(c.Query("command"))
	if command == "" {
		return fiber.NewError(fiber.StatusBadRequest, "command is empty")
	}

	out, err := serverInstance.SendCommand(command)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	resp := struct {
		Output []string `json:"output"`
	}{
		Output: strings.Split(out, "\n"),
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
