package locales

import "github.com/mayswind/ezbookkeeping/pkg/models"

// DefaultLanguage represents the default language
var DefaultLanguage = en

// AllLanguages represents all the supported language
// To add new languages, please refer to https://ezbookkeeping.mayswind.net/translating
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

func IsDecimalSeparatorEqualsDigitGroupingSymbol(decimalSeparator models.DecimalSeparator, digitGroupingSymbol models.DigitGroupingSymbol, locale string) bool {
	if decimalSeparator == models.DECIMAL_SEPARATOR_DEFAULT && digitGroupingSymbol == models.DIGIT_GROUPING_SYMBOL_DEFAULT {
		return false
	}

	if byte(decimalSeparator) == byte(digitGroupingSymbol) {
		return true
	}

	localeTextItems := GetLocaleTextItems(locale)

	if decimalSeparator == models.DECIMAL_SEPARATOR_DEFAULT {
		decimalSeparator = localeTextItems.DefaultTypes.DecimalSeparator
	}

	if digitGroupingSymbol == models.DIGIT_GROUPING_SYMBOL_DEFAULT {
		digitGroupingSymbol = localeTextItems.DefaultTypes.DigitGroupingSymbol
	}

	return byte(decimalSeparator) == byte(digitGroupingSymbol)
}
