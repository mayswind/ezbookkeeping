package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidUsername_ValidUserName(t *testing.T) {
	username := "foobar"
	actualValue := IsValidUsername(username)
	assert.True(t, actualValue)

	username = "--foo_bar--"
	actualValue = IsValidUsername(username)
	assert.True(t, actualValue)
}

func TestIsValidUsername_InvalidUserName(t *testing.T) {
	username := "foo~bar~"
	actualValue := IsValidUsername(username)
	assert.False(t, actualValue)

	username = "012345678901234567890123456789012"
	actualValue = IsValidUsername(username)
	assert.False(t, actualValue)
}

func TestIsValidEmail_ValidEmail(t *testing.T) {
	email := "foo@bar.com"
	actualValue := IsValidEmail(email)
	assert.True(t, actualValue)

	email = "foo@1.2.3.4"
	actualValue = IsValidEmail(email)
	assert.True(t, actualValue)

	email = "foo_bar@foo.bar"
	actualValue = IsValidEmail(email)
	assert.True(t, actualValue)
}

func TestIsValidEmail_InvalidEmail(t *testing.T) {
	email := "foo"
	actualValue := IsValidEmail(email)
	assert.False(t, actualValue)

	email = "@bar"
	actualValue = IsValidEmail(email)
	assert.False(t, actualValue)

	email = "foo@bar"
	actualValue = IsValidEmail(email)
	assert.False(t, actualValue)

	email = "foo@bar."
	actualValue = IsValidEmail(email)
	assert.False(t, actualValue)

	email = "012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789@foobar.com"
	actualValue = IsValidEmail(email)
	assert.False(t, actualValue)
}

func TestIsValidNickName_ValidNickName(t *testing.T) {
	nickname := "0123456789012345678901234567890123456789012345678901234567890123"
	actualValue := IsValidNickName(nickname)
	assert.True(t, actualValue)
}

func TestIsValidNickName_InvalidNickName(t *testing.T) {
	username := "01234567890123456789012345678901234567890123456789012345678901234"
	actualValue := IsValidNickName(username)
	assert.False(t, actualValue)
}

func TestIsValidHexRGBColor_ValidHexRGBColor(t *testing.T) {
	color := "000000"
	actualValue := IsValidHexRGBColor(color)
	assert.True(t, actualValue)

	color = "000"
	actualValue = IsValidHexRGBColor(color)
	assert.True(t, actualValue)

	color = "e0e0e0"
	actualValue = IsValidHexRGBColor(color)
	assert.True(t, actualValue)

	color = "ffffff"
	actualValue = IsValidHexRGBColor(color)
	assert.True(t, actualValue)

	color = "FFFFFF"
	actualValue = IsValidHexRGBColor(color)
	assert.True(t, actualValue)
}

func TestIsValidHexRGBColor_InvalidHexRGBColor(t *testing.T) {
	color := "f"
	actualValue := IsValidHexRGBColor(color)
	assert.False(t, actualValue)

	color = "fffffff"
	actualValue = IsValidHexRGBColor(color)
	assert.False(t, actualValue)

	color = "gggggg"
	actualValue = IsValidHexRGBColor(color)
	assert.False(t, actualValue)

	color = "#ffffff"
	actualValue = IsValidHexRGBColor(color)
	assert.False(t, actualValue)
}

func TestIsValidLongDateTimeFormat_ValidLongDateTimeFormat(t *testing.T) {
	datetime := "2024-09-01 12:34:56"
	actualValue := IsValidLongDateTimeFormat(datetime)
	assert.True(t, actualValue)

	datetime = "2024-10-01 00:00:00"
	actualValue = IsValidLongDateTimeFormat(datetime)
	assert.True(t, actualValue)

	datetime = "9999-12-31 23:59:59"
	actualValue = IsValidLongDateTimeFormat(datetime)
	assert.True(t, actualValue)
}

func TestIsValidLongDateTimeFormat_InvalidLongDateTimeFormat(t *testing.T) {
	datetime := "2024-09-01"
	actualValue := IsValidLongDateTimeFormat(datetime)
	assert.False(t, actualValue)

	datetime = "2024-09-01 12"
	actualValue = IsValidLongDateTimeFormat(datetime)
	assert.False(t, actualValue)

	datetime = "2024-09-01 12:34"
	actualValue = IsValidLongDateTimeFormat(datetime)
	assert.False(t, actualValue)
}

