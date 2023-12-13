package kafka

import (
	"context"
	"github.com/bugrakocabay/konsume/pkg/config"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

// Consumer is the implementation of the MessageQueueConsumer interface for Kafka
type Consumer struct {
	config *config.KafkaConfig
	conn   *kafka.Conn
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.KafkaConfig) *Consumer {
	return &Consumer{
		config: cfg,
	}
}

// Connect creates a connection to Kafka
func (c *Consumer) Connect(ctx context.Context) error {
	slog.Debug("Attempting to connect to Kafka", "brokers", c.config.Brokers, "topic", c.config.Topic)
	var err error
	c.conn, err = kafka.DialLeader(context.Background(), "tcp", c.config.Brokers[0], c.config.Topic, 0)
	if err != nil {
		return err
	}

	return nil
}

// Consume consumes messages from Kafka
func (c *Consumer) Consume(ctx context.Context, queueName string, handler func(msg []byte) error) error {
	slog.Debug("Starting to consume messages from Kafka", "topic", queueName)
	for {
		msg, err := c.conn.ReadMessage(10e6)
		if err != nil {
			slog.Error("Failed to read message from Kafka", "error", err)
			return err
		}
		if err = handler(msg.Value); err != nil {
			slog.Error("Failed to process message", "error", err)
		}
	}
}

// Close closes the connection to Kafka
func (c *Consumer) Close() error {
	slog.Debug("Closing connection to Kafka")
	err := c.conn.Close()
	if err != nil {
		return err
	}
	slog.Debug("Kafka connection closed successfully")
	return nil
}
