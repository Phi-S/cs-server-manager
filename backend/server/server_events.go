package server

import "cs-server-manager/event"

func (s *Instance) OnOutput(handler func(event.PayloadWithData[string])) {
	s.onOutput.Register(handler)
}

func (s *Instance) OnStarting(handler func(event.DefaultPayload)) {
	s.onStarting.Register(handler)
}

func (s *Instance) OnStarted(handler func(event.PayloadWithData[StartParameters])) {
	s.onStarted.Register(handler)
}

func (s *Instance) OnStopped(handler func(event.DefaultPayload)) {
	s.onStopped.Register(handler)
}
func (s *Instance) OnCrashed(handler func(event.PayloadWithData[error])) {
	s.onCrashed.Register(handler)
}
