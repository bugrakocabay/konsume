package config

import (
	"errors"
	"log/slog"
)

var (
	invalidMetricsPortError = errors.New("invalid metrics port")
)

// MetricsConfig is the configuration for the metrics endpoint
type MetricsConfig struct {
	// Enabled is a flag that enables the metrics endpoint
	Enabled bool `yaml:"enabled"`

	// Port is the port that the metrics endpoint will listen on
	Port int `yaml:"port,omitempty"`

	// Path is the path that the metrics endpoint will listen on
	Path string `yaml:"path,omitempty"`
}

// ValidateAll validates the configuration and returns default values if necessary
func (c *MetricsConfig) validateMetrics() error {
	if !c.Enabled {
		return nil
	}
	if c.Port == 0 {
		slog.Debug("No port defined for metrics endpoint, using default port 8080")
		c.Port = 8080
	}
	if c.Port < 0 || c.Port > 65535 {
		return invalidMetricsPortError
	}
	if len(c.Path) == 0 {
		slog.Debug("No path defined for metrics endpoint, using default path /metrics")
		c.Path = "metrics"
	}
	if c.Path[0] != '/' {
		c.Path = "/" + c.Path
	}
	return nil
}
