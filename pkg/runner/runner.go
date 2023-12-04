package runner

import (
	"context"
	"fmt"
	"log"
	"sync"

	"konsume/pkg/config"
	"konsume/pkg/queue"
)

// StartConsumers starts the consumers for all queues
func StartConsumers(cfg *config.Config, consumers map[string]queue.MessageQueueConsumer) error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, qCfg := range cfg.Queues {
		consumer, ok := consumers[qCfg.Provider]
		if !ok {
			return fmt.Errorf("no consumer found for provider: %s", qCfg.Provider)
		}

		wg.Add(1)
		go func(c queue.MessageQueueConsumer, qc *config.QueueConfig) {
			defer wg.Done()
			if err := listen(ctx, c, qc); err != nil {
				log.Printf("Failed to start consumer for queue %s: %s", qc.Name, err)
			}
		}(consumer, qCfg)
	}

	wg.Wait()
	return nil
}

// listen consumes messages from the queue and processes them
func listen(ctx context.Context, consumer queue.MessageQueueConsumer, qCfg *config.QueueConfig) error {
	if err := consumer.Connect(); err != nil {
		return err
	}

	return consumer.Consume(ctx, qCfg.Name, func(msg []byte) error {
		log.Printf("Received message from %s: %s", qCfg.Name, string(msg))
		return nil
	})
}
