package database

import (
	"fmt"
	"log/slog"
	"os"
	"plugin"
	"runtime"

	"github.com/bugrakocabay/konsume/pkg/common"
)

// LoadDatabasePlugin loads the database plugin based on the database type
func LoadDatabasePlugin(dbType string) (Database, error) {
	pluginPath := getPluginPath(dbType)
	if pluginPath == "" {
		return nil, fmt.Errorf("unsupported database type/platform or file doesn't exist: %s", dbType)
	}

	return loadPlugin(pluginPath)
}

// getPluginPath returns the plugin path based on the database type
func getPluginPath(dbType string) string {
	var pluginFile string
	switch dbType {
	case common.DatabaseTypePostgresql:
		pluginFile = "postgres"
	default:
		return ""
	}

	rootPath := getRootPath()
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		slog.Error("Plugin path does not exist", "path", rootPath)
		return ""
	}
	slog.Debug("Loading plugin", "path", rootPath)

	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("%s/%s-linux.so", rootPath, pluginFile)
	case "darwin":
		return fmt.Sprintf("%s/%s-darwin.so", rootPath, pluginFile)
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

func getRootPath() string {
	pluginPathEnv := os.Getenv(common.KonsumePluginPath)
	if pluginPathEnv != "" {
		return pluginPathEnv
	}
	if isRunningInDocker() {
		return "/root/plugins"
	}
	return "./plugins"
}

func isRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}
