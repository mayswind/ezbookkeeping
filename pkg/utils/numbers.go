package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// GetRandomInteger returns a random number, the max parameter represents upper limit
func GetRandomInteger(max int) (int, error) {
	result, err := rand.Int(rand.Reader, big.NewInt(int64(max)))

	if err != nil {
		return 0, err
	}

	return int(result.Int64()), nil
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
