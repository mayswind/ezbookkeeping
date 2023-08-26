package locales

// DefaultLanguage represents the default language
var DefaultLanguage = en

// AllLanguages represents all the supported language
var AllLanguages = map[string]*LocaleInfo{
	"en": {
		Content: en,
	},
	"zh-Hans": {
		Content: zhHans,
	},
}

func GetLocaleTextItems(locale string) *LocaleTextItems {
	localeInfo, exists := AllLanguages[locale]

	if exists {
		return localeInfo.Content
	}

	return DefaultLanguage
}
