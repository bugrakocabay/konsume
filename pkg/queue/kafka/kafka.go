package kafka

import (
	"context"
	"log/slog"

	"github.com/bugrakocabay/konsume/pkg/config"

	"github.com/segmentio/kafka-go"
)

// Consumer is the implementation of the MessageQueueConsumer interface for Kafka
type Consumer struct {
	config *config.KafkaConfig
	reader *kafka.Reader
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.KafkaConfig) *Consumer {
	return &Consumer{
		config: cfg,
	}
}

// Connect creates a connection to Kafka
func (c *Consumer) Connect() error {
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.config.Brokers,
		GroupID: c.config.Group,
		Topic:   c.config.Topic,
	})
	slog.Info("Connected to Kafka", "brokers", c.config.Brokers, "group", c.config.Group, "topic", c.config.Topic)
	return nil
}

// Consume consumes messages from Kafka
func (c *Consumer) Consume(ctx context.Context, queueName string, handler func(msg []byte) error) error {
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return err // Handle error appropriately
		}

		if err = handler(m.Value); err != nil {
			slog.Error("Failed to process message", "error", err)
		}

		if err = c.reader.CommitMessages(ctx, m); err != nil {
			return err
		}
	}
}

// Close closes the connection to Kafka
func (c *Consumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}
