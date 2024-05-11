package config

import (
	"errors"
	"log/slog"
	"time"

	"github.com/bugrakocabay/konsume/pkg/common"
)

var (
	queueNameNotDefinedError       = errors.New("queue name not defined")
	queueProviderNotDefinedError   = errors.New("queue provider not defined")
	queueProviderDoesNotExistError = errors.New("queue provider does not exist in providers list")

	maxRetriesNotDefinedError = errors.New("max retries not defined")
	intervalNotDefinedError   = errors.New("interval not defined")
	invalidStrategyError      = errors.New("invalid strategy")

	routeNameNotDefinedError             = errors.New("route name not defined")
	urlNotDefinedError                   = errors.New("url not defined")
	bodyNotDefinedError                  = errors.New("when using graphql type, body must be defined")
	invalidBodyForGraphQLError           = errors.New("when using graphql type, body must contain query or mutation")
	bodyNotContainsStringForGraphQLError = errors.New("when using graphql type, body must contain string for query or mutation")

	databaseRouteNameNotDefinedError              = errors.New("database route name not defined")
	databaseRouteProviderNotDefinedError          = errors.New("database route provider not defined")
	databaseRouteProviderDoesNotExistError        = errors.New("database route provider does not exist in databases list")
	databaseRouteTableOrCollectionNotDefinedError = errors.New("database route table or collection not defined")
	dataBaseRouteMappingNotDefinedError           = errors.New("database route mapping not defined")
)

// QueueConfig is the main configuration information needed to consume a queue
type QueueConfig struct {
	// Name is the name of the queue
	Name string `yaml:"name" json:"name"`

	// Provider is the provider that will be used to consume the queue
	Provider string `yaml:"provider" json:"provider"`

	// Retry is the retry configuration for the queue
	Retry *RetryConfig `yaml:"retry,omitempty" json:"retry,omitempty"`

	// Routes is the list of routes that will be used to send the messages
	Routes []*RouteConfig `yaml:"routes" json:"routes"`

	// DatabaseRoutes is the list of databases that will be used to store the messages
	DatabaseRoutes []*DatabaseRouteConfig `yaml:"database-routes,omitempty" json:"database-routes,omitempty"`
}

