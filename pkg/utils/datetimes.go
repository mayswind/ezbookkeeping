package utils

import "time"

const (
	longDateTimeFormat = "2006-01-02 15:04:05"
)

// FormatToLongDateTime returns a textual representation of the time value formatted by long date time format
func FormatToLongDateTime(t time.Time) string {
	return t.Format(longDateTimeFormat)
}

// ParseFromLongDateTime parses a formatted string in long date time format
func ParseFromLongDateTime(t string) (time.Time, error) {
	return time.Parse(longDateTimeFormat, t)
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
