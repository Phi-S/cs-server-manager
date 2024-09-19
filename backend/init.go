package main

import (
	"cs-server-manager/config"
	"cs-server-manager/editor"
	"cs-server-manager/event"
	"cs-server-manager/files"
	"cs-server-manager/game_events"
	"cs-server-manager/gvalidator"
	"cs-server-manager/logwrt"
	"cs-server-manager/plugins"
	"cs-server-manager/server"
	"cs-server-manager/start_parameters_json"
	"cs-server-manager/status"
	"cs-server-manager/steamcmd"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func createdRequiredDirs(cfg config.Config) error {
	if err := os.MkdirAll(cfg.DataDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create data dir '%v' %w", cfg.DataDir, err)
	}

	if err := os.MkdirAll(cfg.ServerDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create server dir '%v' %w", cfg.ServerDir, err)
	}

	if err := os.MkdirAll(cfg.SteamcmdDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create steamcmd dir '%v' %w", cfg.SteamcmdDir, err)
	}

	if err := os.MkdirAll(cfg.LogDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log dir '%v' %w", cfg.LogDir, err)
	}

	return nil
}

func createRequiredServices(cfg config.Config) (
	*steamcmd.Instance,
	*server.Instance,
	*start_parameters_json.Instance,
	*logwrt.LogWriter,
	*status.Status,
	*WebSocketServer,
	*game_events.Instance,
	*plugins.Instance,
	*editor.Instance,
	error,
) {
	steamcmdInstance, err := steamcmd.NewInstance(cfg.SteamcmdDir, cfg.ServerDir)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("create steamcmd instance: %w", err)
	}

	serverInstance, err := server.NewInstance(cfg.ServerDir, cfg.CsPort, cfg.SteamcmdDir)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("create server instance: %w", err)
	}

	startParametersJsonPath := filepath.Join(cfg.DataDir, "start-parameters.json")
	startParametersJsonFile, err := start_parameters_json.New(startParametersJsonPath, *server.DefaultStartParameters())
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("create start parameter json instance: %w", err)
	}

	logDir := filepath.Join(cfg.DataDir, "logs")
	userLogWriter, err := logwrt.NewLogWriter(logDir, "user")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("create user log writer: %w", err)
	}

	startParameters, err := startParametersJsonFile.Read()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("read start-parameters.json: %w", err)
	}

	isGameServerInstalled, err := isGameServerInstalled(cfg.ServerDir)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("check if game server is installed: %w", err)
	}

	statusInstance := status.NewStatus(
		isGameServerInstalled,
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
	installedPluginsJsonPath := filepath.Join(cfg.DataDir, "installed-plugin.json")
	csgoDir := filepath.Join(cfg.ServerDir, "game", "csgo")
	if !strings.HasSuffix(csgoDir, string(filepath.Separator)) {
		csgoDir = csgoDir + string(filepath.Separator)
	}
	pluginsInstance, err := plugins.New(csgoDir, pluginsJsonFilePath, installedPluginsJsonPath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("create plugins instance: %w", err)
	}

	editorFilesJsonPath := filepath.Join(cfg.DataDir, "editor-files.json")
	editorInstance, err := editor.New(editorFilesJsonPath, cfg.ServerDir)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("create editor instance: %w", err)
	}

	return steamcmdInstance,
		serverInstance,
		startParametersJsonFile,
		userLogWriter,
		statusInstance,
		webSocketServer,
		&gameEventsInstance,
		pluginsInstance,
		editorInstance,
		nil
}

