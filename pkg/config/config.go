package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	unmarshalConfigFileError = errors.New("failed to unmarshal configuration file")
	configFileNotFoundError  = errors.New("configuration file not found")
	noProvidersDefinedError  = errors.New("no providers defined")
	noQueuesDefinedError     = errors.New("no queues defined")
	formatNotSupportedError  = errors.New("format not supported")
)

// Config is the main configuration struct
type Config struct {
	// Providers is a list of sources that will be connected to
	Providers []*ProviderConfig `yaml:"providers" json:"providers"`

	// Queues is a list of queues that will be consumed
	Queues []*QueueConfig `yaml:"queues" json:"queues"`

	// Debug is a flag that enables debug logging
	Debug bool `yaml:"debug" json:"debug"`

	// Metrics is the configuration for the metrics endpoint
	Metrics *MetricsConfig `yaml:"metrics" json:"metrics"`

	// Log is the format of the log
	Log string `yaml:"log,omitempty" json:"log"`

	// Databases is the configuration for the database connections
	Databases []*DatabaseConfig `yaml:"databases" json:"databases"`
}

// LoadConfig loads the configuration from a file which can be either YAML or JSON.
func LoadConfig() (*Config, error) {
	configPath := os.Getenv("KONSUME_CONFIG_PATH")
	if configPath == "" {
		slog.Debug("No configuration path defined, using default path /config/config.yaml")
		configPath = "/config/config.yaml"
	}

	slog.Info("Loading configuration from", "path", configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, configFileNotFoundError
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	cfg := &Config{}
	switch strings.ToLower(filepath.Ext(configPath)) {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, cfg)
	case ".json":
		err = json.Unmarshal(data, cfg)
	default:
		return nil, formatNotSupportedError
	}

	if err != nil {
		return nil, unmarshalConfigFileError
	}

	slog.Info("Loaded configuration successfully")

	err = cfg.ValidateAll()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// ValidateAll validates the configuration
func (c *Config) ValidateAll() error {
	if len(c.Providers) == 0 {
		return noProvidersDefinedError
	}
	if len(c.Queues) == 0 {
		return noQueuesDefinedError
	}

	for _, p := range c.Providers {
		err := p.validateProvider()
		if err != nil {
			return err
		}
	}

	for _, d := range c.Databases {
		err := validateDatabaseConfig(d)
		if err != nil {
			return err
		}
	}

	for _, q := range c.Queues {
		err := q.validateQueue(c.Providers, c.Databases)
		if err != nil {
			return err
		}
	}

	if c.Metrics != nil {
		err := c.Metrics.validateMetrics()
		if err != nil {
			return err
		}
	}

	if c.Log == "" {
		c.Log = "text"
	}

	slog.Debug("Configuration validated successfully")

	return nil
}
