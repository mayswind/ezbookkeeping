package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubString(t *testing.T) {
	str := "foobar"
	expectedValue := "f"
	actualValue := SubString(str, 0, 1)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "fo"
	actualValue = SubString(str, 0, 2)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "oba"
	actualValue = SubString(str, 2, 3)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "foobar"
	actualValue = SubString(str, 0, 6)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "ba"
	actualValue = SubString(str, -3, 2)
	assert.Equal(t, expectedValue, actualValue)
}

func TestSubString_EmptyStr(t *testing.T) {
	str := ""
	expectedValue := ""
	actualValue := SubString(str, 0, 1)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = ""
	actualValue = SubString(str, 0, 2)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = ""
	actualValue = SubString(str, 2, 3)
	assert.Equal(t, expectedValue, actualValue)
}

func TestSubString_OverBoundary(t *testing.T) {
	str := "foobar"
	expectedValue := ""
	actualValue := SubString(str, 100, 1)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "foobar"
	actualValue = SubString(str, 0, 100)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "obar"
	actualValue = SubString(str, 2, 100)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "bar"
	actualValue = SubString(str, -3, 100)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "foobar"
	actualValue = SubString(str, -100, 100)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = ""
	actualValue = SubString(str, 0, -100)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetRandomString(t *testing.T) {
	actualValue, err := GetRandomString(10)
	assert.Equal(t, nil, err)
	assert.Len(t, actualValue, 10)
}

func TestGetRandomNumberOrLetter(t *testing.T) {
	actualValue, err := GetRandomNumberOrLetter(10)
	assert.Equal(t, nil, err)
	assert.Len(t, actualValue, 10)
}
func TestGetRandomNumberOrLowercaseLetter(t *testing.T) {
	actualValue, err := GetRandomNumberOrLowercaseLetter(10)
	assert.Equal(t, nil, err)
	assert.Len(t, actualValue, 10)
}

func TestMD5Encode(t *testing.T) {
	str := "foobar"
	expectedValue := "3858f62230ac3c915f300c664312c63f"
	actualData := MD5Encode([]byte(str))
	actualValue := fmt.Sprintf("%x", actualData)
	assert.Equal(t, expectedValue, actualValue)

	str = ""
	expectedValue = "d41d8cd98f00b204e9800998ecf8427e"
	actualData = MD5Encode([]byte(str))
	actualValue = fmt.Sprintf("%x", actualData)
	assert.Equal(t, expectedValue, actualValue)
}

func TestEncryptSecretAndDecryptSecret(t *testing.T) {
	str := "foo"
	key := "bar"
	expectedValue := str

	encryptedStr, err := EncryptSecret(str, key)
	actualValue, err := DecryptSecret(encryptedStr, key)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}
