package event

import "time"

// DomainEvent represents a domain event that has occurred in the system.
type DomainEvent interface {
	EventName() string
	OccurredOn() time.Time
}

// EventHandler is a function that handles a domain event.
type EventHandler func(event DomainEvent)

// EventBus defines the interface for publishing and subscribing to domain events.
type EventBus interface {
	Subscribe(eventName string, handler EventHandler)
	Publish(event DomainEvent)
}

// UserCreatedEvent is raised when a new user is created.
type UserCreatedEvent struct {
	UserID    int64
	UserName  string
	Timestamp time.Time
}

func (e *UserCreatedEvent) EventName() string     { return "user.created" }
func (e *UserCreatedEvent) OccurredOn() time.Time { return e.Timestamp }

// UserUpdatedEvent is raised when a user is updated.
type UserUpdatedEvent struct {
	UserID    int64
	UserName  string
	Timestamp time.Time
}

func (e *UserUpdatedEvent) EventName() string     { return "user.updated" }
func (e *UserUpdatedEvent) OccurredOn() time.Time { return e.Timestamp }

// ChatMessageSentEvent is raised when a chat message is sent.
type ChatMessageSentEvent struct {
	SessionID int64
	MessageID int64
	Role      string
	Content   string
	Timestamp time.Time
}

func (e *ChatMessageSentEvent) EventName() string     { return "chat.message.sent" }
func (e *ChatMessageSentEvent) OccurredOn() time.Time { return e.Timestamp }

// DocumentIndexedEvent is raised when a document has been indexed.
type DocumentIndexedEvent struct {
	DocID      int64
	KBID       int64
	ChunkCount int
	Timestamp  time.Time
}

func (e *DocumentIndexedEvent) EventName() string     { return "document.indexed" }
func (e *DocumentIndexedEvent) OccurredOn() time.Time { return e.Timestamp }