// RetryConfig is the main configuration information needed to retry a message
type RetryConfig struct {
	// Enabled is the flag that indicates if the retry is enabled, defaults to false
	Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`

	// MaxRetries is the maximum number of retries for the queue
	MaxRetries int `yaml:"max-retries" json:"max-retries"`

	// Strategy is the retry strategy for the queue, defaults to "fixed"
	Strategy string `yaml:"strategy,omitempty" json:"strategy,omitempty"`

	// Interval is the interval between retries
	Interval time.Duration `yaml:"interval,omitempty" json:"interval,omitempty"`

	// ThresholdStatus is the minimum status code that will trigger a retry, defaults to 500
	ThresholdStatus int `yaml:"threshold-status,omitempty" json:"threshold-status,omitempty"`
}

// RouteConfig is the main configuration information needed to send a message to a service
type RouteConfig struct {
	// Name is the name of the route
	Name string `yaml:"name" json:"name"`

	// URL is the URL of the service
	URL string `yaml:"url" json:"url"`

	// Method is the HTTP method of the request, defaults to "POST"
	Method string `yaml:"method,omitempty" json:"method,omitempty"`

	// Type is the type of the request, defaults to "REST"
	Type string `yaml:"type,omitempty" json:"type,omitempty"`

	// Headers is the list of headers that will be sent with the request
	Headers map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`

	// Body is the body of the request
	Body map[string]interface{} `yaml:"body,omitempty" json:"body,omitempty"`

	// Query is the query string of the request
	Query map[string]string `yaml:"query,omitempty" json:"query,omitempty"`

	// Timeout is the timeout of the request, defaults to 10 seconds
	Timeout time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

// DatabaseRouteConfig is the main configuration information needed to store a message in a database
type DatabaseRouteConfig struct {
	// Name is the name of the database route
	Name string `yaml:"name" json:"name"`

	// Provider is the name of the provider database
	Provider string `yaml:"provider" json:"provider"`

	// Table is the name of the table in a SQL database
	Table string `yaml:"table,omitempty" json:"table,omitempty"`

	// Collection is the name of the collection in a NoSQL database
	Collection string `yaml:"collection,omitempty" json:"collection,omitempty"`

	// Mapping is the mapping of the message to the database
	Mapping map[string]string `yaml:"mapping" json:"mapping"`
}

func (queue *QueueConfig) validateQueue(providers []*ProviderConfig, databases []*DatabaseConfig) error {
	if len(queue.Name) == 0 {
		return queueNameNotDefinedError
	}
	if len(queue.Provider) == 0 {
		return queueProviderNotDefinedError
	}
	providerExists := false
	for _, provider := range providers {
		if provider.Name == queue.Provider {
			providerExists = true
			break
		}
	}
	if !providerExists {
		return queueProviderDoesNotExistError
	}

	if queue.Retry != nil && queue.Retry.Enabled {
		if queue.Retry.MaxRetries == 0 {
			return maxRetriesNotDefinedError
		}
		if queue.Retry.Interval == 0 {
			return intervalNotDefinedError
		}
		if queue.Retry.Strategy == "" {
			slog.Debug("Retry strategy not defined, using default strategy fixed", "queue", queue.Name)
			queue.Retry.Strategy = "fixed"
		}
		if queue.Retry.Strategy != common.RetryStrategyFixed &&
			queue.Retry.Strategy != common.RetryStrategyExpo &&
			queue.Retry.Strategy != common.RetryStrategyRand {
			return invalidStrategyError
		}
		if queue.Retry.ThresholdStatus == 0 {
			slog.Debug("Threshold status not defined, using default status 500", "queue", queue.Name)
			queue.Retry.ThresholdStatus = 500
		}
	}

	if len(queue.Routes) > 0 {
		for _, route := range queue.Routes {
			if len(route.Name) == 0 {
				return routeNameNotDefinedError
			}
			if len(route.URL) == 0 {
				return urlNotDefinedError
			}
			if route.Method == "" {
				slog.Debug("Route method not defined, using default method POST", "route", route.Name)
				route.Method = "POST"
			}
			if route.Type == "" {
				slog.Debug("Route type not defined, using default type REST", "route", route.Name)
				route.Type = common.RouteTypeREST
			}
			if route.Type == common.RouteTypeGraphQL {
				if len(route.Body) == 0 {
					return bodyNotDefinedError
				}
				v1, ok := route.Body["query"]
				v2, ok2 := route.Body["mutation"]
				if !ok && !ok2 {
					return invalidBodyForGraphQLError
				}
				if ok || ok2 {
					if ok {
						_, ok = v1.(string)
					}
					if ok2 {
						_, ok2 = v2.(string)
					}
					if !ok && !ok2 {
						return bodyNotContainsStringForGraphQLError
					}
				}
				if route.Method != "POST" {
					slog.Debug("GraphQL route method is not POST, setting it to POST", "route", route.Name)
					route.Method = "POST"
				}
			}
			if route.Timeout == 0 {
				slog.Debug("Route timeout not defined, using default timeout 10 seconds", "route", route.Name)
				route.Timeout = 10 * time.Second
			}
		}
	}

	if len(queue.DatabaseRoutes) > 0 {
		for _, databaseRoute := range queue.DatabaseRoutes {
			if len(databaseRoute.Name) == 0 {
				return databaseRouteNameNotDefinedError
			}
			if len(databaseRoute.Provider) == 0 {
				return databaseRouteProviderNotDefinedError
			}
			databaseExists := false
			for _, database := range databases {
				if database.Name == databaseRoute.Provider {
					databaseExists = true
					break
				}
			}
			if !databaseExists {
				return databaseRouteProviderDoesNotExistError
			}
			if len(databaseRoute.Table) == 0 && len(databaseRoute.Collection) == 0 {
				return databaseRouteTableOrCollectionNotDefinedError
			}
			if len(databaseRoute.Mapping) == 0 {
				return dataBaseRouteMappingNotDefinedError
			}
		}
	}
	return nil
}
