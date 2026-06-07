package eventbus

import (
	"sync"

	"github.com/kevin-ai/go-kevin/internal/domain/event"
)

// InMemoryEventBus is a thread-safe in-memory implementation of event.EventBus.
type InMemoryEventBus struct {
	handlers map[string][]event.EventHandler
	mu       sync.RWMutex
}

// NewInMemoryEventBus creates a new InMemoryEventBus.
func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]event.EventHandler),
	}
}

// Subscribe registers a handler for the given event name.
func (b *InMemoryEventBus) Subscribe(eventName string, handler event.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

// Publish dispatches an event to all registered handlers concurrently.
func (b *InMemoryEventBus) Publish(e event.DomainEvent) {
	b.mu.RLock()
	handlers, ok := b.handlers[e.EventName()]
	b.mu.RUnlock()

	if ok {
		for _, handler := range handlers {
			go handler(e)
		}
	}
}
