package runner

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"runtime"

	"github.com/bugrakocabay/konsume/pkg/common"
	"github.com/bugrakocabay/konsume/pkg/config"
	"github.com/bugrakocabay/konsume/pkg/database"
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
			messageData, err := util.ParseJSONToMap(msg)
			if err != nil {
				slog.Error("Failed to parse message", "error", err)
				return
			}

			if qCfg.Routes != nil && len(qCfg.Routes) > 0 {
				for _, rCfg := range qCfg.Routes {
					var body []byte

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

			if qCfg.DatabaseRoutes != nil && len(qCfg.DatabaseRoutes) > 0 {
				for _, dbRoute := range qCfg.DatabaseRoutes {
					db, ok := databases[dbRoute.Provider]
					if !ok {
						slog.Error("Database not found", "database", dbRoute.Name)
						continue
					}
					if err = db.Insert(messageData, *dbRoute); err != nil {
						slog.Error("Failed to insert data into database", "error", err)
					}
				}
			}

		}(msg)
		return nil
	})
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

func prepareRESTBody(rCfg *config.RouteConfig, messageData map[string]interface{}) ([]byte, error) {
	return util.ProcessTemplate(rCfg.Body, messageData)
}

func getGraphQLOperation(bodyMap map[string]interface{}) string {
	if operation, ok := bodyMap["query"].(string); ok {
		return operation
	}
	if operation, ok := bodyMap["mutation"].(string); ok {
		return operation
	}
	return ""
}

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
