package mq

import "context"

// Message represents a message
type Message struct {
	Topic   string
	Body    []byte
	Headers map[string]interface{}
}

// Handler is a message handler function
type Handler func(ctx context.Context, msg *Message) error

// MessageQueue is the interface for message queue
type MessageQueue interface {
	// Publish publishes a message to a topic
	Publish(ctx context.Context, topic string, msg *Message) error

	// Subscribe subscribes to a topic
	Subscribe(ctx context.Context, topic string, handler Handler) error

	// Close closes the connection
	Close() error
}
