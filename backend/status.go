package main

import (
    "cs-server-controller/event"
    "cs-server-controller/server"
    "cs-server-controller/steamcmd"
    "encoding/json"
    "github.com/google/uuid"
    "sync"
)

func NewStatus(hostname string, maxPlayerCount uint8, startMap string) *Status {
    return &Status{
        internalStatus: InternalStatus{
            Hostname:       hostname,
            Server:         ServerStatusStopped,
            Steamcmd:       SteamcmdStatusStopped,
            PlayerCount:    0,
            MaxPlayerCount: maxPlayerCount,
            Map:            startMap,
        },
    }
}

type ServerStatus string

const (
    ServerStatusStarted               = "server-status-started"
    ServerStatusStarting              = "server-status-starting"
    ServerStatusStopped  ServerStatus = "server-status-stopped"
    ServerStatusStopping              = "server-status-stopping"
)

type SteamcmdStatus string

const (
    SteamcmdStatusStopped  SteamcmdStatus = "steamcmd-status-stopped"
    SteamcmdStatusUpdating                = "steamcmd-status-updating"
)

type InternalStatus struct {
    Hostname       string         `json:"hostname"`
    Server         ServerStatus   `json:"server"`
    Steamcmd       SteamcmdStatus `json:"steamcmd"`
    PlayerCount    uint8          `json:"player-count"`
    MaxPlayerCount uint8          `json:"max-player-count"`
    Map            string         `json:"map"`
}

type Status struct {
    internalStatus  InternalStatus
    lock            sync.RWMutex
    onStatusChanged event.Instance
}

func (s *Status) OnStatusChanged(handler func(payload event.DefaultPayload)) uuid.UUID {
    return s.onStatusChanged.Register(handler)
}

func (s *Status) Status() InternalStatus {
    s.lock.Lock()
    defer s.lock.Unlock()
    return s.internalStatus
}

func (s *Status) Json() ([]byte, error) {
    s.lock.RLock()
    defer s.lock.RUnlock()

    result, err := json.Marshal(s.internalStatus)
    if err != nil {
        return nil, err
    }

    return result, nil
}

func (s *Status) ChangeStatusOnEvents(serverInstance *server.Instance, steamcmdInstance *steamcmd.Instance) {
    serverInstance.OnOutput(func(pwd event.PayloadWithData[string]) {
        // TODO: check if player joins regex 
        //slog.Debug("serverInstance", "event", "onOutput", "triggeredAtUtc", pwd.TriggeredAtUtc, "output", pwd.Data)
        //    fmt.Println(pwd.TriggeredAtUtc.Format(time.RFC3339Nano) + " | SERVER: " + pwd.Data)
    })

    serverInstance.OnStarting(func(dp event.DefaultPayload) {
        s.lock.Lock()
        s.internalStatus.Server = ServerStatusStarting
        s.lock.Unlock()
        s.onStatusChanged.Trigger()
    })

    serverInstance.OnStarted(func(e event.PayloadWithData[server.StartParameters]) {
        s.lock.Lock()
        s.internalStatus.Server = ServerStatusStarted
        s.internalStatus.Hostname = e.Data.Hostname
        s.internalStatus.MaxPlayerCount = e.Data.MaxPlayers
        s.internalStatus.Map = e.Data.StartMap
        s.lock.Unlock()
        s.onStatusChanged.Trigger()
    })

    serverInstance.OnCrashed(func(pwd event.PayloadWithData[error]) {
        s.lock.Lock()
        s.internalStatus.Server = ServerStatusStopped
        s.lock.Unlock()
        s.onStatusChanged.Trigger()
    })

    serverInstance.OnStopped(func(dp event.DefaultPayload) {
        s.lock.Lock()
        s.internalStatus.Server = ServerStatusStopped
        s.lock.Unlock()
        s.onStatusChanged.Trigger()
    })

    //steamcmdInstance.OnOutput(func(p event.PayloadWithData[string]) {
    //        fmt.Println(p.TriggeredAtUtc.String() + " | steamcmdInstance: " + p.Data)
    //slog.Debug("steamcmdInstance", "event", "onOutput", "triggeredAtUtc", p.TriggeredAtUtc, "output", p.Data)
    //})

    steamcmdInstance.OnStarted(func(p event.DefaultPayload) {
        s.lock.Lock()
        s.internalStatus.Steamcmd = SteamcmdStatusUpdating
        s.lock.Unlock()
        s.onStatusChanged.Trigger()
    })

    steamcmdInstance.OnFinished(func(p event.DefaultPayload) {
        s.lock.Lock()
        s.internalStatus.Steamcmd = SteamcmdStatusStopped
        s.lock.Unlock()
        s.onStatusChanged.Trigger()
    })

    steamcmdInstance.OnCancelled(func(p event.DefaultPayload) {
        s.lock.Lock()
        s.internalStatus.Steamcmd = SteamcmdStatusStopped
        s.lock.Unlock()
        s.onStatusChanged.Trigger()
    })

    steamcmdInstance.OnFailed(func(p event.PayloadWithData[error]) {
        s.lock.Lock()
        s.internalStatus.Steamcmd = SteamcmdStatusStopped
        s.lock.Unlock()
        s.onStatusChanged.Trigger()
    })
}
