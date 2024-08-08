package settings

// ConfigContainer contains the current setting config
type ConfigContainer struct {
	Current *Config
}

// Initialize a config container singleton instance
var (
	Version    string
	CommitHash string
	Container  = &ConfigContainer{}
)

// SetCurrentConfig sets the current config by a given config
func SetCurrentConfig(config *Config) {
	Container.Current = config
}

// GetAfterRegisterNotificationContent returns the notification content displayed each time users register
func (c *ConfigContainer) GetAfterRegisterNotificationContent(userLanguage string, clientLanguage string) string {
	language := userLanguage

	if language == "" {
		language = clientLanguage
	}

	if !c.Current.AfterRegisterNotification.Enabled {
		return ""
	}

	if multiLanguageContent, exists := c.Current.AfterRegisterNotification.MultiLanguageContent[language]; exists {
		return multiLanguageContent
	}

	return c.Current.AfterRegisterNotification.DefaultContent
}

// GetAfterLoginNotificationContent returns the notification content displayed each time users log in
func (c *ConfigContainer) GetAfterLoginNotificationContent(userLanguage string, clientLanguage string) string {
	language := userLanguage

	if language == "" {
		language = clientLanguage
	}

	if !c.Current.AfterLoginNotification.Enabled {
		return ""
	}

	if multiLanguageContent, exists := c.Current.AfterLoginNotification.MultiLanguageContent[language]; exists {
		return multiLanguageContent
	}

	return c.Current.AfterLoginNotification.DefaultContent
}

// GetAfterOpenNotificationContent returns the notification content displayed each time users open the app
func (c *ConfigContainer) GetAfterOpenNotificationContent(userLanguage string, clientLanguage string) string {
	language := userLanguage

	if language == "" {
		language = clientLanguage
	}

	if !c.Current.AfterOpenNotification.Enabled {
		return ""
	}

	if multiLanguageContent, exists := c.Current.AfterOpenNotification.MultiLanguageContent[language]; exists {
		return multiLanguageContent
	}

	return c.Current.AfterOpenNotification.DefaultContent
}
