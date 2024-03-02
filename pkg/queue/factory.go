package queue

import "github.com/bugrakocabay/konsume/pkg/config"

// Factory is a function type that creates a new MessageQueueConsumer based on the provided configuration.
type Factory func(*config.ProviderConfig) (MessageQueueConsumer, error)
