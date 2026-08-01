package utils

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsStringOnlyContainsDigits(t *testing.T) {
	actualValue := IsStringOnlyContainsDigits("0123456789")
	assert.True(t, actualValue)

	actualValue = IsStringOnlyContainsDigits("12345a")
	assert.False(t, actualValue)

	actualValue = IsStringOnlyContainsDigits("12345 ")
	assert.False(t, actualValue)
}

func TestAddInt64(t *testing.T) {
	actualValue, valid := AddInt64(math.MaxInt64-1, 1)
	assert.True(t, valid)
	assert.Equal(t, int64(math.MaxInt64), actualValue)

	actualValue, valid = AddInt64(math.MinInt64+1, -1)
	assert.True(t, valid)
	assert.Equal(t, int64(math.MinInt64), actualValue)

	_, valid = AddInt64(math.MaxInt64, 1)
	assert.False(t, valid)

	_, valid = AddInt64(math.MinInt64, -1)
	assert.False(t, valid)
}

func TestSubtractInt64(t *testing.T) {
	actualValue, valid := SubtractInt64(math.MaxInt64, 0)
	assert.True(t, valid)
	assert.Equal(t, int64(math.MaxInt64), actualValue)

	actualValue, valid = SubtractInt64(math.MinInt64, 0)
	assert.True(t, valid)
	assert.Equal(t, int64(math.MinInt64), actualValue)

	_, valid = SubtractInt64(math.MaxInt64, -1)
	assert.False(t, valid)

	_, valid = SubtractInt64(math.MinInt64, 1)
	assert.False(t, valid)
}

func TestParseFirstConsecutiveNumber(t *testing.T) {
	expectedValue := "￥123.45"
	actualValue, success := ParseFirstConsecutiveNumber(expectedValue)
	assert.True(t, success)
	assert.Equal(t, "123.45", actualValue)

	expectedValue = "$-123.45"
	actualValue, success = ParseFirstConsecutiveNumber(expectedValue)
	assert.True(t, success)
	assert.Equal(t, "-123.45", actualValue)

	expectedValue = "$0.12$123.45"
	actualValue, success = ParseFirstConsecutiveNumber(expectedValue)
	assert.True(t, success)
	assert.Equal(t, "0.12", actualValue)

	expectedValue = "$.12"
	actualValue, success = ParseFirstConsecutiveNumber(expectedValue)
	assert.True(t, success)
	assert.Equal(t, "12", actualValue)

	expectedValue = ""
	actualValue, success = ParseFirstConsecutiveNumber(expectedValue)
	assert.False(t, success)

	expectedValue = "xff"
	actualValue, success = ParseFirstConsecutiveNumber(expectedValue)
	assert.False(t, success)
}

func TestTrimTrailingZerosInDecimal(t *testing.T) {
	expectedValue := "123.45"
	actualValue := TrimTrailingZerosInDecimal("123.45000000000")
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "0.12"
	actualValue = TrimTrailingZerosInDecimal("0.12000000000")
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "0.120000000001"
	actualValue = TrimTrailingZerosInDecimal("0.120000000001")
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = ".12"
	actualValue = TrimTrailingZerosInDecimal(".12000000000")
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "12345000000000"
	actualValue = TrimTrailingZerosInDecimal("12345000000000")
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = ""
	actualValue = TrimTrailingZerosInDecimal("")
	assert.Equal(t, expectedValue, actualValue)
}
