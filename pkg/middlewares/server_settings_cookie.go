package middlewares

import (
	"fmt"
	"strings"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/settings"
)

const settingsCookieName = "ACP_SETTINGS"

// ServerSettingsCookie adds server settings to cookies in response
func ServerSettingsCookie(config *settings.Config) core.MiddlewareHandlerFunc {
	return func(c *core.Context) {
		settingsArr := []string{
			buildBooleanSetting("r", config.EnableUserRegister),
		}

		bundledSettings := strings.Join(settingsArr, "_")
		c.SetCookie(settingsCookieName, bundledSettings, config.TokenExpiredTime, "", "", false, false)

		c.Next()
	}
}

func buildBooleanSetting(key string, value bool) string {
	if value {
		return fmt.Sprintf("%s.1", key)
	} else {
		return fmt.Sprintf("%s.0", key)
	}
}