func TestIsValidLongDateTimeWithoutSecondFormat_ValidLongDateTimeWithoutSecondFormat(t *testing.T) {
	datetime := "2024-09-01 12:34"
	actualValue := IsValidLongDateTimeWithoutSecondFormat(datetime)
	assert.True(t, actualValue)

	datetime = "2024-10-01 00:00"
	actualValue = IsValidLongDateTimeWithoutSecondFormat(datetime)
	assert.True(t, actualValue)

	datetime = "9999-12-31 23:59"
	actualValue = IsValidLongDateTimeWithoutSecondFormat(datetime)
	assert.True(t, actualValue)
}

func TestIsValidLongDateTimeWithoutSecondFormat_InvalidLongDateTimeWithoutSecondFormat(t *testing.T) {
	datetime := "2024-09-01"
	actualValue := IsValidLongDateTimeWithoutSecondFormat(datetime)
	assert.False(t, actualValue)

	datetime = "2024-09-01 12"
	actualValue = IsValidLongDateTimeWithoutSecondFormat(datetime)
	assert.False(t, actualValue)

	datetime = "2024-09-01 12:34:56"
	actualValue = IsValidLongDateTimeWithoutSecondFormat(datetime)
	assert.False(t, actualValue)
}

func TestIsValidLongDateFormat_ValidLongDateFormat(t *testing.T) {
	datetime := "2024-09-01"
	actualValue := IsValidLongDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "9999-12-31"
	actualValue = IsValidLongDateFormat(datetime)
	assert.True(t, actualValue)
}

func TestIsValidLongDateFormat_InvalidLongDateFormat(t *testing.T) {
	datetime := "24-09-01"
	actualValue := IsValidLongDateFormat(datetime)
	assert.False(t, actualValue)

	datetime = "2024-9-1"
	actualValue = IsValidLongDateFormat(datetime)
	assert.False(t, actualValue)

	datetime = "2024-09-1"
	actualValue = IsValidLongDateFormat(datetime)
	assert.False(t, actualValue)

	datetime = "2024-9-01"
	actualValue = IsValidLongDateFormat(datetime)
	assert.False(t, actualValue)

	datetime = "2024-09-01 12"
	actualValue = IsValidLongDateFormat(datetime)
	assert.False(t, actualValue)

	datetime = "2024-09-01 12:34"
	actualValue = IsValidLongDateFormat(datetime)
	assert.False(t, actualValue)

	datetime = "2024-09-01 12:34:56"
	actualValue = IsValidLongDateFormat(datetime)
	assert.False(t, actualValue)
}

func TestIsValidYearMonthDayLongOrShortDateFormat_ValidFormat(t *testing.T) {
	datetime := "2024-09-01"
	actualValue := IsValidYearMonthDayLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "24-09-01"
	actualValue = IsValidYearMonthDayLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "2024-09-1"
	actualValue = IsValidYearMonthDayLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "2024-9-01"
	actualValue = IsValidYearMonthDayLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "2024-9-1"
	actualValue = IsValidYearMonthDayLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "9999-12-31"
	actualValue = IsValidYearMonthDayLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "2024/09/01"
	actualValue = IsValidYearMonthDayLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "2024.09.01"
	actualValue = IsValidYearMonthDayLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "2024'09.01"
	actualValue = IsValidYearMonthDayLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)
}

func TestIsValidMonthDayYearLongOrShortDateFormat_ValidFormat(t *testing.T) {
	datetime := "09-01-2024"
	actualValue := IsValidMonthDayYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "09-01-24"
	actualValue = IsValidMonthDayYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "09-1-2024"
	actualValue = IsValidMonthDayYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "9-01-2024"
	actualValue = IsValidMonthDayYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "9-1-2024"
	actualValue = IsValidMonthDayYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "12-31-9999"
	actualValue = IsValidMonthDayYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "09/01/2024"
	actualValue = IsValidMonthDayYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "09.01.2024"
	actualValue = IsValidMonthDayYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "09/01'2024"
	actualValue = IsValidMonthDayYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)
}

func TestIsValidDayMonthYearLongDateFormat_ValidLongDateFormat(t *testing.T) {
	datetime := "01-09-2024"
	actualValue := IsValidDayMonthYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "01-09-24"
	actualValue = IsValidDayMonthYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "1-09-2024"
	actualValue = IsValidDayMonthYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "01-9-2024"
	actualValue = IsValidDayMonthYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "1-9-2024"
	actualValue = IsValidDayMonthYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "31-12-9999"
	actualValue = IsValidDayMonthYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "01/09/2024"
	actualValue = IsValidDayMonthYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "01.09.2024"
	actualValue = IsValidDayMonthYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)

	datetime = "01/09'2024"
	actualValue = IsValidDayMonthYearLongOrShortDateFormat(datetime)
	assert.True(t, actualValue)
}
