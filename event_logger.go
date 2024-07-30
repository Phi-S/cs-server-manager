package main

import (
    "cs-server-controller/event"
    "cs-server-controller/logwrt"
    "cs-server-controller/server"
    "cs-server-controller/steamcmd"
    "log/slog"
)

func enableEventLogging(s *server.Instance, steamcmd *steamcmd.Instance) {
    s.OnOutput(func(pwd event.PayloadWithData[string]) {
        slog.Debug("server", "event", "onOutput", "triggeredAtUtc", pwd.TriggeredAtUtc, "output", pwd.Data)
        //    fmt.Println(pwd.TriggeredAtUtc.Format(time.RFC3339Nano) + " | SERVER: " + pwd.Data)
    })

    s.OnStarting(func(dp event.DefaultPayload) {
        slog.Debug("server", "event", "onStarting", "triggeredAtUtc", dp.TriggeredAtUtc)
    })

    s.OnStarted(func(dp event.DefaultPayload) {
        slog.Debug("server", "event", "onStarted", "triggeredAtUtc", dp.TriggeredAtUtc)
    })

    s.OnCrashed(func(pwd event.PayloadWithData[error]) {
        slog.Debug("server", "event", "onCrashed", "triggeredAtUtc", pwd.TriggeredAtUtc, "data", pwd.Data)
    })

    s.OnStopped(func(dp event.DefaultPayload) {
        slog.Debug("server", "event", "onStopped", "triggeredAtUtc", dp.TriggeredAtUtc)
    })

    steamcmd.OnOutput(func(p event.PayloadWithData[string]) {
        //        fmt.Println(p.TriggeredAtUtc.String() + " | steamcmd: " + p.Data)
        slog.Debug("steamcmd", "event", "onOutput", "triggeredAtUtc", p.TriggeredAtUtc, "output", p.Data)
    })

    steamcmd.OnStarted(func(p event.DefaultPayload) {
        slog.Debug("steamcmd", "event", "onStarted", "triggeredAtUtc", p.TriggeredAtUtc)
    })

    steamcmd.OnFinished(func(p event.DefaultPayload) {
        slog.Debug("steamcmd", "event", "onFinished", "triggeredAtUtc", p.TriggeredAtUtc)
    })

    steamcmd.OnCancelled(func(p event.DefaultPayload) {
        slog.Debug("steamcmd", "event", "onCancelled", "triggeredAtUtc", p.TriggeredAtUtc)
    })

    steamcmd.OnFailed(func(p event.PayloadWithData[error]) {
        slog.Debug("steamcmd", "event", "onFailed", "triggeredAtUtc", p.TriggeredAtUtc, "data", p.Data)
    })
}

func writeEventsTpUserLogFile(w *logwrt.LogWriter, server *server.Instance, steamcmd *steamcmd.Instance) {
    const systemInfoLogType = "system-info"
    const systemErrorLogType = "system-error"
    const serverLogType = "server"
    const steamcmdLogType = "steamcmd"

    server.OnStarting(func(p event.DefaultPayload) {
        w.WriteLog(p.TriggeredAtUtc, systemInfoLogType, "server starting")
    })

    server.OnStarted(func(p event.DefaultPayload) {
        w.WriteLog(p.TriggeredAtUtc, systemInfoLogType, "server started")
    })

    server.OnStopped(func(p event.DefaultPayload) {
        w.WriteLog(p.TriggeredAtUtc, systemInfoLogType, "server stopped")
    })

    server.OnCrashed(func(p event.PayloadWithData[error]) {
        w.WriteLog(p.TriggeredAtUtc, systemErrorLogType, "server crashed")
    })

    server.OnOutput(func(p event.PayloadWithData[string]) {
        w.WriteLog(p.TriggeredAtUtc, serverLogType, p.Data)
    })

    steamcmd.OnStarted(func(p event.DefaultPayload) {
        w.WriteLog(p.TriggeredAtUtc, systemInfoLogType, "steamcmd update starting")
    })

    steamcmd.OnFinished(func(p event.DefaultPayload) {
        w.WriteLog(p.TriggeredAtUtc, systemInfoLogType, "steamcmd update finished")
    })

    steamcmd.OnCancelled(func(p event.DefaultPayload) {
        w.WriteLog(p.TriggeredAtUtc, systemInfoLogType, "steamcmd update cancelled")
    })

    steamcmd.OnFailed(func(p event.PayloadWithData[error]) {
        w.WriteLog(p.TriggeredAtUtc, systemErrorLogType, "steamcmd update failed")
    })

    steamcmd.OnOutput(func(p event.PayloadWithData[string]) {
        w.WriteLog(p.TriggeredAtUtc, steamcmdLogType, p.Data)
    })
}
