package main

import (
	"cs-server-manager/config"
	"cs-server-manager/constants"
	"cs-server-manager/gvalidator"
	"cs-server-manager/handlers"
	"cs-server-manager/logwrt"
	"cs-server-manager/plugins"
	"cs-server-manager/server"
	"cs-server-manager/start_parameters_json"
	"cs-server-manager/status"
	"cs-server-manager/steamcmd"
	"embed"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "cs-server-manager/docs"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/static"
	"golang.org/x/net/websocket"
)

//go:embed swagger-ui
//go:embed docs
//go:embed web
var dir embed.FS

// @title cs-server-manager API
// @version 1.0
// @schemes http https
// @BasePath /api/v1
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
	startParametersJsonFile *start_parameters_json.Instance,
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

	api := app.Group("/api")

	v1 := api.Group("/v1", func(c fiber.Ctx) error {
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
	v1.Post("/command", handlers.SendCommandHandler)

	v1.Post("/update", handlers.UpdateHandler)
	v1.Post("/update/cancel", handlers.CancelUpdateHandler)

	v1.Get("/settings", handlers.GetSettingsHandler)
	v1.Post("/settings", handlers.UpdateSettingsHandler)

	v1.Get("/plugins", handlers.GetPluginsHandler)
	v1.Post("/plugins", handlers.InstallPluginHandler)
	v1.Delete("/plugins", handlers.UninstallPluginHandler)

	logGroup := v1.Group("/logs", func(c fiber.Ctx) error {
		c.Locals(constants.UserLogWriterKey, userLogWriter)
		return c.Next()
	})

	logGroup.Get("/:count", handlers.LogsHandler)

	v1.Get("/ws", adaptor.HTTPHandler(websocket.Handler(webSocketServer.handleWs)))

	if config.EnableSwagger {
		swagger := api.Group("swagger")

		swagger.Get("/swagger.json", static.New("docs/swagger.json", static.Config{
			FS:     dir,
			Browse: false,
		}))

		if err := mapDir(swagger, "", dir, "swagger-ui"); err != nil {
			log.Fatal("failed to map swagger-ui dir", err)
		}
	}

	if config.EnableWebUi {
		if err := mapDir(app, "", dir, "web"); err != nil {
			log.Fatal("failed to map web dir: ", err)
		}
	}

	log.Fatal(app.Listen(":" + config.HttpPort))
}

func mapDir(router fiber.Router, path string, fs embed.FS, dir string) error {
	router.Get("", static.New(fmt.Sprintf("%v/index.html", dir), static.Config{
		FS:     fs,
		Browse: false,
	}))

	content, err := fs.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read folder content: %w", err)
	}

	for _, entry := range content {
		epath := fmt.Sprintf("%v/%v", path, entry.Name())
		efPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			err := mapDir(router, epath, fs, efPath)
			if err != nil {
				return fmt.Errorf("failed to map sub dir '%v': %w", entry, err)
			}
		} else {
			router.Get(epath, static.New(efPath, static.Config{
				FS:     fs,
				Browse: true,
			}))
		}
	}

	return nil
}
