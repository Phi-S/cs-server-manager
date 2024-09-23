package main

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/Phi-S/cs-server-manager/constants"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
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
			"duration-ms", float64(duration.Nanoseconds())/1e6,
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
			"duration-ms", float64(duration.Nanoseconds())/1e6,
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
