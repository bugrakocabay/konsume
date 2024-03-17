package e2e

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// connectToKafka establishes a connection to Kafka and returns the connection
func connectToKafka(broker, topic string) (*kafka.Conn, error) {
	c, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, 0)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// pushMessageToKafka publishes a message to the specified topic in Kafka
func pushMessageToKafka(c *kafka.Conn, body []byte) error {
	_, err := c.WriteMessages(
		kafka.Message{Value: body},
	)
	if err != nil {
		return err
	}

	return nil
}
