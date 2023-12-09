package queue

import "context"

// MessageQueueConsumer is the interface that each message queue producer should implement
type MessageQueueConsumer interface {
	Connect(ctx context.Context) error
	Consume(ctx context.Context, queueName string, handler func(msg []byte) error) error
	Close() error
}
