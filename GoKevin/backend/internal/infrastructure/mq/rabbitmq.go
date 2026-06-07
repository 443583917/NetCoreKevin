package mq

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ implements MessageQueue using RabbitMQ
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// RabbitMQConfig represents RabbitMQ configuration
type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	VHost    string
}

// NewRabbitMQ creates a new RabbitMQ client
func NewRabbitMQ(cfg RabbitMQConfig) (*RabbitMQ, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.VHost)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("connect to rabbitmq: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("open channel: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

// Publish publishes a message
func (r *RabbitMQ) Publish(ctx context.Context, topic string, msg *Message) error {
	// Declare queue
	_, err := r.channel.QueueDeclare(topic, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	// Publish message
	return r.channel.Publish("", topic, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        msg.Body,
		Headers:     msg.Headers,
	})
}

// Subscribe subscribes to a topic
func (r *RabbitMQ) Subscribe(ctx context.Context, topic string, handler Handler) error {
	// Declare queue
	_, err := r.channel.QueueDeclare(topic, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	// Consume messages
	msgs, err := r.channel.Consume(topic, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	// Process messages in background
	go func() {
		for msg := range msgs {
			m := &Message{
				Topic:   topic,
				Body:    msg.Body,
				Headers: msg.Headers,
			}

			if err := handler(ctx, m); err != nil {
				log.Printf("Handle message error: %v", err)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}()

	return nil
}

// Close closes the connection
func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
