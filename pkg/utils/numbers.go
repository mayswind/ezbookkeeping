package utils

import (
	"crypto/rand"
	"math/big"
)

// GetRandomInteger returns a random number, the max parameter represents upper limit
func GetRandomInteger(max int) (int, error) {
	result, err := rand.Int(rand.Reader, big.NewInt(int64(max)))

	if err != nil {
		return 0, err
	}

	return int(result.Int64()), nil
}
