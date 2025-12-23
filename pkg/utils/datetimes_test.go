package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func TestParseNumericYearMonth(t *testing.T) {
	expectedYear := int32(2024)
	expectedMonth := int32(3)
	actualYear, actualMonth, err := ParseNumericYearMonth("2024-03")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedYear, actualYear)
	assert.Equal(t, expectedMonth, actualMonth)
}

func TestFormatUnixTimeToLongDate(t *testing.T) {
	unixTime := int64(1617228083)
	utcTimezone := time.FixedZone("Test Timezone", 0)      // UTC
	utc8Timezone := time.FixedZone("Test Timezone", 28800) // UTC+8

	expectedValue := "2021-03-31"
	actualValue := FormatUnixTimeToLongDate(unixTime, utcTimezone)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2021-04-01"
	actualValue = FormatUnixTimeToLongDate(unixTime, utc8Timezone)
	assert.Equal(t, expectedValue, actualValue)
}

func TestFormatUnixTimeToLongDateTimeWithTimezone(t *testing.T) {
	unixTime := int64(1617228083)
	utcTimezone := time.FixedZone("Test Timezone", 0)      // UTC
	utc8Timezone := time.FixedZone("Test Timezone", 28800) // UTC+8

	expectedValue := "2021-03-31 22:01:23Z"
	actualValue := FormatUnixTimeToLongDateTimeWithTimezone(unixTime, utcTimezone)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2021-04-01 06:01:23+08:00"
	actualValue = FormatUnixTimeToLongDateTimeWithTimezone(unixTime, utc8Timezone)
	assert.Equal(t, expectedValue, actualValue)
}

func TestFormatUnixTimeToLongDateTimeWithTimezoneRFC3339Format(t *testing.T) {
	unixTime := int64(1617228083)
	utcTimezone := time.FixedZone("Test Timezone", 0)      // UTC
	utc8Timezone := time.FixedZone("Test Timezone", 28800) // UTC+8

	expectedValue := "2021-03-31T22:01:23Z"
	actualValue := FormatUnixTimeToLongDateTimeWithTimezoneRFC3339Format(unixTime, utcTimezone)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2021-04-01T06:01:23+08:00"
	actualValue = FormatUnixTimeToLongDateTimeWithTimezoneRFC3339Format(unixTime, utc8Timezone)
	assert.Equal(t, expectedValue, actualValue)
}

func TestFormatUnixTimeToLongDateTime(t *testing.T) {
	unixTime := int64(1617228083)
	utcTimezone := time.FixedZone("Test Timezone", 0)      // UTC
	utc8Timezone := time.FixedZone("Test Timezone", 28800) // UTC+8

	expectedValue := "2021-03-31 22:01:23"
	actualValue := FormatUnixTimeToLongDateTime(unixTime, utcTimezone)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2021-04-01 06:01:23"
	actualValue = FormatUnixTimeToLongDateTime(unixTime, utc8Timezone)
	assert.Equal(t, expectedValue, actualValue)
}

