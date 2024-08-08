package game_events

import (
	"cs-server-manager/event"
	"log/slog"
	"regexp"
	"strconv"
)

type PlayerConnected struct {
	Name string
	Id   uint8
	Ip   string
	Port uint16
}

type Instance struct {
	onMapChanged         event.InstanceWithData[string]
	onPlayerConnected    event.InstanceWithData[PlayerConnected]
	onPlayerDisconnected event.InstanceWithData[string]
}

func (g *Instance) OnMapChanged(handler func(p event.PayloadWithData[string])) {
	g.onMapChanged.Register(handler)
}

func (g *Instance) OnPlayerConnected(handler func(p event.PayloadWithData[PlayerConnected])) {
	g.onPlayerConnected.Register(handler)
}

func (g *Instance) OnPlayerDisconnected(handler func(p event.PayloadWithData[string])) {
	g.onPlayerDisconnected.Register(handler)
}

func (g *Instance) DetectGameEvent(msg string) {
	g.detectMapChange(msg)
	g.detectPlayerConnected(msg)
	g.detectPlayerDisconnected(msg)
}

func (g *Instance) detectMapChange(msg string) {
	regexExpr := `Host activate: Changelevel \((.+)\)`
	r, err := regexp.Compile(regexExpr)
	if err != nil {
		slog.Error("regex is not valid")
		return
	}

	groups := r.FindStringSubmatch(msg)
	if len(groups) != 2 {
		return
	}

	g.onMapChanged.Trigger(groups[1])
}

func (g *Instance) detectPlayerConnected(msg string) {
	regexExpr := `CServerSideClientBase::Connect\( name='(.+)', userid=(\d), fake=\d, chan->addr=(\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}):(\d{1,5}) \)`
	r, err := regexp.Compile(regexExpr)
	if err != nil {
		slog.Error("regex is not valid", "regex_expr", regexExpr)
		return
	}

	groups := r.FindStringSubmatch(msg)
	if len(groups) != 5 {
		return
	}

	userId, err := strconv.ParseUint(groups[2], 10, 8)
	if err != nil {
		slog.Error("failed to detect player connection. Failed to parse user id", "user_id", groups[1], "error", err)
		return
	}

	port, err := strconv.ParseUint(groups[4], 10, 16)
	if err != nil {
		slog.Error("failed to detect player connection. Failed to parse port", "port", groups[2], "error", err)
	}

	g.onPlayerConnected.Trigger(PlayerConnected{
		Name: groups[1],
		Id:   uint8(userId),
		Ip:   groups[3],
		Port: uint16(port),
	})
}

func (g *Instance) detectPlayerDisconnected(msg string) {
	regexExpr := `SV:  Dropped client '(.+)' from server\(.+`
	r, err := regexp.Compile(regexExpr)
	if err != nil {
		slog.Error("regex is not valid", "regex_expr", regexExpr)
		return
	}

	groups := r.FindStringSubmatch(msg)
	if len(groups) != 2 {
		return
	}

	g.onPlayerDisconnected.Trigger(groups[1])
}
