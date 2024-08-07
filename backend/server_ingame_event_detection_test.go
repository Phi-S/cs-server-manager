package main

import (
	"cs-server-manager/status"
	"testing"
)

func Test_checkForServerEvents_ChangeMap_Ok(t *testing.T) {
	s := status.NewStatus("test", 0, "test")

	msg := "Host activate: Changelevel (de_anubis)"
	DetectMapChange(s, msg)

	if s.Status().Map != "de_anubis" {
		t.Fatalf("test failed with message %v", msg)
	}
}

func Test_checkForServerEvents_ChangeMap_Fail1(t *testing.T) {
	s := status.NewStatus("test", 0, "test")

	msg := "Host activate: keklevel (de_anubis)"
	DetectMapChange(s, msg)

	if s.Status().Map != "test" {
		t.Fatalf("test failed with message %v", msg)
	}
}

func Test_checkForServerEvents_ChangeMap_Fail2(t *testing.T) {
	s := status.NewStatus("test", 0, "test")

	msg := "level (de_anubis)"
	DetectMapChange(s, msg)

	if s.Status().Map != "test" {
		t.Fatalf("test failed with message %v", msg)
	}
}

func Test_checkForServerEvents_PlayerConnected_Ok(t *testing.T) {
	s := status.NewStatus("test", 0, "test")

	msg := `CServerSideClientBase::Connect( name='PhiS > :) < --L', userid=2, fake=0, chan->addr=127.0.0.1:51018 )`
	DetectPlayerConnected(s, msg)

	playerCount := s.Status().PlayerCount
	if playerCount != 1 {
		t.Fatalf("test failed. player count: %v | message %v", playerCount, msg)
	}
}

func Test_checkForServerEvents_PlayerDisconnected_Ok(t *testing.T) {
	s := status.NewStatus("test", 0, "test")

	s.Update(func(sts *status.InternalStatus) { sts.PlayerCount = 2 })

	msg := `SV: Dropped client 'PhiS > :) < --L' from server(59): NETWORK_DISCONNECT_EXITING`
	DetectPlayerDisconnected(s, msg)

	playerCount := s.Status().PlayerCount
	if playerCount != 1 {
		t.Fatalf("test failed. player count: %v | message %v", playerCount, msg)
	}
}
