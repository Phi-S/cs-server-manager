package main

import (
	"cs-server-manager/event"
	"cs-server-manager/game_events"
	"cs-server-manager/logwrt"
	"cs-server-manager/plugins"
	"cs-server-manager/server"
	"cs-server-manager/steamcmd"
	"fmt"
	"log/slog"
	"time"
)

// Handles server and steamcmd events to log them to the console and log file.
// Also sends the events to all connected websocket clients
func logEvents(
	logWriter *logwrt.LogWriter,
	webSocketServer *WebSocketServer,
	serverInstance *server.Instance,
	steamcmdInstance *steamcmd.Instance,
	gameEventsInstance *game_events.Instance,
	pluginsInstance *plugins.Instance,
) {
	handleEvent := func(logType string, timestampUtc time.Time, message string, args ...any) {
		logEntry := logwrt.NewLogEntry(timestampUtc, logType, message)

		args = append(args, "timestamp-utc")
		args = append(args, timestampUtc)
		args = append(args, "message")
		args = append(args, message)

		slog.Debug(logType, args...)

		if err := logWriter.WriteLogEntry(logEntry); err != nil {
			slog.Error("failed to write log entry", "log_entry", logEntry, "error", err)
		}
		if err := webSocketServer.BroadcastLogMessage(logEntry); err != nil {
			slog.Error("failed to broadcast log message", "log_entry", logEntry, "error", err)
		}
	}

	const systemInfoLogType = "system_info"
	const systemErrorLogType = "system_error"
	const serverLogType = "server_log"
	const steamcmdLogType = "steamcmd_log"

	// server
	serverInstance.OnStarting(func(p event.DefaultPayload) {
		handleEvent(systemInfoLogType, p.TriggeredAtUtc, "server starting")
	})

	serverInstance.OnStarted(func(p event.PayloadWithData[server.StartParameters]) {
		handleEvent(systemInfoLogType, p.TriggeredAtUtc, "server started")
	})

	serverInstance.OnStopped(func(p event.DefaultPayload) {
		handleEvent(systemInfoLogType, p.TriggeredAtUtc, "server stopped")
	})

	serverInstance.OnCrashed(func(p event.PayloadWithData[error]) {
		handleEvent(systemErrorLogType, p.TriggeredAtUtc, "server crashed with error", "error", p.Data)
	})

	serverInstance.OnOutput(func(p event.PayloadWithData[string]) {
		handleEvent(serverLogType, p.TriggeredAtUtc, p.Data)
	})

	// steamcmd
	steamcmdInstance.OnStarted(func(p event.DefaultPayload) {
		handleEvent(systemInfoLogType, p.TriggeredAtUtc, "server update started")
	})

	steamcmdInstance.OnFinished(func(p event.DefaultPayload) {
		handleEvent(systemInfoLogType, p.TriggeredAtUtc, "server update finished")
	})

	steamcmdInstance.OnCancelled(func(p event.DefaultPayload) {
		handleEvent(systemInfoLogType, p.TriggeredAtUtc, "server update cancelled")
	})

	steamcmdInstance.OnFailed(func(p event.PayloadWithData[error]) {
		handleEvent(systemErrorLogType, p.TriggeredAtUtc, "server update failed with error", "error", p.Data)
	})

	steamcmdInstance.OnOutput(func(p event.PayloadWithData[string]) {
		handleEvent(steamcmdLogType, p.TriggeredAtUtc, p.Data)
	})

	// game_events
	gameEventsInstance.OnMapChanged(func(p event.PayloadWithData[string]) {
		handleEvent(systemInfoLogType, p.TriggeredAtUtc, fmt.Sprintf("Map changed to %v", p.Data))
	})

	gameEventsInstance.OnPlayerConnected(func(p event.PayloadWithData[game_events.PlayerConnected]) {
		handleEvent(systemInfoLogType, p.TriggeredAtUtc, fmt.Sprintf("New player '%v'(%v) connected from '%v:%v'", p.Data.Name, p.Data.Id, p.Data.Ip, p.Data.Port))
	})

	gameEventsInstance.OnPlayerDisconnected(func(p event.PayloadWithData[string]) {
		handleEvent(systemInfoLogType, p.TriggeredAtUtc, fmt.Sprintf("Player '%v' disconnected", p.Data))
	})

	// plugins
	pluginsInstance.OnPluginInstalled(func(p event.PayloadWithData[plugins.PluginEventsPayload]) {
		handleEvent(systemInfoLogType, p.TriggeredAtUtc, fmt.Sprintf("Plugin '%v(%v)' installed", p.Data.Name, p.Data.Version))
	})

	pluginsInstance.OnPluginUninstalledEvent(func(p event.PayloadWithData[plugins.PluginEventsPayload]) {
		handleEvent(systemInfoLogType, p.TriggeredAtUtc, fmt.Sprintf("Plugin '%v' uninstalled", p.Data.Name))
	})
}
