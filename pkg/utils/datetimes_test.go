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

func TestFormatNumericYearMonthDayToLongDate(t *testing.T) {
	expectedValue := "2026-12-31"
	actualValue := FormatNumericYearMonthDayToLongDate(int32(20261231))
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-01-01"
	actualValue = FormatNumericYearMonthDayToLongDate(int32(20260101))
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

func TestGetDateFirstTimeInTimezone(t *testing.T) {
	date := time.Date(2026, 10, 4, 12, 34, 56, 789, time.UTC)

	expectedValue := "2026-10-04T00:00:00Z"
	actualTime := GetDateFirstTimeInTimezone(date, time.UTC)
	actualValue := actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, time.UTC, actualTime.Location())

	nepalLocation := time.FixedZone("Nepal", 345*60)
	expectedValue = "2026-10-04T00:00:00+05:45"
	actualTime = GetDateFirstTimeInTimezone(date, nepalLocation)
	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, nepalLocation, actualTime.Location())

	expectedValue = "2026-10-04T00:00:00-12:00"
	actualValue = GetDateFirstTimeInTimezone(date, time.FixedZone("UTC-12", -720*60)).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-10-04T00:00:00+14:00"
	actualValue = GetDateFirstTimeInTimezone(date, time.FixedZone("UTC+14", 840*60)).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetDateFirstTimeInTimezone_InputDateInDifferentTimezone(t *testing.T) {
	// Use the input's calendar date, even when its UTC date is different.
	date := time.Date(2026, 1, 1, 0, 30, 0, 0, time.FixedZone("UTC+14", 840*60))
	expectedValue := "2026-01-01T00:00:00Z"
	actualValue := GetDateFirstTimeInTimezone(date, time.UTC).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	date = time.Date(2026, 12, 31, 23, 30, 0, 0, time.FixedZone("UTC-12", -720*60))
	expectedValue = "2026-12-31T00:00:00Z"
	actualValue = GetDateFirstTimeInTimezone(date, time.UTC).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetDateFirstTimeInTimezone_TimezoneWithDST(t *testing.T) {
	berlinLocation, err := time.LoadLocation("Europe/Berlin")
	assert.Equal(t, nil, err)
	lordHoweLocation, err := time.LoadLocation("Australia/Lord_Howe")
	assert.Equal(t, nil, err)

	expectedValue := "2026-03-29T00:00:00+01:00"
	actualValue := GetDateFirstTimeInTimezone(time.Date(2026, 3, 29, 12, 0, 0, 0, time.UTC), berlinLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-03-30T00:00:00+02:00"
	actualValue = GetDateFirstTimeInTimezone(time.Date(2026, 3, 30, 12, 0, 0, 0, time.UTC), berlinLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-10-25T00:00:00+02:00"
	actualValue = GetDateFirstTimeInTimezone(time.Date(2026, 10, 25, 12, 0, 0, 0, time.UTC), berlinLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-10-26T00:00:00+01:00"
	actualValue = GetDateFirstTimeInTimezone(time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC), berlinLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-10-04T00:00:00+10:30"
	actualValue = GetDateFirstTimeInTimezone(time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC), lordHoweLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-10-05T00:00:00+11:00"
	actualValue = GetDateFirstTimeInTimezone(time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC), lordHoweLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetDateFirstTimeInTimezone_TimezoneWithMidnightDST(t *testing.T) {
	santiagoLocation, err := time.LoadLocation("America/Santiago")
	assert.Equal(t, nil, err)
	havanaLocation, err := time.LoadLocation("America/Havana")
	assert.Equal(t, nil, err)
	cairoLocation, err := time.LoadLocation("Africa/Cairo")
	assert.Equal(t, nil, err)

	expectedValue := "2026-09-05T00:00:00-04:00"
	actualValue := GetDateFirstTimeInTimezone(time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC), santiagoLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-09-06T01:00:00-03:00"
	actualValue = GetDateFirstTimeInTimezone(time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC), santiagoLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-09-07T00:00:00-03:00"
	actualValue = GetDateFirstTimeInTimezone(time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC), santiagoLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-03-08T01:00:00-04:00"
	actualValue = GetDateFirstTimeInTimezone(time.Date(2026, 3, 8, 12, 0, 0, 0, time.UTC), havanaLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-04-24T01:00:00+03:00"
	actualValue = GetDateFirstTimeInTimezone(time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC), cairoLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetDateFirstTimeInTimezone_AmbiguousMidnight(t *testing.T) {
	havanaLocation, err := time.LoadLocation("America/Havana")
	assert.Equal(t, nil, err)

	// Select the first occurrence of midnight when it occurs twice.
	expectedValue := "2026-11-01T00:00:00-04:00"
	actualValue := GetDateFirstTimeInTimezone(time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC), havanaLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-11-02T00:00:00-05:00"
	actualValue = GetDateFirstTimeInTimezone(time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC), havanaLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetDateFirstTimeInTimezone_SkippedDate(t *testing.T) {
	apiaLocation, err := time.LoadLocation("Pacific/Apia")
	assert.Equal(t, nil, err)

	expectedValue := "2011-12-29T00:00:00-10:00"
	actualValue := GetDateFirstTimeInTimezone(time.Date(2011, 12, 29, 12, 0, 0, 0, time.UTC), apiaLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	// A skipped date advances to the beginning of the next existing date.
	expectedValue = "2011-12-31T00:00:00+14:00"
	actualValue = GetDateFirstTimeInTimezone(time.Date(2011, 12, 30, 12, 0, 0, 0, time.UTC), apiaLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	actualValue = GetDateFirstTimeInTimezone(time.Date(2011, 12, 31, 12, 0, 0, 0, time.UTC), apiaLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetDateFirstTimeInTimezone_FutureLeapYear(t *testing.T) {
	berlinLocation, err := time.LoadLocation("Europe/Berlin")
	assert.Equal(t, nil, err)
	sydneyLocation, err := time.LoadLocation("Australia/Sydney")
	assert.Equal(t, nil, err)

	expectedValue := "2040-02-29T00:00:00+01:00"
	actualValue := GetDateFirstTimeInTimezone(time.Date(2040, 2, 29, 12, 0, 0, 0, time.UTC), berlinLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2040-12-31T00:00:00+01:00"
	actualValue = GetDateFirstTimeInTimezone(time.Date(2040, 12, 31, 12, 0, 0, 0, time.UTC), berlinLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2040-12-31T00:00:00+11:00"
	actualValue = GetDateFirstTimeInTimezone(time.Date(2040, 12, 31, 12, 0, 0, 0, time.UTC), sydneyLocation).Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateFirstTimeInTimezone_TimezoneWithMidnightDST(t *testing.T) {
	santiagoLocation, err := time.LoadLocation("America/Santiago")
	assert.Equal(t, nil, err)
	havanaLocation, err := time.LoadLocation("America/Havana")
	assert.Equal(t, nil, err)
	cairoLocation, err := time.LoadLocation("Africa/Cairo")
	assert.Equal(t, nil, err)

	// Missing midnight must advance to the first valid time on the requested date.
	expectedValue := "2026-09-06T01:00:00-03:00"
	actualTime, err := ParseFromLongDateFirstTimeInTimezone("2026-09-06", santiagoLocation)
	assert.Equal(t, nil, err)
	actualValue := actualTime.Format(time.RFC3339)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-03-08T01:00:00-04:00"
	actualTime, err = ParseFromLongDateFirstTimeInTimezone("2026-03-08", havanaLocation)
	assert.Equal(t, nil, err)
	actualValue = actualTime.Format(time.RFC3339)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-04-24T01:00:00+03:00"
	actualTime, err = ParseFromLongDateFirstTimeInTimezone("2026-04-24", cairoLocation)
	assert.Equal(t, nil, err)
	actualValue = actualTime.Format(time.RFC3339)
	assert.Equal(t, expectedValue, actualValue)

	// The following date starts at midnight again.
	expectedValue = "2026-09-07T00:00:00-03:00"
	actualTime, err = ParseFromLongDateFirstTimeInTimezone("2026-09-07", santiagoLocation)
	assert.Equal(t, nil, err)
	actualValue = actualTime.Format(time.RFC3339)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateFirstTimeInTimezone_AmbiguousMidnight(t *testing.T) {
	havanaLocation, err := time.LoadLocation("America/Havana")
	assert.Equal(t, nil, err)

	// Midnight occurs twice; the date begins at the first occurrence, still UTC-4.
	expectedValue := "2026-11-01T00:00:00-04:00"
	actualTime, err := ParseFromLongDateFirstTimeInTimezone("2026-11-01", havanaLocation)
	assert.Equal(t, nil, err)
	actualValue := actualTime.Format(time.RFC3339)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateFirstTimeInTimezone_InvalidDate(t *testing.T) {
	_, err := ParseFromLongDateFirstTimeInTimezone("2026-02-30", time.UTC)
	assert.NotEqual(t, nil, err)

	apiaLocation, err := time.LoadLocation("Pacific/Apia")
	assert.Equal(t, nil, err)

	// This entire local date was skipped when Samoa crossed the date line.
	_, err = ParseFromLongDateFirstTimeInTimezone("2011-12-30", apiaLocation)
	assert.Equal(t, errs.ErrParameterInvalid, err)
}

func TestParseFromLongDateLastTimeInTimezone(t *testing.T) {
	shanghaiLocation, err := time.LoadLocation("Asia/Shanghai")
	assert.Equal(t, nil, err)

	expectedValue := "2026-03-29T23:59:59.999999999+08:00"
	actualTime, err := ParseFromLongDateLastTimeInTimezone("2026-03-29", shanghaiLocation)
	assert.Equal(t, nil, err)

	actualValue := actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-12-31T23:59:59.999999999Z"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2026-12-31", time.UTC)
	assert.Equal(t, nil, err)

	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateLastTimeInTimezone_TimezoneWithDST(t *testing.T) {
	berlinLocation, err := time.LoadLocation("Europe/Berlin")
	assert.Equal(t, nil, err)
	sydneyLocation, err := time.LoadLocation("Australia/Sydney")
	assert.Equal(t, nil, err)
	lordHoweLocation, err := time.LoadLocation("Australia/Lord_Howe")
	assert.Equal(t, nil, err)

	// daylight saving time starts (23-hour day)
	expectedValue := "2026-03-29T23:59:59.999999999+02:00"
	actualTime, err := ParseFromLongDateLastTimeInTimezone("2026-03-29", berlinLocation)
	assert.Equal(t, nil, err)

	actualValue := actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	// daylight saving time ends (25-hour day)
	expectedValue = "2026-10-25T23:59:59.999999999+01:00"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2026-10-25", berlinLocation)
	assert.Equal(t, nil, err)

	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	// daylight saving time starts in the southern hemisphere
	expectedValue = "2026-10-04T23:59:59.999999999+11:00"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2026-10-04", sydneyLocation)
	assert.Equal(t, nil, err)

	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	// daylight saving time ends in the southern hemisphere
	expectedValue = "2026-04-05T23:59:59.999999999+10:00"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2026-04-05", sydneyLocation)
	assert.Equal(t, nil, err)

	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	// daylight saving time starts with a half-hour change
	expectedValue = "2026-10-04T23:59:59.999999999+11:00"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2026-10-04", lordHoweLocation)
	assert.Equal(t, nil, err)

	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateLastTimeInTimezone_InvalidDate(t *testing.T) {
	_, err := ParseFromLongDateLastTimeInTimezone("2026-02-30", time.UTC)
	assert.NotEqual(t, nil, err)
}

func TestParseFromLongDateLastTimeInTimezone_TimezoneWithMidnightDST(t *testing.T) {
	santiagoLocation, err := time.LoadLocation("America/Santiago")
	assert.Equal(t, nil, err)
	havanaLocation, err := time.LoadLocation("America/Havana")
	assert.Equal(t, nil, err)
	cairoLocation, err := time.LoadLocation("Africa/Cairo")
	assert.Equal(t, nil, err)

	// The day before missing midnight must end at its own last second.
	expectedValue := "2026-09-05T23:59:59.999999999-04:00"
	actualTime, err := ParseFromLongDateLastTimeInTimezone("2026-09-05", santiagoLocation)
	assert.Equal(t, nil, err)
	actualValue := actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-09-06T23:59:59.999999999-03:00"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2026-09-06", santiagoLocation)
	assert.Equal(t, nil, err)
	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-09-07T23:59:59.999999999-03:00"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2026-09-07", santiagoLocation)
	assert.Equal(t, nil, err)
	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-03-07T23:59:59.999999999-05:00"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2026-03-07", havanaLocation)
	assert.Equal(t, nil, err)
	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-03-08T23:59:59.999999999-04:00"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2026-03-08", havanaLocation)
	assert.Equal(t, nil, err)
	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-04-23T23:59:59.999999999+02:00"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2026-04-23", cairoLocation)
	assert.Equal(t, nil, err)
	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	// Starting at 01:00 must not move the end time into April 25.
	expectedValue = "2026-04-24T23:59:59.999999999+03:00"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2026-04-24", cairoLocation)
	assert.Equal(t, nil, err)
	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateLastTimeInTimezone_AmbiguousMidnight(t *testing.T) {
	havanaLocation, err := time.LoadLocation("America/Havana")
	assert.Equal(t, nil, err)
	cairoLocation, err := time.LoadLocation("Africa/Cairo")
	assert.Equal(t, nil, err)

	// The previous date ends before the first occurrence of repeated midnight.
	expectedValue := "2026-10-31T23:59:59.999999999-04:00"
	actualTime, err := ParseFromLongDateLastTimeInTimezone("2026-10-31", havanaLocation)
	assert.Equal(t, nil, err)
	actualValue := actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2026-11-01T23:59:59.999999999-05:00"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2026-11-01", havanaLocation)
	assert.Equal(t, nil, err)
	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	// Cairo repeats the previous day's 23:00 hour instead of midnight.
	expectedValue = "2026-10-29T23:59:59.999999999+02:00"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2026-10-29", cairoLocation)
	assert.Equal(t, nil, err)
	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateLastTimeInTimezone_SkippedDate(t *testing.T) {
	apiaLocation, err := time.LoadLocation("Pacific/Apia")
	assert.Equal(t, nil, err)

	// The next valid date after December 29 is December 31.
	expectedValue := "2011-12-29T23:59:59.999999999-10:00"
	actualTime, err := ParseFromLongDateLastTimeInTimezone("2011-12-29", apiaLocation)
	assert.Equal(t, nil, err)
	actualValue := actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	_, err = ParseFromLongDateLastTimeInTimezone("2011-12-30", apiaLocation)
	assert.Equal(t, errs.ErrParameterInvalid, err)

	expectedValue = "2011-12-31T23:59:59.999999999+14:00"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2011-12-31", apiaLocation)
	assert.Equal(t, nil, err)
	actualValue = actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromLongDateLastTimeInTimezone_FutureLeapYear(t *testing.T) {
	berlinLocation, err := time.LoadLocation("Europe/Berlin")
	assert.Equal(t, nil, err)
	sydneyLocation, err := time.LoadLocation("Australia/Sydney")
	assert.Equal(t, nil, err)

	// Future timezone bounds must advance correctly past a leap year's final day.
	expectedValue := "2040-12-31T23:59:59.999999999+01:00"
	actualTime, err := ParseFromLongDateLastTimeInTimezone("2040-12-31", berlinLocation)
	assert.Equal(t, nil, err)
	actualValue := actualTime.Format(time.RFC3339Nano)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "2040-12-31T23:59:59.999999999+11:00"
	actualTime, err = ParseFromLongDateLastTimeInTimezone("2040-12-31", sydneyLocation)
	assert.Equal(t, nil, err)
	actualValue = actualTime.Format(time.RFC3339Nano)
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

func TestGetMaxDayOfMonth(t *testing.T) {
	expectedValue := 31
	actualValue := GetMaxDayOfMonth(2023, 1)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = 28
	actualValue = GetMaxDayOfMonth(2023, 2)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = 29
	actualValue = GetMaxDayOfMonth(2024, 2)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = 30
	actualValue = GetMaxDayOfMonth(2023, 4)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = 31
	actualValue = GetMaxDayOfMonth(2023, 12)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = 28
	actualValue = GetMaxDayOfMonth(2100, 2)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetDaysBetweenDates(t *testing.T) {
	shanghaiLocation, err := time.LoadLocation("Asia/Shanghai")
	assert.Equal(t, nil, err)

	startDate := time.Date(2026, 3, 28, 0, 0, 0, 0, shanghaiLocation)
	endDate := time.Date(2026, 3, 30, 0, 0, 0, 0, shanghaiLocation)
	expectedValue := 2
	actualValue := GetDaysBetweenDates(startDate, endDate)
	assert.Equal(t, expectedValue, actualValue)

	// across a leap day
	startDate = time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC)
	endDate = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	expectedValue = 2
	actualValue = GetDaysBetweenDates(startDate, endDate)
	assert.Equal(t, expectedValue, actualValue)

	// across a year boundary
	startDate = time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	endDate = time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedValue = 1
	actualValue = GetDaysBetweenDates(startDate, endDate)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetDaysBetweenDates_FixedTimezone(t *testing.T) {
	timezone := time.FixedZone("Test Timezone", 120*60)
	startDate := time.Date(2026, 3, 28, 0, 0, 0, 0, timezone)
	endDate := time.Date(2026, 3, 30, 0, 0, 0, 0, timezone)
	expectedValue := 2
	actualValue := GetDaysBetweenDates(startDate, endDate)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = -2
	actualValue = GetDaysBetweenDates(endDate, startDate)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetDaysBetweenDates_TimezoneWithDST(t *testing.T) {
	berlinLocation, err := time.LoadLocation("Europe/Berlin")
	assert.Equal(t, nil, err)
	sydneyLocation, err := time.LoadLocation("Australia/Sydney")
	assert.Equal(t, nil, err)
	lordHoweLocation, err := time.LoadLocation("Australia/Lord_Howe")
	assert.Equal(t, nil, err)

	// daylight saving time starts (47 elapsed hours)
	startDate := time.Date(2026, 3, 28, 0, 0, 0, 0, berlinLocation)
	endDate := time.Date(2026, 3, 30, 0, 0, 0, 0, berlinLocation)
	expectedValue := 2
	actualValue := GetDaysBetweenDates(startDate, endDate)
	assert.Equal(t, expectedValue, actualValue)

	// daylight saving time ends (49 elapsed hours)
	startDate = time.Date(2026, 10, 24, 0, 0, 0, 0, berlinLocation)
	endDate = time.Date(2026, 10, 26, 0, 0, 0, 0, berlinLocation)
	expectedValue = 2
	actualValue = GetDaysBetweenDates(startDate, endDate)
	assert.Equal(t, expectedValue, actualValue)

	// across standard time and daylight saving time
	startDate = time.Date(2026, 1, 1, 0, 0, 0, 0, berlinLocation)
	endDate = time.Date(2026, 7, 1, 0, 0, 0, 0, berlinLocation)
	expectedValue = 181
	actualValue = GetDaysBetweenDates(startDate, endDate)
	assert.Equal(t, expectedValue, actualValue)

	// daylight saving time starts in the southern hemisphere
	startDate = time.Date(2026, 10, 3, 0, 0, 0, 0, sydneyLocation)
	endDate = time.Date(2026, 10, 5, 0, 0, 0, 0, sydneyLocation)
	expectedValue = 2
	actualValue = GetDaysBetweenDates(startDate, endDate)
	assert.Equal(t, expectedValue, actualValue)

	// daylight saving time ends in the southern hemisphere
	startDate = time.Date(2026, 4, 4, 0, 0, 0, 0, sydneyLocation)
	endDate = time.Date(2026, 4, 6, 0, 0, 0, 0, sydneyLocation)
	expectedValue = 2
	actualValue = GetDaysBetweenDates(startDate, endDate)
	assert.Equal(t, expectedValue, actualValue)

	// daylight saving time starts with a half-hour change
	startDate = time.Date(2026, 10, 3, 0, 0, 0, 0, lordHoweLocation)
	endDate = time.Date(2026, 10, 5, 0, 0, 0, 0, lordHoweLocation)
	expectedValue = 2
	actualValue = GetDaysBetweenDates(startDate, endDate)
	assert.Equal(t, expectedValue, actualValue)

	// the end date is before the start date across daylight saving time
	startDate = time.Date(2026, 3, 30, 0, 0, 0, 0, berlinLocation)
	endDate = time.Date(2026, 3, 29, 0, 0, 0, 0, berlinLocation)
	expectedValue = -1
	actualValue = GetDaysBetweenDates(startDate, endDate)
	assert.Equal(t, expectedValue, actualValue)

	// the start date and end date are the same
	startDate = time.Date(2026, 3, 29, 0, 0, 0, 0, berlinLocation)
	endDate = time.Date(2026, 3, 29, 0, 0, 0, 0, berlinLocation)
	expectedValue = 0
	actualValue = GetDaysBetweenDates(startDate, endDate)
	assert.Equal(t, expectedValue, actualValue)
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

func TestGetMinimumTimezoneOffsetMinutes_TimezoneWithNegativeDST(t *testing.T) {
	dublinLocation, err := time.LoadLocation("Europe/Dublin")
	assert.Equal(t, nil, err)
	casablancaLocation, err := time.LoadLocation("Africa/Casablanca")
	assert.Equal(t, nil, err)

	referenceUnixTime := time.Now().Unix()
	expectedValue := int16(0)
	actualValue := GetMinimumTimezoneOffsetMinutes(referenceUnixTime, dublinLocation)
	assert.Equal(t, expectedValue, actualValue)

	// the lower offset occurs during Ramadan, outside both January and July
	referenceUnixTime = time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue = int16(0)
	actualValue = GetMinimumTimezoneOffsetMinutes(referenceUnixTime, casablancaLocation)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetMinimumTimezoneOffsetMinutes_YearBoundary(t *testing.T) {
	saoTomeLocation, err := time.LoadLocation("Africa/Sao_Tome")
	assert.Equal(t, nil, err)
	moscowLocation, err := time.LoadLocation("Europe/Moscow")
	assert.Equal(t, nil, err)

	// UTC+0 applies for only the first hour of 2018, before the change to UTC+1
	referenceUnixTime := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue := int16(0)
	actualValue := GetMinimumTimezoneOffsetMinutes(referenceUnixTime, saoTomeLocation)
	assert.Equal(t, expectedValue, actualValue)

	// the lower offset introduced in 2014 must not affect the minimum for 2013
	referenceUnixTime = time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue = int16(240)
	actualValue = GetMinimumTimezoneOffsetMinutes(referenceUnixTime, moscowLocation)
	assert.Equal(t, expectedValue, actualValue)

	referenceUnixTime = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue = int16(180)
	actualValue = GetMinimumTimezoneOffsetMinutes(referenceUnixTime, moscowLocation)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetMinimumTimezoneOffsetMinutes_FutureLeapYear(t *testing.T) {
	berlinLocation, err := time.LoadLocation("Europe/Berlin")
	assert.Equal(t, nil, err)
	sydneyLocation, err := time.LoadLocation("Australia/Sydney")
	assert.Equal(t, nil, err)
	referenceUnixTime := time.Date(2040, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

	// future rules must advance past the last day of a leap year
	expectedValue := int16(60)
	actualValue := GetMinimumTimezoneOffsetMinutes(referenceUnixTime, berlinLocation)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = int16(600)
	actualValue = GetMinimumTimezoneOffsetMinutes(referenceUnixTime, sydneyLocation)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetScheduledTransactionScheduledAt_FixedTimezone(t *testing.T) {
	referenceUnixTime := time.Now().Unix()

	timezone := time.FixedZone("Test Timezone", 0)
	expectedValue := int16(0)
	actualValue := GetScheduledTransactionScheduledAt(referenceUnixTime, timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", 120*60)
	expectedValue = int16(1320)
	actualValue = GetScheduledTransactionScheduledAt(referenceUnixTime, timezone)
	assert.Equal(t, expectedValue, actualValue)

	// legacy templates with no creation time retain the fixed-offset schedule
	actualValue = GetScheduledTransactionScheduledAt(0, timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", 345*60)
	expectedValue = int16(1095)
	actualValue = GetScheduledTransactionScheduledAt(referenceUnixTime, timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", -300*60)
	expectedValue = int16(300)
	actualValue = GetScheduledTransactionScheduledAt(referenceUnixTime, timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", -210*60)
	expectedValue = int16(210)
	actualValue = GetScheduledTransactionScheduledAt(referenceUnixTime, timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", -720*60)
	expectedValue = int16(720)
	actualValue = GetScheduledTransactionScheduledAt(referenceUnixTime, timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", 840*60)
	expectedValue = int16(600)
	actualValue = GetScheduledTransactionScheduledAt(referenceUnixTime, timezone)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetScheduledTransactionScheduledAt_TimezoneWithDST(t *testing.T) {
	referenceUnixTime := time.Now().Unix()

	berlinLocation, err := time.LoadLocation("Europe/Berlin")
	assert.Equal(t, nil, err)
	londonLocation, err := time.LoadLocation("Europe/London")
	assert.Equal(t, nil, err)
	newYorkLocation, err := time.LoadLocation("America/New_York")
	assert.Equal(t, nil, err)
	sydneyLocation, err := time.LoadLocation("Australia/Sydney")
	assert.Equal(t, nil, err)
	lordHoweLocation, err := time.LoadLocation("Australia/Lord_Howe")
	assert.Equal(t, nil, err)

	// the minimum yearly offset uses standard time in the northern hemisphere
	expectedValue := int16(1380)
	actualValue := GetScheduledTransactionScheduledAt(referenceUnixTime, berlinLocation)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = int16(0)
	actualValue = GetScheduledTransactionScheduledAt(referenceUnixTime, londonLocation)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = int16(300)
	actualValue = GetScheduledTransactionScheduledAt(referenceUnixTime, newYorkLocation)
	assert.Equal(t, expectedValue, actualValue)

	// the minimum yearly offset also uses standard time in the southern hemisphere
	expectedValue = int16(840)
	actualValue = GetScheduledTransactionScheduledAt(referenceUnixTime, sydneyLocation)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = int16(810)
	actualValue = GetScheduledTransactionScheduledAt(referenceUnixTime, lordHoweLocation)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetScheduledTransactionScheduledAt_DifferentYears(t *testing.T) {
	moscowLocation, err := time.LoadLocation("Europe/Moscow")
	assert.Equal(t, nil, err)

	// standard time was UTC+3 in 2010
	createdUnixTime := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue := int16(1260)
	actualValue := GetScheduledTransactionScheduledAt(createdUnixTime, moscowLocation)
	assert.Equal(t, expectedValue, actualValue)

	// UTC+4 was used throughout 2012
	createdUnixTime = time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue = int16(1200)
	actualValue = GetScheduledTransactionScheduledAt(createdUnixTime, moscowLocation)
	assert.Equal(t, expectedValue, actualValue)

	// the minimum changed to UTC+3 in October 2014
	createdUnixTime = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue = int16(1260)
	actualValue = GetScheduledTransactionScheduledAt(createdUnixTime, moscowLocation)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetScheduledTransactionScheduledAt_CreationYearBoundary(t *testing.T) {
	moscowLocation, err := time.LoadLocation("Europe/Moscow")
	assert.Equal(t, nil, err)

	// use the UTC creation year even when the local date is already in the next year
	createdUnixTime := time.Date(2013, 12, 31, 23, 59, 59, 0, time.UTC).Unix()
	expectedValue := int16(1200)
	actualValue := GetScheduledTransactionScheduledAt(createdUnixTime, moscowLocation)
	assert.Equal(t, expectedValue, actualValue)

	createdUnixTime = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue = int16(1260)
	actualValue = GetScheduledTransactionScheduledAt(createdUnixTime, moscowLocation)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetScheduledTransactionFinalTime_FixedTimezone(t *testing.T) {
	referenceUnixTime := time.Now().Unix()

	todayFirstUnixTimeInUTC := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC).Unix()

	timezone := time.FixedZone("Test Timezone", 0)
	expectedValue := time.Date(2026, 11, 30, 0, 0, 0, 0, timezone)
	actualValue := GetScheduledTransactionFinalTime(todayFirstUnixTimeInUTC, 0, referenceUnixTime, timezone)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	// legacy templates retain their fixed offset and advance to the next local date
	timezone = time.FixedZone("Test Timezone", 120*60)
	expectedValue = time.Date(2026, 12, 1, 0, 0, 0, 0, timezone)
	actualValue = GetScheduledTransactionFinalTime(todayFirstUnixTimeInUTC, 1320, referenceUnixTime, timezone)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	timezone = time.FixedZone("Test Timezone", 345*60)
	expectedValue = time.Date(2026, 12, 1, 0, 0, 0, 0, timezone)
	actualValue = GetScheduledTransactionFinalTime(todayFirstUnixTimeInUTC, 1095, referenceUnixTime, timezone)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	// negative offsets retain the same local date
	timezone = time.FixedZone("Test Timezone", -300*60)
	expectedValue = time.Date(2026, 11, 30, 0, 0, 0, 0, timezone)
	actualValue = GetScheduledTransactionFinalTime(todayFirstUnixTimeInUTC, 300, referenceUnixTime, timezone)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	timezone = time.FixedZone("Test Timezone", -210*60)
	expectedValue = time.Date(2026, 11, 30, 0, 0, 0, 0, timezone)
	actualValue = GetScheduledTransactionFinalTime(todayFirstUnixTimeInUTC, 210, referenceUnixTime, timezone)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	timezone = time.FixedZone("Test Timezone", -720*60)
	expectedValue = time.Date(2026, 11, 30, 0, 0, 0, 0, timezone)
	actualValue = GetScheduledTransactionFinalTime(todayFirstUnixTimeInUTC, 720, referenceUnixTime, timezone)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	timezone = time.FixedZone("Test Timezone", 840*60)
	expectedValue = time.Date(2026, 12, 1, 0, 0, 0, 0, timezone)
	actualValue = GetScheduledTransactionFinalTime(todayFirstUnixTimeInUTC, 600, referenceUnixTime, timezone)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}
}

func TestGetScheduledTransactionFinalTime_DateBoundary(t *testing.T) {
	referenceUnixTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

	shanghaiLocation, err := time.LoadLocation("Asia/Shanghai")
	assert.Equal(t, nil, err)

	// across a year boundary
	expectedValue := time.Date(2027, 1, 1, 0, 0, 0, 0, shanghaiLocation)
	actualValue := GetScheduledTransactionFinalTime(time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC).Unix(), 960, referenceUnixTime, shanghaiLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	// across a leap day
	expectedValue = time.Date(2024, 2, 29, 0, 0, 0, 0, shanghaiLocation)
	actualValue = GetScheduledTransactionFinalTime(time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC).Unix(), 960, referenceUnixTime, shanghaiLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	expectedValue = time.Date(2024, 3, 1, 0, 0, 0, 0, shanghaiLocation)
	actualValue = GetScheduledTransactionFinalTime(time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC).Unix(), 960, referenceUnixTime, shanghaiLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}
}

func TestGetScheduledTransactionFinalTime_TimezoneWithDST(t *testing.T) {
	referenceUnixTime := time.Now().Unix()

	berlinLocation, err := time.LoadLocation("Europe/Berlin")
	assert.Equal(t, nil, err)
	londonLocation, err := time.LoadLocation("Europe/London")
	assert.Equal(t, nil, err)
	newYorkLocation, err := time.LoadLocation("America/New_York")
	assert.Equal(t, nil, err)

	// during standard time (UTC+1)
	expectedValue := time.Date(2026, 12, 1, 0, 0, 0, 0, berlinLocation)
	actualValue := GetScheduledTransactionFinalTime(time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC).Unix(), 1380, referenceUnixTime, berlinLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	// during daylight saving time (UTC+2), the scheduled UTC minute stays unchanged
	expectedValue = time.Date(2026, 7, 1, 0, 0, 0, 0, berlinLocation)
	actualValue = GetScheduledTransactionFinalTime(time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC).Unix(), 1380, referenceUnixTime, berlinLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	// the reference offset is zero, but the transaction date uses daylight saving time
	expectedValue = time.Date(2026, 7, 1, 0, 0, 0, 0, londonLocation)
	actualValue = GetScheduledTransactionFinalTime(time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC).Unix(), 0, referenceUnixTime, londonLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	// during daylight saving time with a negative offset (UTC-4)
	expectedValue = time.Date(2026, 7, 1, 0, 0, 0, 0, newYorkLocation)
	actualValue = GetScheduledTransactionFinalTime(time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC).Unix(), 300, referenceUnixTime, newYorkLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}
}

func TestGetScheduledTransactionFinalTime_DSTTransition(t *testing.T) {
	referenceUnixTime := time.Now().Unix()

	berlinLocation, err := time.LoadLocation("Europe/Berlin")
	assert.Equal(t, nil, err)

	// midnight before daylight saving time starts (UTC+1)
	expectedValue := time.Date(2026, 3, 29, 0, 0, 0, 0, berlinLocation)
	actualValue := GetScheduledTransactionFinalTime(time.Date(2026, 3, 28, 0, 0, 0, 0, time.UTC).Unix(), 1380, referenceUnixTime, berlinLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	// midnight after daylight saving time starts (UTC+2)
	expectedValue = time.Date(2026, 3, 30, 0, 0, 0, 0, berlinLocation)
	actualValue = GetScheduledTransactionFinalTime(time.Date(2026, 3, 29, 0, 0, 0, 0, time.UTC).Unix(), 1380, referenceUnixTime, berlinLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	// midnight before daylight saving time ends (UTC+2)
	expectedValue = time.Date(2026, 10, 25, 0, 0, 0, 0, berlinLocation)
	actualValue = GetScheduledTransactionFinalTime(time.Date(2026, 10, 24, 0, 0, 0, 0, time.UTC).Unix(), 1380, referenceUnixTime, berlinLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	// midnight after daylight saving time ends (UTC+1)
	expectedValue = time.Date(2026, 10, 26, 0, 0, 0, 0, berlinLocation)
	actualValue = GetScheduledTransactionFinalTime(time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC).Unix(), 1380, referenceUnixTime, berlinLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}
}

func TestGetScheduledTransactionFinalTime_SouthernTimezoneWithDST(t *testing.T) {
	referenceUnixTime := time.Now().Unix()

	sydneyLocation, err := time.LoadLocation("Australia/Sydney")
	assert.Equal(t, nil, err)
	lordHoweLocation, err := time.LoadLocation("Australia/Lord_Howe")
	assert.Equal(t, nil, err)

	// during standard time (UTC+10), using the minimum offset in the creation year
	expectedValue := time.Date(2026, 7, 1, 0, 0, 0, 0, sydneyLocation)
	actualValue := GetScheduledTransactionFinalTime(time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC).Unix(), 840, referenceUnixTime, sydneyLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	// during daylight saving time (UTC+11)
	expectedValue = time.Date(2027, 1, 1, 0, 0, 0, 0, sydneyLocation)
	actualValue = GetScheduledTransactionFinalTime(time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC).Unix(), 840, referenceUnixTime, sydneyLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	// during standard time with a half-hour daylight saving time difference (UTC+10:30)
	expectedValue = time.Date(2026, 7, 1, 0, 0, 0, 0, lordHoweLocation)
	actualValue = GetScheduledTransactionFinalTime(time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC).Unix(), 810, referenceUnixTime, lordHoweLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	// during daylight saving time (UTC+11)
	expectedValue = time.Date(2027, 1, 1, 0, 0, 0, 0, lordHoweLocation)
	actualValue = GetScheduledTransactionFinalTime(time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC).Unix(), 810, referenceUnixTime, lordHoweLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}
}

func TestGetScheduledTransactionFinalTime_CreatedUnixTime(t *testing.T) {
	moscowLocation, err := time.LoadLocation("Europe/Moscow")
	assert.Equal(t, nil, err)
	todayFirstUnixTimeInUTC := time.Date(2014, 11, 30, 0, 0, 0, 0, time.UTC).Unix()

	// a template created in 2012 keeps its UTC+4 reference offset when executed in 2014
	createdUnixTime := time.Date(2012, 6, 1, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue := time.Date(2014, 12, 1, 0, 0, 0, 0, moscowLocation)
	actualValue := GetScheduledTransactionFinalTime(todayFirstUnixTimeInUTC, 1200, createdUnixTime, moscowLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}

	// a template created in 2014 uses UTC+3 for the same transaction date
	createdUnixTime = time.Date(2014, 6, 1, 0, 0, 0, 0, time.UTC).Unix()
	actualValue = GetScheduledTransactionFinalTime(todayFirstUnixTimeInUTC, 1260, createdUnixTime, moscowLocation)
	assert.NotNil(t, actualValue)
	if actualValue != nil {
		assert.Equal(t, expectedValue, *actualValue)
	}
}

func TestGetScheduledTransactionFinalTime_TimezoneWithMidnightDST(t *testing.T) {
	santiagoLocation, err := time.LoadLocation("America/Santiago")
	assert.Equal(t, nil, err)
	havanaLocation, err := time.LoadLocation("America/Havana")
	assert.Equal(t, nil, err)
	cairoLocation, err := time.LoadLocation("Africa/Cairo")
	assert.Equal(t, nil, err)
	createdUnixTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

	scheduledAt := GetScheduledTransactionScheduledAt(createdUnixTime, santiagoLocation)
	expectedValue := "2026-09-06T01:00:00-03:00"
	actualTime := GetScheduledTransactionFinalTime(time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC).Unix(), scheduledAt, createdUnixTime, santiagoLocation)
	assert.NotNil(t, actualTime)
	if actualTime != nil {
		actualValue := actualTime.Format(time.RFC3339)
		assert.Equal(t, expectedValue, actualValue)
	}

	scheduledAt = GetScheduledTransactionScheduledAt(createdUnixTime, havanaLocation)
	expectedValue = "2026-03-08T01:00:00-04:00"
	actualTime = GetScheduledTransactionFinalTime(time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC).Unix(), scheduledAt, createdUnixTime, havanaLocation)
	assert.NotNil(t, actualTime)
	if actualTime != nil {
		actualValue := actualTime.Format(time.RFC3339)
		assert.Equal(t, expectedValue, actualValue)
	}

	// Repeated midnight uses its first occurrence.
	expectedValue = "2026-11-01T00:00:00-04:00"
	actualTime = GetScheduledTransactionFinalTime(time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC).Unix(), scheduledAt, createdUnixTime, havanaLocation)
	assert.NotNil(t, actualTime)
	if actualTime != nil {
		actualValue := actualTime.Format(time.RFC3339)
		assert.Equal(t, expectedValue, actualValue)
	}

	scheduledAt = GetScheduledTransactionScheduledAt(createdUnixTime, cairoLocation)
	expectedValue = "2026-04-24T01:00:00+03:00"
	actualTime = GetScheduledTransactionFinalTime(time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC).Unix(), scheduledAt, createdUnixTime, cairoLocation)
	assert.NotNil(t, actualTime)
	if actualTime != nil {
		actualValue := actualTime.Format(time.RFC3339)
		assert.Equal(t, expectedValue, actualValue)
	}
}

func TestGetScheduledTransactionFinalTime_SkippedDate(t *testing.T) {
	apiaLocation, err := time.LoadLocation("Pacific/Apia")
	assert.Equal(t, nil, err)
	createdUnixTime := time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	scheduledAt := GetScheduledTransactionScheduledAt(createdUnixTime, apiaLocation)

	actualTime := GetScheduledTransactionFinalTime(time.Date(2011, 12, 30, 0, 0, 0, 0, time.UTC).Unix(), scheduledAt, createdUnixTime, apiaLocation)
	assert.Nil(t, actualTime)

	expectedValue := "2011-12-31T00:00:00+14:00"
	actualTime = GetScheduledTransactionFinalTime(time.Date(2011, 12, 31, 0, 0, 0, 0, time.UTC).Unix(), scheduledAt, createdUnixTime, apiaLocation)
	assert.NotNil(t, actualTime)
	if actualTime != nil {
		actualValue := actualTime.Format(time.RFC3339)
		assert.Equal(t, expectedValue, actualValue)
	}
}

func TestGetScheduledTransactionFinalTime_NonMidnight(t *testing.T) {
	timezone := time.FixedZone("Template Timezone", 120*60)
	createdUnixTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	scheduledAt := GetScheduledTransactionScheduledAt(createdUnixTime, timezone) + 15

	// Existing non-midnight schedules retain their wall-clock time.
	expectedValue := "2026-12-01T00:15:00+02:00"
	actualTime := GetScheduledTransactionFinalTime(time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC).Unix(), scheduledAt, createdUnixTime, timezone)
	assert.NotNil(t, actualTime)
	if actualTime != nil {
		actualValue := actualTime.Format(time.RFC3339)
		assert.Equal(t, expectedValue, actualValue)
	}
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

	timezone = time.FixedZone("Test Timezone", -30*60)
	expectedValue = "-00:30"
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

	expectedValue = "-00:30"
	actualValue, err = FormatTimezoneOffsetFromHoursOffset("-0.5")
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

func TestGetStartOfDay(t *testing.T) {
	expectedValue := int64(1617148800) // 2021-03-31 00:00:00 UTC
	actualValue := GetStartOfDay(time.Unix(1617228083, 0).In(time.UTC))
	assert.Equal(t, expectedValue, actualValue.Unix())

	expectedValue = int64(1617206400) // 2021-04-01 00:00:00 UTC+8
	actualValue = GetStartOfDay(time.Unix(1617228083, 0).In(time.FixedZone("Test Timezone", 28800)))
	assert.Equal(t, expectedValue, actualValue.Unix())
}

func TestParseFromUnixTime(t *testing.T) {
	expectedValue := int64(1617228083)
	actualTime := parseFromUnixTime(expectedValue)
	actualValue := actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)
}
