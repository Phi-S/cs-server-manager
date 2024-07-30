package event

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type DefaultPayload struct {
	TriggeredAtUtc time.Time
}

type Event struct {
	lock     sync.Mutex
	handlers map[uuid.UUID]func(DefaultPayload)
}

func (e *Event) Register(handler func(DefaultPayload)) uuid.UUID {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.handlers == nil {
		e.handlers = make(map[uuid.UUID]func(DefaultPayload))
	}

	uuid := uuid.New()
	e.handlers[uuid] = handler
	return uuid
}

func (e *Event) Deregister(handlerUuid uuid.UUID) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.handlers == nil {
		return
	}

	delete(e.handlers, handlerUuid)
}

func (e *Event) Trigger() {
	payload := DefaultPayload{
		TriggeredAtUtc: time.Now().UTC(),
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	if e.handlers == nil {
		return
	}

	wg := sync.WaitGroup{}
	for _, handler := range e.handlers {
		wg.Add(1)
		go func() {
			handler(payload)
			wg.Done()
		}()
	}
	wg.Wait()
}
