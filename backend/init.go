package main

import (
	"cs-server-manager/config"
	"cs-server-manager/event"
	"cs-server-manager/game_events"
	"cs-server-manager/jfile"
	"cs-server-manager/logwrt"
	"cs-server-manager/plugins"
	"cs-server-manager/server"
	"cs-server-manager/status"
	"cs-server-manager/steamcmd"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func createdRequiredDirs(cfg config.Config) error {
	if err := os.MkdirAll(cfg.DataDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create data dir %v %w", cfg.DataDir, err)
	}

	if err := os.MkdirAll(cfg.ServerDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create server dir %v %w", cfg.ServerDir, err)
	}

	if err := os.MkdirAll(cfg.SteamcmdDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create steamcmd dir %v %w", cfg.SteamcmdDir, err)
	}

	if err := os.MkdirAll(cfg.LogDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log dir %v %w", cfg.LogDir, err)
	}

	return nil
}

func createRequiredServices(cfg config.Config) (
	*steamcmd.Instance,
	*server.Instance,
	*jfile.Instance[server.StartParameters],
	*logwrt.LogWriter,
	*status.Status,
	*WebSocketServer,
	*game_events.Instance,
	*plugins.Instance,
	error,
) {

	steamcmdInstance, err := steamcmd.NewInstance(cfg.SteamcmdDir, cfg.ServerDir)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to create steamcmd instance %w", err)
	}

	serverInstance, err := server.NewInstance(cfg.ServerDir, cfg.CsPort, cfg.SteamcmdDir)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to create server instance %w", err)
	}

	startParametersJsonPath := filepath.Join(cfg.DataDir, "start-parameters.json")
	startParametersJsonFile, err := jfile.New[server.StartParameters](startParametersJsonPath, *server.DefaultStartParameters())
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to create new json file service for start-parameter.json %w", err)
	}

	logDir := filepath.Join(cfg.DataDir, "logs")
	userLogWriter, err := logwrt.NewLogWriter(logDir, "user")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to create user log writer %w", err)
	}

	startParameters, err := startParametersJsonFile.Read()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to read start-parameters.json %w", err)
	}

	statusInstance := status.NewStatus(
		startParameters.Hostname,
		cfg.Ip,
		cfg.CsPort,
		startParameters.Password,
		startParameters.MaxPlayers,
		startParameters.StartMap,
	)

	webSocketServer := NewWebSocketServer()

	gameEventsInstance := game_events.Instance{}

	pluginsJsonFilePath := filepath.Join(cfg.DataDir, "plugins.json")
	installedPluginsJsonPath := filepath.Join(cfg.DataDir, "installed-plugins.json")
	csgoDirPath := filepath.Join(cfg.ServerDir, "game", "csgo")
	pluginsInstance, err := plugins.New(csgoDirPath, pluginsJsonFilePath, installedPluginsJsonPath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to create plugins instance %w", err)
	}

	return steamcmdInstance,
		serverInstance,
		startParametersJsonFile,
		userLogWriter,
		statusInstance,
		webSocketServer,
		&gameEventsInstance,
		pluginsInstance,
		nil
}

func registerEvents(
	configInstance config.Config,
	serverInstance *server.Instance,
	steamcmdInstance *steamcmd.Instance,
	startParametersJfileInstance *jfile.Instance[server.StartParameters],
	logWriterInstance *logwrt.LogWriter,
	statusInstance *status.Status,
	webSocketServerInstance *WebSocketServer,
	gameEventsInstance *game_events.Instance,
	pluginsInstance *plugins.Instance,
) {
	logEvents(logWriterInstance, webSocketServerInstance, serverInstance, steamcmdInstance, gameEventsInstance, pluginsInstance)

	// detect game events via server output
	serverInstance.OnOutput(func(p event.PayloadWithData[string]) {
		gameEventsInstance.DetectGameEvent(p.Data)
	})

	// send status update via websocket
	statusInstance.OnStatusChanged(func(p event.PayloadWithData[status.InternalStatus]) {
		if err := webSocketServerInstance.Broadcast("status", p.Data); err != nil {
			slog.Error("failed to send status message", "status", p.Data, "error", err)
		}
	})

	// update status if start parameters get changed (only applies if server is stopped / status gets updated on server start anyway)
	startParametersJfileInstance.OnUpdated(func(data event.PayloadWithData[server.StartParameters]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			if internalStatus.Server == status.ServerStatusStopped {
				internalStatus.Hostname = data.Data.Hostname
				internalStatus.Map = data.Data.StartMap
				internalStatus.MaxPlayerCount = data.Data.MaxPlayers
				internalStatus.Password = data.Data.Password
			}
		})
	})

	serverInstance.OnStarting(func(p event.DefaultPayload) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.Server = status.ServerStatusStarting
		})
	})

	serverInstance.OnStarted(func(e event.PayloadWithData[server.StartParameters]) {
		ip := configInstance.Ip
		if !configInstance.UsedIpFromEnv() {
			publicIp, err := config.GetPublicIp()
			if err != nil {
				slog.Error("failed to get public ip after server started", "error", err)
			} else {
				ip = publicIp
			}
		}

		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.Server = status.ServerStatusStarted
			internalStatus.Hostname = e.Data.Hostname
			internalStatus.MaxPlayerCount = e.Data.MaxPlayers
			internalStatus.Map = e.Data.StartMap
			internalStatus.Ip = ip
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

	// TODO: detect update progress???
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
