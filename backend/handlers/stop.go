package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func StopHandler(c fiber.Ctx) error {
	lock, server, _, err := GetServerSteamcmdInstances(c)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	lock.Lock()
	defer lock.Unlock()

	if !server.IsRunning() {
		return fiber.NewError(http.StatusInternalServerError, "server is not running")
	}

	if err := server.Stop(); err != nil {
		return NewErrorWithInternal(c, fiber.StatusInternalServerError, "failed to stop server", err)
	}

	return c.SendStatus(fiber.StatusOK)
}
