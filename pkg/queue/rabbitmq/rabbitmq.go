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
	dlExchange := "dlx.exchange"
	dlQueue := queueName + ".dlx"

	err := c.channel.ExchangeDeclare(
		dlExchange, // name of the exchange
		"direct",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("error declaring the DLX: %w", err)
	}

	_, err = c.channel.QueueDeclare(
		dlQueue, // name of the queue
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("error declaring the DLQ: %w", err)
	}

	err = c.channel.QueueBind(
		dlQueue,    // queue name
		queueName,  // routing key
		dlExchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error binding the DLQ: %w", err)
	}

	args := amqp.Table{
		"x-dead-letter-exchange":    dlExchange,
		"x-dead-letter-routing-key": queueName,
	}

	_, err = c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments for dead-lettering
	)
	if err != nil {
		return fmt.Errorf("error declaring the main queue: %w", err)
	}

	slog.Debug("Starting to consume messages from RabbitMQ", "queueName", queueName)

	return c.startConsuming(queueName, handler)
}

func (c *Consumer) startConsuming(queueName string, handler func(msg []byte) error) error {
	msgs, err := c.channel.Consume(
		queueName,
		"",    // Consumer tag - Identifier for the consumer
		false, // Auto-Acknowledge, set to false for manual ack
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return fmt.Errorf("error starting to consume: %w", err)
	}

	go func() {
		for d := range msgs {
			err = handler(d.Body)
			if err != nil {
				slog.Error("Failed to process message sending to dead letter exchange", "message", string(d.Body), "error", err)
				d.Nack(false, false)
			} else {
				d.Ack(false)
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
