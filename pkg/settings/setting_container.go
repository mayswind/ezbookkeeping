package settings

// ConfigContainer contains the current setting config
type ConfigContainer struct {
	Current *Config
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
	Container.Current = config
}
