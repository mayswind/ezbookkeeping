package cmd

import (
	"encoding/json"
	"os"

	"github.com/Paxtiny/oscar/pkg/avatars"
	"github.com/Paxtiny/oscar/pkg/core"
	"github.com/Paxtiny/oscar/pkg/datastore"
	"github.com/Paxtiny/oscar/pkg/duplicatechecker"
	"github.com/Paxtiny/oscar/pkg/exchangerates"
	"github.com/Paxtiny/oscar/pkg/llm"
	"github.com/Paxtiny/oscar/pkg/log"
	"github.com/Paxtiny/oscar/pkg/mail"
	"github.com/Paxtiny/oscar/pkg/settings"
	"github.com/Paxtiny/oscar/pkg/storage"
	"github.com/Paxtiny/oscar/pkg/utils"
	"github.com/Paxtiny/oscar/pkg/uuid"
)

func initializeSystem(c *core.CliContext) (*settings.Config, error) {
	var err error
	configFilePath := c.String("conf-path")
	isDisableBootLog := c.Bool("no-boot-log")

	if configFilePath != "" {
		if _, err = os.Stat(configFilePath); err != nil {
			if !isDisableBootLog {
				log.BootErrorf(c, "[initializer.initializeSystem] cannot load configuration from custom config path %s, because file not exists", configFilePath)
			}
			return nil, err
		}

		if !isDisableBootLog {
			log.BootInfof(c, "[initializer.initializeSystem] will loading configuration from custom config path %s", configFilePath)
		}
	} else {
		configFilePath, err = settings.GetDefaultConfigFilePath()

		if err != nil {
			if !isDisableBootLog {
				log.BootErrorf(c, "[initializer.initializeSystem] cannot get default configuration path, because %s", err.Error())
			}
			return nil, err
		}

		if !isDisableBootLog {
			log.BootInfof(c, "[initializer.initializeSystem] will load configuration from default config path %s", configFilePath)
		}
	}

	config, err := settings.LoadConfiguration(configFilePath)

	if err != nil {
		if !isDisableBootLog {
			log.BootErrorf(c, "[initializer.initializeSystem] cannot load configuration, because %s", err.Error())
		}
		return nil, err
	}

	if config.SecretKeyNoSet {
		log.BootWarnf(c, "[initializer.initializeSystem] \"secret_key\" in config file is not set, please change it to keep your user data safe")
	}

	settings.SetCurrentConfig(config)

	err = datastore.InitializeDataStore(config)

	if err != nil {
		if !isDisableBootLog {
			log.BootErrorf(c, "[initializer.initializeSystem] initializes data store failed, because %s", err.Error())
		}
		return nil, err
	}

	err = log.SetLoggerConfiguration(config, isDisableBootLog)

	if err != nil {
		if !isDisableBootLog {
			log.BootErrorf(c, "[initializer.initializeSystem] sets logger configuration failed, because %s", err.Error())
		}
		return nil, err
	}

	err = storage.InitializeStorageContainer(config)

	if err != nil {
		if !isDisableBootLog {
			log.BootErrorf(c, "[initializer.initializeSystem] initializes object storage failed, because %s", err.Error())
		}
		return nil, err
	}

	err = llm.InitializeLargeLanguageModelProvider(config)

	if err != nil {
		if !isDisableBootLog {
			log.BootErrorf(c, "[initializer.initializeSystem] initializes large language model provider failed, because %s", err.Error())
		}
		return nil, err
	}

	err = uuid.InitializeUuidGenerator(config)

	if err != nil {
		if !isDisableBootLog {
			log.BootErrorf(c, "[initializer.initializeSystem] initializes uuid generator failed, because %s", err.Error())
		}
		return nil, err
	}

	err = duplicatechecker.InitializeDuplicateChecker(config)

	if err != nil {
		if !isDisableBootLog {
			log.BootErrorf(c, "[initializer.initializeSystem] initializes duplicate checker failed, because %s", err.Error())
		}
		return nil, err
	}

	err = avatars.InitializeAvatarProvider(config)

	if err != nil {
		if !isDisableBootLog {
			log.BootErrorf(c, "[initializer.initializeSystem] initializes avatar provider failed, because %s", err.Error())
		}
		return nil, err
	}

	err = mail.InitializeMailer(config)

	if err != nil {
		if !isDisableBootLog {
			log.BootErrorf(c, "[initializer.initializeSystem] initializes mailer failed, because %s", err.Error())
		}
		return nil, err
	}

	err = exchangerates.InitializeExchangeRatesDataSource(config)

	if err != nil {
		if !isDisableBootLog {
			log.BootErrorf(c, "[initializer.initializeSystem] initializes exchange rates data source failed, because %s", err.Error())
		}
		return nil, err
	}

	cfgJson, _ := json.Marshal(getConfigWithoutSensitiveData(config))

	if !isDisableBootLog {
		log.BootInfof(c, "[initializer.initializeSystem] has loaded configuration %s", cfgJson)
	}

	return config, nil
}

