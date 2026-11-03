package services

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
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

func initializeScopedTokenServiceTest(t *testing.T) (core.Context, *core.UserTokenClaims, *models.TokenRecord) {
	t.Helper()
	initializeTokenServiceTest(t)
	previousStores := *datastore.Container
	t.Cleanup(func() {
		*datastore.Container = previousStores
	})

	err := datastore.InitializeDataStore(&settings.Config{
		DatabaseConfig: &settings.DatabaseConfig{
			DatabaseType: settings.Sqlite3DbType,
			DatabasePath: filepath.Join(t.TempDir(), "tokens.db"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = datastore.Container.TokenStore.SyncStructs(new(models.TokenRecord))
	if err != nil {
		t.Fatal(err)
	}

	config := Tokens.CurrentConfig()
	config.UserAvatarTokenExpiredTimeDuration = time.Hour
	config.TransactionPictureTokenExpiredTimeDuration = 30 * time.Minute
	c := core.NewNullContext()
	_, claims, record, err := Tokens.createToken(c, &models.User{Uid: 123}, core.USER_TOKEN_TYPE_NORMAL, "", "", 2*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	return c, claims, record
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
	actualValue, err := Tokens.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, []byte{byte(core.NormalTokenScopeAll)})
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
	actualValue, err := Tokens.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
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
	actualValue, err := Tokens.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
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
	actualValue, err := Tokens.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
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
	actualValue, err := Tokens.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
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
	actualValue, err := Tokens.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
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
	actualValue, err := Tokens.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
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
	actualValue, err := Tokens.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
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
	}, core.NONE_DERIVED_TOKEN, nil)

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

	firstValue, err := Tokens.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, []byte{byte(core.NormalTokenScopeAll)})
	assert.Equal(t, nil, err)

	secondValue, err := Tokens.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, []byte{byte(core.NormalTokenScopeAll)})
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

	token, err := Tokens.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
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

	actualValue, actualData, err = Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk.", "00"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidTokenType, err)

	actualValue, actualData, err = Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk.unknown", "00"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidTokenType, err)

	actualValue, actualData, err = Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk.a.p", "00"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidTokenType, err)

	actualValue, actualData, err = Tokens.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak.a", "00"})
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

func TestGenerateEncryptedToken_UserAvatar(t *testing.T) {
	initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_NORMAL,
		Secret:          "0123456789",
	}
	expectedValue := &core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: 2000,
	}
	additionalDataBytes, err := expectedValue.MarshalBinary()
	assert.Equal(t, nil, err)

	token, err := Tokens.generateEncryptedToken(record, core.USER_AVATAR_DERIVED_NORMAL_TOKEN, additionalDataBytes)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(token, "ebk-tk.a-"))

	actualRecord, encryptedData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), strings.Split(token, "-"))
	assert.Equal(t, nil, err)
	assert.Equal(t, record.Uid, actualRecord.Uid)
	assert.Equal(t, record.UserTokenId, actualRecord.UserTokenId)
	assert.Equal(t, record.CreatedUnixTime, actualRecord.CreatedUnixTime)
	assert.Equal(t, record.TokenType, actualRecord.TokenType)

	additionalData, err := utils.AESGCMDecrypt(utils.MD5Encode([]byte(record.Secret)), encryptedData)
	assert.Equal(t, nil, err)

	actualValue, err := Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 1000, []core.NormalTokenScope{
		core.NormalTokenScopeUserAvatar,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseNormalTokenAdditionalData_UserAvatarInvalidScope(t *testing.T) {
	additionalData, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: 2000,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	actualValue, err := Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 1000, nil)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)

	actualValue, err = Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 1000, []core.NormalTokenScope{core.NormalTokenScopeTransactionPicture})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
}

func TestParseNormalTokenAdditionalData_UserAvatarExpired(t *testing.T) {
	additionalData, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: 2000,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	actualValue, err := Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 2000, []core.NormalTokenScope{core.NormalTokenScopeUserAvatar})
	assert.Equal(t, nil, err)
	assert.Equal(t, &core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: 2000,
	}, actualValue)

	actualValue, err = Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 2001, []core.NormalTokenScope{core.NormalTokenScopeUserAvatar})
	assert.Equal(t, errs.ErrCurrentTokenExpired, err)
	assert.Nil(t, actualValue)
}

