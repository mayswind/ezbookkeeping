package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"strings"
	"unicode"

	"golang.org/x/crypto/pbkdf2"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

const (
	availableCharacters                      = "!#$&()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_abcdefghijklmnopqrstuvwxyz{|}~"
	availableNumberAndLetters                = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	availableNumberAndLowercaseLetters       = "0123456789abcdefghijklmnopqrstuvwxyz"
	availableCharactersLength                = len(availableCharacters)
	availableNumberAndLettersLength          = len(availableNumberAndLetters)
	availableNumberAndLowercaseLettersLength = len(availableNumberAndLowercaseLetters)
	base62Characters                         = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base62ChunkSize                          = 10
	base62ChunkBase                          = 839299365868340224 // 62^10
)

// SubString returns part of the source string according to start index and length
func SubString(str string, start int, length int) string {
	chars := []rune(str)
	realLength := len(chars)
	end := 0

	if start < 0 {
		start = realLength + start
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

// ContainsAnyString returns whether the specified string contains any string of sub string slice
func ContainsAnyString(s string, substrs []string) bool {
	for i := 0; i < len(substrs); i++ {
		if strings.Index(s, substrs[i]) >= 0 {
			return true
		}
	}

	return false
}

// GetFirstLowerCharString returns the source string parameter, but makes the first character lower case
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

// ContainsOnlyOneRune returns the source string only contains one character
func ContainsOnlyOneRune(s string, r rune) bool {
	if len(s) < 1 {
		return false
	}

	for i := 0; i < len(s); i++ {
		if rune(s[i]) != r {
			return false
		}
	}

	return true
}

// GetRandomString returns a random string of which length is n
func GetRandomString(n int) (string, error) {
	var result = make([]byte, n)

	for i := 0; i < n; i++ {
		index, err := GetRandomInteger(availableCharactersLength)

		if err != nil {
			return "", err
		}

		result[i] = availableCharacters[index]
	}

	return string(result), nil
}

// GetRandomNumberOrLetter returns a random string which only contains number or letter characters
func GetRandomNumberOrLetter(n int) (string, error) {
	var result = make([]byte, n)

	for i := 0; i < n; i++ {
		index, err := GetRandomInteger(availableNumberAndLettersLength)

		if err != nil {
			return "", err
		}

		result[i] = availableNumberAndLetters[index]
	}

	return string(result), nil
}

// GetRandomNumberOrLowercaseLetter returns a random string which only contains number or letter characters
func GetRandomNumberOrLowercaseLetter(n int) (string, error) {
	var result = make([]byte, n)

	for i := 0; i < n; i++ {
		index, err := GetRandomInteger(availableNumberAndLowercaseLettersLength)

		if err != nil {
			return "", err
		}

		result[i] = availableNumberAndLowercaseLetters[index]
	}

	return string(result), nil
}

// MD5Encode returns a hashed string by md5
func MD5Encode(data []byte) []byte {
	m := md5.New()
	m.Write(data)
	return m.Sum(nil)
}

// MD5EncodeToString returns a hashed string by md5
func MD5EncodeToString(data []byte) string {
	hash := MD5Encode(data)
	return hex.EncodeToString(hash)
}

// MD5EncodeToStringWithUidAndSalt returns a hashed string by md5 with uid and salt
func MD5EncodeToStringWithUidAndSalt(data []byte, uid int64, salt string) string {
	hash := MD5Encode(append(data, []byte(salt+Int64ToString(uid))...))
	hash = MD5Encode(append(hash, []byte(Int64ToString(uid)+salt)...))
	return hex.EncodeToString(hash)
}

// AESGCMEncrypt returns a encrypted string by aes-gcm
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

// AESGCMDecrypt returns a decrypted string by aes-gcm
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

	if len(ciphertext)-nonceSize <= 0 {
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

// EncodePassword returns a encoded password
func EncodePassword(password string, salt string) string {
	encodedPassword := pbkdf2.Key([]byte(password), []byte(salt), 10000, 48, sha256.New) // 256^48 = 64^64
	return strings.TrimRight(base64.StdEncoding.EncodeToString(encodedPassword), "=")
}

// EncryptSecretToBase64 returns an encrypted secret encoded as Base64
func EncryptSecretToBase64(secret string, key string) (string, error) {
	encryptedSecret, err := AESGCMEncrypt(MD5Encode([]byte(key)), []byte(secret)) // md5encode make the aes key's length to 16

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedSecret), nil
}

// EncryptSecretToBase62 returns an encrypted secret encoded as Base62
func EncryptSecretToBase62(secret string, key string) (string, error) {
	encryptedSecret, err := AESGCMEncrypt(MD5Encode([]byte(key)), []byte(secret)) // md5encode make the aes key's length to 16

	if err != nil {
		return "", err
	}

	return encodeBase62(encryptedSecret), nil
}

// EncryptSecretToHex returns an encrypted secret encoded as hexadecimal
func EncryptSecretToHex(secret string, key string) (string, error) {
	encryptedSecret, err := AESGCMEncrypt(MD5Encode([]byte(key)), []byte(secret)) // md5encode make the aes key's length to 16

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encryptedSecret), nil
}

// DecryptSecretFromBase64 returns the decrypted secret from Base64 ciphertext
func DecryptSecretFromBase64(encyptedSecret string, key string) (string, error) {
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

// DecryptSecretFromBase62 returns the decrypted secret from Base62 ciphertext
func DecryptSecretFromBase62(encryptedSecret string, key string) (string, error) {
	encryptedData, err := decodeBase62(encryptedSecret)

	if err != nil {
		return "", err
	}

	secret, err := AESGCMDecrypt(MD5Encode([]byte(key)), encryptedData)

	if err != nil {
		return "", err
	}

	return string(secret), nil
}

// DecryptSecretFromHex returns the decrypted secret from hexadecimal ciphertext
func DecryptSecretFromHex(encryptedSecret string, key string) (string, error) {
	encryptedData, err := hex.DecodeString(encryptedSecret)

	if err != nil {
		return "", err
	}

	secret, err := AESGCMDecrypt(MD5Encode([]byte(key)), encryptedData)

	if err != nil {
		return "", err
	}

	return string(secret), nil
}

func encodeBase62(data []byte) string {
	leadingZeros := 0

	for leadingZeros < len(data) && data[leadingZeros] == 0 {
		leadingZeros++
	}

	if leadingZeros == len(data) {
		return strings.Repeat("0", leadingZeros)
	}

	var value, base, remainder big.Int
	value.SetBytes(data)
	base.SetUint64(base62ChunkBase)
	result := make([]byte, len(data)*8/5+1)
	index := len(result)

	for value.Sign() > 0 {
		value.QuoRem(&value, &base, &remainder)
		chunk := remainder.Uint64()

		for i := 0; i < base62ChunkSize; i++ {
			index--
			result[index] = base62Characters[chunk%62]
			chunk /= 62

			if value.Sign() == 0 && chunk == 0 {
				break
			}
		}
	}

	for i := 0; i < leadingZeros; i++ {
		index--
		result[index] = '0'
	}

	return string(result[index:])
}

func decodeBase62(str string) ([]byte, error) {
	if len(str) < 1 {
		return nil, errs.ErrCiphertextInvalid
	}

	var value, base, chunkValue big.Int

	for i := 0; i < len(str); {
		var chunk uint64
		power := uint64(1)
		end := min(i+base62ChunkSize, len(str))

		for ; i < end; i++ {
			var digit uint64

			switch {
			case str[i] >= '0' && str[i] <= '9':
				digit = uint64(str[i] - '0')
			case str[i] >= 'A' && str[i] <= 'Z':
				digit = uint64(str[i]-'A') + 10
			case str[i] >= 'a' && str[i] <= 'z':
				digit = uint64(str[i]-'a') + 36
			default:
				return nil, errs.ErrCiphertextInvalid
			}

			chunk = chunk*62 + digit
			power *= 62
		}

		base.SetUint64(power)
		chunkValue.SetUint64(chunk)
		value.Mul(&value, &base)
		value.Add(&value, &chunkValue)
	}

	leadingZeros := 0

	for leadingZeros < len(str) && str[leadingZeros] == '0' {
		leadingZeros++
	}

	return append(make([]byte, leadingZeros), value.Bytes()...), nil
}
