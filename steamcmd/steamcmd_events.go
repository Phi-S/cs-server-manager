package steamcmd

import "cs-server-controller/event"

func (s *SteamcmdInstance) OnOutput(handler func(event.PayloadWithData[string])) {
	s.onOutput.Register(handler)
}

func (s *SteamcmdInstance) OnStarted(handler func(event.DefaultPayload)) {
	s.onStarted.Register(handler)
}

func (s *SteamcmdInstance) OnFinished(handler func(event.DefaultPayload)) {
	s.onFinished.Register(handler)
}

func (s *SteamcmdInstance) OnCancelled(handler func(event.DefaultPayload)) {
	s.onCancelled.Register(handler)
}
func (s *SteamcmdInstance) OnFailed(handler func(event.PayloadWithData[error])) {
	s.onFailed.Register(handler)
}
