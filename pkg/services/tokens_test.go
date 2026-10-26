package services

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func initializeTokenServiceTest(t *testing.T) {
	t.Helper()
	previousConfig := settings.Container.GetCurrentConfig()
	settings.SetCurrentConfig(&settings.Config{SecretKey: "token-test-global-key"})
	t.Cleanup(func() {
		settings.SetCurrentConfig(previousConfig)
	})
}

func encryptTokenIdentityForTest(t *testing.T, uid, userTokenId, createdUnixTime int64, tokenType core.TokenType, encryptedData []byte) string {
	t.Helper()
	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, uid)
	_ = binary.Write(buf, binary.BigEndian, userTokenId)
	_ = binary.Write(buf, binary.BigEndian, createdUnixTime)
	_ = binary.Write(buf, binary.BigEndian, tokenType)
	_, _ = buf.Write(encryptedData)

	token, err := utils.EncryptSecretToHex(buf.String(), Tokens.CurrentConfig().SecretKey)

	if err != nil {
		t.Fatal(err)
	}

	return token
}

func TestGenerateTokenId(t *testing.T) {
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
	}

	expectedValue := "123:789:456"
	actualValue := Tokens.GenerateTokenId(record)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromTokenId(t *testing.T) {
	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
	}
	actualValue, err := Tokens.ParseFromTokenId("123:789:456")

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromTokenId_InvalidFormat(t *testing.T) {
	actualValue, err := Tokens.ParseFromTokenId("")
	assert.Nil(t, actualValue)
	assert.Equal(t, errs.ErrInvalidTokenId, err)

	actualValue, err = Tokens.ParseFromTokenId("123:456")
	assert.Nil(t, actualValue)
	assert.Equal(t, errs.ErrInvalidTokenId, err)

	actualValue, err = Tokens.ParseFromTokenId("123:789:456:extra")
	assert.Nil(t, actualValue)
	assert.Equal(t, errs.ErrInvalidTokenId, err)
}

func TestParseFromTokenId_InvalidUserId(t *testing.T) {
	actualValue, err := Tokens.ParseFromTokenId("invalid:789:456")
	assert.Nil(t, actualValue)
	assert.Equal(t, errs.ErrInvalidTokenId, err)
}

func TestParseFromTokenId_InvalidCreationTime(t *testing.T) {
	actualValue, err := Tokens.ParseFromTokenId("123:invalid:456")
	assert.Nil(t, actualValue)
	assert.Equal(t, errs.ErrInvalidTokenId, err)
}

func TestParseFromTokenId_InvalidUserTokenId(t *testing.T) {
	actualValue, err := Tokens.ParseFromTokenId("123:789:invalid")
	assert.Nil(t, actualValue)
	assert.Equal(t, errs.ErrInvalidTokenId, err)
}

func TestParseFromTokenId_InvalidOverflow(t *testing.T) {
	actualValue, err := Tokens.ParseFromTokenId("123:789:9223372036854775808")
	assert.Nil(t, actualValue)
	assert.Equal(t, errs.ErrInvalidTokenId, err)
}

func TestGenerateEncryptedToken_Normal(t *testing.T) {
	initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_NORMAL,
		Secret:          "0123456789",
	}
	actualValue, err := Tokens.generateEncryptedToken(record)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-tk-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, encryptedData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)
	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_NORMAL,
	}
	assert.Equal(t, expectedValue, actualRecord)

	additionalData, err := utils.AESGCMDecrypt(utils.MD5Encode([]byte(record.Secret)), encryptedData)
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte{byte(core.NormalTokenScopeAll)}, additionalData)

	var actualAdditionalData core.NormalTokenAdditionalData
	err = actualAdditionalData.UnmarshalBinary(additionalData)
	assert.Equal(t, nil, err)
	assert.Equal(t, core.NormalTokenScopeAll, actualAdditionalData.Scope)
}

func TestGenerateEncryptedToken_Require2FA(t *testing.T) {
	initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_REQUIRE_2FA,
		Secret:          "0123456789",
	}
	actualValue, err := Tokens.generateEncryptedToken(record)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-2fa-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, encryptedData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_REQUIRE_2FA,
	}
	assert.Equal(t, expectedValue, actualRecord)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromHex(parts[2], Tokens.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 25, len(plaintext))
}

func TestGenerateEncryptedToken_EmailVerify(t *testing.T) {
	initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_EMAIL_VERIFY,
		Secret:          "0123456789",
	}
	actualValue, err := Tokens.generateEncryptedToken(record)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-vfy-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, encryptedData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_EMAIL_VERIFY,
	}
	assert.Equal(t, expectedValue, actualRecord)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromHex(parts[2], Tokens.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 25, len(plaintext))
}

