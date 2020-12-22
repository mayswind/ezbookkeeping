package utils

import "strconv"

func Int32ToString(num int) string {
	return strconv.Itoa(num)
}

func StringToInt32(str string) (int, error) {
	return strconv.Atoi(str)
}

func StringTryToInt32(str string, defaultValue int) int {
	num, err := StringToInt32(str)

	if err != nil {
		return defaultValue
	}

	return num
}

func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func StringToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func StringTryToInt64(str string, defaultValue int64) int64 {
	num, err := StringToInt64(str)

	if err != nil {
		return defaultValue
	}

	return num
}
