package main

import (
	"cs-server-manager/constants"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"log/slog"
	"runtime/debug"
	"time"
)

func logMiddleware(c fiber.Ctx) error {
	startTime := time.Now().UTC()
	err := c.Next()

	duration := time.Now().UTC().Sub(startTime)
	internalError, _ := c.Locals(constants.InternalErrorKey).(error)

	statusCode := fiber.StatusInternalServerError
	responseMessage := ""

	var e *fiber.Error
	if errors.As(err, &e) {
		statusCode = e.Code
		responseMessage = e.Message
	}

	if err == nil {
		slog.Info("request finished",
			"request-id", requestid.FromContext(c),
			"method", c.Method(),
			"path", c.Path(),
			"query", c.Request().URI().QueryString(),
			"ip", c.IP(),
			"port", c.Port(),
			"status", c.Response().StatusCode(),
			"duration", duration,
		)
	} else {
		slog.Error("request finished with error",
			"request-id", requestid.FromContext(c),
			"method", c.Method(),
			"path", c.Path(),
			"query", c.Request().URI().QueryString(),
			"ip", c.IP(),
			"port", c.Port(),
			"status", statusCode,
			"duration", duration,
			"response-message", responseMessage,
			"internal-error", internalError,
			"error", err,
		)
	}
	return err
}

func panicHandler(c fiber.Ctx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("handler paniced: %v | %s", r, debug.Stack())
			}
		}
	}()

	return c.Next()
}