func TestGenerateEncryptedToken_PasswordReset(t *testing.T) {
	initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_PASSWORD_RESET,
		Secret:          "0123456789",
	}
	actualValue, err := Tokens.generateEncryptedToken(record)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-rst-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, encryptedData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_PASSWORD_RESET,
	}
	assert.Equal(t, expectedValue, actualRecord)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromHex(parts[2], Tokens.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 25, len(plaintext))
}

func TestGenerateEncryptedToken_MCP(t *testing.T) {
	initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_MCP,
		Secret:          "0123456789",
	}
	actualValue, err := Tokens.generateEncryptedToken(record)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-mcp-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, encryptedData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_MCP,
	}
	assert.Equal(t, expectedValue, actualRecord)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromHex(parts[2], Tokens.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 25, len(plaintext))
}

func TestGenerateEncryptedToken_OAuth2RequireVerify(t *testing.T) {
	initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_OAUTH2_CALLBACK_REQUIRE_VERIFY,
		Secret:          "0123456789",
	}
	actualValue, err := Tokens.generateEncryptedToken(record)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-olk-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, encryptedData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_OAUTH2_CALLBACK_REQUIRE_VERIFY,
	}
	assert.Equal(t, expectedValue, actualRecord)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromHex(parts[2], Tokens.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 25, len(plaintext))
}

func TestGenerateEncryptedToken_OAuth2Callback(t *testing.T) {
	initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_OAUTH2_CALLBACK,
		Secret:          "0123456789",
	}
	actualValue, err := Tokens.generateEncryptedToken(record)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-oauth2-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, encryptedData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_OAUTH2_CALLBACK,
	}
	assert.Equal(t, expectedValue, actualRecord)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromHex(parts[2], Tokens.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 25, len(plaintext))
}

func TestGenerateEncryptedToken_API(t *testing.T) {
	initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_API,
		Secret:          "0123456789",
	}
	actualValue, err := Tokens.generateEncryptedToken(record)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-ak-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, encryptedData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_API,
	}
	assert.Equal(t, expectedValue, actualRecord)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromHex(parts[2], Tokens.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 25, len(plaintext))
}

func TestGenerateEncryptedToken_InvalidType(t *testing.T) {
	initializeTokenServiceTest(t)
	actualValue, err := Tokens.generateEncryptedToken(&models.TokenRecord{
		TokenType: core.TokenType(255),
	})

	assert.Equal(t, "", actualValue)
	assert.Equal(t, errs.ErrCurrentInvalidTokenType, err)
}

func TestGenerateEncryptedToken_RandomNonce(t *testing.T) {
	initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_NORMAL,
		Secret:          "0123456789",
	}

	firstValue, err := Tokens.generateEncryptedToken(record)
	assert.Equal(t, nil, err)

	secondValue, err := Tokens.generateEncryptedToken(record)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, firstValue, secondValue)
}

func TestGenerateEncryptedToken_BinaryLayout(t *testing.T) {
	initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             0x0102030405060708,
		UserTokenId:     0x1112131415161718,
		CreatedUnixTime: 0x2122232425262728,
		TokenType:       core.USER_TOKEN_TYPE_API,
	}

	token, err := Tokens.generateEncryptedToken(record)
	assert.Equal(t, nil, err)

	plaintext, err := utils.DecryptSecretFromHex(strings.Split(token, "-")[2], Tokens.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	expectedValue := []byte{
		1, 2, 3, 4, 5, 6, 7, 8,
		17, 18, 19, 20, 21, 22, 23, 24,
		33, 34, 35, 36, 37, 38, 39, 40,
		8,
	}
	assert.Equal(t, expectedValue, []byte(plaintext))
}

func TestParseEncryptedTokenId_InvalidFormat(t *testing.T) {
	initializeTokenServiceTest(t)
	actualValue, actualData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), nil)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	actualValue, actualData, err = Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	actualValue, actualData, err = Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"wrong", "tk", "00"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	actualValue, actualData, err = Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", "00", "extra"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_UnknownType(t *testing.T) {
	initializeTokenServiceTest(t)
	actualValue, actualData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "unknown", "00"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidTokenType, err)
}

