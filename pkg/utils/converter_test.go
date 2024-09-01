package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntToString(t *testing.T) {
	expectedValue := "-123456789"
	actualValue := IntToString(-123456789)
	assert.Equal(t, expectedValue, actualValue)
}

func TestStringToInt(t *testing.T) {
	expectedValue := -123456789
	actualValue, err := StringToInt("-123456789")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestStringToInt_InvalidNumber(t *testing.T) {
	_, err := StringToInt("")
	assert.NotEqual(t, nil, err)

	_, err = StringToInt("null")
	assert.NotEqual(t, nil, err)
}

func TestStringToInt32(t *testing.T) {
	expectedValue := int32(-123456789)
	actualValue, err := StringToInt32("-123456789")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestStringToInt32_OutOfRange(t *testing.T) {
	_, err := StringToInt32("2147483648")
	assert.NotEqual(t, nil, err)

	_, err = StringToInt32("-2147483649")
	assert.NotEqual(t, nil, err)
}

func TestStringToInt32_InvalidNumber(t *testing.T) {
	_, err := StringToInt32("")
	assert.NotEqual(t, nil, err)

	_, err = StringToInt32("null")
	assert.NotEqual(t, nil, err)
}

func TestStringTryToInt32_InvalidNumber(t *testing.T) {
	expectedValue := -1
	actualValue := StringTryToInt("", -1)
	assert.Equal(t, expectedValue, actualValue)

	actualValue = StringTryToInt("null", -1)
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

func TestFormatAmount(t *testing.T) {
	expectedValue := "0.00"
	actualValue := FormatAmount(0)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "0.00"
	actualValue = FormatAmount(-0)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "0.01"
	actualValue = FormatAmount(1)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "-0.01"
	actualValue = FormatAmount(-1)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "0.10"
	actualValue = FormatAmount(10)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "-0.10"
	actualValue = FormatAmount(-10)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "0.12"
	actualValue = FormatAmount(12)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "-0.12"
	actualValue = FormatAmount(-12)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "1.23"
	actualValue = FormatAmount(123)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "-1.23"
	actualValue = FormatAmount(-123)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "12.34"
	actualValue = FormatAmount(1234)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "-12.34"
	actualValue = FormatAmount(-1234)
	assert.Equal(t, expectedValue, actualValue)
}
