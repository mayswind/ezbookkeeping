package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
