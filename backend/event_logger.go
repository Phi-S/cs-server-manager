package main

import (
    "cs-server-controller/event"
    "cs-server-controller/logwrt"
    "cs-server-controller/server"
    "cs-server-controller/steamcmd"
    "log/slog"
)

func enableEventLogging(serverInstance *server.Instance, steamcmdInstance *steamcmd.Instance) {
    serverInstance.OnOutput(func(pwd event.PayloadWithData[string]) {
        slog.Debug("serverInstance", "event", "onOutput", "triggeredAtUtc", pwd.TriggeredAtUtc, "output", pwd.Data)
        //    fmt.Println(pwd.TriggeredAtUtc.Format(time.RFC3339Nano) + " | SERVER: " + pwd.Data)
    })

    serverInstance.OnStarting(func(dp event.DefaultPayload) {
        slog.Debug("serverInstance", "event", "onStarting", "triggeredAtUtc", dp.TriggeredAtUtc)
    })

    serverInstance.OnStarted(func(dp event.PayloadWithData[server.StartParameters]) {
        slog.Debug("serverInstance", "event", "onStarted", "triggeredAtUtc", dp.TriggeredAtUtc, "data", dp.Data)
    })

    serverInstance.OnCrashed(func(pwd event.PayloadWithData[error]) {
        slog.Debug("serverInstance", "event", "onCrashed", "triggeredAtUtc", pwd.TriggeredAtUtc, "data", pwd.Data)
    })

    serverInstance.OnStopped(func(dp event.DefaultPayload) {
        slog.Debug("serverInstance", "event", "onStopped", "triggeredAtUtc", dp.TriggeredAtUtc)
    })

    steamcmdInstance.OnOutput(func(p event.PayloadWithData[string]) {
        //        fmt.Println(p.TriggeredAtUtc.String() + " | steamcmdInstance: " + p.Data)
        slog.Debug("steamcmdInstance", "event", "onOutput", "triggeredAtUtc", p.TriggeredAtUtc, "output", p.Data)
    })

    steamcmdInstance.OnStarted(func(p event.DefaultPayload) {
        slog.Debug("steamcmdInstance", "event", "onStarted", "triggeredAtUtc", p.TriggeredAtUtc)
    })

    steamcmdInstance.OnFinished(func(p event.DefaultPayload) {
        slog.Debug("steamcmdInstance", "event", "onFinished", "triggeredAtUtc", p.TriggeredAtUtc)
    })

    steamcmdInstance.OnCancelled(func(p event.DefaultPayload) {
        slog.Debug("steamcmdInstance", "event", "onCancelled", "triggeredAtUtc", p.TriggeredAtUtc)
    })

    steamcmdInstance.OnFailed(func(p event.PayloadWithData[error]) {
        slog.Debug("steamcmdInstance", "event", "onFailed", "triggeredAtUtc", p.TriggeredAtUtc, "data", p.Data)
    })
}

