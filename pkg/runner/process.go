package runner

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/bugrakocabay/konsume/pkg/common"
	"github.com/bugrakocabay/konsume/pkg/config"
	"github.com/bugrakocabay/konsume/pkg/database"
	"github.com/bugrakocabay/konsume/pkg/metrics"
	"github.com/bugrakocabay/konsume/pkg/queue"
	"github.com/bugrakocabay/konsume/pkg/requester"
	"github.com/bugrakocabay/konsume/pkg/util"
)

// listenAndProcess listens the queue and processes the messages
func listenAndProcess(consumer queue.MessageQueueConsumer, qCfg *config.QueueConfig, mCfg *config.MetricsConfig, databases map[string]database.Database) error {
	semaphore := make(chan struct{}, runtime.NumCPU()*2)

	return consumer.Consume(qCfg.Name, func(msg []byte) error {
		slog.Info("Received a message", "queue", qCfg.Name, "message", string(msg))
		semaphore <- struct{}{}
		go func(msg []byte) {
			defer func() {
				<-semaphore
				if r := recover(); r != nil {
					slog.Error("Recovered from panic in message processing", "error", r)
				}
			}()
			processMessage(msg, qCfg, mCfg, databases)
		}(msg)

		return nil
	})
}

// processMessage processes the message by sending requests and inserting data into databases
func processMessage(msg []byte, qCfg *config.QueueConfig, mCfg *config.MetricsConfig, databases map[string]database.Database) {
	messageData, err := util.ParseJSONToMap(msg)
	if err != nil {
		slog.Error("Failed to parse message", "error", err)
		return
	}
	handleRoutes(qCfg, messageData, msg, mCfg)
	handleDatabaseRoutes(qCfg, messageData, databases)
}

// handleRoutes sends requests to the routes defined in the queue config
func handleRoutes(qCfg *config.QueueConfig, messageData map[string]interface{}, msg []byte, mCfg *config.MetricsConfig) {
	if qCfg.Routes == nil {
		return
	}
	for _, rCfg := range qCfg.Routes {
		var body []byte
		var err error
		if len(rCfg.Body) > 0 {
			body, err = prepareRequestBody(rCfg, messageData)
			if err != nil {
				slog.Error("Failed to prepare request body", "error", err)
				continue
			}
		} else {
			body = msg
		}
		rCfg.URL = appendQueryParams(rCfg.URL, rCfg.Query)
		rqstr := requester.NewRequester(rCfg.URL, rCfg.Method, body, rCfg.Headers)
		sendRequestWithStrategy(qCfg, rCfg, mCfg, rqstr)
	}
}

// handleDatabaseRoutes inserts data into the databases defined in the queue config
func handleDatabaseRoutes(qCfg *config.QueueConfig, messageData map[string]interface{}, databases map[string]database.Database) {
	if qCfg.DatabaseRoutes == nil {
		return
	}
	for _, dbRoute := range qCfg.DatabaseRoutes {
		db, ok := databases[dbRoute.Provider]
		if !ok {
			slog.Error("Database not found", "database", dbRoute.Name)
			continue
		}
		if err := db.Insert(messageData, *dbRoute); err != nil {
			slog.Error("Failed to insert data into database", "error", err)
		}
	}
}

// sendRequestWithStrategy attempts to send an HTTP request and retries based on the provided configuration
func sendRequestWithStrategy(qCfg *config.QueueConfig, rCfg *config.RouteConfig, mCfg *config.MetricsConfig, requester requester.HTTPRequester) {
	resp, err := requester.SendRequest(mCfg, rCfg.Timeout)
	if err != nil {
		slog.Error("Error occurred while sending request", "route", rCfg.Name, "error", err)
	}
	if resp != nil {
		body, err := util.ReadRequestBody(resp)
		if err != nil {
			slog.Error("Failed to read response body", "route", rCfg.Name, "error", err)
		}
		slog.Info("Received a response from", "route", rCfg.Name, "status", resp.StatusCode, "response", body)
	} else {
		slog.Error("Received an empty response", "route", rCfg.Name)
	}
	if shouldRetry(resp, qCfg.Retry) {
		retryRequest(qCfg, rCfg, mCfg, requester)
	}
	metrics.MessagesConsumed.Inc()
}

// shouldRetry determines whether a request should be retried based on the response and retry configuration
func shouldRetry(resp *http.Response, retryConfig *config.RetryConfig) bool {
	if retryConfig == nil || !retryConfig.Enabled {
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
		if err != nil {
			slog.Error("Error occurred while retrying the request", "route", rCfg.Name, "error", err)
		}
		if resp != nil {
			body, err := util.ReadRequestBody(resp)
			if err != nil {
				slog.Error("Failed to read response body", "route", rCfg.Name, "error", err)
			}
			slog.Info("Received a response from retry", "route", rCfg.Name, "status", resp.StatusCode, "response", body)
			if !shouldRetry(resp, qCfg.Retry) {
				return
			}
		} else {
			slog.Error("Received an empty response from retry", "route", rCfg.Name)
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

// prepareRequestBody prepares the request body according to the route type
func prepareRequestBody(rCfg *config.RouteConfig, messageData map[string]interface{}) ([]byte, error) {
	if rCfg.Type == common.RouteTypeGraphQL {
		return prepareGraphQLBody(rCfg, messageData)
	} else {
		return prepareRESTBody(rCfg, messageData)
	}
}

// prepareGraphQLBody creates a graphql request body from the given route config and message data
func prepareGraphQLBody(rCfg *config.RouteConfig, messageData map[string]interface{}) ([]byte, error) {
	graphqlOperation := getGraphQLOperation(rCfg.Body)
	if graphqlOperation == "" {
		return nil, fmt.Errorf("no query or mutation found in graphql body")
	}
	bodyStr, err := util.ProcessGraphQLTemplate(graphqlOperation, messageData)
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]string{"query": bodyStr})
}

// prepareRESTBody creates a rest request body from the given route config and message data
func prepareRESTBody(rCfg *config.RouteConfig, messageData map[string]interface{}) ([]byte, error) {
	return util.ProcessTemplate(rCfg.Body, messageData)
}

// getGraphQLOperation returns the query or mutation from the graphql request body
func getGraphQLOperation(bodyMap map[string]interface{}) string {
	if operation, ok := bodyMap["query"].(string); ok {
		return operation
	}
	if operation, ok := bodyMap["mutation"].(string); ok {
		return operation
	}
	return ""
}

// appendQueryParams appends the query parameters to the given url
func appendQueryParams(url string, queryParams map[string]string) string {
	if len(queryParams) == 0 {
		return url
	}
	url += "?"
	for key, value := range queryParams {
		url += fmt.Sprintf("%s=%s&", key, value)
	}
	return url[:len(url)-1]
}
