package api

import (
	"fmt"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const ezbookkeepingServerSettingsGlobalVariableName = "EZBOOKKEEPING_SERVER_SETTINGS"
const ezbookkeepingServerSettingsGlobalVariableFullName = "window." + ezbookkeepingServerSettingsGlobalVariableName
const ezbookkeepingServerSettingsJavascriptFileHeader = ezbookkeepingServerSettingsGlobalVariableFullName +
	"=" + ezbookkeepingServerSettingsGlobalVariableFullName + "||{};\n"

// ServerSettingsApi represents server settings api
type ServerSettingsApi struct {
	ApiUsingConfig
}

// Initialize a server settings api singleton instance
var (
	ServerSettings = &ServerSettingsApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
	}
)

// ServerSettingsJavascriptHandler returns the javascript contains server settings
func (a *ServerSettingsApi) ServerSettingsJavascriptHandler(c *core.WebContext) ([]byte, string, *errs.Error) {
	config := a.CurrentConfig()
	builder := &strings.Builder{}
	builder.WriteString(ezbookkeepingServerSettingsJavascriptFileHeader)

	a.appendBooleanSetting(builder, "r", config.EnableUserRegister)
	a.appendBooleanSetting(builder, "f", config.EnableUserForgetPassword)
	a.appendBooleanSetting(builder, "v", config.EnableUserVerifyEmail)
	a.appendBooleanSetting(builder, "p", config.EnableTransactionPictures)
	a.appendBooleanSetting(builder, "s", config.EnableScheduledTransaction)
	a.appendBooleanSetting(builder, "e", config.EnableDataExport)
	a.appendBooleanSetting(builder, "i", config.EnableDataImport)

	if config.EnableMCPServer {
		a.appendBooleanSetting(builder, "mcp", config.EnableMCPServer)
	}

	if config.LoginPageTips.Enabled {
		a.appendMultiLanguageTipSetting(builder, "lpt", config.LoginPageTips)
	}

	a.appendStringSetting(builder, "m", config.MapProvider)

	if config.EnableMapDataFetchProxy &&
		(config.MapProvider == settings.OpenStreetMapProvider ||
			config.MapProvider == settings.OpenStreetMapHumanitarianStyleProvider ||
			config.MapProvider == settings.OpenTopoMapProvider ||
			config.MapProvider == settings.OPNVKarteMapProvider ||
			config.MapProvider == settings.CyclOSMMapProvider ||
			config.MapProvider == settings.CartoDBMapProvider ||
			config.MapProvider == settings.TomTomMapProvider ||
			config.MapProvider == settings.TianDiTuProvider ||
			config.MapProvider == settings.CustomProvider) {
		a.appendBooleanSetting(builder, "mp", config.EnableMapDataFetchProxy)
	}

	if config.MapProvider == settings.CustomProvider {
		a.appendStringSetting(builder, "cmzl", fmt.Sprintf("%d-%d-%d", config.CustomMapTileServerMinZoomLevel, config.CustomMapTileServerMaxZoomLevel, config.CustomMapTileServerDefaultZoomLevel))

		if !config.EnableMapDataFetchProxy {
			a.appendStringSetting(builder, "cmsu", config.CustomMapTileServerTileLayerUrl)

			if config.CustomMapTileServerAnnotationLayerUrl != "" {
				a.appendStringSetting(builder, "cmau", config.CustomMapTileServerAnnotationLayerUrl)
			}
		} else {
			if config.CustomMapTileServerAnnotationLayerUrl != "" {
				a.appendBooleanSetting(builder, "cmap", config.EnableMapDataFetchProxy)
			}
		}
	}

	if config.MapProvider == settings.TomTomMapProvider && config.TomTomMapAPIKey != "" && !config.EnableMapDataFetchProxy {
		a.appendStringSetting(builder, "tmak", config.TomTomMapAPIKey)
	}

	if config.MapProvider == settings.TianDiTuProvider && config.TianDiTuAPIKey != "" && !config.EnableMapDataFetchProxy {
		a.appendStringSetting(builder, "tdak", config.TianDiTuAPIKey)
	}

	if config.MapProvider == settings.GoogleMapProvider && config.GoogleMapAPIKey != "" {
		a.appendStringSetting(builder, "gmak", config.GoogleMapAPIKey)
	}

	if config.MapProvider == settings.BaiduMapProvider && config.BaiduMapAK != "" {
		a.appendStringSetting(builder, "bmak", config.BaiduMapAK)
	}

	if config.MapProvider == settings.AmapProvider && config.AmapApplicationKey != "" {
		a.appendStringSetting(builder, "amak", config.AmapApplicationKey)
	}

	if config.MapProvider == settings.AmapProvider && config.AmapSecurityVerificationMethod != "" {
		a.appendStringSetting(builder, "amsv", config.AmapSecurityVerificationMethod)

		if config.AmapSecurityVerificationMethod == settings.AmapSecurityVerificationExternalProxyMethod {
			a.appendStringSetting(builder, "amep", config.AmapApiExternalProxyUrl)
		}

		if config.AmapSecurityVerificationMethod == settings.AmapSecurityVerificationPlainTextMethod {
			a.appendStringSetting(builder, "amas", config.AmapApplicationSecret)
		}
	}

	if config.ExchangeRatesRequestTimeoutExceedDefaultValue {
		a.appendIntegerSetting(builder, "errt", int(config.ExchangeRatesRequestTimeout))
	}

	return []byte(builder.String()), "", nil
}

