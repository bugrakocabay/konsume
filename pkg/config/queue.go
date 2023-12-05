package config

import (
	"errors"
	"time"

	"konsume/pkg/common"
)

var (
	queueNameNotDefinedError       = errors.New("queue name not defined")
	queueProviderNotDefinedError   = errors.New("queue provider not defined")
	queueProviderDoesNotExistError = errors.New("queue provider does not exist in providers list")

	maxRetriesNotDefinedError = errors.New("max retries not defined")
	intervalNotDefinedError   = errors.New("interval not defined")
	invalidStrategyError      = errors.New("invalid strategy")

	noRoutesDefinedError     = errors.New("no routes defined")
	routeNameNotDefinedError = errors.New("route name not defined")
	urlNotDefinedError       = errors.New("url not defined")
)

// QueueConfig is the main configuration information needed to consume a queue
type QueueConfig struct {
	// Name is the name of the queue
	Name string `yaml:"name"`

	// Provider is the provider that will be used to consume the queue
	Provider string `yaml:"provider"`

	// Retry is the retry configuration for the queue
	Retry *RetryConfig `yaml:"retry,omitempty"`

	// Routes is the list of routes that will be used to send the messages
	Routes []*RouteConfig `yaml:"routes"`
}

// RetryConfig is the main configuration information needed to retry a message
type RetryConfig struct {
	// MaxRetries is the maximum number of retries for the queue
	MaxRetries int `yaml:"max_retries"`

	// Strategy is the retry strategy for the queue, defaults to "fixed"
	Strategy string `yaml:"strategy,omitempty"`

	// Interval is the interval between retries
	Interval time.Duration `yaml:"interval"`

	// ThresholdStatus is the minimum status code that will trigger a retry, defaults to 500
	ThresholdStatus int `yaml:"threshold_status,omitempty"`
}

// RouteConfig is the main configuration information needed to send a message to a service
type RouteConfig struct {
	// Name is the name of the route
	Name string `yaml:"name"`

	// URL is the URL of the service
	URL string `yaml:"url"`

	// Method is the HTTP method of the request, defaults to "POST"
	Method string `yaml:"method,omitempty"`

	// Type is the type of the request, defaults to "REST"
	Type string `yaml:"type,omitempty"`

	// Headers is the list of headers that will be sent with the request
	Headers map[string]string `yaml:"headers,omitempty"`

	// Body is the body of the request
	Body map[string]interface{} `yaml:"body,omitempty"`

	// Query is the query string of the request
	Query map[string]string `yaml:"query,omitempty"`

	// Timeout is the timeout of the request, defaults to 10 seconds
	Timeout time.Duration `yaml:"timeout,omitempty"`
}

func (queue *QueueConfig) validateQueue(providers []*ProviderConfig) error {
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

	if queue.Retry != nil {
		if queue.Retry.MaxRetries == 0 {
			return maxRetriesNotDefinedError
		}
		if queue.Retry.Interval == 0 {
			return intervalNotDefinedError
		}
		if queue.Retry.Strategy == "" {
			queue.Retry.Strategy = "fixed"
		}
		if queue.Retry.Strategy != common.RetryStrategyFixed &&
			queue.Retry.Strategy != common.RetryStrategyExpo &&
			queue.Retry.Strategy != common.RetryStrategyRand {
			return invalidStrategyError
		}
		if queue.Retry.ThresholdStatus == 0 {
			queue.Retry.ThresholdStatus = 500
		}
	}

	if len(queue.Routes) == 0 {
		return noRoutesDefinedError
	} else {
		for _, route := range queue.Routes {
			if len(route.Name) == 0 {
				return routeNameNotDefinedError
			}
			if len(route.URL) == 0 {
				return urlNotDefinedError
			}
			if route.Method == "" {
				route.Method = "POST"
			}
			if route.Type == "" {
				route.Type = "REST"
			}
			if route.Timeout == 0 {
				route.Timeout = 10 * time.Second
			}
		}
	}
	return nil
}
