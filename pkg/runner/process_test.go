package runner

import (
	"errors"
	"strings"
	"testing"

	"github.com/bugrakocabay/konsume/pkg/config"
)

func TestListenAndProcess(t *testing.T) {
	qCfg := &config.QueueConfig{Name: "testQueue"}

	mockConsumer := &MockMessageQueueConsumer{
		ConnectFunc: func() error { return nil },
		ConsumeFunc: func(queueName string, handler func(msg []byte) error) error { return nil },
	}
	err := listenAndProcess(mockConsumer, qCfg, nil, nil)
	if err != nil {
		t.Errorf("listenAndProcess() error = %v, wantErr %v", err, nil)
	}
	if !mockConsumer.ConsumeCalled {
		t.Errorf("Expected Consume to be called, but it was not")
	}
}

func TestListenAndProcess_ConnectFails(t *testing.T) {
	qCfg := &config.QueueConfig{Name: "testQueue"}

	mockConsumer := &MockMessageQueueConsumer{
		ConnectFunc: func() error { return errors.New("connection failed") },
	}
	err := listenAndProcess(mockConsumer, qCfg, nil, nil)
	if err == nil {
		t.Error("Expected an error when connection fails, but got nil")
	}
}

func TestListenAndProcess_ConsumptionFails(t *testing.T) {
	qCfg := &config.QueueConfig{Name: "testQueue"}

	mockConsumer := &MockMessageQueueConsumer{
		ConnectFunc: func() error { return nil },
		ConsumeFunc: func(queueName string, handler func(msg []byte) error) error {
			return errors.New("consumption failed")
		},
	}
	err := listenAndProcess(mockConsumer, qCfg, nil, nil)
	if err == nil {
		t.Error("Expected an error when consumption fails, but got nil")
	}
}

