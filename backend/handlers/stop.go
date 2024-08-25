package handlers

import (
	"github.com/gofiber/fiber/v3"
)

// StopHandler
// @Summary 	Stop the server
// @Description Stops the server of if the server is not running, returns 200 OK
// @Tags        server
// @Success     200
// @Failure     400  {object}  handlers.ErrorResponse
// @Failure     500  {object}  handlers.ErrorResponse
// @Router      /stop [post]
func StopHandler(c fiber.Ctx) error {
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
