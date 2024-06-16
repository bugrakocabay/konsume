package rabbitmq

import (
	"fmt"
	"log/slog"

	"github.com/bugrakocabay/konsume/pkg/config"
	"github.com/bugrakocabay/konsume/pkg/queue"

	amqp "github.com/rabbitmq/amqp091-go"
)

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

// NewConsumerFactory returns a new RabbitMQ consumer based on the provided configuration.
func NewConsumerFactory(cfg *config.ProviderConfig) (queue.MessageQueueConsumer, error) {
	return NewConsumer(cfg.AMQPConfig), nil
}

// Connect creates a connection to RabbitMQ and a channel
func (c *Consumer) Connect() error {
	slog.Debug("Attempting to connect to RabbitMQ", "host", c.config.Host, "port", c.config.Port)
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
func (c *Consumer) Consume(queueName string, handler func(msg []byte) error) error {
	slog.Debug("Starting to consume messages from RabbitMQ", "queueName", queueName)
	msgs, err := c.channel.Consume(
		queueName,
		"konsume", // Consumer tag - Identifier for the consumer
		false,     // Auto-Acknowledge, set to false for manual ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return fmt.Errorf("error starting to consume: %w", err)
	}

	go func() {
		for d := range msgs {
			err = handler(d.Body)
			if err != nil {
				slog.Error("Failed to process message sending to dead letter exchange", "message", string(d.Body), "error", err)
				err = d.Nack(false, false)
				if err != nil {
					slog.Error("Failed to nack the message", "error", err)
				}
			} else {
				err = d.Ack(false)
				if err != nil {
					slog.Error("Failed to ack the message", "error", err)
				}
			}
		}
	}()

	return nil
}

// Close closes the connection to RabbitMQ
func (c *Consumer) Close() error {
	slog.Debug("Closing RabbitMQ connection and channel")
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
	slog.Debug("RabbitMQ connection and channel closed successfully")
	return nil
}
