package events

import (
	"sync"
)

type Handler func(Event)

type EventBus struct {
	handlers map[Action][]Handler
	mu       sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[Action][]Handler),
		mu:       sync.RWMutex{},
	}
}
func (b *EventBus) Subscribe(action Action, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[action] = append(b.handlers[action], handler)
}

func (b *EventBus) Publish(event Event) {
	b.mu.RLock()
	handlers, exists := b.handlers[event.Action]
	defer b.mu.RUnlock()
	if !exists {
		return
	}
	for _, h := range handlers {
		go h(event)
	}

}