func TestFormatYearMonthDayToLongDateTime(t *testing.T) {
	expectedValue := "2025-06-01 00:00:00"
	actualValue, err := FormatYearMonthDayToLongDateTime("25", "06", "01")
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2025-06-01 00:00:00"
	actualValue, err = FormatYearMonthDayToLongDateTime("25", "6", "1")
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "1990-06-01 00:00:00"
	actualValue, err = FormatYearMonthDayToLongDateTime("90", "06", "01")
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestFormatUnixTimeToLongDateTimeWithoutSecond(t *testing.T) {
	unixTime := int64(1617228083)
	utcTimezone := time.FixedZone("Test Timezone", 0)      // UTC
	utc8Timezone := time.FixedZone("Test Timezone", 28800) // UTC+8

	expectedValue := "2021-03-31 22:01"
	actualValue := FormatUnixTimeToLongDateTimeWithoutSecond(unixTime, utcTimezone)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2021-04-01 06:01"
	actualValue = FormatUnixTimeToLongDateTimeWithoutSecond(unixTime, utc8Timezone)
	assert.Equal(t, expectedValue, actualValue)
}

func TestFormatUnixTimeToYearMonth(t *testing.T) {
	unixTime := int64(1617228083)
	utcTimezone := time.FixedZone("Test Timezone", 0)      // UTC
	utc8Timezone := time.FixedZone("Test Timezone", 28800) // UTC+8

	expectedValue := "2021-03"
	actualValue := FormatUnixTimeToYearMonth(unixTime, utcTimezone)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2021-04"
	actualValue = FormatUnixTimeToYearMonth(unixTime, utc8Timezone)
	assert.Equal(t, expectedValue, actualValue)
}

func TestFormatUnixTimeToNumericYearMonth(t *testing.T) {
	unixTime := int64(1617228083)
	utcTimezone := time.FixedZone("Test Timezone", 0)      // UTC
	utc8Timezone := time.FixedZone("Test Timezone", 28800) // UTC+8

	expectedValue := int32(202103)
	actualValue := FormatUnixTimeToNumericYearMonth(unixTime, utcTimezone)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = int32(202104)
	actualValue = FormatUnixTimeToNumericYearMonth(unixTime, utc8Timezone)
	assert.Equal(t, expectedValue, actualValue)
}

func TestFormatUnixTimeToNumericYearMonthDay(t *testing.T) {
	unixTime := int64(1617228083)
	utcTimezone := time.FixedZone("Test Timezone", 0)      // UTC
	utc8Timezone := time.FixedZone("Test Timezone", 28800) // UTC+8

	expectedValue := int32(20210331)
	actualValue := FormatUnixTimeToNumericYearMonthDay(unixTime, utcTimezone)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = int32(20210401)
	actualValue = FormatUnixTimeToNumericYearMonthDay(unixTime, utc8Timezone)
	assert.Equal(t, expectedValue, actualValue)
}

func TestFormatUnixTimeToNumericLocalDateTime(t *testing.T) {
	unixTime := int64(1617228083)
	utcTimezone := time.FixedZone("Test Timezone", 0)      // UTC
	utc8Timezone := time.FixedZone("Test Timezone", 28800) // UTC+8

	expectedValue := int64(20210331220123)
	actualValue := FormatUnixTimeToNumericLocalDateTime(unixTime, utcTimezone)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = int64(20210401060123)
	actualValue = FormatUnixTimeToNumericLocalDateTime(unixTime, utc8Timezone)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetMinUnixTimeWithSameLocalDateTime(t *testing.T) {
	expectedValue := int64(1690797600)
	actualValue := GetMinUnixTimeWithSameLocalDateTime(1690819200, 480)
	assert.Equal(t, expectedValue, actualValue)

	actualValue = GetMinUnixTimeWithSameLocalDateTime(1690873200, -420)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetMaxUnixTimeWithSameLocalDateTime(t *testing.T) {
	expectedValue := int64(1690891200)
	actualValue := GetMaxUnixTimeWithSameLocalDateTime(1690819200, 480)
	assert.Equal(t, expectedValue, actualValue)

	actualValue = GetMaxUnixTimeWithSameLocalDateTime(1690873200, -420)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateFirstTime(t *testing.T) {
	expectedValue := int64(1690819200)
	actualTime, err := ParseFromLongDateFirstTime("2023-08-01", 480)
	assert.Equal(t, nil, err)

	actualValue := actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateLastTime(t *testing.T) {
	expectedValue := int64(1690905599)
	actualTime, err := ParseFromLongDateLastTime("2023-08-01", 480)
	assert.Equal(t, nil, err)

	actualValue := actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateTimeToMinUnixTime(t *testing.T) {
	expectedValue := int64(1690797600)
	actualTime, err := ParseFromLongDateTimeToMinUnixTime("2023-08-01 00:00:00")
	assert.Equal(t, nil, err)

	actualValue := actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateTimeToMaxUnixTime(t *testing.T) {
	expectedValue := int64(1690891200)
	actualTime, err := ParseFromLongDateTimeToMaxUnixTime("2023-08-01 00:00:00")
	assert.Equal(t, nil, err)

	actualValue := actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateTimeInFixedUtcOffset(t *testing.T) {
	expectedValue := int64(1617228083)
	actualTime, err := ParseFromLongDateTimeInFixedUtcOffset("2021-04-01 06:01:23", 480)
	assert.Equal(t, nil, err)

	actualValue := actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateTimeInTimeZone(t *testing.T) {
	londonLocation, err := time.LoadLocation("Europe/London")
	assert.Equal(t, nil, err)

	// during standard time (UTC+0)
	expectedValue := int64(1577858483)
	actualTime, err := ParseFromLongDateTimeInTimeZone("2020-01-01 06:01:23", londonLocation)
	assert.Equal(t, nil, err)

	actualValue := actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)

	// during daylight saving time (UTC+1)
	expectedValue = int64(1619845283)
	actualTime, err = ParseFromLongDateTimeInTimeZone("2021-05-01 06:01:23", londonLocation)
	assert.Equal(t, nil, err)

	actualValue = actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateTimeWithTimezone(t *testing.T) {
	expectedValue := int64(1617238883)
	actualTime, err := ParseFromLongDateTimeWithTimezone("2021-04-01 06:01:23+05:00")
	assert.Equal(t, nil, err)

	actualValue := actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateTimeWithTimezone2(t *testing.T) {
	expectedValue := int64(1617238883)
	actualTime, err := ParseFromLongDateTimeWithTimezone2("2021-04-01 06:01:23 +0500")
	assert.Equal(t, nil, err)

	actualValue := actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateTimeWithTimezoneRFC3339Format(t *testing.T) {
	expectedValue := int64(1617238883)
	actualTime, err := ParseFromLongDateTimeWithTimezoneRFC3339Format("2021-04-01T06:01:23+05:00")
	assert.Equal(t, nil, err)

	actualValue := actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateTimeWithoutSecondInFixedUtcOffset(t *testing.T) {
	expectedValue := int64(1691947440)
	actualTime, err := ParseFromLongDateTimeWithoutSecondInFixedUtcOffset("2023-08-13 17:24", 0)
	assert.Equal(t, nil, err)

	actualValue := actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromShortDateTimeInFixedUtcOffset(t *testing.T) {
	expectedValue := int64(1617228083)
	actualTime, err := ParseFromShortDateTimeInFixedUtcOffset("2021-4-1 6:1:23", 480)
	assert.Equal(t, nil, err)

	actualValue := actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromElapsedSeconds(t *testing.T) {
	expectedValue := "00:00:00"
	actualValue, err := ParseFromElapsedSeconds(0)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "00:00:09"
	actualValue, err = ParseFromElapsedSeconds(9)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "00:01:08"
	actualValue, err = ParseFromElapsedSeconds(68)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "01:00:07"
	actualValue, err = ParseFromElapsedSeconds(3607)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "23:59:59"
	actualValue, err = ParseFromElapsedSeconds(86399)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromElapsedSeconds_InvalidTime(t *testing.T) {
	_, err := ParseFromElapsedSeconds(-1)
	assert.NotEqual(t, nil, err)

	_, err = ParseFromElapsedSeconds(86400)
	assert.NotEqual(t, nil, err)
}

func TestIsUnixTimeEqualsYearAndMonth(t *testing.T) {
	actualValue := IsUnixTimeEqualsYearAndMonth(1691947440, time.UTC, 2023, 8)
	assert.Equal(t, true, actualValue)

	actualValue = IsUnixTimeEqualsYearAndMonth(1690847999, time.UTC, 2023, 8)
	assert.Equal(t, false, actualValue)
}

func TestGetTimezoneOffsetMinutes_FixedTimezone(t *testing.T) {
	timezone := time.FixedZone("Test Timezone", 120*60)
	expectedValue := int16(120)
	actualValue := GetTimezoneOffsetMinutes(time.Now().Unix(), timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", 345*60)
	expectedValue = int16(345)
	actualValue = GetTimezoneOffsetMinutes(time.Now().Unix(), timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", -720*60)
	expectedValue = int16(-720)
	actualValue = GetTimezoneOffsetMinutes(time.Now().Unix(), timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", 0)
	expectedValue = int16(0)
	actualValue = GetTimezoneOffsetMinutes(time.Now().Unix(), timezone)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetTimezoneOffsetMinutes_TimezoneWithDST(t *testing.T) {
	londonLocation, err := time.LoadLocation("Europe/London")
	assert.Equal(t, nil, err)

	// during standard time (UTC+0)
	expectedValue := int16(0)
	actualValue := GetTimezoneOffsetMinutes(1577858483, londonLocation) // 2020-01-01 06:01:23 +00:00
	assert.Equal(t, expectedValue, actualValue)

	// during daylight saving time (UTC+1)
	expectedValue = int16(60)
	actualValue = GetTimezoneOffsetMinutes(1619845283, londonLocation) // 2021-05-01 06:01:23 +01:00
	assert.Equal(t, expectedValue, actualValue)
}

func TestFormatTimezoneOffset_FixedTimezone(t *testing.T) {
	timezone := time.FixedZone("Test Timezone", 120*60)
	expectedValue := "+02:00"
	actualValue := FormatTimezoneOffset(time.Now().Unix(), timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", 345*60)
	expectedValue = "+05:45"
	actualValue = FormatTimezoneOffset(time.Now().Unix(), timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", -720*60)
	expectedValue = "-12:00"
	actualValue = FormatTimezoneOffset(time.Now().Unix(), timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", -150*60)
	expectedValue = "-02:30"
	actualValue = FormatTimezoneOffset(time.Now().Unix(), timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", 0)
	expectedValue = "+00:00"
	actualValue = FormatTimezoneOffset(time.Now().Unix(), timezone)
	assert.Equal(t, expectedValue, actualValue)
}

func TestFormatTimezoneOffset_TimezoneWithDST(t *testing.T) {
	londonLocation, err := time.LoadLocation("Europe/London")
	assert.Equal(t, nil, err)

	// during standard time (UTC+0)
	expectedValue := "+00:00"
	actualValue := FormatTimezoneOffset(1577858483, londonLocation) // 2020-01-01 06:01:23 +00:00
	assert.Equal(t, expectedValue, actualValue)

	// during daylight saving time (UTC+1)
	expectedValue = "+01:00"
	actualValue = FormatTimezoneOffset(1619845283, londonLocation) // 2021-05-01 06:01:23 +01:00
	assert.Equal(t, expectedValue, actualValue)
}

func TestFormatTimezoneOffsetFromHoursOffset(t *testing.T) {
	expectedValue := "+02:00"
	actualValue, err := FormatTimezoneOffsetFromHoursOffset("2")
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "+05:45"
	actualValue, err = FormatTimezoneOffsetFromHoursOffset("+5.75")
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "-12:00"
	actualValue, err = FormatTimezoneOffsetFromHoursOffset("-12.00")
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "-02:30"
	actualValue, err = FormatTimezoneOffsetFromHoursOffset("-2.5")
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "+00:00"
	actualValue, err = FormatTimezoneOffsetFromHoursOffset("0")
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestFormatTimezoneOffsetFromHoursOffset_InvalidHoursOffset(t *testing.T) {
	_, err := FormatTimezoneOffsetFromHoursOffset("")
	assert.EqualError(t, err, errs.ErrFormatInvalid.Message)

	_, err = FormatTimezoneOffsetFromHoursOffset("+")
	assert.EqualError(t, err, errs.ErrFormatInvalid.Message)

	_, err = FormatTimezoneOffsetFromHoursOffset("-")
	assert.EqualError(t, err, errs.ErrFormatInvalid.Message)

	_, err = FormatTimezoneOffsetFromHoursOffset("a")
	assert.EqualError(t, err, errs.ErrFormatInvalid.Message)
}

func TestParseFromTimezoneOffset(t *testing.T) {
	expectedValue := time.FixedZone("Timezone", 120*60)
	actualValue, err := ParseFromTimezoneOffset("+02:00")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = time.FixedZone("Timezone", 345*60)
	actualValue, err = ParseFromTimezoneOffset("+05:45")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = time.FixedZone("Timezone", -720*60)
	actualValue, err = ParseFromTimezoneOffset("-12:00")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = time.FixedZone("Timezone", -150*60)
	actualValue, err = ParseFromTimezoneOffset("-02:30")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = time.FixedZone("Timezone", 0)
	actualValue, err = ParseFromTimezoneOffset("+00:00")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)

	actualValue, err = ParseFromTimezoneOffset("00:00")
	assert.NotEqual(t, nil, err)

	actualValue, err = ParseFromTimezoneOffset("0")
	assert.NotEqual(t, nil, err)

	actualValue, err = ParseFromTimezoneOffset("1000")
	assert.NotEqual(t, nil, err)
}

func TestGetMinTransactionTimeFromUnixTime(t *testing.T) {
	expectedValue := int64(1617228083000)
	actualValue := GetMinTransactionTimeFromUnixTime(1617228083)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetMaxTransactionTimeFromUnixTime(t *testing.T) {
	expectedValue := int64(1617228083999)
	actualValue := GetMaxTransactionTimeFromUnixTime(1617228083)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetUnixTimeFromTransactionTime(t *testing.T) {
	expectedValue := int64(1617228083)
	actualValue := GetUnixTimeFromTransactionTime(1617228083999)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetTransactionTimeRangeByYearMonth(t *testing.T) {
	expectedMinValue := int64(1704016800000)
	expectedMaxValue := int64(1706788799999)
	actualMinValue, actualMaxValue, err := GetTransactionTimeRangeByYearMonth(2024, 1)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedMinValue, actualMinValue)
	assert.Equal(t, expectedMaxValue, actualMaxValue)
}

func TestParseFromUnixTime(t *testing.T) {
	expectedValue := int64(1617228083)
	actualTime := parseFromUnixTime(expectedValue)
	actualValue := actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)
}
