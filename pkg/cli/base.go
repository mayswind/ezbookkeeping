package cli

import "github.com/mayswind/ezbookkeeping/pkg/settings"

// CliUsingConfig represents a cli that need to use config
type CliUsingConfig struct {
	container *settings.ConfigContainer
}

// CurrentConfig returns the current config
func (l *CliUsingConfig) CurrentConfig() *settings.Config {
	return l.container.GetCurrentConfig()
}
