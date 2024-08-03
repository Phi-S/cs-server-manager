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

type InstanceWithData[T any] struct {
	lock     sync.Mutex
	handlers map[uuid.UUID]func(PayloadWithData[T])
}

func (e *InstanceWithData[T]) Register(handler func(PayloadWithData[T])) uuid.UUID {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.handlers == nil {
		e.handlers = make(map[uuid.UUID]func(PayloadWithData[T]))
	}

	newUUID := uuid.New()
	e.handlers[newUUID] = handler
	return newUUID
}

func (e *InstanceWithData[T]) Deregister(handlerUuid uuid.UUID) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.handlers == nil {
		return
	}

	delete(e.handlers, handlerUuid)
}

func (e *InstanceWithData[T]) Trigger(dataIn T) {
	p := PayloadWithData[T]{
		DefaultPayload: DefaultPayload{
			TriggeredAtUtc: time.Now().UTC(),
		},
		Data: dataIn,
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	if e.handlers == nil || len(e.handlers) == 0 {
		return
	}

	wg := sync.WaitGroup{}
	for _, handler := range e.handlers {
		wg.Add(1)
		go func() {
			handler(p)
			wg.Done()
		}()
	}
	wg.Wait()
}