func TestGenerateEncryptedToken_TransactionPicture(t *testing.T) {
	initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_NORMAL,
		Secret:          "0123456789",
	}
	expectedValue := &core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeTransactionPicture,
		DerivedTokenExpiresAt: 2000,
	}
	additionalDataBytes, err := expectedValue.MarshalBinary()
	assert.Equal(t, nil, err)

	token, err := Tokens.generateEncryptedToken(record, core.TRANSACTION_PICTURE_DERIVED_NORMAL_TOKEN, additionalDataBytes)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(token, "ebk-tk.p-"))

	actualRecord, encryptedData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), strings.Split(token, "-"))
	assert.Equal(t, nil, err)
	assert.Equal(t, record.Uid, actualRecord.Uid)
	assert.Equal(t, record.UserTokenId, actualRecord.UserTokenId)
	assert.Equal(t, record.CreatedUnixTime, actualRecord.CreatedUnixTime)
	assert.Equal(t, record.TokenType, actualRecord.TokenType)

	additionalData, err := utils.AESGCMDecrypt(utils.MD5Encode([]byte(record.Secret)), encryptedData)
	assert.Equal(t, nil, err)

	actualValue, err := Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 1000, []core.NormalTokenScope{core.NormalTokenScopeTransactionPicture})
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseNormalTokenAdditionalData_TransactionPictureInvalidScope(t *testing.T) {
	additionalData, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeTransactionPicture,
		DerivedTokenExpiresAt: 2000,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	actualValue, err := Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 1000, nil)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)

	actualValue, err = Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 1000, []core.NormalTokenScope{core.NormalTokenScopeUserAvatar})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
}

func TestParseNormalTokenAdditionalData_TransactionPictureExpired(t *testing.T) {
	additionalData, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeTransactionPicture,
		DerivedTokenExpiresAt: 2000,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	actualValue, err := Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 2000, []core.NormalTokenScope{core.NormalTokenScopeTransactionPicture})
	assert.Equal(t, nil, err)
	assert.Equal(t, &core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeTransactionPicture,
		DerivedTokenExpiresAt: 2000,
	}, actualValue)

	actualValue, err = Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 2001, []core.NormalTokenScope{core.NormalTokenScopeTransactionPicture})
	assert.Equal(t, errs.ErrCurrentTokenExpired, err)
	assert.Nil(t, actualValue)
}

func TestGenerateEncryptedToken_ScopedAdditionalDataInvalidSecret(t *testing.T) {
	initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_NORMAL,
		Secret:          "0123456789",
	}
	additionalDataBytes, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: 2000,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	token, err := Tokens.generateEncryptedToken(record, core.USER_AVATAR_DERIVED_NORMAL_TOKEN, additionalDataBytes)
	assert.Equal(t, nil, err)

	_, encryptedData, err := Tokens.parseEncryptedTokenId(core.NewNullContext(), strings.Split(token, "-"))
	assert.Equal(t, nil, err)

	_, err = utils.AESGCMDecrypt(utils.MD5Encode([]byte("other-secret")), encryptedData)
	assert.NotNil(t, err)

	encryptedData[len(encryptedData)-1] ^= 1
	_, err = utils.AESGCMDecrypt(utils.MD5Encode([]byte(record.Secret)), encryptedData)
	assert.NotNil(t, err)
}

func TestParseNormalTokenAdditionalData_All(t *testing.T) {
	expectedValue := &core.NormalTokenAdditionalData{
		Scope: core.NormalTokenScopeAll,
	}
	actualValue, err := Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), []byte{1}, 2000, nil)

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseNormalTokenAdditionalData_InvalidScope(t *testing.T) {
	additionalData, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeInvalid,
		DerivedTokenExpiresAt: 3000,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	actualValue, err := Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 2000, []core.NormalTokenScope{core.NormalTokenScopeInvalid})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)

	additionalData, err = (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScope(255),
		DerivedTokenExpiresAt: 3000,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	actualValue, err = Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 2000, []core.NormalTokenScope{core.NormalTokenScope(255)})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
}

