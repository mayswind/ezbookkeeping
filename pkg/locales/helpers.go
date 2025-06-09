package locales

import "github.com/mayswind/ezbookkeeping/pkg/core"

// GetLocaleTextItems returns the locale text items for the specified locale
func GetLocaleTextItems(locale string) *LocaleTextItems {
	localeInfo, exists := AllLanguages[locale]

	if exists {
		return localeInfo.Content
	}

	return DefaultLanguage
}

// IsDecimalSeparatorEqualsDigitGroupingSymbol returns whether the decimal separator equals to the digit grouping symbol in the specified locale
func IsDecimalSeparatorEqualsDigitGroupingSymbol(decimalSeparator core.DecimalSeparator, digitGroupingSymbol core.DigitGroupingSymbol, locale string) bool {
	if decimalSeparator == core.DECIMAL_SEPARATOR_DEFAULT && digitGroupingSymbol == core.DIGIT_GROUPING_SYMBOL_DEFAULT {
		return false
	}

	if byte(decimalSeparator) == byte(digitGroupingSymbol) {
		return true
	}

	localeTextItems := GetLocaleTextItems(locale)

	if decimalSeparator == core.DECIMAL_SEPARATOR_DEFAULT {
		decimalSeparator = localeTextItems.DefaultTypes.DecimalSeparator
	}

	if digitGroupingSymbol == core.DIGIT_GROUPING_SYMBOL_DEFAULT {
		digitGroupingSymbol = localeTextItems.DefaultTypes.DigitGroupingSymbol
	}

	return byte(decimalSeparator) == byte(digitGroupingSymbol)
}
