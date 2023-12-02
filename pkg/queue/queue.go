package queue

// MessageQueueConsumer is the interface that each message queue producer should implement
type MessageQueueConsumer interface {
	Connect() error
	Consume(queueName string, handler func(msg []byte) error) error
	Close() error
}
