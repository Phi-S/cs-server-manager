package handlers

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterUpdate(r fiber.Router) {
	r.Post("/update", startUpdateHandler)
	r.Post("/update/cancel", cancelUpdateHandler)
}

// @Summary				Start server update
// @Tags         		update
// @Success     		200
// @Failure				400  {object}  handlers.ErrorResponse
// @Failure				500  {object}  handlers.ErrorResponse
// @Router       		/update [post]
func startUpdateHandler(c fiber.Ctx) error {
	lock, serverInstance, steamcmdInstance, err := GetServerSteamcmdInstances(c)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	lock.Lock()
	defer lock.Unlock()

	if steamcmdInstance.IsRunning() {
		return fiber.NewError(fiber.StatusInternalServerError, "another update is already running")
	}

	if serverInstance.IsRunning() {
		return fiber.NewError(fiber.StatusInternalServerError, "can not update server. Server is still running")
	}

	if err := steamcmdInstance.Update(false); err != nil {
		return NewErrorWithInternal(c, fiber.StatusInternalServerError, "failed to update server", err)
	}

	return c.SendStatus(fiber.StatusAccepted)
}

// @Summary				Cancel the server update
// @Description 		Cancel the currently running server update or if no update is currently running, returns 200 OK
// @Tags         		update
// @Success     		200
// @Failure				400  {object}  handlers.ErrorResponse
// @Failure				500  {object}  handlers.ErrorResponse
// @Router       		/update/cancel [post]
func cancelUpdateHandler(c fiber.Ctx) error {
	_, _, steamcmdInstance, err := GetServerSteamcmdInstances(c)
	if err != nil {
		return NewInternalServerErrorWithInternal(c, err)
	}

	if !steamcmdInstance.IsRunning() {
		return c.SendStatus(fiber.StatusOK)
	}

	if err := steamcmdInstance.Cancel(); err != nil {
		return NewErrorWithInternal(c, fiber.StatusInternalServerError, "failed to cancel server update", err)
	}

	return c.SendStatus(fiber.StatusOK)
}
