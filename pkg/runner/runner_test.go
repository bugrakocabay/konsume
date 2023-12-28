package runner

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/bugrakocabay/konsume/pkg/common"
	"github.com/bugrakocabay/konsume/pkg/config"
	"github.com/bugrakocabay/konsume/pkg/queue"
)

type MockHTTPRequester struct {
	MockResponse *http.Response
	MockError    error
	CallCount    int
}

func (m *MockHTTPRequester) SendRequest() (*http.Response, error) {
	m.CallCount++

	if m.MockError != nil {
		return nil, m.MockError
	}
	return m.MockResponse, nil
}

type MockMessageQueueConsumer struct {
	ConnectFunc   func() error
	ConsumeFunc   func(queueName string, handler func(msg []byte) error) error
	CloseFunc     func() error
	ConnectCalled bool
	ConsumeCalled bool
	CloseCalled   bool
}

func (m *MockMessageQueueConsumer) Connect(ctx context.Context) error {
	m.ConnectCalled = true
	if m.ConnectFunc != nil {
		return m.ConnectFunc()
	}
	return errors.New("Connect not implemented")
}

func (m *MockMessageQueueConsumer) Consume(ctx context.Context, queueName string, handler func(msg []byte) error) error {
	m.ConsumeCalled = true
	if m.ConsumeFunc != nil {
		return m.ConsumeFunc(queueName, handler)
	}
	return errors.New("Consume not implemented")
}

func (m *MockMessageQueueConsumer) Close() error {
	m.CloseCalled = true
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return errors.New("Close not implemented")
}

func TestStartConsumers(t *testing.T) {
	cfg := &config.Config{
		Queues:    []*config.QueueConfig{{Name: "testQueue", Provider: "rabbitmq"}},
		Providers: []*config.ProviderConfig{{Name: "rabbitmq", Type: "amqp"}},
	}

	mockConsumer := &MockMessageQueueConsumer{
		ConnectFunc: func() error { return nil },
		ConsumeFunc: func(queueName string, handler func(msg []byte) error) error { return nil },
	}

	consumers := map[string]queue.MessageQueueConsumer{"rabbitmq": mockConsumer}
	providerMap := make(map[string]*config.ProviderConfig)
	providerMap["rabbitmq"] = &config.ProviderConfig{Name: "rabbitmq", Type: "amqp"}

	err := StartConsumers(cfg, consumers, providerMap)
	if err != nil {
		t.Errorf("StartConsumers() error = %v, wantErr %v", err, nil)
	}
	if !mockConsumer.ConnectCalled {
		t.Errorf("Expected Connect to be called, but it was not")
	}
	if !mockConsumer.ConsumeCalled {
		t.Errorf("Expected Consume to be called, but it was not")
	}
}

func TestStartConsumersMultipleQueues(t *testing.T) {
	cfg := &config.Config{
		Queues: []*config.QueueConfig{
			{Name: "validQueue", Provider: "rabbitmq"},
			{Name: "invalidQueue", Provider: "unknown"},
		},
		Providers: []*config.ProviderConfig{{Name: "rabbitmq", Type: "amqp"}},
	}

	mockConsumer := &MockMessageQueueConsumer{
		ConnectFunc: func() error { return nil },
		ConsumeFunc: func(queueName string, handler func(msg []byte) error) error { return nil },
	}

	consumers := map[string]queue.MessageQueueConsumer{"rabbitmq": mockConsumer}
	providerMap := make(map[string]*config.ProviderConfig)
	providerMap["rabbitmq"] = &config.ProviderConfig{Name: "rabbitmq", Type: "amqp"}

	err := StartConsumers(cfg, consumers, providerMap)
	if err == nil || !strings.Contains(err.Error(), "no consumer found for provider: unknown") {
		t.Errorf("Expected error for missing provider, got %v", err)
	}
}

func TestConnectWithRetry(t *testing.T) {
	tests := []struct {
		name              string
		connectFunc       func() error
		retryCount        int
		expectedRetry     int
		expectConnectCall bool
		expectError       bool
	}{
		{
			name: "Successful connection without retries",
			connectFunc: func() error {
				return nil
			},
			retryCount:        1,
			expectedRetry:     0,
			expectConnectCall: true,
			expectError:       false,
		},
		{
			name: "Fail to connect and exhaust retries",
			connectFunc: func() error {
				return errors.New("connection failed")
			},
			retryCount:        1,
			expectedRetry:     1,
			expectConnectCall: true,
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConsumer := &MockMessageQueueConsumer{
				ConnectFunc: tt.connectFunc,
			}
			ctx := context.Background()
			cfg := &config.ProviderConfig{
				Retry: tt.retryCount,
			}

			err := connectWithRetry(ctx, mockConsumer, cfg)

			if (err != nil) != tt.expectError {
				t.Errorf("connectWithRetry() error = %v, wantErr %v", err, tt.expectError)
			}
			if mockConsumer.ConnectCalled != tt.expectConnectCall {
				t.Errorf("Expected Connect to be called = %v, got %v", tt.expectConnectCall, mockConsumer.ConnectCalled)
			}
		})
	}
}

