package main

import (
	"cs-server-manager/config"
	"cs-server-manager/constants"
	"cs-server-manager/gvalidator"
	"cs-server-manager/handlers"
	"cs-server-manager/jfile"
	"cs-server-manager/logwrt"
	"cs-server-manager/plugins"
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
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/static"
	"golang.org/x/net/websocket"
)

func main() {
	configureLogger()

	if err := gvalidator.RegisterCustomTags(); err != nil {
		slog.Error("FATAL: gvalidator: failed to register custom tags", "error", err)
		os.Exit(1)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error("FATAL: failed to get config", "error", err)
		os.Exit(1)
	}

	if err := createdRequiredDirs(cfg); err != nil {
		slog.Error("FATAL: failed to create required directories", "error", err)
		os.Exit(1)
	}

	steamcmdInstance,
		serverInstance,
		startParametersJsonFileHandler,
		userLogWriter,
		statusInstance,
		webSocketServerInstance,
		gameEventsInstance,
		pluginsInstance,
		err := createRequiredServices(cfg)
	if err != nil {
		slog.Error("FATAL: failed to create required services", "error", err)
		os.Exit(1)
	}

	// linking up all services via events
	registerEvents(
		cfg,
		serverInstance,
		steamcmdInstance,
		startParametersJsonFileHandler,
		userLogWriter,
		statusInstance,
		webSocketServerInstance,
		gameEventsInstance,
		pluginsInstance,
	)

	defer func() {
		_ = steamcmdInstance.Cancel()
		steamcmdInstance.Close()

		_ = serverInstance.Stop()
		serverInstance.Close()

		userLogWriter.Close()
	}()

	// this lock is used to prevent collision between the server and steamcmd instance
	// Fox example the lock is used to prevent the server from being started while a steamcmd updated is getting started at the same time.
	// This can occur if two http request are coming in at the same time and the internal status of the steamcmd and/or server instances is not yet updated
	ServerSteamcmdLock := sync.Mutex{}

	startApi(
		cfg,
		&ServerSteamcmdLock,
		serverInstance,
		steamcmdInstance,
		startParametersJsonFileHandler,
		statusInstance,
		userLogWriter,
		webSocketServerInstance,
		pluginsInstance,
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

func startApi(
	config config.Config,
	ServerSteamcmdLock *sync.Mutex,
	serverInstance *server.Instance,
	steamcmdInstance *steamcmd.Instance,
	startParametersJsonFile *jfile.Instance[server.StartParameters],
	status *status.Status,
	userLogWriter *logwrt.LogWriter,
	webSocketServer *WebSocketServer,
	pluginsInstance *plugins.Instance,
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

	app.Use(cors.New())
	app.Use(requestid.New())
	app.Use(logMiddleware)
	app.Use(panicHandler)

	v1 := app.Group("/api/v1", func(c fiber.Ctx) error {
		c.Locals(constants.ConfigKey, config)
		c.Locals(constants.ServerSteamcmdLockKey, ServerSteamcmdLock)
		c.Locals(constants.ServerInstanceKey, serverInstance)
		c.Locals(constants.SteamCmdInstanceKey, steamcmdInstance)
		c.Locals(constants.StartParametersJsonFileKey, startParametersJsonFile)
		c.Locals(constants.StatusKey, status)
		c.Locals(constants.PluginsKey, pluginsInstance)
		return c.Next()
	})

	v1.Get("/status", handlers.StatusHandler)

	v1.Post("/start", handlers.StartHandler)
	v1.Post("/stop", handlers.StopHandler)
	v1.Post("/send-command", handlers.SendCommandHandler)

	v1.Post("/update", handlers.UpdateHandler)
	v1.Post("/update/cancel", handlers.CancelUpdateHandler)

	v1.Get("/settings", handlers.GetSettingsHandler)
	v1.Post("/settings", handlers.UpdateSettingsHandler)

	v1.Get("/plugins", handlers.GetAvailablePluginsHandler)
	v1.Post("/plugins/install", handlers.InstallPluginHandler)
	v1.Post("/plugins/uninstall", handlers.UninstallPluginHandler)

	logGroup := v1.Group("/log", func(c fiber.Ctx) error {
		c.Locals(constants.UserLogWriterKey, userLogWriter)
		return c.Next()
	})

	logGroup.Get("/:countOrSince", handlers.LogsHandler)

	v1.Get("/ws", adaptor.HTTPHandler(websocket.Handler(webSocketServer.handleWs)))

	if !config.DisableWebsite {
		app.Get("/*", static.New("./dist"))
	}

	log.Fatal(app.Listen(":" + config.HttpPort))
}
