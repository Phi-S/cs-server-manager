package handlers

import (
	"cs-server-manager/constants"
	"cs-server-manager/server"
	"cs-server-manager/steamcmd"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v3"
)

type InternalError struct {
}

func NewErrorWithInternal(c fiber.Ctx, code int, message string, internalError error) error {
	c.Locals(constants.InternalErrorKey, internalError)
	return fiber.NewError(code, message)
}

func NewInternalServerErrorWithInternal(c fiber.Ctx, err error) error {
	return NewErrorWithInternal(c, fiber.StatusInternalServerError, "internal error", err)
}

func GetFromLocals[T any](c fiber.Ctx, key any) (T, error) {
	var zeroResult T

	value, ok := c.Locals(key).(T)
	if !ok {
		return zeroResult, fmt.Errorf("failed to get %T from locals", key)
	}

	return value, nil
}

func GetServerSteamcmdInstances(c fiber.Ctx) (*sync.Mutex, *server.Instance, *steamcmd.Instance, error) {

	lock, err := GetFromLocals[*sync.Mutex](c, constants.ServerSteamcmdLockKey)
	if err != nil {
		return nil, nil, nil, err
	}

	server, err := GetFromLocals[*server.Instance](c, constants.ServerInstanceKey)
	if err != nil {
		return nil, nil, nil, err
	}

	steamcmd, err := GetFromLocals[*steamcmd.Instance](c, constants.SteamCmdInstanceKey)
	if err != nil {
		return nil, nil, nil, err
	}

	return lock, server, steamcmd, nil
}
