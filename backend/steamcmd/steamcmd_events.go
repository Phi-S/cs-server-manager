package steamcmd

import "cs-server-controller/event"

func (s *Instance) OnOutput(handler func(event.PayloadWithData[string])) {
	s.onOutput.Register(handler)
}

func (s *Instance) OnStarted(handler func(event.DefaultPayload)) {
	s.onStarted.Register(handler)
}

func (s *Instance) OnFinished(handler func(event.DefaultPayload)) {
	s.onFinished.Register(handler)
}

func (s *Instance) OnCancelled(handler func(event.DefaultPayload)) {
	s.onCancelled.Register(handler)
}
func (s *Instance) OnFailed(handler func(event.PayloadWithData[error])) {
	s.onFailed.Register(handler)
}
