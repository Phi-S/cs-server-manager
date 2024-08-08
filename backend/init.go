package main

import (
	"cs-server-manager/config"
	"cs-server-manager/game_events"
	"cs-server-manager/jfile"
	"cs-server-manager/logwrt"
	"cs-server-manager/server"
	"cs-server-manager/status"
	"cs-server-manager/steamcmd"
	"log/slog"
	"os"
	"path/filepath"
)

func createdRequiredDirs(cfg config.Config) {
	if err := os.MkdirAll(cfg.DataDir, os.ModePerm); err != nil {
		slog.Error("failed to create data directory", "path", cfg.DataDir, "error", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(cfg.ServerDir, os.ModePerm); err != nil {
		slog.Error("failed to create server directory", "path", cfg.ServerDir, "error", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(cfg.SteamcmdDir, os.ModePerm); err != nil {
		slog.Error("failed to create steamcmd directory", "path", cfg.SteamcmdDir, "error", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(cfg.LogDir, os.ModePerm); err != nil {
		slog.Error("failed to create log directory", "path", cfg.LogDir, "error", err)
		os.Exit(1)
	}
}

func createRequiredServices(cfg config.Config) (
	*steamcmd.Instance,
	*server.Instance,
	*jfile.Instance[server.StartParameters],
	*logwrt.LogWriter,
	*status.Status,
	*WebSocketServer,
	*game_events.Instance,
) {
	steamcmdInstance, err := steamcmd.NewInstance(cfg.SteamcmdDir, cfg.ServerDir)
	if err != nil {
		slog.Error("failed to create new steamcmd instance", "error", err)
		os.Exit(1)
	}

	serverInstance, err := server.NewInstance(cfg.ServerDir, cfg.CsPort, cfg.SteamcmdDir)
	if err != nil {
		slog.Error("failed to create new server instance", "error", err)
		os.Exit(1)
	}

	startParametersJsonPath := filepath.Join(cfg.DataDir, "start-parameters.json")
	startParametersJsonFile, err := jfile.New[server.StartParameters](startParametersJsonPath, *server.DefaultStartParameters())
	if err != nil {
		slog.Error("failed to create new json file service for start-parameter.json", "path", startParametersJsonPath, "error", err)
		os.Exit(1)
	}

	logDir := filepath.Join(cfg.DataDir, "logs")
	userLogWriter, err := logwrt.NewLogWriter(logDir, "user")
	if err != nil {
		slog.Error("failed to create user log write", "error", err)
		os.Exit(1)
	}

	startParameters, err := startParametersJsonFile.Read()
	if err != nil {
		slog.Error("failed to read start-parameter.json", "error", err)
		os.Exit(1)
	}

	statusInstance := status.NewStatus(startParameters.Hostname, startParameters.MaxPlayers, startParameters.StartMap)
	webSocketServer := NewWebSocketServer()

	gameEventsInstance := game_events.Instance{}

	return steamcmdInstance, serverInstance, startParametersJsonFile, userLogWriter, statusInstance, webSocketServer, &gameEventsInstance
}
