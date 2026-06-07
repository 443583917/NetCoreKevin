package eventbus

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kevin-ai/go-kevin/internal/domain/event"
)

type TestEvent struct {
	name string
}

func (e *TestEvent) EventName() string     { return "test.event" }
func (e *TestEvent) OccurredOn() time.Time { return time.Now() }

func TestInMemoryEventBusSubscribeAndPublish(t *testing.T) {
	bus := NewInMemoryEventBus()

	var received atomic.Bool
	handler := func(e event.DomainEvent) {
		received.Store(true)
	}

	bus.Subscribe("test.event", handler)

	testEvent := &TestEvent{name: "test"}
	bus.Publish(testEvent)

	time.Sleep(100 * time.Millisecond)
	assert.True(t, received.Load(), "handler should have been called")
}

func TestInMemoryEventBusMultipleHandlers(t *testing.T) {
	bus := NewInMemoryEventBus()

	var count atomic.Int32
	handler1 := func(e event.DomainEvent) {
		count.Add(1)
	}
	handler2 := func(e event.DomainEvent) {
		count.Add(1)
	}

	bus.Subscribe("test.event", handler1)
	bus.Subscribe("test.event", handler2)

	testEvent := &TestEvent{name: "test"}
	bus.Publish(testEvent)

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int32(2), count.Load(), "both handlers should have been called")
}

func TestInMemoryEventBusNoMatchingHandlers(t *testing.T) {
	bus := NewInMemoryEventBus()

	var called atomic.Bool
	handler := func(e event.DomainEvent) {
		called.Store(true)
	}

	bus.Subscribe("other.event", handler)

	testEvent := &TestEvent{name: "test"}
	bus.Publish(testEvent)

	time.Sleep(100 * time.Millisecond)
	assert.False(t, called.Load(), "handler should not be called for different event name")
}
