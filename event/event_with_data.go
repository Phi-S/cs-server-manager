package event

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type PayloadWithData[T any] struct {
	DefaultPayload
	Data T
}

type EventWithData[T any] struct {
	lock     sync.Mutex
	handlers map[uuid.UUID]func(PayloadWithData[T])
}

func (e *EventWithData[T]) Register(handler func(PayloadWithData[T])) uuid.UUID {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.handlers == nil {
		e.handlers = make(map[uuid.UUID]func(PayloadWithData[T]))
	}

	uuid := uuid.New()
	e.handlers[uuid] = handler
	return uuid
}

func (e *EventWithData[T]) Deregister(handlerUuid uuid.UUID) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.handlers == nil {
		return
	}

	delete(e.handlers, handlerUuid)
}

func (e *EventWithData[T]) Trigger(dataIn T) {
	p := PayloadWithData[T]{
		DefaultPayload: DefaultPayload{
			TriggeredAtUtc: time.Now().UTC(),
		},
		Data: dataIn,
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	if e.handlers == nil {
		return
	}

	for _, handler := range e.handlers {
		handler(p)
	}
}
