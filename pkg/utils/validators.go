package utils

import "regexp"

var (
	usernamePattern                    = regexp.MustCompile("^(?i)[a-z0-9_-]+$")
	emailPattern                       = regexp.MustCompile("^(?i)(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])$")
	hexRGBColorPattern                 = regexp.MustCompile("^(?i)([0-9a-f]{6}|[0-9a-f]{3})$")
	longDateTimePattern                = regexp.MustCompile("^([1-9][0-9]{3})-(0[1-9]|1[0-2])-(0[1-9]|1[0-9]|2[0-9]|3[01]) ([0-1][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$")
	longDateTimeWithoutSecondPattern   = regexp.MustCompile("^([1-9][0-9]{3})-(0[1-9]|1[0-2])-(0[1-9]|1[0-9]|2[0-9]|3[01]) ([0-1][0-9]|2[0-3]):([0-5][0-9])$")
	longDatePattern                    = regexp.MustCompile("^([1-9][0-9]{3})-(0[1-9]|1[0-2])-(0[1-9]|1[0-9]|2[0-9]|3[01])$")
	longOrShortYearMonthDayDatePattern = regexp.MustCompile("^(([1-9][0-9])?[0-9]{2})[-/.']([1-9]|0[1-9]|1[0-2])[-/.']([1-9]|0[1-9]|1[0-9]|2[0-9]|3[01])$")
	longOrShortMonthDayYearDatePattern = regexp.MustCompile("^([1-9]|0[1-9]|1[0-2])[-/.']([1-9]|0[1-9]|1[0-9]|2[0-9]|3[01])[-/.'](([1-9][0-9])?[0-9]{2})$")
	longOrShortDayMonthYearDatePattern = regexp.MustCompile("^([1-9]|0[1-9]|1[0-9]|2[0-9]|3[01])[-/.']([1-9]|0[1-9]|1[0-2])[-/.'](([1-9][0-9])?[0-9]{2})$")
)

// IsValidUsername reports whether username is valid
func IsValidUsername(username string) bool {
	return len(username) <= 32 && usernamePattern.MatchString(username)
}

// IsValidEmail reports whether email is valid
func IsValidEmail(email string) bool {
	return len(email) <= 100 && emailPattern.MatchString(email)
}

// IsValidNickName reports whether user nick name is valid
func IsValidNickName(nickname string) bool {
	return len(nickname) <= 64
}

// IsValidHexRGBColor reports whether color is valid
func IsValidHexRGBColor(color string) bool {
	return hexRGBColorPattern.MatchString(color)
}

// IsValidLongDateTimeFormat reports whether long date time is valid format
func IsValidLongDateTimeFormat(datetime string) bool {
	return longDateTimePattern.MatchString(datetime)
}

// IsValidLongDateTimeWithoutSecondFormat reports long date time without seconds is valid format
func IsValidLongDateTimeWithoutSecondFormat(datetime string) bool {
	return longDateTimeWithoutSecondPattern.MatchString(datetime)
}

// IsValidLongDateFormat reports long date is valid format
func IsValidLongDateFormat(date string) bool {
	return longDatePattern.MatchString(date)
}

// IsValidYearMonthDayLongOrShortDateFormat reports long date is valid format
func IsValidYearMonthDayLongOrShortDateFormat(date string) bool {
	return longOrShortYearMonthDayDatePattern.MatchString(date)
}

// IsValidMonthDayYearLongOrShortDateFormat reports long date is valid format
func IsValidMonthDayYearLongOrShortDateFormat(date string) bool {
	return longOrShortMonthDayYearDatePattern.MatchString(date)
}

// IsValidDayMonthYearLongOrShortDateFormat reports long date is valid format
func IsValidDayMonthYearLongOrShortDateFormat(date string) bool {
	return longOrShortDayMonthYearDatePattern.MatchString(date)
}
