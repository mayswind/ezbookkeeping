package utils

import (
	"fmt"
	"time"
)

const (
	longDateTimeFormat              = "2006-01-02 15:04:05"
	longDateTimeWithoutSecondFormat = "2006-01-02 15:04"
	shortDateTimeFormat             = "2006-1-2 15:4:5"
	yearMonthDateTimeFormat         = "2006-01"
)

// FormatUnixTimeToLongDateTimeInServerTimezone returns a textual representation of the unix time formatted by long date time format
func FormatUnixTimeToLongDateTimeInServerTimezone(unixTime int64) string {
	return parseFromUnixTime(unixTime).Format(longDateTimeFormat)
}

// FormatUnixTimeToLongDateTimeWithoutSecond returns a textual representation of the unix time formatted by long date time format (no second)
func FormatUnixTimeToLongDateTimeWithoutSecond(unixTime int64, timezone *time.Location) string {
	t := parseFromUnixTime(unixTime)

	if timezone != nil {
		t = t.In(timezone)
	}

	return t.Format(longDateTimeWithoutSecondFormat)
}

// FormatUnixTimeToYearMonth returns year and month of specified unix time
func FormatUnixTimeToYearMonth(unixTime int64, timezone *time.Location) string {
	t := parseFromUnixTime(unixTime)

	if timezone != nil {
		t = t.In(timezone)
	}

	return t.Format(yearMonthDateTimeFormat)
}

// ParseFromLongDateTime parses a formatted string in long date time format
func ParseFromLongDateTime(t string, utcOffset int16) (time.Time, error) {
	timezone := time.FixedZone("Timezone", int(utcOffset)*60)
	return time.ParseInLocation(longDateTimeFormat, t, timezone)
}

// ParseFromShortDateTime parses a formatted string in short date time format
func ParseFromShortDateTime(t string, utcOffset int16) (time.Time, error) {
	timezone := time.FixedZone("Timezone", int(utcOffset)*60)
	return time.ParseInLocation(shortDateTimeFormat, t, timezone)
}

// FormatTimezoneOffset returns "+/-HH:MM" format of timezone
func FormatTimezoneOffset(timezone *time.Location) string {
	_, tzOffset := time.Now().In(timezone).Zone()
	tzMinutesOffset := tzOffset / 60

	sign := "+"
	hourAbsOffset := tzMinutesOffset / 60
	minuteAbsOffset := tzMinutesOffset % 60

	if hourAbsOffset < 0 {
		sign = "-"
		hourAbsOffset = -hourAbsOffset
		minuteAbsOffset = -minuteAbsOffset
	}

	return fmt.Sprintf("%s%02d:%02d", sign, hourAbsOffset, minuteAbsOffset)
}

// GetMinTransactionTimeFromUnixTime returns the minimum transaction time from unix time
func GetMinTransactionTimeFromUnixTime(unixTime int64) int64 {
	return unixTime * 1000
}

// GetMaxTransactionTimeFromUnixTime returns the maximum transaction time from unix time
func GetMaxTransactionTimeFromUnixTime(unixTime int64) int64 {
	return unixTime*1000 + 999
}

// GetUnixTimeFromTransactionTime returns unix time from the transaction time
func GetUnixTimeFromTransactionTime(transactionTime int64) int64 {
	return transactionTime / 1000
}

// parseFromUnixTime parses a unix time and returns a golang time struct
func parseFromUnixTime(unixTime int64) time.Time {
	return time.Unix(unixTime, 0)
}
