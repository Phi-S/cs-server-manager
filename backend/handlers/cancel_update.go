package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func CancelUpdateHandler(c fiber.Ctx) error {
	_, _, steamcmdInstance, err := GetServerSteamcmdInstances(c)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	if !steamcmdInstance.IsRunning() {
		return fiber.NewError(http.StatusInternalServerError, "no update is running")
	}

	if err := steamcmdInstance.Cancel(); err != nil {
		return NewErrorWithInternal(c, fiber.StatusInternalServerError, "failed to cancel server update", err)
	}

	return c.SendStatus(fiber.StatusOK)
}
