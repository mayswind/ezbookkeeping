package settings

// ConfigContainer contains the current setting config
type ConfigContainer struct {
	current *Config
}

// Initialize a config container singleton instance
var (
	Container = &ConfigContainer{}
)

// NewConfigContainer returns a new config container with the given config
func NewConfigContainer(config *Config) *ConfigContainer {
	return &ConfigContainer{
		current: config,
	}
}

// SetCurrentConfig sets the current config by a given config
func SetCurrentConfig(config *Config) {
	Container.current = config
}

// GetCurrentConfig returns the current config
func (c *ConfigContainer) GetCurrentConfig() *Config {
	return c.current
}
