package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt32ToString(t *testing.T) {
	expectedValue := "-123456789"
	actualValue := Int32ToString(-123456789)
	assert.Equal(t, expectedValue, actualValue)
}

func TestStringToInt32(t *testing.T) {
	expectedValue := -123456789
	actualValue, err := StringToInt32("-123456789")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestStringToInt32_InvalidNumber(t *testing.T) {
	_, err := StringToInt32("")
	assert.NotEqual(t, nil, err)

	_, err = StringToInt32("null")
	assert.NotEqual(t, nil, err)
}

func TestStringTryToInt32_InvalidNumber(t *testing.T) {
	expectedValue := -1
	actualValue := StringTryToInt32("", -1)
	assert.Equal(t, expectedValue, actualValue)

	actualValue = StringTryToInt32("null", -1)
	assert.Equal(t, expectedValue, actualValue)
}

func TestInt64ToString(t *testing.T) {
	expectedValue := "-123456789012345"
	actualValue := Int64ToString(-123456789012345)
	assert.Equal(t, expectedValue, actualValue)
}

func TestInt64ArrayToStringArray(t *testing.T) {
	expectedValue := []string{"0", "1", "-123456789012345", "12345678", "1234567890123456"}
	actualValue := Int64ArrayToStringArray([]int64{0, 1, -123456789012345, 12345678, 1234567890123456})
	assert.EqualValues(t, expectedValue, actualValue)
}

func TestStringToInt64(t *testing.T) {
	expectedValue := int64(-123456789012345)
	actualValue, err := StringToInt64("-123456789012345")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestStringToInt64_InvalidNumber(t *testing.T) {
	_, err := StringToInt64("")
	assert.NotEqual(t, nil, err)

	_, err = StringToInt64("null")
	assert.NotEqual(t, nil, err)
}

func TestStringArrayToInt64Array(t *testing.T) {
	expectedValue := []int64{0, 1, -123456789012345, 12345678, 1234567890123456}
	actualValue, err := StringArrayToInt64Array([]string{"0", "1", "-123456789012345", "12345678", "1234567890123456"})
	assert.Equal(t, nil, err)
	assert.EqualValues(t, expectedValue, actualValue)
}

func TestStringArrayToInt64Array_InvalidNumber(t *testing.T) {
	_, err := StringArrayToInt64Array([]string{"0", "1", "", "12345678", "1234567890123456"})
	assert.NotEqual(t, nil, err)

	_, err = StringArrayToInt64Array([]string{"0", "1", "null", "12345678", "1234567890123456"})
	assert.NotEqual(t, nil, err)
}

func TestStringTryToInt64_InvalidNumber(t *testing.T) {
	expectedValue := int64(-1)
	actualValue := StringTryToInt64("", -1)
	assert.Equal(t, expectedValue, actualValue)

	actualValue = StringTryToInt64("null", -1)
	assert.Equal(t, expectedValue, actualValue)
}

func TestFloat64ToString(t *testing.T) {
	expectedValue := "-123456789.123456"
	actualValue := Float64ToString(-123456789.123456)
	assert.Equal(t, expectedValue, actualValue)
}

func TestStringToFloat64(t *testing.T) {
	expectedValue := -123456789.123456
	actualValue, err := StringToFloat64("-123456789.123456")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}
