package database

import (
	"plugin"

	"github.com/bugrakocabay/konsume/pkg/common"
)

// LoadDatabasePlugin loads the database plugin based on the database type
func LoadDatabasePlugin(dbType string) (Database, error) {
	switch dbType {
	case common.DatabaseTypePostgresql:
		const path = "/root/plugins/postgres.so"
		return loadPlugin(path)
	default:
		return nil, nil
	}
}

// loadPlugin loads the plugin from the given path
func loadPlugin(path string) (Database, error) {
	plug, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	symbol, err := plug.Lookup("Plugin")
	if err != nil {
		return nil, err
	}
	postgresPlugin, ok := symbol.(Database)
	if !ok {
		return nil, err
	}
	return postgresPlugin, nil
}
