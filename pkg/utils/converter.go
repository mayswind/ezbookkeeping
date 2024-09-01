package utils

import (
	"strconv"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// IntToString returns the textual representation of this number
func IntToString(num int) string {
	return strconv.Itoa(num)
}

// StringToInt parses a textual representation of the number to int
func StringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

// StringTryToInt parses a textual representation of the number to int if str is valid,
// or returns the default value
func StringTryToInt(str string, defaultValue int) int {
	num, err := StringToInt(str)

	if err != nil {
		return defaultValue
	}

	return num
}

// StringToInt32 parses a textual representation of the number to int32
func StringToInt32(str string) (int32, error) {
	val, err := strconv.ParseInt(str, 10, 32)

	if err != nil {
		return 0, err
	}

	return int32(val), nil
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

// Float64ToString returns the textual representation of this number
func Float64ToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}

// StringToFloat64 parses a textual representation of the number to float64
func StringToFloat64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

// FormatAmount returns a textual representation of amount
func FormatAmount(amount int64) string {
	displayAmount := Int64ToString(amount)
	negative := displayAmount[0] == '-'

	if negative {
		displayAmount = displayAmount[1:]
	}

	integer := SubString(displayAmount, 0, len(displayAmount)-2)
	decimals := SubString(displayAmount, -2, 2)

	if integer == "" {
		integer = "0"
	}

	if len(decimals) == 0 {
		decimals = "00"
	} else if len(decimals) == 1 {
		decimals = "0" + decimals
	}

	if negative {
		return "-" + integer + "." + decimals
	}

	return integer + "." + decimals
}

// ParseAmount parses a textual representation of amount
func ParseAmount(amount string) (int64, error) {
	if len(amount) < 1 {
		return 0, nil
	}

	sign := int64(1)

	if amount[0] == '-' {
		amount = amount[1:]
		sign = -1
	}

	if len(amount) < 1 {
		return 0, errs.ErrNumberInvalid
	}

	items := strings.Split(amount, ".")

	if len(items) > 2 {
		return 0, errs.ErrNumberInvalid
	}

	var err error
	integer := int64(0)
	decimals := int64(0)

	if len(items[0]) > 0 {
		integer, err = StringToInt64(items[0])

		if err != nil {
			return 0, err
		}

		if integer < 0 {
			return 0, errs.ErrNumberInvalid
		}
	}

	if len(items) == 2 {
		if len(items[1]) > 2 {
			return 0, errs.ErrNumberInvalid
		}

		decimals, err = StringToInt64(items[1])

		if err != nil {
			return 0, err
		}

		if decimals < 0 {
			return 0, errs.ErrNumberInvalid
		}

		if len(items[1]) == 1 {
			decimals = decimals * 10
		}
	}

	return sign*integer*100 + sign*decimals, nil
}
