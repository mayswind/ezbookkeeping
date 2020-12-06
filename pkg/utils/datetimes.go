package utils

import "time"

const LongDateTimeFormat = "2006-01-02 15:04:05"

func FormatToLongDateTime(t time.Time) string {
	return t.Format(LongDateTimeFormat)
}

func ParseFromLongDateTime(t string) (time.Time, error) {
	return time.Parse(LongDateTimeFormat, t)
}
