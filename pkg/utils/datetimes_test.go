package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseNumericYearMonth(t *testing.T) {
	expectedYear := int32(2024)
	expectedMonth := int32(3)
	actualYear, actualMonth, err := ParseNumericYearMonth("2024-03")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedYear, actualYear)
	assert.Equal(t, expectedMonth, actualMonth)
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

func TestParseFromLongDateTime(t *testing.T) {
	expectedValue := int64(1617228083)
	actualTime, err := ParseFromLongDateTime("2021-04-01 06:01:23", 480)
	assert.Equal(t, nil, err)

	actualValue := actualTime.Unix()
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

func TestParseFromLongDateTimeWithoutSecond(t *testing.T) {
	expectedValue := int64(1691947440)
	actualTime, err := ParseFromLongDateTimeWithoutSecond("2023-08-13 17:24", 0)
	assert.Equal(t, nil, err)

	actualValue := actualTime.Unix()
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromShortDateTime(t *testing.T) {
	expectedValue := int64(1617228083)
	actualTime, err := ParseFromShortDateTime("2021-4-1 6:1:23", 480)
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

func TestGetTimezoneOffsetMinutes(t *testing.T) {
	timezone := time.FixedZone("Test Timezone", 120*60)
	expectedValue := int16(120)
	actualValue := GetTimezoneOffsetMinutes(timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", 345*60)
	expectedValue = int16(345)
	actualValue = GetTimezoneOffsetMinutes(timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", -720*60)
	expectedValue = int16(-720)
	actualValue = GetTimezoneOffsetMinutes(timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", 0)
	expectedValue = int16(0)
	actualValue = GetTimezoneOffsetMinutes(timezone)
	assert.Equal(t, expectedValue, actualValue)
}

func TestFormatTimezoneOffset(t *testing.T) {
	timezone := time.FixedZone("Test Timezone", 120*60)
	expectedValue := "+02:00"
	actualValue := FormatTimezoneOffset(timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", 345*60)
	expectedValue = "+05:45"
	actualValue = FormatTimezoneOffset(timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", -720*60)
	expectedValue = "-12:00"
	actualValue = FormatTimezoneOffset(timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", -150*60)
	expectedValue = "-02:30"
	actualValue = FormatTimezoneOffset(timezone)
	assert.Equal(t, expectedValue, actualValue)

	timezone = time.FixedZone("Test Timezone", 0)
	expectedValue = "+00:00"
	actualValue = FormatTimezoneOffset(timezone)
	assert.Equal(t, expectedValue, actualValue)
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
