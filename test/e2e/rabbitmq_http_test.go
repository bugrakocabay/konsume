package e2e

import (
	"fmt"
	"os"
	"testing"
	"time"

	konsume "github.com/bugrakocabay/konsume/cmd"
	"github.com/bugrakocabay/konsume/pkg/common"
	"github.com/bugrakocabay/konsume/pkg/config"
)

const (
	host          = "localhost"
	port          = 5672
	username      = "user"
	password      = "password"
	testQueueName = "test-queue"
)

type TestCase struct {
	Description    string
	KonsumeConfig  *config.Config
	SetupMessage   SetupMessage
	ExpectedResult []HTTPRequestExpectation
}

type SetupMessage struct {
	QueueName string
	Message   []byte
}

type HTTPRequestExpectation struct {
	URL    string
	Method string
	Body   string
}

func TestKonsumeWithRabbitMQHTTP(t *testing.T) {
	mockServer, url, requestCapture := setupMockServer(t)
	defer mockServer.Close()

	tests := []TestCase{
		{
			Description: "Test with single message",
			KonsumeConfig: &config.Config{
				Providers: []*config.ProviderConfig{
					{
						Name: "rabbit-queue",
						Type: "rabbitmq",
						AMQPConfig: &config.AMQPConfig{
							Host:     host,
							Port:     port,
							Username: username,
							Password: password,
						},
					},
				},
				Queues: []*config.QueueConfig{
					{
						Name:     testQueueName + "-1",
						Provider: "rabbit-queue",
						Routes: []*config.RouteConfig{
							{
								Name: "test-route",
								URL:  fmt.Sprintf("%s/200", url),
							},
						},
					},
				},
			},
			SetupMessage: SetupMessage{
				QueueName: testQueueName + "-1",
				Message:   []byte("{\"id\": 0, \"name\": \"test\"}"),
			},
			ExpectedResult: []HTTPRequestExpectation{
				{
					URL:    "/200",
					Body:   "{\"id\": 0, \"name\": \"test\"}",
					Method: "POST",
				},
			},
		},
		{
			Description: "Test with single message dynamic body",
			KonsumeConfig: &config.Config{
				Providers: []*config.ProviderConfig{
					{
						Name: "rabbit-queue",
						Type: "rabbitmq",
						AMQPConfig: &config.AMQPConfig{
							Host:     host,
							Port:     port,
							Username: username,
							Password: password,
						},
					},
				},
				Queues: []*config.QueueConfig{
					{
						Name:     testQueueName + "-2",
						Provider: "rabbit-queue",
						Routes: []*config.RouteConfig{
							{
								Name: "test-route",
								URL:  fmt.Sprintf("%s/200", url),
								Body: map[string]interface{}{
									"some-id":   "{{id}}",
									"some-name": "{{name}}",
								},
							},
						},
					},
				},
			},
			SetupMessage: SetupMessage{
				QueueName: testQueueName + "-2",
				Message:   []byte("{\"id\": 1, \"name\": \"test\"}"),
			},
			ExpectedResult: []HTTPRequestExpectation{
				{
					URL:    "/200",
					Body:   "{\"some-id\":1,\"some-name\":\"test\"}",
					Method: "POST",
				},
			},
		},
		{
			Description: "Test with single message fixed retry strategy",
			KonsumeConfig: &config.Config{
				Providers: []*config.ProviderConfig{
					{
						Name: "rabbit-queue",
						Type: "rabbitmq",
						AMQPConfig: &config.AMQPConfig{
							Host:     host,
							Port:     port,
							Username: username,
							Password: password,
						},
					},
				},
				Queues: []*config.QueueConfig{
					{
						Name:     testQueueName + "-3",
						Provider: "rabbit-queue",
						Retry: &config.RetryConfig{
							Enabled:         true,
							MaxRetries:      2,
							Strategy:        common.RetryStrategyFixed,
							ThresholdStatus: 500,
							Interval:        1 * time.Second,
						},
						Routes: []*config.RouteConfig{
							{
								Name: "test-route",
								URL:  fmt.Sprintf("%s/500", url),
							},
						},
					},
				},
			},
			SetupMessage: SetupMessage{
				QueueName: testQueueName + "-3",
				Message:   []byte("{\"id\": 1, \"name\": \"test\"}"),
			},
			ExpectedResult: []HTTPRequestExpectation{
				{
					URL:    "/500",
					Body:   "{\"id\": 1, \"name\": \"test\"}",
					Method: "POST",
				},
				{
					URL:    "/500",
					Body:   "{\"id\": 1, \"name\": \"test\"}",
					Method: "POST",
				},
				{
					URL:    "/500",
					Body:   "{\"id\": 1, \"name\": \"test\"}",
					Method: "POST",
				},
			},
		},
		{
			Description: "Test with single message expo retry strategy",
			KonsumeConfig: &config.Config{
				Providers: []*config.ProviderConfig{
					{
						Name: "rabbit-queue",
						Type: "rabbitmq",
						AMQPConfig: &config.AMQPConfig{
							Host:     host,
							Port:     port,
							Username: username,
							Password: password,
						},
					},
				},
				Queues: []*config.QueueConfig{
					{
						Name:     testQueueName + "-4",
						Provider: "rabbit-queue",
						Retry: &config.RetryConfig{
							Enabled:         true,
							MaxRetries:      2,
							Strategy:        common.RetryStrategyExpo,
							ThresholdStatus: 500,
							Interval:        1 * time.Second,
						},
						Routes: []*config.RouteConfig{
							{
								Name: "test-route",
								URL:  fmt.Sprintf("%s/500", url),
							},
						},
					},
				},
			},
			SetupMessage: SetupMessage{
				QueueName: testQueueName + "-4",
				Message:   []byte("{\"id\": 1, \"name\": \"test\"}"),
			},
			ExpectedResult: []HTTPRequestExpectation{
				{
					URL:    "/500",
					Body:   "{\"id\": 1, \"name\": \"test\"}",
					Method: "POST",
				},
				{
					URL:    "/500",
					Body:   "{\"id\": 1, \"name\": \"test\"}",
					Method: "POST",
				},
				{
					URL:    "/500",
					Body:   "{\"id\": 1, \"name\": \"test\"}",
					Method: "POST",
				},
			},
		},
		{
			Description: "Test with single message rand retry strategy",
			KonsumeConfig: &config.Config{
				Providers: []*config.ProviderConfig{
					{
						Name: "rabbit-queue",
						Type: "rabbitmq",
						AMQPConfig: &config.AMQPConfig{
							Host:     host,
							Port:     port,
							Username: username,
							Password: password,
						},
					},
				},
				Queues: []*config.QueueConfig{
					{
						Name:     testQueueName + "-5",
						Provider: "rabbit-queue",
						Retry: &config.RetryConfig{
							Enabled:         true,
							MaxRetries:      2,
							Strategy:        common.RetryStrategyRand,
							ThresholdStatus: 500,
							Interval:        1 * time.Second,
						},
						Routes: []*config.RouteConfig{
							{
								Name: "test-route",
								URL:  fmt.Sprintf("%s/500", url),
							},
						},
					},
				},
			},
			SetupMessage: SetupMessage{
				QueueName: testQueueName + "-5",
				Message:   []byte("{\"id\": 1, \"name\": \"test\"}"),
			},
			ExpectedResult: []HTTPRequestExpectation{
				{
					URL:    "/500",
					Body:   "{\"id\": 1, \"name\": \"test\"}",
					Method: "POST",
				},
				{
					URL:    "/500",
					Body:   "{\"id\": 1, \"name\": \"test\"}",
					Method: "POST",
				},
				{
					URL:    "/500",
					Body:   "{\"id\": 1, \"name\": \"test\"}",
					Method: "POST",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			requestCapture.ReceivedRequests = nil
			// Setting up the config file
			configFilePath, cleanup := writeConfigToFile(test.KonsumeConfig)
			defer cleanup()
			os.Setenv("KONSUME_CONFIG_PATH", configFilePath)

			// Running konsume and waiting for it to consume the message
			go konsume.Execute()
			time.Sleep(2 * time.Second)

			// Pushing the message to the queue
			connString := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port)
			conn, ch, err := connectToRabbitMQ(connString)
			if err != nil {
				t.Fatalf("Failed to connect to RabbitMQ: %v", err)
			}
			defer conn.Close()
			defer ch.Close()
			err = pushMessageToQueue(ch, test.SetupMessage.QueueName, test.SetupMessage.Message)
			if err != nil {
				t.Fatalf("Failed to push message to queue: %v", err)
			}
			sleep(test)

			// Checking the captured requests
			requestCapture.Mutex.Lock()
			defer requestCapture.Mutex.Unlock()
			if len(requestCapture.ReceivedRequests) != len(test.ExpectedResult) {
				t.Fatalf("Expected %d HTTP requests, but got %d", len(test.ExpectedResult), len(requestCapture.ReceivedRequests))
			}
			for i, expectedRequest := range test.ExpectedResult {
				if requestCapture.ReceivedRequests[i].URL != expectedRequest.URL {
					t.Errorf("Expected URL: %s, but got: %s", expectedRequest.URL, requestCapture.ReceivedRequests[i].URL)
				}
				if requestCapture.ReceivedRequests[i].Method != expectedRequest.Method {
					t.Errorf("Expected method: %s, but got: %s", expectedRequest.Method, requestCapture.ReceivedRequests[i].Method)
				}
				if requestCapture.ReceivedRequests[i].Body != expectedRequest.Body {
					t.Errorf("Expected body: %s, but got: %s", expectedRequest.Body, requestCapture.ReceivedRequests[i].Body)
				}
			}
		})
	}
}
