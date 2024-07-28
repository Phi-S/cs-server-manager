package main

import (
	"cs-server-controller/config"
	"cs-server-controller/event"
	"cs-server-controller/handlers"
	"cs-server-controller/httpex/errorwrp"
	"cs-server-controller/middleware"
	"cs-server-controller/server"
	"cs-server-controller/steamcmd"
	"cs-server-controller/user_logs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	configureLogger()
	config, err := config.GetConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err := os.MkdirAll(config.DataDir, os.ModePerm); err != nil {
		slog.Error("failed to create data directory", "dir", config.DataDir, "error", err)
		os.Exit(1)
	}

	var steamcmdDir = filepath.Join(config.DataDir, "steamcmd")
	var serverDir = filepath.Join(config.DataDir, "server")

	// this lock is used to prevent collision between the server and steamcmd instance
	// Fox example the lock is used to prevent the server from being started while being updated.
	// This can occur if two http request are coming in at the same time
	ServerSteamcmdLock := sync.Mutex{}

	steamcmdInstance, err := steamcmd.NewInstance(steamcmdDir, serverDir, true)
	if err != nil {
		slog.Error("failed to create new steamcmd instance", "error", err)
		os.Exit(1)
	}
	defer steamcmdInstance.Cancel()
	defer steamcmdInstance.Close()

	serverInstance, err := server.NewInstance(serverDir, config.CsPort, true)
	if err != nil {
		slog.Error("failed to create new server instance", "error", err)
		os.Exit(1)
	}
	defer serverInstance.Stop()
	defer serverInstance.Close()

	logDir := filepath.Join(config.DataDir, "logs")
	if err := os.MkdirAll(logDir, 0777); err != nil {
		slog.Error("failed to create log dir", "logDir", logDir, "error", err)
		os.Exit(1)
	}

	userLogsWriter, err := user_logs.NewLogWriter(logDir, "user-logs")
	if err != nil {
		slog.Error("failed to create user logs write", "error", err)
		os.Exit(1)
	}
	defer userLogsWriter.Close()

	logServerEvents(userLogsWriter, serverInstance)
	logSteamcmdEvents(userLogsWriter, steamcmdInstance)

	main := http.NewServeMux()

	v1 := http.NewServeMux()

	errorwrp.GET(v1, "/status", handlers.StatusHandler)

	errorwrp.POST(v1, "/start", handlers.StartHandler)
	errorwrp.POST(v1, "/stop", handlers.StopHandler)
	errorwrp.POST(v1, "/send-command", handlers.SendCommandHandler)

	errorwrp.POST(v1, "/update", handlers.UpdateHandler)
	errorwrp.POST(v1, "/cancel-update", handlers.CancelUpdateHandler)

	errorwrp.GET(v1, "/logs", handlers.LogHandler)

	main.Handle("/v1/", middleware.ContextValues(
		http.StripPrefix("/v1", v1),
		config,
		&ServerSteamcmdLock,
		serverInstance,
		steamcmdInstance,
		userLogsWriter,
	))

	slog.Info("listening at port "+config.HttpPort, "port", config.HttpPort)
	slog.Error("Failed to start http server", "error",
		http.ListenAndServe(":"+config.HttpPort, middleware.TraceId(
			middleware.Logging(
				middleware.Recover(
					main,
				),
			),
		)),
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

func logServerEvents(w *user_logs.LogWriter, s *server.ServerInstance) {
	s.OnServerStarting(func(dp event.DefaultPayload) {
		w.WriteSystemInfoLog(dp.TriggeredAtUtc, "server starting")
	})

	s.OnServerStarted(func(dp event.DefaultPayload) {
		w.WriteSystemInfoLog(dp.TriggeredAtUtc, "server started")
	})

	s.OnServerStopped(func(dp event.DefaultPayload) {
		w.WriteSystemInfoLog(dp.TriggeredAtUtc, "server stopped")
	})

	s.OnServerCrashed(func(dp event.PayloadWithData[error]) {
		w.WriteSystemErrorLog(dp.TriggeredAtUtc, "server crashed")
	})

	s.OnOutput(func(dp event.PayloadWithData[string]) {
		w.WriteServerLog(dp.TriggeredAtUtc, dp.Data)
	})
}

func logSteamcmdEvents(w *user_logs.LogWriter, s *steamcmd.SteamcmdInstance) {
	s.OnStarted(func(dp event.DefaultPayload) {
		w.WriteSystemInfoLog(dp.TriggeredAtUtc, "steamcmd update starting")
	})

	s.OnFinished(func(dp event.DefaultPayload) {
		w.WriteSystemInfoLog(dp.TriggeredAtUtc, "steamcmd update finished")
	})

	s.OnCancelled(func(dp event.DefaultPayload) {
		w.WriteSystemInfoLog(dp.TriggeredAtUtc, "steamcmd update cancelled")
	})

	s.OnFailed(func(dp event.PayloadWithData[error]) {
		w.WriteSystemErrorLog(dp.TriggeredAtUtc, "steamcmd update failed")
	})

	s.OnOutput(func(dp event.PayloadWithData[string]) {
		w.WriteSteamcmdLog(dp.TriggeredAtUtc, dp.Data)
	})

}
