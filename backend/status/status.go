package status

import (
	"cs-server-manager/event"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
)

func NewStatus(hostname string, ip string, port string, password string, maxPlayerCount uint8, startMap string) *Status {
	instance := Status{
		internalStatus: &InternalStatus{
			Hostname:       hostname,
			Server:         ServerStatusStopped,
			Steamcmd:       SteamcmdStatusStopped,
			PlayerCount:    0,
			MaxPlayerCount: maxPlayerCount,
			Map:            startMap,
			Ip:             ip,
			Port:           port,
			Password:       password,
		},
	}

	return &instance
}

type ServerStatus string

const (
	ServerStatusStarted  ServerStatus = "server-status-started"
	ServerStatusStarting ServerStatus = "server-status-starting"
	ServerStatusStopped  ServerStatus = "server-status-stopped"
	ServerStatusStopping ServerStatus = "server-status-stopping"
)

type SteamcmdStatus string

const (
	SteamcmdStatusStopped  SteamcmdStatus = "steamcmd-status-stopped"
	SteamcmdStatusUpdating SteamcmdStatus = "steamcmd-status-updating"
)

type InternalStatus struct {
	Hostname       string         `json:"hostname"`
	Server         ServerStatus   `json:"server"`
	Steamcmd       SteamcmdStatus `json:"steamcmd"`
	PlayerCount    uint8          `json:"player_count"`
	MaxPlayerCount uint8          `json:"max_player_count"`
	Map            string         `json:"map"`
	Ip             string         `json:"ip"`
	Port           string         `json:"port"`
	Password       string         `json:"password"`
}

type Status struct {
	internalStatus  *InternalStatus
	lock            sync.RWMutex
	onStatusChanged event.InstanceWithData[InternalStatus]
}

func (s *Status) OnStatusChanged(handler func(payload event.PayloadWithData[InternalStatus])) uuid.UUID {
	return s.onStatusChanged.Register(handler)
}

func (s *Status) Status() InternalStatus {
	s.lock.Lock()
	defer s.lock.Unlock()
	return *s.internalStatus
}

func (s *Status) Json() ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	result, err := json.Marshal(*s.internalStatus)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Status) Update(ch func(internalStatus *InternalStatus)) {
	var localCopy InternalStatus

	s.lock.Lock()
	ch(s.internalStatus)
	localCopy = *s.internalStatus
	s.lock.Unlock()

	s.onStatusChanged.Trigger(localCopy)
}
