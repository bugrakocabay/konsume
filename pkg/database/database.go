package database

import "github.com/bugrakocabay/konsume/pkg/config"

// Database is an interface that defines the methods that a database should implement
type Database interface {
	Connect(connectionString, dbName string) error
	Insert(data map[string]interface{}, dbRouteConfig config.DatabaseRouteConfig) error
	Close() error
}
