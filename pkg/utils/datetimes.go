package utils

import "time"

const (
	longDateTimeFormat              = "2006-01-02 15:04:05"
	longDateTimeWithoutSecondFormat = "2006-01-02 15:04"
)

// FormatToLongDateTime returns a textual representation of the time value formatted by long date time format
func FormatToLongDateTime(t time.Time) string {
	return t.Format(longDateTimeFormat)
}

// FormatToLongDateTimeWithoutSecond returns a textual representation of the time value formatted by long date time format (no second)
func FormatToLongDateTimeWithoutSecond(t time.Time) string {
	return t.Format(longDateTimeWithoutSecondFormat)
}

// ParseFromUnixTime parses a unix time and returns a golang time struct
func ParseFromUnixTime(unixTime int64) time.Time {
	return time.Unix(unixTime, 0)
}

// ParseFromLongDateTime parses a formatted string in long date time format
func ParseFromLongDateTime(t string, utcOffset int16) (time.Time, error) {
	timezone := time.FixedZone("Timezone", int(utcOffset)*60)
	return time.ParseInLocation(longDateTimeFormat, t, timezone)
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