func getConfigWithoutSensitiveData(config *settings.Config) *settings.Config {
	clonedConfig := &settings.Config{}
	err := utils.Clone(config, clonedConfig)

	if err != nil {
		return config
	}

	if clonedConfig.DatabaseConfig.DatabasePassword != "" {
		clonedConfig.DatabaseConfig.DatabasePassword = "****"
	}

	if clonedConfig.SMTPConfig.SMTPPasswd != "" {
		clonedConfig.SMTPConfig.SMTPPasswd = "****"
	}

	if clonedConfig.MinIOConfig.SecretAccessKey != "" {
		clonedConfig.MinIOConfig.SecretAccessKey = "****"
	}

	if clonedConfig.SecretKey != "" {
		clonedConfig.SecretKey = "****"
	}

	if clonedConfig.AmapApplicationSecret != "" {
		clonedConfig.AmapApplicationSecret = "****"
	}

	if clonedConfig.WebDAVConfig != nil && clonedConfig.WebDAVConfig.Password != "" {
		clonedConfig.WebDAVConfig.Password = "****"
	}

	if clonedConfig.ReceiptImageRecognitionLLMConfig != nil {
		if clonedConfig.ReceiptImageRecognitionLLMConfig.OpenAIAPIKey != "" {
			clonedConfig.ReceiptImageRecognitionLLMConfig.OpenAIAPIKey = "****"
		}

		if clonedConfig.ReceiptImageRecognitionLLMConfig.OpenAICompatibleAPIKey != "" {
			clonedConfig.ReceiptImageRecognitionLLMConfig.OpenAICompatibleAPIKey = "****"
		}

		if clonedConfig.ReceiptImageRecognitionLLMConfig.AnthropicCompatibleAPIKey != "" {
			clonedConfig.ReceiptImageRecognitionLLMConfig.AnthropicCompatibleAPIKey = "****"
		}

		if clonedConfig.ReceiptImageRecognitionLLMConfig.AnthropicAPIKey != "" {
			clonedConfig.ReceiptImageRecognitionLLMConfig.AnthropicAPIKey = "****"
		}

		if clonedConfig.ReceiptImageRecognitionLLMConfig.OpenRouterAPIKey != "" {
			clonedConfig.ReceiptImageRecognitionLLMConfig.OpenRouterAPIKey = "****"
		}

		if clonedConfig.ReceiptImageRecognitionLLMConfig.LMStudioToken != "" {
			clonedConfig.ReceiptImageRecognitionLLMConfig.LMStudioToken = "****"
		}

		if clonedConfig.ReceiptImageRecognitionLLMConfig.GoogleAIAPIKey != "" {
			clonedConfig.ReceiptImageRecognitionLLMConfig.GoogleAIAPIKey = "****"
		}
	}

	if clonedConfig.OAuth2ClientSecret != "" {
		clonedConfig.OAuth2ClientSecret = "****"
	}

	return clonedConfig
}
