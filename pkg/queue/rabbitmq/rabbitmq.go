package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"

	"konsume/pkg/config"
	"konsume/pkg/queue"

	amqp "github.com/rabbitmq/amqp091-go"
)

var _ queue.MessageQueueConsumer = (*Consumer)(nil)

// Consumer is the implementation of the MessageQueueConsumer interface for RabbitMQ
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *config.AMQPConfig
}

// NewConsumer creates a new RabbitMQ consumer
func NewConsumer(cfg *config.AMQPConfig) *Consumer {
	return &Consumer{
		config: cfg,
	}
}

// Connect creates a connection to RabbitMQ and a channel
func (c *Consumer) Connect() error {
	var err error
	cfg := c.config
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	c.conn, err = amqp.Dial(connectionString)
	if err != nil {
		return err
	}
	c.channel, err = c.conn.Channel()
	if err != nil {
		return err
	}
	slog.Info("Connected to RabbitMQ", "host", cfg.Host, "port", cfg.Port)

	return nil
}

// Consume consumes messages from RabbitMQ
func (c *Consumer) Consume(ctx context.Context, queueName string, handler func(msg []byte) error) error {
	_, err := c.channel.QueueDeclare(
		queueName, // Name of the queue
		true,      // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return err
	}
	msgs, err := c.channel.Consume(
		queueName,
		"",    // Consumer tag - Identifier for the consumer
		true,  // Auto-Acknowledge
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			if err = handler(d.Body); err != nil {
				slog.Error("Failed to process message", "error", err)
			}
		}
	}()

	return nil
}

// Close closes the connection to RabbitMQ
func (c *Consumer) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
