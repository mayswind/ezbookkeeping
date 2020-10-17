package settings

type ConfigContainer struct {
	Current *Config
}

var (
	Container = &ConfigContainer{}
)

func SetCurrentConfig(config *Config) {
	Container.Current = config
}
