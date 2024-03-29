package config

import (
	"errors"

	"github.com/bugrakocabay/konsume/pkg/common"
)

var (
	providerNameNotDefinedError = errors.New("provider name not defined")
	providerTypeNotDefinedError = errors.New("provider type not defined")
	invalidProviderTypeError    = errors.New("invalid provider type")

	amqpConfigNotDefinedError   = errors.New("amqp config not defined")
	amqpHostNotDefinedError     = errors.New("amqp host not defined")
	amqpPortNotDefinedError     = errors.New("amqp port not defined")
	amqpUsernameNotDefinedError = errors.New("amqp username not defined")
	amqpPasswordNotDefinedError = errors.New("amqp password not defined")

	kafkaConfigNotDefinedError = errors.New("kafka config not defined")
	brokersNotDefinedError     = errors.New("brokers not defined")
	topicNotDefinedError       = errors.New("topic not defined")
	groupNotDefinedError       = errors.New("group not defined")

	stompConfigNotDefinedError   = errors.New("stomp config not defined")
	stompHostNotDefinedError     = errors.New("stomp host not defined")
	stompPortNotDefinedError     = errors.New("stomp port not defined")
	stompUsernameNotDefinedError = errors.New("stomp username not defined")
	stompPasswordNotDefinedError = errors.New("stomp password not defined")
)

// ProviderConfig is the main configuration information needed to connect to a provider
type ProviderConfig struct {
	// Name is the name of the provider
	Name string `yaml:"name" json:"name"`

	// Type is the type of the provider, such as "amqp" or "kafka"
	Type string `yaml:"type" json:"type"`

	// Retry is the number of times to retry connecting to the provider
	Retry int `yaml:"retry,omitempty" json:"retry,omitempty"`

	// AMQPConfig is the configuration for the AMQP provider
	AMQPConfig *AMQPConfig `yaml:"amqp-config,omitempty" json:"amqp-config,omitempty"`

	// KafkaConfig is the configuration for the Kafka provider
	KafkaConfig *KafkaConfig `yaml:"kafka-config,omitempty" json:"kafka-config,omitempty"`

	// StompMQConfig is the configuration for the ActiveMQ provider
	StompMQConfig *StompConfig `yaml:"stomp-config,omitempty" json:"stomp-config,omitempty"`
}

// AMQPConfig is the main configuration information needed to connect to an AMQP provider
type AMQPConfig struct {
	// Host is the host of the queue
	Host string `yaml:"host,omitempty" json:"host,omitempty"`

	// Port is the port of the queue
	Port int `yaml:"port,omitempty" json:"port,omitempty"`

	// Username is the username of the queue
	Username string `yaml:"username,omitempty" json:"username,omitempty"`

	// Password is the password of the queue
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
}

// KafkaConfig is the main configuration information needed to connect to a Kafka provider
type KafkaConfig struct {
	// Brokers is a list of brokers that will be consumed
	Brokers []string `yaml:"brokers,omitempty" json:"brokers,omitempty"`

	// Topic is a list of topics that will be consumed
	Topic string `yaml:"topic,omitempty" json:"topic,omitempty"`

	// Group is the consumer group that will be used
	Group string `yaml:"group,omitempty" json:"group,omitempty"`
}

type StompConfig struct {
	// Host is the host of the queue
	Host string `yaml:"host,omitempty" json:"host,omitempty"`

	// Port is the port of the queue
	Port int `yaml:"port,omitempty" json:"port,omitempty"`

	// Username is the username of the queue
	Username string `yaml:"username,omitempty" json:"username,omitempty"`

	// Password is the password of the queue
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
}

// ValidateProvider validates the ProviderConfig struct
func (p *ProviderConfig) validateProvider() error {
	if len(p.Name) == 0 {
		return providerNameNotDefinedError
	}

	if len(p.Type) == 0 {
		return providerTypeNotDefinedError
	}

	if p.Type != common.QueueSourceRabbitMQ && p.Type != common.QueueSourceKafka &&
		p.Type != common.QueueSourceActiveMQ {
		return invalidProviderTypeError
	}

	if p.Type == common.QueueSourceRabbitMQ {
		if p.AMQPConfig == nil {
			return amqpConfigNotDefinedError
		}

		err := p.AMQPConfig.validateAMQPConfig()
		if err != nil {
			return err
		}
	}

	if p.Type == common.QueueSourceKafka {
		if p.KafkaConfig == nil {
			return kafkaConfigNotDefinedError
		}

		err := p.KafkaConfig.validateKafkaConfig()
		if err != nil {
			return err
		}
	}

	if p.Type == common.QueueSourceActiveMQ {
		if p.StompMQConfig == nil {
			return stompConfigNotDefinedError
		}
		err := p.StompMQConfig.validateStompConfig()
		if err != nil {
			return err
		}
	}

	return nil
}

// validateAMQPConfig validates the AMQPConfig struct
func (a *AMQPConfig) validateAMQPConfig() error {
	if len(a.Host) == 0 {
		return amqpHostNotDefinedError
	}

	if a.Port == 0 {
		return amqpPortNotDefinedError
	}

	if len(a.Username) == 0 {
		return amqpUsernameNotDefinedError
	}

	if len(a.Password) == 0 {
		return amqpPasswordNotDefinedError
	}

	return nil
}

// validateKafkaConfig validates the KafkaConfig struct
func (k *KafkaConfig) validateKafkaConfig() error {
	if len(k.Brokers) == 0 {
		return brokersNotDefinedError
	}

	if len(k.Topic) == 0 {
		return topicNotDefinedError
	}

	if len(k.Group) == 0 {
		return groupNotDefinedError
	}

	return nil
}

func (s *StompConfig) validateStompConfig() error {
	if len(s.Host) == 0 {
		return stompHostNotDefinedError
	}

	if s.Port == 0 {
		return stompPortNotDefinedError
	}

	if len(s.Username) == 0 {
		return stompUsernameNotDefinedError
	}

	if len(s.Password) == 0 {
		return stompPasswordNotDefinedError
	}

	return nil
}
