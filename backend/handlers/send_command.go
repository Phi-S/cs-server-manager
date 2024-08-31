package handlers

import (
	"cs-server-manager/constants"
	"cs-server-manager/server"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type CommandRequest struct {
	Command string `json:"command" validate:"required,lt=128"`
}

type CommandResponse struct {
	Output []string `json:"output"`
}

// SendCommandHandler
// @Summary				Send game-server command
// @Tags         		server
// @Accept       		json
// @Produce 			json
// @Param		 		command body CommandRequest true "This command will be executed on the game server"
// @Success     		200  {object}  CommandResponse
// @Failure				400  {object}  handlers.ErrorResponse
// @Failure				500  {object}  handlers.ErrorResponse
// @Router       		/send-command [post]
func SendCommandHandler(c fiber.Ctx) error {
	serverInstance, err := GetFromLocals[*server.Instance](c, constants.ServerInstanceKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	if !serverInstance.IsRunning() {
		return fiber.NewError(fiber.StatusInternalServerError, "server is not running")
	}

	var commandRequest CommandRequest
	if err := c.Bind().JSON(&commandRequest); err != nil {
		return NewErrorWithInternal(c, fiber.StatusBadRequest, "request is not valid", err)
	}
	if commandRequest.Command == "" {
		return fiber.NewError(fiber.StatusBadRequest, "command is empty")
	}

	out, err := serverInstance.SendCommand(commandRequest.Command)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	resp := CommandResponse{
		Output: strings.Split(out, "\n"),
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