func writeEventToLogFileAndWebSocket(
    logWriter *logwrt.LogWriter,
    webSocketserver *WebSocketServer,
    serverInstance *server.Instance,
    steamcmdInstance *steamcmd.Instance,
) {
    const systemInfoLogType = "system-info"
    const systemErrorLogType = "system-error"
    const serverLogType = "serverInstance"
    const steamcmdLogType = "steamcmdInstance"

    serverInstance.OnStarting(func(p event.DefaultPayload) {
        logEntry := logwrt.NewLogEntry(p.TriggeredAtUtc, systemInfoLogType, "serverInstance starting")
        if err := logWriter.WriteLogEntry(logEntry); err != nil {
            slog.Error("failed to write log entry", "log-entry", logEntry, "error", err)
        }
        if err := webSocketserver.BroadcastLogMessage(logEntry); err != nil {
            slog.Error("failed to brodcast log message", "log-entry", logEntry, "error", err)
        }
    })

    serverInstance.OnStarted(func(p event.PayloadWithData[server.StartParameters]) {
        logEntry := logwrt.NewLogEntry(p.TriggeredAtUtc, systemInfoLogType, "serverInstance started")
        if err := logWriter.WriteLogEntry(logEntry); err != nil {
            slog.Error("failed to write log entry", "log-entry", logEntry, "error", err)
        }
        if err := webSocketserver.BroadcastLogMessage(logEntry); err != nil {
            slog.Error("failed to brodcast log message", "log-entry", logEntry, "error", err)
        }
    })

    serverInstance.OnStopped(func(p event.DefaultPayload) {
        logEntry := logwrt.NewLogEntry(p.TriggeredAtUtc, systemInfoLogType, "serverInstance stopped")
        if err := logWriter.WriteLogEntry(logEntry); err != nil {
            slog.Error("failed to write log entry", "log-entry", logEntry, "error", err)
        }
        if err := webSocketserver.BroadcastLogMessage(logEntry); err != nil {
            slog.Error("failed to brodcast log message", "log-entry", logEntry, "error", err)
        }

    })

    serverInstance.OnCrashed(func(p event.PayloadWithData[error]) {
        logEntry := logwrt.NewLogEntry(p.TriggeredAtUtc, systemErrorLogType, "serverInstance crashed")
        if err := logWriter.WriteLogEntry(logEntry); err != nil {
            slog.Error("failed to write log entry", "log-entry", logEntry, "error", err)
        }
        if err := webSocketserver.BroadcastLogMessage(logEntry); err != nil {
            slog.Error("failed to brodcast log message", "log-entry", logEntry, "error", err)
        }

    })

    serverInstance.OnOutput(func(p event.PayloadWithData[string]) {
        logEntry := logwrt.NewLogEntry(p.TriggeredAtUtc, serverLogType, p.Data)
        if err := logWriter.WriteLogEntry(logEntry); err != nil {
            slog.Error("failed to write log entry", "log-entry", logEntry, "error", err)
        }
        if err := webSocketserver.BroadcastLogMessage(logEntry); err != nil {
            slog.Error("failed to brodcast log message", "log-entry", logEntry, "error", err)
        }
    })

    steamcmdInstance.OnStarted(func(p event.DefaultPayload) {
        logEntry := logwrt.NewLogEntry(p.TriggeredAtUtc, systemInfoLogType, "steamcmdInstance update starting")
        if err := logWriter.WriteLogEntry(logEntry); err != nil {
            slog.Error("failed to write log entry", "log-entry", logEntry, "error", err)
        }
        if err := webSocketserver.BroadcastLogMessage(logEntry); err != nil {
            slog.Error("failed to brodcast log message", "log-entry", logEntry, "error", err)
        }
    })

    steamcmdInstance.OnFinished(func(p event.DefaultPayload) {
        logEntry := logwrt.NewLogEntry(p.TriggeredAtUtc, systemInfoLogType, "steamcmdInstance update finished")
        if err := logWriter.WriteLogEntry(logEntry); err != nil {
            slog.Error("failed to write log entry", "log-entry", logEntry, "error", err)
        }
        if err := webSocketserver.BroadcastLogMessage(logEntry); err != nil {
            slog.Error("failed to brodcast log message", "log-entry", logEntry, "error", err)
        }
    })

    steamcmdInstance.OnCancelled(func(p event.DefaultPayload) {
        logEntry := logwrt.NewLogEntry(p.TriggeredAtUtc, systemInfoLogType, "steamcmdInstance update cancelled")
        if err := logWriter.WriteLogEntry(logEntry); err != nil {
            slog.Error("failed to write log entry", "log-entry", logEntry, "error", err)
        }
        if err := webSocketserver.BroadcastLogMessage(logEntry); err != nil {
            slog.Error("failed to brodcast log message", "log-entry", logEntry, "error", err)
        }
    })

    steamcmdInstance.OnFailed(func(p event.PayloadWithData[error]) {
        logEntry := logwrt.NewLogEntry(p.TriggeredAtUtc, systemErrorLogType, "steamcmdInstance update failed")
        if err := logWriter.WriteLogEntry(logEntry); err != nil {
            slog.Error("failed to write log entry", "log-entry", logEntry, "error", err)
        }
        if err := webSocketserver.BroadcastLogMessage(logEntry); err != nil {
            slog.Error("failed to brodcast log message", "log-entry", logEntry, "error", err)
        }
    })

    steamcmdInstance.OnOutput(func(p event.PayloadWithData[string]) {
        logEntry := logwrt.NewLogEntry(p.TriggeredAtUtc, steamcmdLogType, p.Data)
        if err := logWriter.WriteLogEntry(logEntry); err != nil {
            slog.Error("failed to write log entry", "log-entry", logEntry, "error", err)
        }
        if err := webSocketserver.BroadcastLogMessage(logEntry); err != nil {
            slog.Error("failed to brodcast log message", "log-entry", logEntry, "error", err)
        }
    })
}
