package main

import (
	"cs-server-manager/config"
	"cs-server-manager/constants"
	"cs-server-manager/event"
	globalvalidator "cs-server-manager/global_validator"
	"cs-server-manager/handlers"
	jsonfile "cs-server-manager/jsonfile"
	"cs-server-manager/logwrt"
	"cs-server-manager/server"
	"cs-server-manager/status"
	"cs-server-manager/steamcmd"
	"errors"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/static"
	"golang.org/x/net/websocket"
)

func main() {
	configureLogger()
	globalvalidator.Init()

	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// this will kill the application if unsuccessful
	createdRequiredDirs(cfg)

	// this lock is used to prevent collision between the server and steamcmd instance
	// Fox example the lock is used to prevent the server from being started while an steamcmd updated is getting started at the same time.
	// This can occur if two http request are coming in at the same time and the internal status of the steamcmd and/or server instances is not yet updated
	ServerSteamcmdLock := sync.Mutex{}

	steamcmdInstance,
		serverInstance,
		startParametersJsonFileHandler,
		userLogWriter,
		statusInstance,
		webSocketServer,
		gameEventsInstance := createRequiredServices(cfg)

	serverInstance.OnOutput(func(p event.PayloadWithData[string]) {
		gameEventsInstance.DetectGameEvent(p.Data)
	})

	statusEventHandler(statusInstance, serverInstance, steamcmdInstance, gameEventsInstance)
	statusInstance.OnStatusChanged(func(p event.PayloadWithData[status.InternalStatus]) {
		if err := webSocketServer.Broadcast("status", p.Data); err != nil {
			slog.Error("failed to send status message", "status", p.Data, "error", err)
		}
	})

	defer func() {
		_ = steamcmdInstance.Cancel()
		steamcmdInstance.Close()

		_ = serverInstance.Stop()
		serverInstance.Close()

		userLogWriter.Close()
	}()

	logEvents(userLogWriter, webSocketServer, serverInstance, steamcmdInstance, gameEventsInstance)

	////////////////////
	StartApi(
		cfg,
		&ServerSteamcmdLock,
		serverInstance,
		steamcmdInstance,
		startParametersJsonFileHandler,
		statusInstance,
		userLogWriter,
		webSocketServer,
	)
}

func configureLogger() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(time.Now().UTC().Format(time.RFC3339Nano))
			}
			return a
		},
	})

	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}

func StartApi(
	config config.Config,
	ServerSteamcmdLock *sync.Mutex,
	serverInstance *server.Instance,
	steamcmdInstance *steamcmd.Instance,
	startParametersJsonFile *jsonfile.JsonFile[server.StartParameters],
	status *status.Status,
	userLogWriter *logwrt.LogWriter,
	webSocketServer *WebSocketServer,
) {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			requestId := requestid.FromContext(c)
			if requestId == "" {
				slog.Warn("while handling error the request id was not set")
			}

			code := fiber.StatusInternalServerError
			msg := "unknown error"

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
				msg = e.Message
			}

			resp := struct {
				Status    int    `json:"status"`
				Message   string `json:"message"`
				RequestId string `json:"request_id"`
			}{
				Status:    code,
				Message:   msg,
				RequestId: requestId,
			}

			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return c.Status(code).JSON(resp)
		},
	})

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(requestid.New())
	app.Use(func(c fiber.Ctx) error {
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
				"path", c.Path(),
				"query", c.Request().URI().QueryString(),
				"ip", c.IP(),
				"port", c.Port(),
				"status", statusCode,
				"duration", duration,
				"response-message", responseMessage,
				"internal-error", internalError,
			)
		}
		return err
	})

	v1 := app.Group("/api/v1", func(c fiber.Ctx) error {
		c.Locals(constants.ConfigKey, config)
		c.Locals(constants.ServerSteamcmdLockKey, ServerSteamcmdLock)
		c.Locals(constants.ServerInstanceKey, serverInstance)
		c.Locals(constants.SteamCmdInstanceKey, steamcmdInstance)
		c.Locals(constants.StartParametersJsonFileKey, startParametersJsonFile)
		c.Locals(constants.StatusKey, status)
		return c.Next()
	})

	v1.Get("/status", handlers.StatusHandler)

	v1.Post("/start", handlers.StartHandler)
	v1.Post("/stop", handlers.StopHandler)
	v1.Post("/send-command", handlers.SendCommandHandler)

	v1.Post("/update", handlers.UpdateHandler)
	v1.Post("/update/cancel", handlers.CancelUpdateHandler)

	logGroup := v1.Group("/log", func(c fiber.Ctx) error {
		c.Locals(constants.UserLogWriterKey, userLogWriter)
		return c.Next()
	})

	logGroup.Get("/:countOrSince", handlers.LogsHandler)

	v1.Get("/ws", adaptor.HTTPHandler(websocket.Handler(webSocketServer.handleWs)))

	app.Get("/*", static.New("./dist"))

	log.Fatal(app.Listen(":" + config.HttpPort))
}
