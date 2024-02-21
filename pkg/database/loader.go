package database

import (
	"fmt"
	"plugin"
	"runtime"

	"github.com/bugrakocabay/konsume/pkg/common"
)

// LoadDatabasePlugin loads the database plugin based on the database type
func LoadDatabasePlugin(dbType string) (Database, error) {
	pluginPath := getPluginPath(dbType)
	if pluginPath == "" {
		return nil, fmt.Errorf("unsupported database type or platform: %s", dbType)
	}

	return loadPlugin(pluginPath)
}

func getPluginPath(dbType string) string {
	var pluginFile string
	switch dbType {
	case common.DatabaseTypePostgresql:
		pluginFile = "postgres"
	default:
		return ""
	}

	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("/root/plugins/%s-linux.so", pluginFile)
	case "darwin":
		return fmt.Sprintf("/root/plugins/%s-darwin.so", pluginFile)
	default:
		return ""
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