func TestStartConsumersNoQueues(t *testing.T) {
	cfg := &config.Config{}

	err := StartConsumers(cfg, nil, nil)
	if err != nil {
		t.Errorf("Expected no error for no queues, got %v", err)
	}
}

func TestSendRequestWithStrategy(t *testing.T) {
	log.SetOutput(io.Discard)

	tests := []struct {
		name          string
		route         *config.RouteConfig
		mockResponse  *http.Response
		mockError     error
		retryEnabled  bool
		retryStrategy string
		expectedCalls int
		maxRetries    int
		interval      time.Duration
	}{
		{
			name: "should return success when all checks are met",
			route: &config.RouteConfig{
				Name:   "TestRoute",
				URL:    "http://localhost:8080",
				Method: "GET",
				Type:   "REST",
			},
			mockResponse: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("Success")),
			},
			retryEnabled:  false,
			retryStrategy: "",
			expectedCalls: 1,
		},
		{
			name: "should return success when all checks are met and retry is enabled",
			route: &config.RouteConfig{
				Name:   "TestRoute",
				URL:    "http://localhost:8080",
				Method: "GET",
				Type:   "REST",
			},
			mockResponse: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("Success")),
			},
			retryEnabled:  true,
			retryStrategy: "",
			expectedCalls: 1,
		},
		{
			name: "should run fixed strategy when retry is enabled and retry strategy is fixed",
			route: &config.RouteConfig{
				Name:   "TestRoute",
				URL:    "http://localhost:8080",
				Method: "GET",
				Type:   "REST",
			},
			mockResponse: &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
			},
			retryEnabled:  true,
			retryStrategy: common.RetryStrategyFixed,
			expectedCalls: 3,
			maxRetries:    2,
			interval:      1 * time.Millisecond,
		},
		{
			name: "should run expo strategy when retry is enabled and retry strategy is expo",
			route: &config.RouteConfig{
				Name:   "TestRoute",
				URL:    "http://localhost:8080",
				Method: "GET",
				Type:   "REST",
			},
			mockResponse: &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
			},
			retryEnabled:  true,
			retryStrategy: common.RetryStrategyExpo,
			expectedCalls: 3,
			maxRetries:    2,
			interval:      1 * time.Millisecond,
		}, {
			name: "should run random strategy when retry is enabled and retry strategy is random",
			route: &config.RouteConfig{
				Name:   "TestRoute",
				URL:    "http://localhost:8080",
				Method: "GET",
				Type:   "REST",
			},
			mockResponse: &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
			},
			retryEnabled:  true,
			retryStrategy: common.RetryStrategyRand,
			expectedCalls: 3,
			maxRetries:    2,
			interval:      1 * time.Millisecond,
		},
		{
			name: "should return error when retry is enabled and retry strategy is invalid",
			route: &config.RouteConfig{
				Name:   "TestRoute",
				URL:    "http://localhost:8080",
				Method: "GET",
				Type:   "REST",
			},
			mockResponse: &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
			},
			retryEnabled:  true,
			retryStrategy: "invalid",
			expectedCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTPRequester := &MockHTTPRequester{
				MockResponse: tt.mockResponse,
				MockError:    tt.mockError,
			}

			qCfg := &config.QueueConfig{
				Name: "testQueue",
				Retry: &config.RetryConfig{
					Strategy:        tt.retryStrategy,
					MaxRetries:      tt.maxRetries,
					Interval:        tt.interval,
					ThresholdStatus: 500,
				},
			}
			startTime := time.Now()
			sendRequestWithStrategy(qCfg, tt.route, mockHTTPRequester)
			duration := time.Since(startTime)

			if mockHTTPRequester.CallCount != tt.expectedCalls {
				t.Errorf("Expected %d calls to SendRequest, got %d", tt.expectedCalls, mockHTTPRequester.CallCount)
			}
			if qCfg.Retry != nil && qCfg.Retry.Enabled {
				if duration < qCfg.Retry.Interval*time.Duration(qCfg.Retry.MaxRetries) {
					t.Errorf("Expected duration to be greater than %d, got %d", qCfg.Retry.Interval*time.Duration(qCfg.Retry.MaxRetries), duration)
				}
			}
		})
	}
}