func registerEvents(
	configInstance config.Config,
	serverInstance *server.Instance,
	steamcmdInstance *steamcmd.Instance,
	startParametersJfileInstance *start_parameters_json.Instance,
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
			slog.Error("after status changed: send status message", "status", p.Data, "error", err)
		}
	})

	startParametersJfileInstance.OnUpdated(func(data event.PayloadWithData[server.StartParameters]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			// update status if start parameters get changed (only applies if server is stopped / status gets updated on server start anyway)
			if internalStatus.State == status.Idle {
				internalStatus.Hostname = data.Data.Hostname
				internalStatus.Map = data.Data.StartMap
				internalStatus.MaxPlayerCount = data.Data.MaxPlayers
				internalStatus.Password = data.Data.Password
			}
		})
	})

	serverInstance.OnStarting(func(p event.DefaultPayload) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			ip, err := configInstance.GetCurrentIp()
			if err != nil {
				slog.Error("after starting: get current ip", "error", err)
				ip = internalStatus.Ip
			}

			internalStatus.State = status.ServerStarting
			internalStatus.Ip = ip
		})
	})

	serverInstance.OnStarted(func(e event.PayloadWithData[server.StartParameters]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.State = status.ServerStarted
			internalStatus.Hostname = e.Data.Hostname
			internalStatus.MaxPlayerCount = e.Data.MaxPlayers
			internalStatus.Map = e.Data.StartMap
			internalStatus.Password = e.Data.Password
		})
	})

	serverInstance.OnCrashed(func(p event.PayloadWithData[error]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.State = status.Idle

			startParametersJson, err := startParametersJfileInstance.Read()
			if err != nil {
				slog.Error("after server crashed: read start parameters json", "error", err)
				return
			}

			internalStatus.Hostname = startParametersJson.Hostname
			internalStatus.MaxPlayerCount = startParametersJson.MaxPlayers
			internalStatus.Map = startParametersJson.StartMap
			internalStatus.Password = startParametersJson.Password
		})
	})

	serverInstance.OnStopped(func(p event.DefaultPayload) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.State = status.Idle

			startParametersJson, err := startParametersJfileInstance.Read()
			if err != nil {
				slog.Error("after server stopped: read start parameters json", "error", err)
				return
			}

			internalStatus.Hostname = startParametersJson.Hostname
			internalStatus.MaxPlayerCount = startParametersJson.MaxPlayers
			internalStatus.Map = startParametersJson.StartMap
			internalStatus.Password = startParametersJson.Password
		})
	})

	steamcmdInstance.OnStarted(func(p event.DefaultPayload) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.State = status.SteamcmdUpdating
		})
	})

	steamcmdInstance.OnFinished(func(p event.DefaultPayload) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.State = status.Idle
			isServerInstalled, err := isGameServerInstalled(configInstance.ServerDir)
			if err != nil {
				slog.Warn("after steamcmd finished: check if game server is installed", "error", err)
				return
			}

			internalStatus.IsGameServerInstalled = isServerInstalled
		})
	})

	steamcmdInstance.OnCancelled(func(p event.DefaultPayload) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.State = status.Idle

			isServerInstalled, err := isGameServerInstalled(configInstance.ServerDir)
			if err != nil {
				slog.Warn("after steamcmd go canceled: check if game server is installed", "error", err)
				return
			}

			internalStatus.IsGameServerInstalled = isServerInstalled
		})
	})

	steamcmdInstance.OnFailed(func(p event.PayloadWithData[error]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.State = status.Idle

			isServerInstalled, err := isGameServerInstalled(configInstance.ServerDir)
			if err != nil {
				slog.Warn("after steamcmd failed: check if game server is installed", "error", err)
				return
			}

			internalStatus.IsGameServerInstalled = isServerInstalled
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

	//plugins
	pluginsInstance.OnPluginInstalling(func(p event.PayloadWithData[plugins.PluginEventsPayload]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.State = status.PluginInstalling
		})
	})

	pluginsInstance.OnPluginInstalled(func(p event.PayloadWithData[plugins.PluginEventsPayload]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.State = status.Idle
		})
	})

	pluginsInstance.OnPluginInstallationFailedEvent(func(p event.PayloadWithData[plugins.PluginEventsPayload]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.State = status.Idle
		})
	})

	pluginsInstance.OnPluginUninstallingEvent(func(p event.PayloadWithData[plugins.PluginEventsPayload]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.State = status.PluginUninstalling
		})
	})

	pluginsInstance.OnPluginUninstalledEvent(func(p event.PayloadWithData[plugins.PluginEventsPayload]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.State = status.Idle
		})
	})

	pluginsInstance.OnPluginUninstallFailedEvent(func(p event.PayloadWithData[plugins.PluginEventsPayload]) {
		statusInstance.Update(func(internalStatus *status.InternalStatus) {
			internalStatus.State = status.Idle
		})
	})
}

func isGameServerInstalled(serverDir string) (bool, error) {
	if err := gvalidator.Instance().Var(serverDir, "dir"); err != nil {
		return false, nil
	}

	size, err := files.GetDirSize(serverDir)
	if err != nil {
		return false, fmt.Errorf("failed to get serverDir '%v' size: %w", serverDir, err)
	}

	gib := size / 1024 / 1024 / 1024
	if gib < 30 {
		return false, nil
	}

	csgoDir := filepath.Join(serverDir, "game", "csgo")
	if err := gvalidator.Instance().Var(csgoDir, "dir"); err != nil {
		return false, nil
	}

	if err := gvalidator.Instance().Var(filepath.Join(csgoDir, "pak01_001.vpk"), "file"); err != nil {
		return false, nil
	}

	if err := gvalidator.Instance().Var(filepath.Join(csgoDir, "pak01_001.vpk"), "file"); err != nil {
		return false, nil
	}

	if err := gvalidator.Instance().Var(filepath.Join(csgoDir, "pak01_002.vpk"), "file"); err != nil {
		return false, nil
	}

	cfgFolder := filepath.Join(csgoDir, "cfg")
	if err := gvalidator.Instance().Var(csgoDir, "dir"); err != nil {
		return false, nil
	}

	if err := gvalidator.Instance().Var(filepath.Join(cfgFolder, "gamemode_competitive.cfg"), "file"); err != nil {
		return false, nil
	}

	if err := gvalidator.Instance().Var(filepath.Join(cfgFolder, "gamemode_deathmatch.cfg"), "file"); err != nil {
		return false, nil
	}

	return true, nil

}
