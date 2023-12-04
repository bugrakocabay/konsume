package runner

import (
	"errors"
	"strings"
	"testing"

	"konsume/pkg/config"
	"konsume/pkg/queue"
)

type MockMessageQueueConsumer struct {
	ConnectFunc   func() error
	ConsumeFunc   func(queueName string, handler func(msg []byte) error) error
	CloseFunc     func() error
	ConnectCalled bool
	ConsumeCalled bool
	CloseCalled   bool
}

func (m *MockMessageQueueConsumer) Connect() error {
	m.ConnectCalled = true
	if m.ConnectFunc != nil {
		return m.ConnectFunc()
	}
	return errors.New("Connect not implemented")
}

func (m *MockMessageQueueConsumer) Consume(queueName string, handler func(msg []byte) error) error {
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

	err := StartConsumers(cfg, consumers)
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

func TestListen(t *testing.T) {
	qCfg := &config.QueueConfig{Name: "testQueue"}

	mockConsumer := &MockMessageQueueConsumer{
		ConnectFunc: func() error { return nil },
		ConsumeFunc: func(queueName string, handler func(msg []byte) error) error { return nil },
	}

	err := listen(mockConsumer, qCfg)
	if err != nil {
		t.Errorf("listen() error = %v, wantErr %v", err, nil)
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

	err := StartConsumers(cfg, consumers)
	if err == nil || !strings.Contains(err.Error(), "no consumer found for provider: unknown") {
		t.Errorf("Expected error for missing provider, got %v", err)
	}
}

func TestStartConsumersNoQueues(t *testing.T) {
	cfg := &config.Config{}

	err := StartConsumers(cfg, nil)
	if err != nil {
		t.Errorf("Expected no error for no queues, got %v", err)
	}
}
