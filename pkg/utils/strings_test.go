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

func TestContainsAnyString(t *testing.T) {
	actualValue := ContainsAnyString("test", []string{"test"})
	assert.Equal(t, true, actualValue)

	actualValue = ContainsAnyString("test", []string{"st"})
	assert.Equal(t, true, actualValue)

	actualValue = ContainsAnyString("test", []string{"tt", "tet", "tst", "est"})
	assert.Equal(t, true, actualValue)

	actualValue = ContainsAnyString("test", []string{"tt", "tet", "tst"})
	assert.Equal(t, false, actualValue)
}

func TestGetFirstLowerCharString(t *testing.T) {
	expectedValue := "fooBar"
	actualValue := GetFirstLowerCharString("fooBar")
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = "fooBar"
	actualValue = GetFirstLowerCharString("FooBar")
	assert.Equal(t, expectedValue, actualValue)
}

func TestContainsOnlyOneRune(t *testing.T) {
	actualValue := ContainsOnlyOneRune("-------", '-')
	assert.Equal(t, true, actualValue)

	actualValue = ContainsOnlyOneRune(" -------", '-')
	assert.Equal(t, false, actualValue)
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

func TestMD5EncodeToString(t *testing.T) {
	str := "foobar"
	expectedValue := "3858f62230ac3c915f300c664312c63f"
	actualValue := MD5EncodeToString([]byte(str))
	assert.Equal(t, expectedValue, actualValue)

	str = ""
	expectedValue = "d41d8cd98f00b204e9800998ecf8427e"
	actualValue = MD5EncodeToString([]byte(str))
	assert.Equal(t, expectedValue, actualValue)
}

func TestMD5EncodeToStringWithUidAndSalt(t *testing.T) {
	str := "foobar"
	uid := int64(12345)
	salt := "salt"
	expectedValue := "aeb360a70b82f7ca8aac8d66cd326069"
	actualValue := MD5EncodeToStringWithUidAndSalt([]byte(str), uid, salt)
	assert.Equal(t, expectedValue, actualValue)

	str = ""
	expectedValue = "cb4439f409d0cac03f7e60fb0f7c12f9"
	actualValue = MD5EncodeToStringWithUidAndSalt([]byte(str), uid, salt)
	assert.Equal(t, expectedValue, actualValue)
}

func TestEncodePassword(t *testing.T) {
	password := "foobar"
	salt := "salt"
	expectedValue := "QrpKShMygoe4Ym4ibnA7cNDzCcSonBkgFl69IrtnDmp3oROft3/Td/DNXjsweosa"
	actualValue := EncodePassword(password, salt)
	assert.Equal(t, expectedValue, actualValue)
}

func TestEncryptSecretToBase64AndDecryptSecretFromBase64(t *testing.T) {
	str := "foo"
	key := "bar"
	expectedValue := str

	encryptedStr, err := EncryptSecretToBase64(str, key)
	assert.Equal(t, nil, err)
	actualValue, err := DecryptSecretFromBase64(encryptedStr, key)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestEncryptSecretToHexAndDecryptSecretFromHex(t *testing.T) {
	str := "foo"
	key := "bar"
	expectedValue := str

	encryptedStr, err := EncryptSecretToHex(str, key)
	assert.Equal(t, nil, err)
	assert.Regexp(t, `^[0-9a-f]+$`, encryptedStr)
	actualValue, err := DecryptSecretFromHex(encryptedStr, key)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestEncryptSecretToHexAndDecryptSecretFromHex_EmptyString(t *testing.T) {
	str := ""
	key := "bar"
	expectedValue := str

	encryptedStr, err := EncryptSecretToHex(str, key)
	assert.Equal(t, nil, err)
	actualValue, err := DecryptSecretFromHex(encryptedStr, key)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestEncryptSecretToHexAndDecryptSecretFromHex_Unicode(t *testing.T) {
	str := "图片:token/context?"
	key := "bar"
	expectedValue := str

	encryptedStr, err := EncryptSecretToHex(str, key)
	assert.Equal(t, nil, err)
	actualValue, err := DecryptSecretFromHex(encryptedStr, key)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDecryptSecretFromHex_WrongKey(t *testing.T) {
	encryptedStr, err := EncryptSecretToHex("foo", "bar")
	assert.Equal(t, nil, err)

	_, err = DecryptSecretFromHex(encryptedStr, "wrong-key")
	assert.NotNil(t, err)
}

func TestDecryptSecretFromHex_TamperedCiphertext(t *testing.T) {
	encryptedStr, err := EncryptSecretToHex("foo", "bar")
	assert.Equal(t, nil, err)

	// Change the authentication tag while keeping valid hexadecimal encoding.
	lastCharacter := "0"
	if encryptedStr[len(encryptedStr)-1:] == lastCharacter {
		lastCharacter = "1"
	}
	_, err = DecryptSecretFromHex(encryptedStr[:len(encryptedStr)-1]+lastCharacter, "bar")
	assert.NotNil(t, err)
}

func TestDecryptSecretFromHex_EmptyString(t *testing.T) {
	_, err := DecryptSecretFromHex("", "bar")
	assert.NotNil(t, err)
}

func TestDecryptSecretFromHex_OddLength(t *testing.T) {
	_, err := DecryptSecretFromHex("0", "bar")
	assert.NotNil(t, err)
}

func TestDecryptSecretFromHex_InvalidHex(t *testing.T) {
	_, err := DecryptSecretFromHex("zz", "bar")
	assert.NotNil(t, err)
}

func TestDecryptSecretFromHex_ShortCiphertext(t *testing.T) {
	_, err := DecryptSecretFromHex("00", "bar")
	assert.NotNil(t, err)
}
