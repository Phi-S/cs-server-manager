package main

import (
	"cs-server-controller/event"
	"cs-server-controller/logwrt"
	"cs-server-controller/server"
	"cs-server-controller/steamcmd"
)

func writeEventsTpUserLogFile(w *logwrt.LogWriter, server *server.ServerInstance, steamcmd *steamcmd.SteamcmdInstance) {
	const systemInfoLogType = "system-info"
	const systemErrorLogType = "system-error"
	const serverLogType = "server"
	const steamcmdLogType = "steamcmd"

	server.OnServerStarting(func(p event.DefaultPayload) {
		w.WriteLog(p.TriggeredAtUtc, systemInfoLogType, "server starting")
	})

	server.OnServerStarted(func(p event.DefaultPayload) {
		w.WriteLog(p.TriggeredAtUtc, systemInfoLogType, "server started")
	})

	server.OnServerStopped(func(p event.DefaultPayload) {
		w.WriteLog(p.TriggeredAtUtc, systemInfoLogType, "server stopped")
	})

	server.OnServerCrashed(func(p event.PayloadWithData[error]) {
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
