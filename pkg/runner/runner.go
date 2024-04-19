package runner

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/bugrakocabay/konsume/pkg/config"
	"github.com/bugrakocabay/konsume/pkg/database"
	"github.com/bugrakocabay/konsume/pkg/queue"
)

// StartConsumers starts the consumers for all queues
func StartConsumers(cfg *config.Config, consumers map[string]queue.MessageQueueConsumer, providers map[string]*config.ProviderConfig, databases map[string]database.Database) error {
	var wg sync.WaitGroup

	for _, qCfg := range cfg.Queues {
		consumer, ok := consumers[qCfg.Provider]
		if !ok {
			return fmt.Errorf("no consumer found for provider: %s", qCfg.Provider)
		}

		providerCfg, ok := providers[qCfg.Provider]
		if !ok {
			return fmt.Errorf("no provider config found for provider: %s", qCfg.Provider)
		}

		wg.Add(1)
		go func(c queue.MessageQueueConsumer, qc *config.QueueConfig, pc *config.ProviderConfig) {
			defer wg.Done()
			if err := connectProviderWithRetry(c, pc); err != nil {
				return
			}
			if err := listenAndProcess(c, qc, cfg.Metrics, databases); err != nil {
				slog.Error("Failed to start consumer for", "queue", qc.Name, "error", err)
			}
		}(consumer, qCfg, providerCfg)
	}

	wg.Wait()
	return nil
}

// StopConsumers kills the connections of the consumers
func StopConsumers(consumers map[string]queue.MessageQueueConsumer, databases map[string]database.Database) {
	for _, c := range consumers {
		c.Close()
	}
	for _, db := range databases {
		db.Close()
	}
}

// connectProviderWithRetry tries to connect to the queue with the given consumer
func connectProviderWithRetry(consumer queue.MessageQueueConsumer, cfg *config.ProviderConfig) error {
	var err error
	err = consumer.Connect()
	if err != nil {
		slog.Error("Failed to connect to queue", "error", err)
		if cfg.Retry > 0 {
			for i := 1; i <= cfg.Retry; i++ {
				time.Sleep(time.Duration(5) * time.Second)
				slog.Info("Retrying to connect", "retry", i)
				err = consumer.Connect()
				if err == nil {
					break
				}
			}
		}
		return err
	}
	return err
}
