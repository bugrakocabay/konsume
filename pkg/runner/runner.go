package runner

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/bugrakocabay/konsume/pkg/common"
	"github.com/bugrakocabay/konsume/pkg/config"
	"github.com/bugrakocabay/konsume/pkg/queue"
	"github.com/bugrakocabay/konsume/pkg/requester"
	"github.com/bugrakocabay/konsume/pkg/util"
)

// StartConsumers starts the consumers for all queues
func StartConsumers(cfg *config.Config, consumers map[string]queue.MessageQueueConsumer, providers map[string]*config.ProviderConfig) error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
			if err := connectWithRetry(ctx, c, pc); err != nil {
				return
			}
			if err := listenAndProcess(ctx, c, qc); err != nil {
				log.Printf("Failed to start consumer for queue %s: %s", qc.Name, err)
			}
		}(consumer, qCfg, providerCfg)
	}

	wg.Wait()
	return nil
}

// connectWithRetry tries to connect to the queue with the given consumer
func connectWithRetry(ctx context.Context, consumer queue.MessageQueueConsumer, cfg *config.ProviderConfig) error {
	var err error
	err = consumer.Connect(ctx)
	if err != nil {
		slog.Error("Failed to connect to queue", "error", err)
		if cfg.Retry > 0 {
			for i := 1; i <= cfg.Retry; i++ {
				time.Sleep(time.Duration(5) * time.Second)
				slog.Info("Retrying to connect", "retry", i)
				err = consumer.Connect(ctx)
				if err == nil {
					break
				}
			}
		}
		return err
	}
	return err
}

// listenAndProcess consumes messages from the queue and processes them
func listenAndProcess(ctx context.Context, consumer queue.MessageQueueConsumer, qCfg *config.QueueConfig) error {
	return consumer.Consume(ctx, qCfg.Name, func(msg []byte) error {
		slog.Info("Received a message", "queue", qCfg.Name, "message", string(msg))
		var (
			messageData map[string]interface{}
			err         error
			body        []byte
		)
		for _, rCfg := range qCfg.Routes {
			if len(rCfg.Body) > 0 {
				messageData, err = util.ParseJSONToMap(msg)
				if err != nil {
					slog.Error("Failed to parse message", "error", err)
				}
				body, err = util.ProcessTemplate(rCfg.Body, messageData)
				if err != nil {
					slog.Error("Failed to process template", "error", err)
					body = msg
				}
			} else {
				body = msg
			}
			rqstr := requester.NewRequester(rCfg.URL, rCfg.Method, body, rCfg.Headers)
			sendRequestWithStrategy(qCfg, rCfg, body, rqstr)
		}
		return nil
	})
}

// sendRequestWithStrategy sends the request to the given endpoint and makes use of the given strategy
func sendRequestWithStrategy(qCfg *config.QueueConfig, rCfg *config.RouteConfig, msg []byte, requester requester.HTTPRequester) {
	resp := requester.SendRequest()
	slog.Info("Received a response from", "route", rCfg.Name, "status", resp.StatusCode)
	retry := qCfg.Retry
	if retry != nil && resp.StatusCode >= retry.ThresholdStatus {
		for i := 1; i <= retry.MaxRetries; i++ {
			slog.Info("Retrying", "route", rCfg.Name, "strategy", retry.Strategy, "retry", i)
			switch retry.Strategy {
			case common.RetryStrategyFixed:
				time.Sleep(retry.Interval)
			case common.RetryStrategyExpo:
				time.Sleep(retry.Interval * time.Duration(i))
			case common.RetryStrategyRand:
				time.Sleep(time.Duration(rand.Intn(int(retry.Interval))))
			default:
				slog.Error("Invalid retry strategy", "strategy", retry.Strategy)
				break
			}
			resp = requester.SendRequest()
			slog.Info("Received a response from", "route", rCfg.Name, "status", resp.StatusCode)
			if resp.StatusCode < retry.ThresholdStatus {
				break
			}
		}
	}
}
