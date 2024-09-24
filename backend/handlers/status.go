package handlers

import (
	"github.com/Phi-S/cs-server-manager/constants"
	"github.com/Phi-S/cs-server-manager/status"

	"github.com/gofiber/fiber/v3"
)

func RegisterStatus(r fiber.Router) {
	r.Get("/status", statusHandler)
}

// @Summary      Get the current status of the server
// @Tags         server
// @Produce      json
// @Success      200  {object}  status.InternalStatus
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /status [get]
func statusHandler(c fiber.Ctx) error {
	statusInstance, err := GetFromLocals[*status.Status](c, constants.StatusKey)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(statusInstance.Status())
}