func (a *ServerSettingsApi) appendStringSetting(builder *strings.Builder, key string, value string) {
	builder.WriteString(ezbookkeepingServerSettingsGlobalVariableFullName)
	builder.WriteString("[")
	a.appendEncodedString(builder, key)
	builder.WriteString("]=")

	a.appendEncodedString(builder, value)

	builder.WriteString(";\n")
}

func (a *ServerSettingsApi) appendMultiLanguageTipSetting(builder *strings.Builder, key string, value settings.TipConfig) {
	builder.WriteString(ezbookkeepingServerSettingsGlobalVariableFullName)
	builder.WriteString("[")
	a.appendEncodedString(builder, key)
	builder.WriteString("]={\n")

	builder.WriteString("'default'")
	builder.WriteRune(':')
	a.appendEncodedString(builder, value.DefaultContent)

	for languageTag, content := range value.MultiLanguageContent {
		builder.WriteString(",\n")
		a.appendEncodedString(builder, languageTag)
		builder.WriteRune(':')
		a.appendEncodedString(builder, content)
	}

	builder.WriteString("\n};\n")
}

func (a *ServerSettingsApi) appendBooleanSetting(builder *strings.Builder, key string, value bool) {
	builder.WriteString(ezbookkeepingServerSettingsGlobalVariableFullName)
	builder.WriteString("[")
	a.appendEncodedString(builder, key)
	builder.WriteString("]=")

	if value {
		builder.WriteRune('1')
	} else {
		builder.WriteRune('0')
	}

	builder.WriteString(";\n")
}

func (a *ServerSettingsApi) appendIntegerSetting(builder *strings.Builder, key string, value int) {
	builder.WriteString(ezbookkeepingServerSettingsGlobalVariableFullName)
	builder.WriteString("[")
	a.appendEncodedString(builder, key)
	builder.WriteString("]=")
	builder.WriteString(utils.IntToString(value))
	builder.WriteString(";\n")
}

func (a *ServerSettingsApi) appendEncodedString(builder *strings.Builder, content string) {
	builder.WriteRune('\'')
	runes := []rune(content)

	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		case '\\':
			builder.WriteRune('\\')
			builder.WriteRune('\\')
		case '\'':
			builder.WriteRune('\\')
			builder.WriteRune('\'')
		case '\n':
			builder.WriteRune('\\')
			builder.WriteRune('n')
		case '\r':
			builder.WriteRune('\\')
			builder.WriteRune('r')
		case '\t':
			builder.WriteRune('\\')
			builder.WriteRune('t')
		case '\f':
			builder.WriteRune('\\')
			builder.WriteRune('f')
		case '\b':
			builder.WriteRune('\\')
			builder.WriteRune('b')
		default:
			builder.WriteRune(runes[i])
		}
	}

	builder.WriteRune('\'')
}
