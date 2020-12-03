package cmd

import (
	"encoding/json"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/mayswind/lab/pkg/datastore"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/settings"
	"github.com/mayswind/lab/pkg/utils"
	"github.com/mayswind/lab/pkg/uuid"
)

func initializeSystem(c *cli.Context) (*settings.Config, error) {
	var err error
	configFilePath := c.String("conf-path")

	if configFilePath != "" {
		if _, err = os.Stat(configFilePath); err != nil {
			log.BootErrorf("[initializer.initializeSystem] cannot load configuration from custom config path %s, because file not exists", configFilePath)
			return nil, err
		}

		log.BootInfof("[initializer.initializeSystem] will loading configuration from custom config path %s", configFilePath)
	} else {
		configFilePath, err = settings.GetDefaultConfigFilePath()

		if err != nil {
			log.BootErrorf("[initializer.initializeSystem] cannot get default configuration path, because %s", err.Error())
			return nil, err
		}

		log.BootInfof("[initializer.initializeSystem] will load configuration from default config path %s", configFilePath)
	}

	config, err := settings.LoadConfiguration(configFilePath)

	if err != nil {
		log.BootErrorf("[initializer.initializeSystem] cannot load configuration, because %s", err.Error())
		return nil, err
	}

	settings.SetCurrentConfig(config)

	err = datastore.InitializeDataStore(config)

	if err != nil {
		log.BootErrorf("[initializer.initializeSystem] initializes data store failed, because %s", err.Error())
		return nil, err
	}

	err = log.SetLoggerConfiguration(config)

	if err != nil {
		log.BootErrorf("[initializer.initializeSystem] sets logger configuration failed, because %s", err.Error())
		return nil, err
	}

	err = uuid.InitializeUuidGenerator(config)

	if err != nil {
		log.BootErrorf("[initializer.initializeSystem] initializes uuid generator failed, because %s", err.Error())
		return nil, err
	}

	cfgJson, _ := json.Marshal(getConfigWithNoSensitiveData(config))
	log.BootInfof("[initializer.initializeSystem] has loaded configuration %s", cfgJson)

	return config, nil
}

func getConfigWithNoSensitiveData(config *settings.Config) *settings.Config {
	clonedConfig := &settings.Config{}
	err := utils.Clone(config, clonedConfig)

	if err != nil {
		return config
	}

	clonedConfig.DatabaseConfig.DatabasePassword = "****"
	clonedConfig.SecretKey = "****"

	return clonedConfig
}
