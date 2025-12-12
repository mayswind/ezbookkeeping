package settings

import (
	"fmt"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

// ConfigContainer contains the current setting config
type ConfigContainer struct {
	current *Config
}

// Initialize a config container singleton instance
var (
	Version    string
	CommitHash string
	BuildTime  string
	Container  = &ConfigContainer{}
)

// SetCurrentConfig sets the current config by a given config
func SetCurrentConfig(config *Config) {
	Container.current = config
}

// GetCurrentConfig returns the current config
func (c *ConfigContainer) GetCurrentConfig() *Config {
	return c.current
}

func GetUserAgent() string {
	if Version == "" {
		return core.ApplicationName
	}

	return fmt.Sprintf("%s/%s", core.ApplicationName, Version)
}
