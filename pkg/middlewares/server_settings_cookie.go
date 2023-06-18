package middlewares

import (
	"encoding/base64"
	"fmt"
	"net/url"
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

		if config.MapProvider == settings.OpenStreetMapProvider && config.EnableMapDataFetchProxy {
			settingsArr = append(settingsArr, buildBooleanSetting("mp", config.EnableMapDataFetchProxy))
		}

		if config.MapProvider == settings.GoogleMapProvider && config.GoogleMapAPIKey != "" {
			settingsArr = append(settingsArr, buildEncodedStringSetting("gmak", config.GoogleMapAPIKey))
		}

		if config.MapProvider == settings.BaiduMapProvider && config.BaiduMapAK != "" {
			settingsArr = append(settingsArr, buildEncodedStringSetting("bmak", config.BaiduMapAK))
		}

		if config.MapProvider == settings.AmapProvider && config.AmapApplicationKey != "" {
			settingsArr = append(settingsArr, buildEncodedStringSetting("amak", config.AmapApplicationKey))
		}

		if config.MapProvider == settings.AmapProvider && config.AmapSecurityVerificationMethod != "" {
			settingsArr = append(settingsArr, buildStringSetting("amsv", strings.Replace(config.AmapSecurityVerificationMethod, "_", "", -1)))

			if config.AmapSecurityVerificationMethod == settings.AmapSecurityVerificationExternalProxyMethod {
				settingsArr = append(settingsArr, buildEncodedStringSetting("amep", config.AmapApiExternalProxyUrl))
			}

			if config.AmapSecurityVerificationMethod == settings.AmapSecurityVerificationPlainTextMethod {
				settingsArr = append(settingsArr, buildEncodedStringSetting("amas", config.AmapApplicationSecret))
			}
		}

		bundledSettings := strings.Join(settingsArr, "_")
		c.SetCookie(settingsCookieName, bundledSettings, int(config.TokenExpiredTime), "", "", false, false)

		c.Next()
	}
}

func buildStringSetting(key string, value string) string {
	return fmt.Sprintf("%s.%s", key, value)
}

func buildEncodedStringSetting(key string, value string) string {
	urlEncodedValue := url.QueryEscape(value)
	base64Value := base64.StdEncoding.EncodeToString([]byte(urlEncodedValue))
	return fmt.Sprintf("%s.%s", key, base64Value)
}

func buildBooleanSetting(key string, value bool) string {
	if value {
		return fmt.Sprintf("%s.1", key)
	} else {
		return fmt.Sprintf("%s.0", key)
	}
}
