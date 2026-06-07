package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserCreatedEvent(t *testing.T) {
	event := &UserCreatedEvent{
		UserID:    1,
		UserName:  "testuser",
		Timestamp: time.Now(),
	}

	assert.Equal(t, "user.created", event.EventName())
	assert.Equal(t, int64(1), event.UserID)
	assert.False(t, event.OccurredOn().IsZero())
}

func TestUserUpdatedEvent(t *testing.T) {
	event := &UserUpdatedEvent{
		UserID:    2,
		UserName:  "updateduser",
		Timestamp: time.Now(),
	}

	assert.Equal(t, "user.updated", event.EventName())
	assert.Equal(t, int64(2), event.UserID)
	assert.False(t, event.OccurredOn().IsZero())
}

func TestChatMessageSentEvent(t *testing.T) {
	event := &ChatMessageSentEvent{
		SessionID: 10,
		MessageID: 20,
		Role:      "user",
		Content:   "hello",
		Timestamp: time.Now(),
	}

	assert.Equal(t, "chat.message.sent", event.EventName())
	assert.Equal(t, int64(10), event.SessionID)
	assert.Equal(t, int64(20), event.MessageID)
	assert.False(t, event.OccurredOn().IsZero())
}

func TestDocumentIndexedEvent(t *testing.T) {
	event := &DocumentIndexedEvent{
		DocID:      5,
		KBID:       3,
		ChunkCount: 100,
		Timestamp:  time.Now(),
	}

	assert.Equal(t, "document.indexed", event.EventName())
	assert.Equal(t, int64(5), event.DocID)
	assert.Equal(t, 100, event.ChunkCount)
	assert.False(t, event.OccurredOn().IsZero())
}

func TestDomainEventInterface(t *testing.T) {
	events := []DomainEvent{
		&UserCreatedEvent{UserID: 1, Timestamp: time.Now()},
		&UserUpdatedEvent{UserID: 2, Timestamp: time.Now()},
		&ChatMessageSentEvent{SessionID: 1, Timestamp: time.Now()},
		&DocumentIndexedEvent{DocID: 1, Timestamp: time.Now()},
	}

	for _, e := range events {
		assert.NotEmpty(t, e.EventName())
		assert.False(t, e.OccurredOn().IsZero())
	}
}
