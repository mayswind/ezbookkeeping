package utils

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strings"
)

var (
	numberPattern = regexp.MustCompile("(-?\\d+)(\\.\\d+)?")
)

// IsStringOnlyContainsDigits returns whether the specified string only contains digit characters
func IsStringOnlyContainsDigits(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] < '0' || str[i] > '9' {
			return false
		}
	}

	return true
}

// GetRandomInteger returns a random number, the max parameter represents upper limit
func GetRandomInteger(max int) (int, error) {
	result, err := rand.Int(rand.Reader, big.NewInt(int64(max)))

	if err != nil {
		return 0, err
	}

	return int(result.Int64()), nil
}

// ParseFirstConsecutiveNumber returns the first consecutive number in the specified string
func ParseFirstConsecutiveNumber(str string) (string, bool) {
	result := numberPattern.FindAllString(str, 1)

	if len(result) > 0 {
		return result[0], true
	} else {
		return "", false
	}
}

// TrimTrailingZerosInDecimal returns a textual number without trailing zeros in decimal
func TrimTrailingZerosInDecimal(num string) string {
	if len(num) < 1 {
		return num
	}

	dotPosition := strings.Index(num, ".")

	if dotPosition < 0 {
		return num
	}

	lastNonZeroPosition := len(num)

	for i := len(num) - 1; i > dotPosition+1; i-- {
		if num[i] == '0' {
			lastNonZeroPosition = i
		} else {
			break
		}
	}

	if lastNonZeroPosition >= len(num) {
		return num
	}

	return num[0:lastNonZeroPosition]
}
