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
	"path/filepath"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/requestid"
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

	if err := os.MkdirAll(cfg.DataDir, os.ModePerm); err != nil {
		slog.Error("failed to create data directory", "dir", cfg.DataDir, "error", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(cfg.ServerDir, os.ModePerm); err != nil {
		slog.Error("failed to create server directory", "dir", cfg.ServerDir, "error", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(cfg.SteamcmdDir, os.ModePerm); err != nil {
		slog.Error("failed to create steamcmd directory", "dir", cfg.SteamcmdDir, "error", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(cfg.LogDir, os.ModePerm); err != nil {
		slog.Error("failed to create log directory", "dir", cfg.LogDir, "error", err)
		os.Exit(1)
	}

	// this lock is used to prevent collision between the server and steamcmd instance
	// Fox example the lock is used to prevent the server from being started while being updated.
	// This can occur if two http request are coming in at the same time
	ServerSteamcmdLock := sync.Mutex{}

	steamcmdInstance, err := steamcmd.NewInstance(cfg.SteamcmdDir, cfg.ServerDir)
	if err != nil {
		slog.Error("failed to create new steamcmd instance", "error", err)
		os.Exit(1)
	}
	defer func() {
		_ = steamcmdInstance.Cancel()
		steamcmdInstance.Close()
	}()

	serverInstance, err := server.NewInstance(cfg.ServerDir, cfg.CsPort, cfg.SteamcmdDir)
	if err != nil {
		slog.Error("failed to create new server instance", "error", err)
		os.Exit(1)
	}
	defer func() {
		_ = serverInstance.Stop()
		serverInstance.Close()
	}()

	startParametersJsonFile, err := jsonfile.New[server.StartParameters](filepath.Join(cfg.DataDir, "start-parameters.json"), *server.DefaultStartParameters())
	if err != nil {
		slog.Error("failed to create new json file service for start-parameter.json", "error", err)
		os.Exit(1)
	}

	startParameters, err := startParametersJsonFile.Read()
	if err != nil {
		slog.Error("failed to read start-parameter.json", "error", err)
		os.Exit(1)
	}

	logDir := filepath.Join(cfg.DataDir, "logs")
	userLogWriter, err := logwrt.NewLogWriter(logDir, "user-logs")
	if err != nil {
		slog.Error("failed to create user logs write", "error", err)
		os.Exit(1)
	}
	defer userLogWriter.Close()

	status := status.NewStatus(startParameters.Hostname, startParameters.MaxPlayers, startParameters.StartMap)
	status.ChangeStatusOnEvents(serverInstance, steamcmdInstance)

	webSocketServer := NewWebSocketServer()
	status.OnStatusChanged(func(payload event.DefaultPayload) {
		currentStatus := status.Status()
		if err := webSocketServer.Broadcast("status", currentStatus); err != nil {
			slog.Error("failed to send status message", "status", currentStatus, "error", err)
		}
	})

	writeEventToLogFileAndWebSocket(userLogWriter, webSocketServer, serverInstance, steamcmdInstance)
	enableEventLogging(serverInstance, steamcmdInstance)

	////////////////////
	StartApi(
		cfg,
		&ServerSteamcmdLock,
		serverInstance,
		steamcmdInstance,
		startParametersJsonFile,
		status,
		userLogWriter,
		webSocketServer,
	)
}

func configureLogger() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
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
				RequestId string `json:"request-id"`
			}{
				Status:    code,
				Message:   msg,
				RequestId: requestId,
			}

			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return c.Status(code).JSON(resp)
		},
	})

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
	app.Use(cors.New())

	v1 := app.Group("/v1", func(c fiber.Ctx) error {
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

	log.Fatal(app.Listen(":" + config.HttpPort))
}

/*

	// log
	log := http.NewServeMux()
	v1.Handle("/log/", func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, ctxex.UserLogWriterKey, userLogWriter)
			http.StripPrefix("/log", log).ServeHTTP(w, r.WithContext(ctx))
		})
	}())

	errorwrp.GET(log, "/last", handlers.LogsHandler)
	errorwrp.GET(log, "/since", handlers.LogsSinceHandler)
	errorwrp.GET(log, "/files", handlers.LogFilesHandler)
	errorwrp.GET(log, "/file", handlers.LogFileContentHandler)
*/
