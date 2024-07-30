package server

import "cs-server-controller/event"

func (s *Instance) OnOutput(handler func(event.PayloadWithData[string])) {
    s.onOutput.Register(handler)
}

func (s *Instance) OnServerStarting(handler func(event.DefaultPayload)) {
    s.onServerStarting.Register(handler)
}

func (s *Instance) OnServerStarted(handler func(event.DefaultPayload)) {
    s.onServerStarted.Register(handler)
}

func (s *Instance) OnServerStopped(handler func(event.DefaultPayload)) {
    s.onServerStopped.Register(handler)
}
func (s *Instance) OnServerCrashed(handler func(event.PayloadWithData[error])) {
    s.onServerCrashed.Register(handler)
}
