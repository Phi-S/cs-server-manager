package main

import (
	"cs-server-manager/event"
	"cs-server-manager/game_events"
	"cs-server-manager/server"
	"cs-server-manager/status"
	"cs-server-manager/steamcmd"
)

func statusEventHandler(statusInstance *status.Status, serverInstance *server.Instance, steamcmdInstance *steamcmd.Instance, gameEventsInstance *game_events.Instance) {
	// server
	serverInstance.OnStarting(func(p event.DefaultPayload) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.Server = status.ServerStatusStarting
		})
	})

	serverInstance.OnStarted(func(e event.PayloadWithData[server.StartParameters]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.Server = status.ServerStatusStarted
			internalStatus.Hostname = e.Data.Hostname
			internalStatus.MaxPlayerCount = e.Data.MaxPlayers
			internalStatus.Map = e.Data.StartMap
		})
	})

	serverInstance.OnCrashed(func(p event.PayloadWithData[error]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.Server = status.ServerStatusStopped
		})
	})

	serverInstance.OnStopped(func(p event.DefaultPayload) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.Server = status.ServerStatusStopped
		})
	})

	// steamcmd

	//steamcmdInstance.OnOutput(func(p event.PayloadWithData[string]) {
	//        fmt.Println(p.TriggeredAtUtc.String() + " | steamcmdInstance: " + p.Data)
	//slog.Debug("steamcmdInstance", "event", "onOutput", "triggeredAtUtc", p.TriggeredAtUtc, "output", p.Data)
	//})

	steamcmdInstance.OnStarted(func(p event.DefaultPayload) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.Steamcmd = status.SteamcmdStatusUpdating
		})
	})

	steamcmdInstance.OnFinished(func(p event.DefaultPayload) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.Steamcmd = status.SteamcmdStatusStopped
		})
	})

	steamcmdInstance.OnCancelled(func(p event.DefaultPayload) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.Steamcmd = status.SteamcmdStatusStopped
		})
	})

	steamcmdInstance.OnFailed(func(p event.PayloadWithData[error]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.Steamcmd = status.SteamcmdStatusStopped
		})
	})

	//game_events
	gameEventsInstance.OnMapChanged(func(p event.PayloadWithData[string]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.Map = p.Data
		})
	})

	gameEventsInstance.OnPlayerConnected(func(p event.PayloadWithData[game_events.PlayerConnected]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.PlayerCount++
		})
	})

	gameEventsInstance.OnPlayerDisconnected(func(p event.PayloadWithData[string]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.PlayerCount--
		})
	})
}