func TestParseEncryptedTokenId_InvalidCiphertext(t *testing.T) {
	initializeTokenServiceTest(t)
	actualValue, actualData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", "not-hex"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	actualValue, actualData, err = Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", "00"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_ShortIdentity(t *testing.T) {
	initializeTokenServiceTest(t)
	ciphertext, err := utils.EncryptSecretToHex(string(make([]byte, 24)), Tokens.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)

	actualValue, actualData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	ciphertext, err = utils.EncryptSecretToHex("", Tokens.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)

	actualValue, actualData, err = Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_InvalidUserId(t *testing.T) {
	initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, 0, 456, 789, core.USER_TOKEN_TYPE_NORMAL, nil)
	actualValue, actualData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})

	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_InvalidUserTokenId(t *testing.T) {
	initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, 123, 0, 789, core.USER_TOKEN_TYPE_NORMAL, nil)
	actualValue, actualData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})

	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_InvalidCreationTime(t *testing.T) {
	initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, 123, 456, 0, core.USER_TOKEN_TYPE_NORMAL, nil)
	actualValue, actualData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})

	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_InvalidNegativeUserId(t *testing.T) {
	initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, -1, 456, 789, core.USER_TOKEN_TYPE_NORMAL, nil)
	actualValue, actualData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})

	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_InvalidNegativeUserTokenId(t *testing.T) {
	initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, 123, -1, 789, core.USER_TOKEN_TYPE_NORMAL, nil)
	actualValue, actualData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})

	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_InvalidNegativeCreationTime(t *testing.T) {
	initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, 123, 456, -1, core.USER_TOKEN_TYPE_NORMAL, nil)
	actualValue, actualData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})

	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_TypeMismatch(t *testing.T) {
	initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, 123, 456, 789, core.USER_TOKEN_TYPE_API, nil)
	actualValue, actualData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})

	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_AdditionalData(t *testing.T) {
	initializeTokenServiceTest(t)
	expectedData := []byte{0, 255, ':', '-', 1}
	ciphertext := encryptTokenIdentityForTest(t, 123, 456, 789, core.USER_TOKEN_TYPE_API, expectedData)
	actualValue, actualData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak", ciphertext})
	assert.Equal(t, nil, err)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_API,
	}
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, expectedData, actualData)
}

func TestParseEncryptedTokenId_MaximumIds(t *testing.T) {
	initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, math.MaxInt64, math.MaxInt64, math.MaxInt64, core.USER_TOKEN_TYPE_API, nil)
	actualValue, _, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak", ciphertext})
	assert.Equal(t, nil, err)

	expectedValue := &models.TokenRecord{
		Uid:             math.MaxInt64,
		UserTokenId:     math.MaxInt64,
		CreatedUnixTime: math.MaxInt64,
		TokenType:       core.USER_TOKEN_TYPE_API,
	}
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseEncryptedTokenId_WrongKey(t *testing.T) {
	initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, 123, 456, 789, core.USER_TOKEN_TYPE_API, nil)
	Tokens.CurrentConfig().SecretKey = "different-global-key"

	actualValue, actualData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak", ciphertext})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_TamperedCiphertext(t *testing.T) {
	initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, 123, 456, 789, core.USER_TOKEN_TYPE_API, nil)
	data, err := hex.DecodeString(ciphertext)
	assert.Equal(t, nil, err)

	data[len(data)-1] ^= 1
	actualValue, actualData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak", hex.EncodeToString(data)})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestGetUserAgent(t *testing.T) {
	webContext := core.WrapWebContext(&gin.Context{Request: httptest.NewRequest("GET", "http://localhost/", nil)}, nil)
	webContext.Request.Header.Set("User-Agent", "test-client")
	expectedValue := "test-client"
	actualValue := Tokens.getUserAgent(webContext)
	assert.Equal(t, expectedValue, actualValue)

	webContext.Request.Header.Set("User-Agent", strings.Repeat("x", models.TokenMaxUserAgentLength+1))
	expectedValue = strings.Repeat("x", models.TokenMaxUserAgentLength)
	actualValue = Tokens.getUserAgent(webContext)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetUserAgent_Empty(t *testing.T) {
	actualValue := Tokens.getUserAgent(nil)
	assert.Equal(t, "", actualValue)

	actualValue = Tokens.getUserAgent(&core.WebContext{Context: &gin.Context{}})
	assert.Equal(t, "", actualValue)

	webContext := core.WrapWebContext(&gin.Context{Request: httptest.NewRequest("GET", "http://localhost/", nil)}, nil)
	actualValue = Tokens.getUserAgent(webContext)
	assert.Equal(t, "", actualValue)
}

func TestGetUserTokenId(t *testing.T) {
	actualValue := Tokens.getUserTokenId()
	assert.Greater(t, actualValue, int64(0))

	secondValue := Tokens.getUserTokenId()
	assert.Greater(t, secondValue, int64(0))
	assert.NotEqual(t, actualValue, secondValue)
}
