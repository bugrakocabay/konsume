package runner

import (
	"context"
	"fmt"
	"konsume/pkg/util"
	"log"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"konsume/pkg/common"
	"konsume/pkg/config"
	"konsume/pkg/queue"
	"konsume/pkg/requester"
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
			if err := listenAndProcess(ctx, c, qc); err != nil {
				log.Printf("Failed to start consumer for queue %s: %s", qc.Name, err)
			}
		}(consumer, qCfg)
	}

	wg.Wait()
	return nil
}

// listenAndProcess consumes messages from the queue and processes them
func listenAndProcess(ctx context.Context, consumer queue.MessageQueueConsumer, qCfg *config.QueueConfig) error {
	if err := consumer.Connect(); err != nil {
		return err
	}

	return consumer.Consume(ctx, qCfg.Name, func(msg []byte) error {
		log.Printf("Received message from %s: %s", qCfg.Name, string(msg))
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
					return err
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