func TestListenAndProcess_SuccessfulConsumption(t *testing.T) {
	qCfg := &config.QueueConfig{Name: "testQueue"}

	mockConsumer := &MockMessageQueueConsumer{
		ConnectFunc: func() error { return nil },
		ConsumeFunc: func(queueName string, handler func(msg []byte) error) error {
			return handler([]byte("{\"key\":\"value\"}"))
		},
	}
	err := listenAndProcess(mockConsumer, qCfg, nil, nil)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestListenAndProcess_InvalidMessageFormat(t *testing.T) {
	qCfg := &config.QueueConfig{Name: "testQueue"}
	handlerCalled := false

	mockConsumer := &MockMessageQueueConsumer{
		ConnectFunc: func() error { return nil },
		ConsumeFunc: func(queueName string, handler func(msg []byte) error) error {
			handlerCalled = true
			return handler([]byte("invalid message"))
		},
	}
	_ = listenAndProcess(mockConsumer, qCfg, nil, nil) // Error is not expected to be returned

	if !handlerCalled {
		t.Error("Expected handler to be called, but it was not")
	}
}

func TestListenAndProcess_RouteHandling(t *testing.T) {
	mockConsumer := &MockMessageQueueConsumer{
		ConsumeFunc: func(queueName string, handler func(msg []byte) error) error {
			return handler([]byte("{\"key\":\"value\"}"))
		},
	}

	qCfg1 := &config.QueueConfig{Name: "testQueue"}
	qCfg1.Routes = []*config.RouteConfig{
		{
			Body:   map[string]interface{}{"key": "value"},
			Method: "GET",
			URL:    "http://localhost/test",
			Query:  map[string]string{"param": "value"},
		},
	}
	err := listenAndProcess(mockConsumer, qCfg1, nil, nil)
	if err != nil {
		t.Errorf("listenAndProcess() with non-empty body returned error: %v", err)
	}

	qCfg2 := &config.QueueConfig{Name: "testQueue"}
	qCfg2.Routes = []*config.RouteConfig{
		{
			Body:   nil,
			Method: "GET",
			URL:    "http://localhost/test",
			Query:  map[string]string{"param": "value"},
		},
	}
	err = listenAndProcess(mockConsumer, qCfg2, nil, nil)
	if err != nil {
		t.Errorf("listenAndProcess() with empty body returned error: %v", err)
	}
}

func TestListenAndProcess_PrepareRequestBody_Success(t *testing.T) {
	qCfg := &config.QueueConfig{
		Name: "testQueue",
		Routes: []*config.RouteConfig{
			{
				Body:   map[string]interface{}{"key": "value"},
				Method: "GET",
				URL:    "http://localhost/test",
			},
		},
	}

	mockConsumer := &MockMessageQueueConsumer{
		ConsumeFunc: func(queueName string, handler func(msg []byte) error) error {
			return handler([]byte("{\"key\":\"value\"}"))
		},
	}

	err := listenAndProcess(mockConsumer, qCfg, nil, nil)
	if err != nil {
		t.Errorf("listenAndProcess() with valid body returned error: %v", err)
	}
}

func TestPrepareRequestBody(t *testing.T) {
	messageData := map[string]interface{}{"key1": "value1"}

	restRouteConfig := config.RouteConfig{Type: "rest", Body: map[string]interface{}{"key": "{{key1}}"}}
	restBody, err := prepareRequestBody(&restRouteConfig, messageData)
	if err != nil || string(restBody) != `{"key":"value1"}` {
		t.Errorf("prepareRequestBody REST failed, got: %s, error: %v", string(restBody), err)
	}

	graphqlRouteConfig := config.RouteConfig{Type: "graphql", Body: map[string]interface{}{"query": "query { test(key: {{key1}}) }"}}
	graphqlBody, err := prepareRequestBody(&graphqlRouteConfig, messageData)
	expectedGraphQL := `{"query":"query { test(key: \"value1\") }"}`
	if err != nil || string(graphqlBody) != expectedGraphQL {
		t.Errorf("prepareRequestBody GraphQL failed, got: %s, expected: %s, error: %v", string(graphqlBody), expectedGraphQL, err)
	}
}

func TestPrepareGraphQLBody(t *testing.T) {
	messageData := map[string]interface{}{"key1": "value1"}
	routeConfig := config.RouteConfig{Body: map[string]interface{}{"query": "query { test(key: {{key1}}) }"}}

	body, err := prepareGraphQLBody(&routeConfig, messageData)
	expectedGraphQL := `{"query":"query { test(key: \"value1\") }"}`
	if err != nil || string(body) != expectedGraphQL {
		t.Errorf("prepareGraphQLBody failed, got: %s, expected: %s, error: %v", string(body), expectedGraphQL, err)
	}
}

func TestPrepareRESTBody(t *testing.T) {
	messageData := map[string]interface{}{"key1": "value1"}
	routeConfig := config.RouteConfig{Body: map[string]interface{}{"key": "{{key1}}"}}

	body, err := prepareRESTBody(&routeConfig, messageData)
	if err != nil || string(body) != `{"key":"value1"}` {
		t.Errorf("prepareRESTBody failed, got: %s, error: %v", string(body), err)
	}
}

func TestGetGraphQLOperation(t *testing.T) {
	bodyMapQuery := map[string]interface{}{"query": "query { test }"}
	bodyMapMutation := map[string]interface{}{"mutation": "mutation { addTest }"}

	if getGraphQLOperation(bodyMapQuery) != "query { test }" {
		t.Errorf("getGraphQLOperation query failed")
	}
	if getGraphQLOperation(bodyMapMutation) != "mutation { addTest }" {
		t.Errorf("getGraphQLOperation mutation failed")
	}
}

func TestAppendQueryParams(t *testing.T) {
	url := "http://localhost:8080"
	queryParams := map[string]string{"key1": "value1", "key2": "value2"}

	urlWithQueryParams := appendQueryParams(url, queryParams)
	if strings.Contains(urlWithQueryParams, "key1=value1") == false || strings.Contains(urlWithQueryParams, "key2=value2") == false {
		t.Errorf("appendQueryParams failed, urlWithQueryParams does not contain queryParams")
	}
}

func TestAppendQueryParams_EmptyQueryParams(t *testing.T) {
	url := "http://localhost:8080"
	queryParams := map[string]string{}

	urlWithQueryParams := appendQueryParams(url, queryParams)
	if urlWithQueryParams != url {
		t.Errorf("appendQueryParams failed, urlWithQueryParams should be equal to url")
	}
}
