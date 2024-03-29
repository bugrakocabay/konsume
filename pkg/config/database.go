package config

import (
	"errors"

	"github.com/bugrakocabay/konsume/pkg/common"
)

var (
	databaseNameNotDefinedError             = errors.New("database name not defined")
	databaseTypeNotDefinedError             = errors.New("database type not defined")
	databaseTypeInvalidError                = errors.New("database type invalid")
	databaseConnectionStringNotDefinedError = errors.New("database connection string not defined")
	databaseDatabaseNotDefinedError         = errors.New("database database not defined")
)

// DatabaseConfig is the configuration for the database connections
type DatabaseConfig struct {
	// Name is the name of the database
	Name string `yaml:"name" json:"name"`

	// Type is the type of the database
	Type string `yaml:"type" json:"type"`

	// ConnectionString is the connection string for the database
	ConnectionString string `yaml:"connection-string" json:"connection-string"`

	// Retry is the amount of times the connection should be retried
	Retry int `yaml:"retry,omitempty" json:"retry,omitempty"`

	// Database is the database name for the database
	Database string `yaml:"database,omitempty" json:"database,omitempty"`
}

func validateDatabaseConfig(database *DatabaseConfig) error {
	if len(database.Name) == 0 {
		return databaseNameNotDefinedError
	}
	if len(database.Type) == 0 {
		return databaseTypeNotDefinedError
	}
	if len(database.ConnectionString) == 0 {
		return databaseConnectionStringNotDefinedError
	}
	if database.Type != common.DatabaseTypePostgresql && database.Type != common.DatabaseTypeMongoDB {
		return databaseTypeInvalidError
	}
	if database.Type == common.DatabaseTypeMongoDB && len(database.Database) == 0 {
		return databaseDatabaseNotDefinedError
	}
	return nil
}
