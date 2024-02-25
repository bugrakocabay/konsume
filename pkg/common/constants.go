package common

const (
	QueueSourceRabbitMQ = "rabbitmq"
	QueueSourceKafka    = "kafka"
	QueueSourceActiveMQ = "activemq"
)

const (
	RetryStrategyFixed = "fixed"
	RetryStrategyExpo  = "expo"
	RetryStrategyRand  = "random"
)

const (
	RouteTypeREST    = "REST"
	RouteTypeGraphQL = "graphql"
)

const (
	DatabaseTypePostgresql = "postgresql"
	DatabaseTypeMongoDB    = "mongodb"
)

const (
	KonsumeConfigPath = "KONSUME_CONFIG_PATH"
	KonsumePluginPath = "KONSUME_PLUGIN_PATH"
)
