package game_events_test

import (
	"testing"

	"github.com/Phi-S/cs-server-manager/event"
	"github.com/Phi-S/cs-server-manager/game_events"
)

func TestDetectGameEvent_ChangeMap_Ok(t *testing.T) {
	msg := "Host activate: Changelevel (de_anubis)"

	ge := game_events.Instance{}
	eventTriggered := false
	ge.OnMapChanged(func(p event.PayloadWithData[string]) { eventTriggered = true })

	ge.DetectGameEvent(msg)
	if eventTriggered == false {
		t.Fatalf("test failed. OnMapChanged not triggered with message %v", msg)
	}
}

func TestDetectGameEvent_ChangeMap_Fail1(t *testing.T) {
	msg := "Host activate: keklevel (de_anubis)"

	ge := game_events.Instance{}
	eventTriggered := false
	ge.OnMapChanged(func(p event.PayloadWithData[string]) { eventTriggered = true })

	ge.DetectGameEvent(msg)
	if eventTriggered == true {
		t.Fatalf("test failed. OnMapChanged triggered with invalid message %v", msg)
	}
}

func TestDetectGameEvent_ChangeMap_Fail2(t *testing.T) {
	msg := "level (de_anubis)"

	ge := game_events.Instance{}
	eventTriggered := false
	ge.OnMapChanged(func(p event.PayloadWithData[string]) { eventTriggered = true })

	ge.DetectGameEvent(msg)
	if eventTriggered == true {
		t.Fatalf("test failed. OnMapChanged triggered with invalid message %v", msg)
	}
}

func TestDetectGameEvent_PlayerConnected_Ok(t *testing.T) {
	msg := `CServerSideClientBase::Connect( name='PhiS > :) < --L', userid=2, fake=0, chan->addr=127.0.0.1:51018 )`

	ge := game_events.Instance{}
	eventTriggered := false
	ge.OnPlayerConnected(func(p event.PayloadWithData[game_events.PlayerConnected]) { eventTriggered = true })

	ge.DetectGameEvent(msg)
	if eventTriggered == false {
		t.Fatalf("test failed. OnPlayerConnected not triggered: Message %v", msg)
	}
}

func TestDetectGameEvent_PlayerDisconnected_Ok(t *testing.T) {
	msg := `SV:  Dropped client 'PhiS > :) < --L' from server(59): NETWORK_DISCONNECT_EXITING`

	ge := game_events.Instance{}
	eventTriggered := false
	ge.OnPlayerDisconnected(func(p event.PayloadWithData[string]) { eventTriggered = true })

	ge.DetectGameEvent(msg)
	if eventTriggered == false {
		t.Fatalf("test failed. OnPlayerConnected not triggered. Message: %v", msg)
	}
}
