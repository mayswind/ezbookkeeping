package utils

import "strconv"

// Int32ToString returns the textual representation of this number
func Int32ToString(num int) string {
	return strconv.Itoa(num)
}

// StringToInt32 parses a textual representation of the number to int32
func StringToInt32(str string) (int, error) {
	return strconv.Atoi(str)
}

// StringTryToInt32 parses a textual representation of the number to int32 if str is valid,
// or returns the default value
func StringTryToInt32(str string, defaultValue int) int {
	num, err := StringToInt32(str)

	if err != nil {
		return defaultValue
	}

	return num
}

// Int64ToString returns the textual representation of this number
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

// Int64ArrayToStringArray returns a array of textual representation of these numbers
func Int64ArrayToStringArray(num []int64) []string {
	ret := make([]string, 0, len(num))

	for i := 0; i < len(num); i++ {
		ret = append(ret, Int64ToString(num[i]))
	}

	return ret
}

// StringToInt64 parses a textual representation of the number to int64
func StringToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

// StringArrayToInt64Array parses a series textual representations of the numbers to int64 array
func StringArrayToInt64Array(strs []string) ([]int64, error) {
	ret := make([]int64, 0, len(strs))

	for i := 0; i < len(strs); i++ {
		val, err := StringToInt64(strs[i])

		if err != nil {
			return nil, err
		}

		ret = append(ret, val)
	}

	return ret, nil
}

// StringTryToInt64 parses a textual representation of the number to int64 if str is valid,
// or returns the default value
func StringTryToInt64(str string, defaultValue int64) int64 {
	num, err := StringToInt64(str)

	if err != nil {
		return defaultValue
	}

	return num
}