func TestDeriveUserAvatarToken_CurrentSession(t *testing.T) {
	c, claims, record := initializeScopedTokenServiceTest(t)
	expectedValue := time.Now().Add(time.Hour).Unix()
	token, err := Tokens.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(token, "ebk-tk.a-"))

	actualValue, actualRecord, err := Tokens.ParseToken(c, token, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, nil, err)
	assert.Equal(t, claims.UserTokenId, actualValue.UserTokenId)
	assert.Equal(t, claims.Uid, actualValue.Uid)
	assert.Equal(t, record, actualRecord)
	assert.InDelta(t, expectedValue, actualValue.ExpiresAt, 2)

	_, _, err = Tokens.ParseToken(c, token)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	_, _, err = Tokens.ParseToken(c, token, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	mismatchedToken := strings.Replace(token, "ebk-tk.a-", "ebk-tk.p-", 1)
	_, _, err = Tokens.ParseToken(c, mismatchedToken, core.NormalTokenScopeUserAvatar, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	missingDerivedTokenType := strings.Replace(token, "ebk-tk.a-", "ebk-tk-", 1)
	actualValue, _, err = Tokens.ParseToken(c, missingDerivedTokenType, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)

	normalToken, err := Tokens.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, []byte{byte(core.NormalTokenScopeAll)})
	assert.Equal(t, nil, err)
	mismatchedToken = strings.Replace(normalToken, "ebk-tk-", "ebk-tk.a-", 1)
	_, _, err = Tokens.ParseToken(c, mismatchedToken)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestDeriveUserAvatarToken_SessionExpiration(t *testing.T) {
	c, claims, record := initializeScopedTokenServiceTest(t)
	claims.ExpiresAt = time.Now().Unix() + 60
	expectedValue := claims.ExpiresAt
	token, err := Tokens.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, nil, err)

	actualValue, _, err := Tokens.ParseToken(c, token, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue.ExpiresAt)
}

func TestParseToken_UserAvatarExpired(t *testing.T) {
	c, _, record := initializeScopedTokenServiceTest(t)
	additionalDataBytes, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: time.Now().Unix() - 1,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	token, err := Tokens.generateEncryptedToken(record, core.USER_AVATAR_DERIVED_NORMAL_TOKEN, additionalDataBytes)
	assert.Equal(t, nil, err)

	actualValue, _, err := Tokens.ParseToken(c, token, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, errs.ErrCurrentTokenExpired, err)
	assert.Nil(t, actualValue)
}

func TestDeriveUserAvatarToken_RevokedSession(t *testing.T) {
	c, claims, record := initializeScopedTokenServiceTest(t)
	token, err := Tokens.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, nil, err)

	err = Tokens.DeleteTokenByClaims(c, claims)
	assert.Equal(t, nil, err)

	actualValue, _, err := Tokens.ParseToken(c, token, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
}

func TestDeriveTransactionPictureToken_CurrentSession(t *testing.T) {
	c, claims, record := initializeScopedTokenServiceTest(t)
	expectedValue := time.Now().Add(30 * time.Minute).Unix()
	token, err := Tokens.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(token, "ebk-tk.p-"))

	actualValue, actualRecord, err := Tokens.ParseToken(c, token, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, nil, err)
	assert.Equal(t, claims.UserTokenId, actualValue.UserTokenId)
	assert.Equal(t, claims.Uid, actualValue.Uid)
	assert.Equal(t, record, actualRecord)
	assert.InDelta(t, expectedValue, actualValue.ExpiresAt, 2)

	_, _, err = Tokens.ParseToken(c, token)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	_, _, err = Tokens.ParseToken(c, token, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	mismatchedToken := strings.Replace(token, "ebk-tk.p-", "ebk-tk.a-", 1)
	_, _, err = Tokens.ParseToken(c, mismatchedToken, core.NormalTokenScopeUserAvatar, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	missingDerivedTokenType := strings.Replace(token, "ebk-tk.p-", "ebk-tk-", 1)
	actualValue, _, err = Tokens.ParseToken(c, missingDerivedTokenType, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)

	normalToken, err := Tokens.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, []byte{byte(core.NormalTokenScopeAll)})
	assert.Equal(t, nil, err)
	mismatchedToken = strings.Replace(normalToken, "ebk-tk-", "ebk-tk.p-", 1)
	_, _, err = Tokens.ParseToken(c, mismatchedToken)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestDeriveTransactionPictureToken_SessionExpiration(t *testing.T) {
	c, claims, record := initializeScopedTokenServiceTest(t)
	claims.ExpiresAt = time.Now().Unix() + 60
	expectedValue := claims.ExpiresAt
	token, err := Tokens.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, nil, err)

	actualValue, _, err := Tokens.ParseToken(c, token, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue.ExpiresAt)
}

