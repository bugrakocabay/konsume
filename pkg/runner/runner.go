package runner

import (
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/bugrakocabay/konsume/pkg/common"
	"github.com/bugrakocabay/konsume/pkg/config"
	"github.com/bugrakocabay/konsume/pkg/metrics"
	"github.com/bugrakocabay/konsume/pkg/queue"
	"github.com/bugrakocabay/konsume/pkg/requester"
)

// StartConsumers starts the consumers for all queues
func StartConsumers(cfg *config.Config, consumers map[string]queue.MessageQueueConsumer, providers map[string]*config.ProviderConfig) error {
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
			if err := connectWithRetry(c, pc); err != nil {
				return
			}
			if err := listenAndProcess(c, qc, cfg.Metrics); err != nil {
				log.Printf("Failed to start consumer for queue %s: %s", qc.Name, err)
			}
		}(consumer, qCfg, providerCfg)
	}

	wg.Wait()
	return nil
}

// StopConsumers kills the connections of the consumers
func StopConsumers(consumers map[string]queue.MessageQueueConsumer) {
	for _, c := range consumers {
		c.Close()
	}
}

// connectWithRetry tries to connect to the queue with the given consumer
func connectWithRetry(consumer queue.MessageQueueConsumer, cfg *config.ProviderConfig) error {
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

// sendRequestWithStrategy attempts to send an HTTP request and retries based on the provided configuration
func sendRequestWithStrategy(qCfg *config.QueueConfig, rCfg *config.RouteConfig, mCfg *config.MetricsConfig, requester requester.HTTPRequester) {
	resp, err := requester.SendRequest(mCfg, rCfg.Timeout)
	if err != nil || shouldRetry(resp, qCfg.Retry) {
		if resp != nil && resp.StatusCode != 0 {
			slog.Info("Received a response from", "route", rCfg.Name, "status", resp.StatusCode)
		} else {
			slog.Error("Failed to send request", "route", rCfg.Name, "error", err)
		}
		retryRequest(qCfg, rCfg, mCfg, requester)
	} else {
		if resp != nil {
			slog.Info("Received a response from", "route", rCfg.Name, "status", resp.StatusCode)
		}
	}
	metrics.MessagesConsumed.Inc()
}

// shouldRetry determines whether a request should be retried based on the response and retry configuration
func shouldRetry(resp *http.Response, retryConfig *config.RetryConfig) bool {
	if retryConfig == nil && !retryConfig.Enabled {
		return false
	}
	return resp == nil || resp.StatusCode >= retryConfig.ThresholdStatus
}

// retryRequest handles the retry logic for a request, attempting retries as configured
func retryRequest(qCfg *config.QueueConfig, rCfg *config.RouteConfig, mCfg *config.MetricsConfig, requester requester.HTTPRequester) {
	for i := 1; i <= qCfg.Retry.MaxRetries; i++ {
		slog.Info("Retrying request", "route", rCfg.Name, "retry", i)
		time.Sleep(calculateRetryInterval(qCfg.Retry, i))
		resp, err := requester.SendRequest(mCfg, rCfg.Timeout)
		if err == nil && !shouldRetry(resp, qCfg.Retry) {
			if resp != nil {
				slog.Info("Received a response from", "route", rCfg.Name, "status", resp.StatusCode)
			}
			return
		} else {
			slog.Error("Failed to send request", "route", rCfg.Name, "error", err)
		}
	}
}

// calculateRetryInterval computes the time to wait before a retry attempt based on the retry strategy
func calculateRetryInterval(retryConfig *config.RetryConfig, attempt int) time.Duration {
	switch retryConfig.Strategy {
	case common.RetryStrategyFixed:
		return retryConfig.Interval
	case common.RetryStrategyExpo:
		return retryConfig.Interval * time.Duration(attempt)
	case common.RetryStrategyRand:
		return time.Duration(rand.Intn(int(retryConfig.Interval)))
	default:
		slog.Error("Invalid retry strategy", "strategy", retryConfig.Strategy)
		return 0
	}
}
