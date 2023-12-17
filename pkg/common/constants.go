package common

const (
	QueueSourceRabbitMQ = "rabbitmq"
	QueueSourceKafka    = "kafka"
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
