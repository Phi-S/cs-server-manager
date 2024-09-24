package handlers

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/Phi-S/cs-server-manager/constants"
	"github.com/Phi-S/cs-server-manager/server"
	"github.com/Phi-S/cs-server-manager/steamcmd"
	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v3"
)

type ErrorResponse struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
}

func NewErrorValidation(c fiber.Ctx, err error) error {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return NewInternalServerErrorWithInternal(c, err)
	}

	if _, ok := err.(validator.ValidationErrors); ok {
		msg := ""
		for _, valErr := range err.(validator.ValidationErrors) {
			msg += fmt.Sprintf("Validation of '%v' failed on tag '%v'",
				valErr.Field(),
				valErr.ActualTag(),
			)

			if valErr.Param() != "" && valErr.ActualTag() != valErr.Param() {
				msg += fmt.Sprintf(" %v", valErr.Param())
			}
			msg += "!\n"
		}

		msg, _ = strings.CutSuffix(msg, "\n")
		return NewErrorWithInternal(c, fiber.StatusBadRequest, msg, err)
	}

	return NewInternalServerErrorWithInternal(c, err)
}

func NewErrorWithInternal(c fiber.Ctx, code int, message string, internalError error) error {
	c.Locals(constants.InternalErrorKey, internalError)
	return fiber.NewError(code, message)
}

func NewErrorWithMessage(c fiber.Ctx, code int, message string) error {
	c.Locals(constants.InternalErrorKey, errors.New(message))
	return fiber.NewError(code, message)
}

func NewInternalServerErrorWithInternal(c fiber.Ctx, err error) error {
	return NewErrorWithInternal(c, fiber.StatusInternalServerError, "internal error", err)
}

func GetFromLocals[T any](c fiber.Ctx, key any) (T, error) {
	var zeroResult T

	value := c.Locals(key)
	if value == nil {
		return zeroResult, fmt.Errorf("key %T not found in locals", key)
	}

	parsedValue, ok := value.(T)
	if !ok {
		return zeroResult, fmt.Errorf("failed to parse value in locals with key %T", key)
	}

	return parsedValue, nil
}

func GetServerSteamcmdInstances(c fiber.Ctx) (*sync.Mutex, *server.Instance, *steamcmd.Instance, error) {

	lock, err := GetFromLocals[*sync.Mutex](c, constants.ServerSteamcmdLockKey)
	if err != nil {
		return nil, nil, nil, err
	}

	serverInstance, err := GetFromLocals[*server.Instance](c, constants.ServerInstanceKey)
	if err != nil {
		return nil, nil, nil, err
	}

	steamcmdInstance, err := GetFromLocals[*steamcmd.Instance](c, constants.SteamCmdInstanceKey)
	if err != nil {
		return nil, nil, nil, err
	}

	return lock, serverInstance, steamcmdInstance, nil
}