func TestParseToken_TransactionPictureExpired(t *testing.T) {
	c, _, record := initializeScopedTokenServiceTest(t)
	additionalDataBytes, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeTransactionPicture,
		DerivedTokenExpiresAt: time.Now().Unix() - 1,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	token, err := Tokens.generateEncryptedToken(record, core.TRANSACTION_PICTURE_DERIVED_NORMAL_TOKEN, additionalDataBytes)
	assert.Equal(t, nil, err)

	actualValue, _, err := Tokens.ParseToken(c, token, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, errs.ErrCurrentTokenExpired, err)
	assert.Nil(t, actualValue)
}

func TestDeriveTransactionPictureToken_RevokedSession(t *testing.T) {
	c, claims, record := initializeScopedTokenServiceTest(t)
	token, err := Tokens.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, nil, err)

	err = Tokens.DeleteTokenByClaims(c, claims)
	assert.Equal(t, nil, err)

	actualValue, _, err := Tokens.ParseToken(c, token, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
}

func TestDeriveUserAvatarToken_InvalidClaims(t *testing.T) {
	initializeTokenServiceTest(t)
	expectedValue := ""
	actualValue, err := Tokens.DeriveUserAvatarToken(core.NewNullContext(), nil, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	actualValue, err = Tokens.DeriveUserAvatarToken(core.NewNullContext(), &core.UserTokenClaims{
		Type: core.USER_TOKEN_TYPE_API,
	}, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	actualValue, err = Tokens.DeriveUserAvatarToken(core.NewNullContext(), &core.UserTokenClaims{
		Type:        core.USER_TOKEN_TYPE_NORMAL,
		UserTokenId: "invalid",
	}, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDeriveUserAvatarToken_ExpiredSession(t *testing.T) {
	c, claims, record := initializeScopedTokenServiceTest(t)
	claims.ExpiresAt = time.Now().Unix() - 1
	expectedValue := ""
	actualValue, err := Tokens.DeriveUserAvatarToken(c, claims, record)

	assert.Equal(t, errs.ErrTokenExpired, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDeriveUserAvatarToken_InvalidRecordType(t *testing.T) {
	c, _, _ := initializeScopedTokenServiceTest(t)
	_, claims, record, err := Tokens.createToken(c, &models.User{
		Uid: 123,
	}, core.USER_TOKEN_TYPE_API, "", "", time.Hour)
	assert.Equal(t, nil, err)

	claims.Type = core.USER_TOKEN_TYPE_NORMAL
	expectedValue := ""
	actualValue, err := Tokens.DeriveUserAvatarToken(c, claims, record)

	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDeriveTransactionPictureToken_InvalidClaims(t *testing.T) {
	initializeTokenServiceTest(t)
	expectedValue := ""
	actualValue, err := Tokens.DeriveTransactionPictureToken(core.NewNullContext(), nil, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	actualValue, err = Tokens.DeriveTransactionPictureToken(core.NewNullContext(), &core.UserTokenClaims{
		Type: core.USER_TOKEN_TYPE_API,
	}, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	actualValue, err = Tokens.DeriveTransactionPictureToken(core.NewNullContext(), &core.UserTokenClaims{
		Type:        core.USER_TOKEN_TYPE_NORMAL,
		UserTokenId: "invalid",
	}, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDeriveTransactionPictureToken_ExpiredSession(t *testing.T) {
	c, claims, record := initializeScopedTokenServiceTest(t)
	claims.ExpiresAt = time.Now().Unix() - 1
	expectedValue := ""
	actualValue, err := Tokens.DeriveTransactionPictureToken(c, claims, record)

	assert.Equal(t, errs.ErrTokenExpired, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDeriveTransactionPictureToken_InvalidRecordType(t *testing.T) {
	c, _, _ := initializeScopedTokenServiceTest(t)
	_, claims, record, err := Tokens.createToken(c, &models.User{
		Uid: 123,
	}, core.USER_TOKEN_TYPE_API, "", "", time.Hour)
	assert.Equal(t, nil, err)

	claims.Type = core.USER_TOKEN_TYPE_NORMAL
	expectedValue := ""
	actualValue, err := Tokens.DeriveTransactionPictureToken(c, claims, record)

	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseNormalTokenAdditionalData_InvalidData(t *testing.T) {
	actualValue, err := Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), nil, 2000, nil)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)

	actualValue, err = Tokens.parseNormalTokenAdditionalData(core.NewNullContext(), []byte{1, 0}, 2000, nil)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
}

func TestDeriveUserAvatarToken_InvalidRecord(t *testing.T) {
	initializeTokenServiceTest(t)
	c := core.NewNullContext()
	claims := &core.UserTokenClaims{
		Uid:         123,
		UserTokenId: "456",
		Type:        core.USER_TOKEN_TYPE_NORMAL,
		IssuedAt:    789,
		ExpiresAt:   time.Now().Add(time.Hour).Unix(),
	}
	record := &models.TokenRecord{
		Uid:             claims.Uid,
		UserTokenId:     456,
		TokenType:       claims.Type,
		CreatedUnixTime: claims.IssuedAt,
		ExpiredUnixTime: claims.ExpiresAt,
	}
	expectedValue := ""
	actualValue, err := Tokens.DeriveUserAvatarToken(c, claims, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	record.Uid++
	actualValue, err = Tokens.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	record.Uid = claims.Uid
	record.UserTokenId++
	actualValue, err = Tokens.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	record.UserTokenId = 456
	record.CreatedUnixTime++
	actualValue, err = Tokens.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDeriveUserAvatarToken_RecordExpiration(t *testing.T) {
	c, claims, record := initializeScopedTokenServiceTest(t)
	record.ExpiredUnixTime = time.Now().Unix() + 60
	expectedValue := record.ExpiredUnixTime
	token, err := Tokens.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, nil, err)

	actualValue, _, err := Tokens.ParseToken(c, token, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue.ExpiresAt)

	record.ExpiredUnixTime = time.Now().Unix() - 1
	token, err = Tokens.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, errs.ErrTokenExpired, err)
	assert.Equal(t, "", token)
}

func TestDeriveTransactionPictureToken_InvalidRecord(t *testing.T) {
	initializeTokenServiceTest(t)
	c := core.NewNullContext()
	claims := &core.UserTokenClaims{
		Uid:         123,
		UserTokenId: "456",
		Type:        core.USER_TOKEN_TYPE_NORMAL,
		IssuedAt:    789,
		ExpiresAt:   time.Now().Add(time.Hour).Unix(),
	}
	record := &models.TokenRecord{
		Uid:             claims.Uid,
		UserTokenId:     456,
		TokenType:       claims.Type,
		CreatedUnixTime: claims.IssuedAt,
		ExpiredUnixTime: claims.ExpiresAt,
	}
	expectedValue := ""
	actualValue, err := Tokens.DeriveTransactionPictureToken(c, claims, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	record.Uid++
	actualValue, err = Tokens.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	record.Uid = claims.Uid
	record.UserTokenId++
	actualValue, err = Tokens.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	record.UserTokenId = 456
	record.CreatedUnixTime++
	actualValue, err = Tokens.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDeriveTransactionPictureToken_RecordExpiration(t *testing.T) {
	c, claims, record := initializeScopedTokenServiceTest(t)
	record.ExpiredUnixTime = time.Now().Unix() + 60
	expectedValue := record.ExpiredUnixTime
	token, err := Tokens.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, nil, err)

	actualValue, _, err := Tokens.ParseToken(c, token, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue.ExpiresAt)

	record.ExpiredUnixTime = time.Now().Unix() - 1
	token, err = Tokens.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, errs.ErrTokenExpired, err)
	assert.Equal(t, "", token)
}

func TestParseToken_ReturnsTokenRecord(t *testing.T) {
	c, _, _ := initializeScopedTokenServiceTest(t)
	token, claims, expectedValue, err := Tokens.createToken(c, &models.User{
		Uid: 123,
	}, core.USER_TOKEN_TYPE_NORMAL, "test-agent", "test-context", time.Hour)
	assert.Equal(t, nil, err)

	actualClaims, actualValue, err := Tokens.ParseToken(c, token)
	assert.Equal(t, nil, err)
	assert.Equal(t, claims.UserTokenId, actualClaims.UserTokenId)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, "test-context", actualValue.GetContext())
}
