package server

import "cs-server-controller/event"

func (s *ServerInstance) OnOutput(handler func(event.PayloadWithData[string])) {
	s.onOutput.Register(handler)
}

func (s *ServerInstance) OnServerStarting(handler func(event.DefaultPayload)) {
	s.onServerStarting.Register(handler)
}

func (s *ServerInstance) OnServerStarted(handler func(event.DefaultPayload)) {
	s.onServerStarted.Register(handler)
}

func (s *ServerInstance) OnServerStopped(handler func(event.DefaultPayload)) {
	s.onServerStopped.Register(handler)
}
func (s *ServerInstance) OnServerCrashed(handler func(event.PayloadWithData[error])) {
	s.onServerCrashed.Register(handler)
}
