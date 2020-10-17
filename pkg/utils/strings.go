package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"unicode"

	"golang.org/x/crypto/pbkdf2"

	"github.com/mayswind/lab/pkg/errs"
)

const (
	CHARACTERS                = "!#$&()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_abcdefghijklmnopqrstuvwxyz{|}~"
	NUMBER_AND_LETTERS        = "0123456789abcdefghijklmnopqrstuvwxyz"
	CHARACTERS_LENGTH         = len(CHARACTERS)
	NUMBER_AND_LETTERS_LENGTH = len(NUMBER_AND_LETTERS)
)

func SubString(str string, start int, length int) string {
	chars := []rune(str)
	realLength := len(chars)
	end := 0

	if start < 0 {
		start = realLength - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}

	if start > realLength {
		start = realLength
	}

	if end < 0 {
		end = 0
	}

	if end > realLength {
		end = realLength
	}

	return string(chars[start:end])
}

func GetFirstLowerCharString(s string) string {
	if s == "" {
		return s
	}

	chars := []rune(s)

	if unicode.IsLower(chars[0]) {
		return s
	}

	chars[0] = unicode.ToLower(chars[0])
	return string(chars)
}

func GetRandomString(n int) (string, error) {
	var result = make([]byte, n)

	for i := 0; i < n; i++ {
		index, err := GetRandomInteger(CHARACTERS_LENGTH)

		if err != nil {
			return "", err
		}

		result[i] = CHARACTERS[index]
	}

	return string(result), nil
}

func GetRandomNumberOrLetter(n int) (string, error) {
	var result = make([]byte, n)

	for i := 0; i < n; i++ {
		index, err := GetRandomInteger(NUMBER_AND_LETTERS_LENGTH)

		if err != nil {
			return "", err
		}

		result[i] = NUMBER_AND_LETTERS[index]
	}

	return string(result), nil
}

func MD5Encode(data []byte) []byte {
	m := md5.New()
	m.Write(data)
	return m.Sum(nil)
}

func AESGCMEncrypt(key []byte, plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())

	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plainText, nil)
	result := append(nonce, ciphertext...)

	return result, nil
}

func AESGCMDecrypt(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()

	if len(ciphertext) - nonceSize <= 0 {
		return nil, errs.ErrCiphertextInvalid
	}

	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	plainText, err := aesgcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func EncodePassword(password string, salt string) string {
	encodedPassword := pbkdf2.Key([]byte(password), []byte(salt), 10000, 48, sha256.New) // 256^48 = 64^64
	return strings.TrimRight(base64.StdEncoding.EncodeToString(encodedPassword), "=")
}

func EncyptSecret(secret string, key string) (string, error) {
	encyptedSecret, err := AESGCMEncrypt(MD5Encode([]byte(key)), []byte(secret)) // md5encode make the aes key's length to 16

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encyptedSecret), nil
}

func DecryptSecret(encyptedSecret string, key string) (string, error) {
	encyptedData, err := base64.StdEncoding.DecodeString(encyptedSecret)

	if err != nil {
		return "", err
	}

	secret, err := AESGCMDecrypt(MD5Encode([]byte(key)), []byte(encyptedData))

	if err != nil {
		return "", err
	}

	return string(secret), nil
}
