package activemq

import (
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/bugrakocabay/konsume/pkg/config"
	"github.com/bugrakocabay/konsume/pkg/queue"

	"github.com/go-stomp/stomp/v3"
)

type Consumer struct {
	conn   *stomp.Conn
	sub    *stomp.Subscription
	config *config.StompConfig
}

func NewConsumer(cfg *config.StompConfig) *Consumer {
	return &Consumer{
		config: cfg,
	}
}

// NewConsumerFactory returns a new RabbitMQ consumer based on the provided configuration.
func NewConsumerFactory(cfg *config.ProviderConfig) (queue.MessageQueueConsumer, error) {
	return NewConsumer(cfg.StompMQConfig), nil
}

func (c *Consumer) Connect() error {
	slog.Debug("Attempting to connect to ActiveMQ", "host", c.config.Host, "port", c.config.Port)
	var err error
	var options = []func(*stomp.Conn) error{
		stomp.ConnOpt.HeartBeat(2*time.Hour, 2*time.Hour),
		stomp.ConnOpt.HeartBeatError(5 * time.Minute),
		stomp.ConnOpt.Login(c.config.Username, c.config.Password),
	}
	var connectionString = net.JoinHostPort(c.config.Host, strconv.Itoa(c.config.Port))

	c.conn, err = stomp.Dial("tcp", connectionString, options...)
	if err != nil {
		return err
	}
	slog.Info("Connected to ActiveMQ", "host", c.config.Host, "port", c.config.Port)
	return nil
}

func (c *Consumer) Consume(queueName string, handler func(msg []byte) error) error {
	slog.Debug("Starting to consume messages from ActiveMQ", "queueName", queueName)
	sub, err := c.conn.Subscribe(queueName, stomp.AckAuto)
	if err != nil {
		return err
	}
	go func() {
		for {
			m, err := sub.Read()
			if err != nil {
				slog.Error("Failed to read message from ActiveMQ", "error", err)
			}
			if err = handler(m.Body); err != nil {
				slog.Error("Failed to process message", "error", err)
			}
		}
	}()
	return nil
}

func (c *Consumer) Close() error {
	slog.Debug("Closing connection to ActiveMQ")
	err := c.sub.Unsubscribe()
	if err != nil {
		return err
	}
	err = c.conn.Disconnect()
	if err != nil {
		return err
	}
	slog.Debug("ActiveMQ connection closed successfully")
	return nil
}
