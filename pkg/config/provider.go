package config

import "errors"

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
	topicsNotDefinedError      = errors.New("topics not defined")
	groupNotDefinedError       = errors.New("group not defined")
)

// ProviderConfig is the main configuration information needed to connect to a provider
type ProviderConfig struct {
	// Name is the name of the provider
	Name string `yaml:"name"`

	// Type is the type of the provider, such as "amqp" or "kafka"
	Type string `yaml:"type"`

	// AMQPConfig is the configuration for the AMQP provider
	AMQPConfig *AMQPConfig `yaml:"amqp-config,omitempty"`

	// KafkaConfig is the configuration for the Kafka provider
	KafkaConfig *KafkaConfig `yaml:"kafka-config,omitempty"`
}

// AMQPConfig is the main configuration information needed to connect to an AMQP provider
type AMQPConfig struct {
	// Host is the host of the queue
	Host string `yaml:"host,omitempty"`

	// Port is the port of the queue
	Port int `yaml:"port,omitempty"`

	// Username is the username of the queue
	Username string `yaml:"username,omitempty"`

	// Password is the password of the queue
	Password string `yaml:"password,omitempty"`
}

// KafkaConfig is the main configuration information needed to connect to a Kafka provider
type KafkaConfig struct {
	// Brokers is a list of brokers that will be consumed
	Brokers []string `yaml:"brokers,omitempty"`

	// Topics is a list of topics that will be consumed
	Topics []string `yaml:"topics,omitempty"`

	// Group is the consumer group that will be used
	Group string `yaml:"group,omitempty"`
}

// ValidateProvider validates the ProviderConfig struct
func (p *ProviderConfig) validateProvider() error {
	if len(p.Name) == 0 {
		return providerNameNotDefinedError
	}

	if len(p.Type) == 0 {
		return providerTypeNotDefinedError
	}

	if p.Type != "amqp" && p.Type != "kafka" {
		return invalidProviderTypeError
	}

	if p.Type == "amqp" {
		if p.AMQPConfig == nil {
			return amqpConfigNotDefinedError
		}

		err := p.AMQPConfig.validateAMQPConfig()
		if err != nil {
			return err
		}
	}

	if p.Type == "kafka" {
		if p.KafkaConfig == nil {
			return kafkaConfigNotDefinedError
		}

		err := p.KafkaConfig.validateKafkaConfig()
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

	if len(k.Topics) == 0 {
		return topicsNotDefinedError
	}

	if len(k.Group) == 0 {
		return groupNotDefinedError
	}

	return nil
}
