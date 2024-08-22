package handlers

import (
	"cs-server-manager/constants"
	"cs-server-manager/status"

	"github.com/gofiber/fiber/v3"
)

// StatusHandler Status
// @Summary      Get the current status of the server
// @Tags         server
// @Produce      json
// @Success      200  {object}  status.InternalStatus
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /status [get]
func StatusHandler(c fiber.Ctx) error {
	statusInstance, err := GetFromLocals[*status.Status](c, constants.StatusKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(statusInstance.Status())
}
