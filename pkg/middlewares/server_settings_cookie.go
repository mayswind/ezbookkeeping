package middlewares

import (
	"fmt"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const settingsCookieName = "ebk_server_settings"

// ServerSettingsCookie adds server settings to cookies in response
func ServerSettingsCookie(config *settings.Config) core.MiddlewareHandlerFunc {
	return func(c *core.Context) {
		settingsArr := []string{
			buildBooleanSetting("r", config.EnableUserRegister),
			buildBooleanSetting("e", config.EnableDataExport),
			buildStringSetting("m", config.MapProvider),
		}

		if config.EnableMapDataFetchProxy {
			settingsArr = append(settingsArr, buildBooleanSetting("mp", config.EnableMapDataFetchProxy))
		}

		if config.GoogleMapAPIKey != "" {
			settingsArr = append(settingsArr, buildStringSetting("gmak", config.GoogleMapAPIKey))
		}

		if config.BaiduMapAK != "" {
			settingsArr = append(settingsArr, buildStringSetting("bmak", config.BaiduMapAK))
		}

		if config.AMapApplicationKey != "" {
			settingsArr = append(settingsArr, buildStringSetting("amak", config.AMapApplicationKey))
		}

		if config.AMapSecurityVerificationMethod != "" {
			settingsArr = append(settingsArr, buildStringSetting("amsv", config.AMapSecurityVerificationMethod))

			if config.AMapSecurityVerificationMethod == settings.AmapSecurityVerificationPlainMethod {
				settingsArr = append(settingsArr, buildStringSetting("amas", config.AMapApplicationSecret))
			}
		}

		bundledSettings := strings.Join(settingsArr, "_")
		c.SetCookie(settingsCookieName, bundledSettings, int(config.TokenExpiredTime), "", "", false, false)

		c.Next()
	}
}

func buildStringSetting(key string, value string) string {
	return fmt.Sprintf("%s.%s", key, strings.Replace(value, ".", "-", -1))
}

func buildBooleanSetting(key string, value bool) string {
	if value {
		return fmt.Sprintf("%s.1", key)
	} else {
		return fmt.Sprintf("%s.0", key)
	}
}
