package event

import (
	"sync"
	"time"
)

type DefaultPayload struct {
	Time time.Time
}

type PayloadWithData[T any] struct {
	DefaultPayload DefaultPayload
	Data           T
}

func GetPayload[T any](Data T) *PayloadWithData[T] {
	return &PayloadWithData[T]{
		DefaultPayload{Time: time.Now().UTC()},
		Data,
	}
}

func GetMessagePayload(Message string) *PayloadWithData[string] {
	return GetPayload[string](Message)
}

type EventWithData[T any] struct {
	mutext   sync.Mutex
	handlers []func(PayloadWithData[T])
}

func (e *EventWithData[T]) Register(handler func(PayloadWithData[T])) {
	e.mutext.Lock()
	defer e.mutext.Unlock()
	e.handlers = append(e.handlers, handler)
}

func (e *EventWithData[T]) Trigger(Data T) {
	e.mutext.Lock()
	defer e.mutext.Unlock()

	payload := PayloadWithData[T]{
		Data: Data,
	}

	for _, handler := range e.handlers {
		go handler(payload)
	}
}

type Event struct {
	mutext   sync.Mutex
	handlers []func(DefaultPayload)
}

func (e *Event) Register(handler func(DefaultPayload)) {
	e.mutext.Lock()
	defer e.mutext.Unlock()
	e.handlers = append(e.handlers, handler)
}

func (e *Event) Trigger() {
	e.mutext.Lock()
	defer e.mutext.Unlock()
	for _, handler := range e.handlers {
		go handler(DefaultPayload{Time: time.Now().UTC()})
	}
}
