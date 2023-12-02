package runner

import (
	"fmt"
	"log"

	"konsume/pkg/config"
	"konsume/pkg/queue"
)

// StartConsumers starts the consumers for all queues
func StartConsumers(cfg *config.Config, consumers map[string]queue.MessageQueueConsumer) error {
	for _, qCfg := range cfg.Queues {
		consumer, ok := consumers[qCfg.Provider]
		if !ok {
			return fmt.Errorf("no consumer found for provider: %s", qCfg.Provider)
		}

		err := listen(consumer, qCfg)
		if err != nil {
			return fmt.Errorf("failed to start consumer for queue %s: %s", qCfg.Name, err)
		}
	}
	return nil
}

// listen consumes messages from the queue and processes them
func listen(consumer queue.MessageQueueConsumer, qCfg *config.QueueConfig) error {
	err := consumer.Connect()
	if err != nil {
		return err
	}

	return consumer.Consume(qCfg.Name, func(msg []byte) error {
		log.Printf("Received message from %s: %s", qCfg.Name, string(msg))
		return nil
	})
}
