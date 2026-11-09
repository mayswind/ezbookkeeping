package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"math"
	"net/http/httptest"
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

type tokenServiceTestDatabase struct {
	datastore.Database
	records map[[3]int64]models.TokenRecord
}

type tokenServiceTestSession struct {
	datastore.DatabaseSession
	database *tokenServiceTestDatabase
	query    any
	args     []any
}

func (db *tokenServiceTestDatabase) NewSession(c core.Context) datastore.DatabaseSession {
	return &tokenServiceTestSession{
		database: db,
	}
}

func (db *tokenServiceTestDatabase) DoTransaction(c core.Context, fn func(datastore.DatabaseSession) error) error {
	return fn(db.NewSession(c))
}

func (s *tokenServiceTestSession) Where(query any, args ...any) datastore.DatabaseSession {
	s.query = query
	s.args = args
	return s
}

func (s *tokenServiceTestSession) Limit(limit int, start ...int) datastore.DatabaseSession {
	return s
}

func (s *tokenServiceTestSession) Get(beans ...any) (bool, error) {
	if len(beans) != 1 {
		return false, nil
	}

	tokenRecord, ok := beans[0].(*models.TokenRecord)

	if !ok || tokenRecord == nil {
		return false, nil
	}

	key, ok := s.getTokenRecordKey()

	if !ok {
		return false, nil
	}

	record, exists := s.database.records[key]

	if exists {
		*tokenRecord = record
	}

	return exists, nil
}

func (s *tokenServiceTestSession) Insert(beans ...any) (int64, error) {
	for _, bean := range beans {
		record, ok := bean.(*models.TokenRecord)

		if !ok || record == nil {
			return 0, nil
		}
	}

	for _, bean := range beans {
		record := bean.(*models.TokenRecord)
		key := [3]int64{record.Uid, record.UserTokenId, record.CreatedUnixTime}
		s.database.records[key] = *record
	}

	return int64(len(beans)), nil
}

func (s *tokenServiceTestSession) Delete(beans ...any) (int64, error) {
	if len(beans) != 1 {
		return 0, nil
	}

	record, ok := beans[0].(*models.TokenRecord)

	if !ok || record == nil {
		return 0, nil
	}

	key, ok := s.getTokenRecordKey()

	if !ok {
		return 0, nil
	}

	if _, exists := s.database.records[key]; !exists {
		return 0, nil
	}

	delete(s.database.records, key)
	return 1, nil
}

func (s *tokenServiceTestSession) getTokenRecordKey() ([3]int64, bool) {
	var key [3]int64
	if s.query != "uid=? AND user_token_id=? AND created_unix_time=?" || len(s.args) != len(key) {
		return key, false
	}

	for i, arg := range s.args {
		value, ok := arg.(int64)

		if !ok {
			return key, false
		}

		key[i] = value
	}

	return key, true
}

func initializeTokenServiceTest(t *testing.T) *TokenService {
	t.Helper()
	return &TokenService{
		ServiceUsingConfig: ServiceUsingConfig{
			container: settings.NewConfigContainer(&settings.Config{SecretKey: "token-test-global-key"}),
		},
	}
}

func initializeScopedTokenServiceTest(t *testing.T) (*TokenService, core.Context, *core.UserTokenClaims, *models.TokenRecord) {
	t.Helper()
	service := initializeTokenServiceTest(t)
	tokenStore, err := datastore.NewDataStore(&tokenServiceTestDatabase{
		records: make(map[[3]int64]models.TokenRecord),
	})
	if err != nil {
		t.Fatal(err)
	}

	service.ServiceUsingDB = ServiceUsingDB{
		container: &datastore.DataStoreContainer{
			TokenStore: tokenStore,
		},
	}

	config := service.CurrentConfig()
	config.UserAvatarTokenExpiredTimeDuration = time.Hour
	config.UserCustomIconTokenExpiredTimeDuration = 30 * time.Minute
	config.TransactionPictureTokenExpiredTimeDuration = 30 * time.Minute

	c := core.NewNullContext()

	_, claims, record, err := service.createToken(c, &models.User{Uid: 123}, core.USER_TOKEN_TYPE_NORMAL, "", "", 2*time.Hour)

	if err != nil {
		t.Fatal(err)
	}

	return service, c, claims, record
}

func encryptTokenIdentityForTest(t *testing.T, service *TokenService, uid, userTokenId, createdUnixTime int64, tokenType core.TokenType, encryptedData []byte) string {
	t.Helper()
	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, uid)
	_ = binary.Write(buf, binary.BigEndian, userTokenId)
	_ = binary.Write(buf, binary.BigEndian, createdUnixTime)
	_ = binary.Write(buf, binary.BigEndian, tokenType)
	_, _ = buf.Write(make([]byte, 12))
	_, _ = buf.Write(encryptedData)

	token, err := utils.EncryptSecretToBase62(buf.String(), service.CurrentConfig().SecretKey)

	if err != nil {
		t.Fatal(err)
	}

	return token
}

func TestGenerateTokenId(t *testing.T) {
	service := &TokenService{}
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
	}

	expectedValue := "123:789:456"
	actualValue := service.GenerateTokenId(record)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromTokenId(t *testing.T) {
	service := &TokenService{}
	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
	}
	actualValue, err := service.ParseFromTokenId("123:789:456")

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseFromTokenId_InvalidFormat(t *testing.T) {
	service := &TokenService{}
	actualValue, err := service.ParseFromTokenId("")
	assert.Nil(t, actualValue)
	assert.Equal(t, errs.ErrInvalidTokenId, err)

	actualValue, err = service.ParseFromTokenId("123:456")
	assert.Nil(t, actualValue)
	assert.Equal(t, errs.ErrInvalidTokenId, err)

	actualValue, err = service.ParseFromTokenId("123:789:456:extra")
	assert.Nil(t, actualValue)
	assert.Equal(t, errs.ErrInvalidTokenId, err)
}

func TestParseFromTokenId_InvalidUserId(t *testing.T) {
	service := &TokenService{}
	actualValue, err := service.ParseFromTokenId("invalid:789:456")
	assert.Nil(t, actualValue)
	assert.Equal(t, errs.ErrInvalidTokenId, err)
}

func TestParseFromTokenId_InvalidCreationTime(t *testing.T) {
	service := &TokenService{}
	actualValue, err := service.ParseFromTokenId("123:invalid:456")
	assert.Nil(t, actualValue)
	assert.Equal(t, errs.ErrInvalidTokenId, err)
}

func TestParseFromTokenId_InvalidUserTokenId(t *testing.T) {
	service := &TokenService{}
	actualValue, err := service.ParseFromTokenId("123:789:invalid")
	assert.Nil(t, actualValue)
	assert.Equal(t, errs.ErrInvalidTokenId, err)
}

func TestParseFromTokenId_InvalidOverflow(t *testing.T) {
	service := &TokenService{}
	actualValue, err := service.ParseFromTokenId("123:789:9223372036854775808")
	assert.Nil(t, actualValue)
	assert.Equal(t, errs.ErrInvalidTokenId, err)
}

func TestGenerateEncryptedToken_Normal(t *testing.T) {
	service := initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_NORMAL,
		Secret:          "0123456789",
	}
	actualValue, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, []byte{byte(core.NormalTokenScopeAll)})
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-tk-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, actualSignature, encryptedData, err := service.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)
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
	service := initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_REQUIRE_2FA,
		Secret:          "0123456789",
	}
	actualValue, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-2fa-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, actualSignature, encryptedData, err := service.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_REQUIRE_2FA,
	}
	assert.Equal(t, expectedValue, actualRecord)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromBase62(parts[2], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 37, len(plaintext))
}

func TestGenerateEncryptedToken_EmailVerify(t *testing.T) {
	service := initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_EMAIL_VERIFY,
		Secret:          "0123456789",
	}
	actualValue, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-vfy-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, actualSignature, encryptedData, err := service.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_EMAIL_VERIFY,
	}
	assert.Equal(t, expectedValue, actualRecord)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromBase62(parts[2], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 37, len(plaintext))
}

func TestGenerateEncryptedToken_PasswordReset(t *testing.T) {
	service := initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_PASSWORD_RESET,
		Secret:          "0123456789",
	}
	actualValue, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-rst-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, actualSignature, encryptedData, err := service.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_PASSWORD_RESET,
	}
	assert.Equal(t, expectedValue, actualRecord)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromBase62(parts[2], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 37, len(plaintext))
}

func TestGenerateEncryptedToken_MCP(t *testing.T) {
	service := initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_MCP,
		Secret:          "0123456789",
	}
	actualValue, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-mcp-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, actualSignature, encryptedData, err := service.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_MCP,
	}
	assert.Equal(t, expectedValue, actualRecord)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromBase62(parts[2], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 37, len(plaintext))
}

func TestGenerateEncryptedToken_OAuth2RequireVerify(t *testing.T) {
	service := initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_OAUTH2_CALLBACK_REQUIRE_VERIFY,
		Secret:          "0123456789",
	}
	actualValue, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-olk-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, actualSignature, encryptedData, err := service.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_OAUTH2_CALLBACK_REQUIRE_VERIFY,
	}
	assert.Equal(t, expectedValue, actualRecord)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromBase62(parts[2], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 37, len(plaintext))
}

func TestGenerateEncryptedToken_OAuth2Callback(t *testing.T) {
	service := initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_OAUTH2_CALLBACK,
		Secret:          "0123456789",
	}
	actualValue, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-oauth2-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, actualSignature, encryptedData, err := service.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_OAUTH2_CALLBACK,
	}
	assert.Equal(t, expectedValue, actualRecord)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromBase62(parts[2], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 37, len(plaintext))
}

func TestGenerateEncryptedToken_API(t *testing.T) {
	service := initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_API,
		Secret:          "0123456789",
	}
	actualValue, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, "ebk-ak-"))

	parts := strings.Split(actualValue, "-")
	actualRecord, actualSignature, encryptedData, err := service.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)

	expectedValue := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_API,
	}
	assert.Equal(t, expectedValue, actualRecord)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromBase62(parts[2], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 37, len(plaintext))
}

func TestGenerateEncryptedToken_InvalidType(t *testing.T) {
	service := initializeTokenServiceTest(t)
	actualValue, err := service.generateEncryptedToken(&models.TokenRecord{
		TokenType: core.TokenType(255),
	}, core.NONE_DERIVED_TOKEN, nil)

	assert.Equal(t, "", actualValue)
	assert.Equal(t, errs.ErrCurrentInvalidTokenType, err)
}

func TestGenerateEncryptedToken_RandomNonce(t *testing.T) {
	service := initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_NORMAL,
		Secret:          "0123456789",
	}

	firstValue, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, []byte{byte(core.NormalTokenScopeAll)})
	assert.Equal(t, nil, err)

	secondValue, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, []byte{byte(core.NormalTokenScopeAll)})
	assert.Equal(t, nil, err)
	assert.NotEqual(t, firstValue, secondValue)
}

func TestGenerateEncryptedToken_BinaryLayout(t *testing.T) {
	service := initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             0x0102030405060708,
		UserTokenId:     0x1112131415161718,
		CreatedUnixTime: 0x2122232425262728,
		TokenType:       core.USER_TOKEN_TYPE_API,
		Secret:          "0123456789",
	}

	token, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)

	plaintext, err := utils.DecryptSecretFromBase62(strings.Split(token, "-")[2], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	expectedValue := []byte{
		1, 2, 3, 4, 5, 6, 7, 8,
		17, 18, 19, 20, 21, 22, 23, 24,
		33, 34, 35, 36, 37, 38, 39, 40,
		8,
	}
	signature := hmac.New(sha1.New, []byte(record.Secret))
	_, _ = signature.Write(expectedValue)
	expectedValue = append(expectedValue, signature.Sum(nil)[:12]...)
	assert.Equal(t, expectedValue, []byte(plaintext))
}

func TestParseEncryptedTokenId_InvalidFormat(t *testing.T) {
	service := initializeTokenServiceTest(t)
	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), nil)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), []string{"wrong", "tk", "00"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", "00", "extra"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_UnknownType(t *testing.T) {
	service := initializeTokenServiceTest(t)
	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "unknown", "00"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidTokenType, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk.", "00"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidTokenType, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk.unknown", "00"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidTokenType, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk.a.p", "00"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidTokenType, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak.a", "00"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidTokenType, err)
}

func TestParseEncryptedTokenId_InvalidCiphertext(t *testing.T) {
	service := initializeTokenServiceTest(t)
	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", "not-base62"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", "00"})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_InvalidUserId(t *testing.T) {
	service := initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, service, 0, 456, 789, core.USER_TOKEN_TYPE_NORMAL, nil)
	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})

	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_InvalidUserTokenId(t *testing.T) {
	service := initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, service, 123, 0, 789, core.USER_TOKEN_TYPE_NORMAL, nil)
	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})

	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_InvalidCreationTime(t *testing.T) {
	service := initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, service, 123, 456, 0, core.USER_TOKEN_TYPE_NORMAL, nil)
	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})

	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_InvalidNegativeUserId(t *testing.T) {
	service := initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, service, -1, 456, 789, core.USER_TOKEN_TYPE_NORMAL, nil)
	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})

	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_InvalidNegativeUserTokenId(t *testing.T) {
	service := initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, service, 123, -1, 789, core.USER_TOKEN_TYPE_NORMAL, nil)
	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})

	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_InvalidNegativeCreationTime(t *testing.T) {
	service := initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, service, 123, 456, -1, core.USER_TOKEN_TYPE_NORMAL, nil)
	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})

	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_TypeMismatch(t *testing.T) {
	service := initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, service, 123, 456, 789, core.USER_TOKEN_TYPE_API, nil)
	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "tk", ciphertext})

	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_AdditionalData(t *testing.T) {
	service := initializeTokenServiceTest(t)
	expectedData := []byte{0, 255, ':', '-', 1}
	ciphertext := encryptTokenIdentityForTest(t, service, 123, 456, 789, core.USER_TOKEN_TYPE_API, expectedData)
	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak", ciphertext})
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)
	assert.Equal(t, make([]byte, 12), actualSignature)

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
	service := initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, service, math.MaxInt64, math.MaxInt64, math.MaxInt64, core.USER_TOKEN_TYPE_API, nil)
	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak", ciphertext})
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)
	assert.Empty(t, actualData)

	expectedValue := &models.TokenRecord{
		Uid:             math.MaxInt64,
		UserTokenId:     math.MaxInt64,
		CreatedUnixTime: math.MaxInt64,
		TokenType:       core.USER_TOKEN_TYPE_API,
	}
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseEncryptedTokenId_WrongKey(t *testing.T) {
	service := initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, service, 123, 456, 789, core.USER_TOKEN_TYPE_API, nil)
	service.CurrentConfig().SecretKey = "different-global-key"

	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak", ciphertext})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_TamperedCiphertext(t *testing.T) {
	service := initializeTokenServiceTest(t)
	ciphertext := encryptTokenIdentityForTest(t, service, 123, 456, 789, core.USER_TOKEN_TYPE_API, nil)
	lastCharacter := "0"

	if strings.HasSuffix(ciphertext, lastCharacter) {
		lastCharacter = "1"
	}

	ciphertext = ciphertext[:len(ciphertext)-1] + lastCharacter
	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak", ciphertext})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_Base62TypeMismatch(t *testing.T) {
	service := initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_API,
		Secret:          "0123456789",
	}
	token, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)

	parts := strings.Split(token, "-")
	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)
	assert.Equal(t, record.TokenType, actualValue.TokenType)
	assert.Empty(t, actualData)

	parts[1] = "tk"
	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), parts)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_Base62ShortIdentity(t *testing.T) {
	service := initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             123,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_API,
		Secret:          "0123456789",
	}
	token, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)
	plaintext, err := utils.DecryptSecretFromBase62(strings.Split(token, "-")[2], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)

	ciphertext, err := utils.EncryptSecretToBase62("", service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)

	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak", ciphertext})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	ciphertext, err = utils.EncryptSecretToBase62(plaintext[:23], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak", ciphertext})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	ciphertext, err = utils.EncryptSecretToBase62(plaintext[:24], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak", ciphertext})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	ciphertext, err = utils.EncryptSecretToBase62(plaintext[:25], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak", ciphertext})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	ciphertext, err = utils.EncryptSecretToBase62(plaintext[:36], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), []string{"ebk", "ak", ciphertext})
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseEncryptedTokenId_Base62InvalidIdentity(t *testing.T) {
	service := initializeTokenServiceTest(t)
	record := &models.TokenRecord{
		Uid:             0,
		UserTokenId:     456,
		CreatedUnixTime: 789,
		TokenType:       core.USER_TOKEN_TYPE_API,
		Secret:          "0123456789",
	}
	token, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)

	actualValue, actualSignature, actualData, err := service.parseEncryptedTokenId(core.NewNullContext(), strings.Split(token, "-"))
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	record.Uid = -1
	token, err = service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), strings.Split(token, "-"))
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	record.Uid = 123
	record.UserTokenId = 0
	token, err = service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), strings.Split(token, "-"))
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	record.UserTokenId = -1
	token, err = service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), strings.Split(token, "-"))
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	record.UserTokenId = 456
	record.CreatedUnixTime = 0
	token, err = service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), strings.Split(token, "-"))
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	record.CreatedUnixTime = -1
	token, err = service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)

	actualValue, actualSignature, actualData, err = service.parseEncryptedTokenId(core.NewNullContext(), strings.Split(token, "-"))
	assert.Nil(t, actualValue)
	assert.Nil(t, actualSignature)
	assert.Nil(t, actualData)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestParseToken_Base62Signature(t *testing.T) {
	service, c, _, _ := initializeScopedTokenServiceTest(t)
	token, claims, record, err := service.createToken(c, &models.User{Uid: 123}, core.USER_TOKEN_TYPE_API, "", "", time.Hour)
	assert.Equal(t, nil, err)

	actualValue, actualRecord, err := service.ParseToken(c, token)
	assert.Equal(t, nil, err)
	assert.Equal(t, claims, actualValue)
	assert.Equal(t, record, actualRecord)

	parts := strings.Split(token, "-")
	_, actualSignature, encryptedData, err := service.parseEncryptedTokenId(c, parts)
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)
	assert.Empty(t, encryptedData)

	plaintext, err := utils.DecryptSecretFromBase62(parts[2], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	data := []byte(plaintext)
	data[25] ^= 1
	parts[2], err = utils.EncryptSecretToBase62(string(data), service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)

	actualValue, actualRecord, err = service.ParseToken(c, strings.Join(parts, "-"))
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualRecord)

	wrongSecretRecord := *record
	wrongSecretRecord.Secret = "wrong-secret"
	token, err = service.generateEncryptedToken(&wrongSecretRecord, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)

	actualValue, actualRecord, err = service.ParseToken(c, token)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualRecord)

	signature := hmac.New(sha1.New, []byte(record.Secret))
	_, _ = signature.Write([]byte(plaintext[:25]))
	data = append([]byte(plaintext[:25]), signature.Sum(nil)...)
	parts[2], err = utils.EncryptSecretToBase62(string(data), service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	actualValue, actualRecord, err = service.ParseToken(c, strings.Join(parts, "-"))
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualRecord)
}

func TestParseToken_Base62SignatureTokenType(t *testing.T) {
	service, c, _, _ := initializeScopedTokenServiceTest(t)
	token, _, record, err := service.createToken(c, &models.User{Uid: 123}, core.USER_TOKEN_TYPE_API, "", "", time.Hour)
	assert.Equal(t, nil, err)

	parts := strings.Split(token, "-")
	plaintext, err := utils.DecryptSecretFromBase62(parts[2], service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)
	data := []byte(plaintext)
	data[24] = byte(core.USER_TOKEN_TYPE_MCP)
	parts[1] = "mcp"
	parts[2], err = utils.EncryptSecretToBase62(string(data), service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)

	updatedRecord := *record
	updatedRecord.TokenType = core.USER_TOKEN_TYPE_MCP
	err = service.createTokenRecord(c, &updatedRecord)
	assert.Equal(t, nil, err)

	actualIdentity, _, _, err := service.parseEncryptedTokenId(c, parts)
	assert.Equal(t, nil, err)
	assert.Equal(t, updatedRecord.TokenType, actualIdentity.TokenType)
	actualValue, actualRecord, err := service.ParseToken(c, strings.Join(parts, "-"))
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualRecord)

	token, err = service.generateEncryptedToken(&updatedRecord, core.NONE_DERIVED_TOKEN, nil)
	assert.Equal(t, nil, err)
	actualValue, actualRecord, err = service.ParseToken(c, token)
	assert.Equal(t, nil, err)
	assert.Equal(t, updatedRecord.TokenType, actualValue.Type)
	assert.Equal(t, &updatedRecord, actualRecord)
}

func TestParseToken_LegacyHexRejected(t *testing.T) {
	service, c, _, _ := initializeScopedTokenServiceTest(t)
	_, _, record, err := service.createToken(c, &models.User{Uid: 123}, core.USER_TOKEN_TYPE_API, "", "", time.Hour)
	assert.Equal(t, nil, err)

	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, record.Uid)
	_ = binary.Write(buf, binary.BigEndian, record.UserTokenId)
	_ = binary.Write(buf, binary.BigEndian, record.CreatedUnixTime)
	_ = binary.Write(buf, binary.BigEndian, record.TokenType)
	ciphertext, err := utils.EncryptSecretToHex(buf.String(), service.CurrentConfig().SecretKey)
	assert.Equal(t, nil, err)

	actualIdentity, actualSignature, encryptedData, err := service.parseEncryptedTokenId(c, []string{"ebk", "ak", ciphertext})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualIdentity)
	assert.Nil(t, actualSignature)
	assert.Nil(t, encryptedData)

	actualValue, actualRecord, err := service.ParseToken(c, "ebk-ak-"+ciphertext)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualRecord)
}

func TestGetUserAgent(t *testing.T) {
	service := &TokenService{}
	webContext := core.WrapWebContext(&gin.Context{Request: httptest.NewRequest("GET", "http://localhost/", nil)}, nil)
	webContext.Request.Header.Set("User-Agent", "test-client")
	expectedValue := "test-client"
	actualValue := service.getUserAgent(webContext)
	assert.Equal(t, expectedValue, actualValue)

	webContext.Request.Header.Set("User-Agent", strings.Repeat("x", models.TokenMaxUserAgentLength+1))
	expectedValue = strings.Repeat("x", models.TokenMaxUserAgentLength)
	actualValue = service.getUserAgent(webContext)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetUserAgent_Empty(t *testing.T) {
	service := &TokenService{}
	actualValue := service.getUserAgent(nil)
	assert.Equal(t, "", actualValue)

	actualValue = service.getUserAgent(&core.WebContext{Context: &gin.Context{}})
	assert.Equal(t, "", actualValue)

	webContext := core.WrapWebContext(&gin.Context{Request: httptest.NewRequest("GET", "http://localhost/", nil)}, nil)
	actualValue = service.getUserAgent(webContext)
	assert.Equal(t, "", actualValue)
}

func TestGetUserTokenId(t *testing.T) {
	service := &TokenService{}
	actualValue := service.getUserTokenId()
	assert.Greater(t, actualValue, int64(0))

	secondValue := service.getUserTokenId()
	assert.Greater(t, secondValue, int64(0))
	assert.NotEqual(t, actualValue, secondValue)
}

func TestGenerateEncryptedToken_UserAvatar(t *testing.T) {
	service := initializeTokenServiceTest(t)
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

	token, err := service.generateEncryptedToken(record, core.USER_AVATAR_DERIVED_NORMAL_TOKEN, additionalDataBytes)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(token, "ebk-tk.a-"))

	actualRecord, actualSignature, encryptedData, err := service.parseEncryptedTokenId(core.NewNullContext(), strings.Split(token, "-"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)
	assert.Equal(t, record.Uid, actualRecord.Uid)
	assert.Equal(t, record.UserTokenId, actualRecord.UserTokenId)
	assert.Equal(t, record.CreatedUnixTime, actualRecord.CreatedUnixTime)
	assert.Equal(t, record.TokenType, actualRecord.TokenType)

	additionalData, err := utils.AESGCMDecrypt(utils.MD5Encode([]byte(record.Secret)), encryptedData)
	assert.Equal(t, nil, err)

	actualValue, err := service.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 1000, []core.NormalTokenScope{
		core.NormalTokenScopeUserAvatar,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseNormalTokenAdditionalData_UserAvatarInvalidScope(t *testing.T) {
	service := &TokenService{}
	additionalData, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: 2000,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	actualValue, err := service.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 1000, nil)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)

	actualValue, err = service.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 1000, []core.NormalTokenScope{core.NormalTokenScopeTransactionPicture})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
}

func TestParseNormalTokenAdditionalData_UserAvatarExpired(t *testing.T) {
	service := &TokenService{}
	additionalData, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: 2000,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	actualValue, err := service.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 2000, []core.NormalTokenScope{core.NormalTokenScopeUserAvatar})
	assert.Equal(t, nil, err)
	assert.Equal(t, &core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: 2000,
	}, actualValue)

	actualValue, err = service.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 2001, []core.NormalTokenScope{core.NormalTokenScopeUserAvatar})
	assert.Equal(t, errs.ErrCurrentTokenExpired, err)
	assert.Nil(t, actualValue)
}

func TestGenerateEncryptedToken_TransactionPicture(t *testing.T) {
	service := initializeTokenServiceTest(t)
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

	token, err := service.generateEncryptedToken(record, core.TRANSACTION_PICTURE_DERIVED_NORMAL_TOKEN, additionalDataBytes)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(token, "ebk-tk.p-"))

	actualRecord, actualSignature, encryptedData, err := service.parseEncryptedTokenId(core.NewNullContext(), strings.Split(token, "-"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)
	assert.Equal(t, record.Uid, actualRecord.Uid)
	assert.Equal(t, record.UserTokenId, actualRecord.UserTokenId)
	assert.Equal(t, record.CreatedUnixTime, actualRecord.CreatedUnixTime)
	assert.Equal(t, record.TokenType, actualRecord.TokenType)

	additionalData, err := utils.AESGCMDecrypt(utils.MD5Encode([]byte(record.Secret)), encryptedData)
	assert.Equal(t, nil, err)

	actualValue, err := service.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 1000, []core.NormalTokenScope{core.NormalTokenScopeTransactionPicture})
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseNormalTokenAdditionalData_TransactionPictureInvalidScope(t *testing.T) {
	service := &TokenService{}
	additionalData, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeTransactionPicture,
		DerivedTokenExpiresAt: 2000,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	actualValue, err := service.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 1000, nil)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)

	actualValue, err = service.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 1000, []core.NormalTokenScope{core.NormalTokenScopeUserAvatar})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
}

func TestParseNormalTokenAdditionalData_TransactionPictureExpired(t *testing.T) {
	service := &TokenService{}
	additionalData, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeTransactionPicture,
		DerivedTokenExpiresAt: 2000,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	actualValue, err := service.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 2000, []core.NormalTokenScope{core.NormalTokenScopeTransactionPicture})
	assert.Equal(t, nil, err)
	assert.Equal(t, &core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeTransactionPicture,
		DerivedTokenExpiresAt: 2000,
	}, actualValue)

	actualValue, err = service.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 2001, []core.NormalTokenScope{core.NormalTokenScopeTransactionPicture})
	assert.Equal(t, errs.ErrCurrentTokenExpired, err)
	assert.Nil(t, actualValue)
}

func TestGenerateEncryptedToken_ScopedAdditionalDataInvalidSecret(t *testing.T) {
	service := initializeTokenServiceTest(t)
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

	token, err := service.generateEncryptedToken(record, core.USER_AVATAR_DERIVED_NORMAL_TOKEN, additionalDataBytes)
	assert.Equal(t, nil, err)

	_, actualSignature, encryptedData, err := service.parseEncryptedTokenId(core.NewNullContext(), strings.Split(token, "-"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualSignature, 12)

	_, err = utils.AESGCMDecrypt(utils.MD5Encode([]byte("other-secret")), encryptedData)
	assert.NotNil(t, err)

	encryptedData[len(encryptedData)-1] ^= 1
	_, err = utils.AESGCMDecrypt(utils.MD5Encode([]byte(record.Secret)), encryptedData)
	assert.NotNil(t, err)
}

func TestParseNormalTokenAdditionalData_All(t *testing.T) {
	service := &TokenService{}
	expectedValue := &core.NormalTokenAdditionalData{
		Scope: core.NormalTokenScopeAll,
	}
	actualValue, err := service.parseNormalTokenAdditionalData(core.NewNullContext(), []byte{1}, 2000, nil)

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseNormalTokenAdditionalData_InvalidScope(t *testing.T) {
	service := &TokenService{}
	additionalData, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeInvalid,
		DerivedTokenExpiresAt: 3000,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	actualValue, err := service.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 2000, []core.NormalTokenScope{core.NormalTokenScopeInvalid})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)

	additionalData, err = (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScope(255),
		DerivedTokenExpiresAt: 3000,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	actualValue, err = service.parseNormalTokenAdditionalData(core.NewNullContext(), additionalData, 2000, []core.NormalTokenScope{core.NormalTokenScope(255)})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
}

func TestDeriveUserAvatarToken_CurrentSession(t *testing.T) {
	service, c, claims, record := initializeScopedTokenServiceTest(t)
	expectedValue := time.Now().Add(time.Hour).Unix()
	token, err := service.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(token, "ebk-tk.a-"))

	actualValue, actualRecord, err := service.ParseToken(c, token, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, nil, err)
	assert.Equal(t, claims.UserTokenId, actualValue.UserTokenId)
	assert.Equal(t, claims.Uid, actualValue.Uid)
	assert.Equal(t, record, actualRecord)
	assert.InDelta(t, expectedValue, actualValue.ExpiresAt, 2)

	_, _, err = service.ParseToken(c, token)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	_, _, err = service.ParseToken(c, token, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	mismatchedToken := strings.Replace(token, "ebk-tk.a-", "ebk-tk.p-", 1)
	_, _, err = service.ParseToken(c, mismatchedToken, core.NormalTokenScopeUserAvatar, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	missingDerivedTokenType := strings.Replace(token, "ebk-tk.a-", "ebk-tk-", 1)
	actualValue, _, err = service.ParseToken(c, missingDerivedTokenType, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)

	normalToken, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, []byte{byte(core.NormalTokenScopeAll)})
	assert.Equal(t, nil, err)
	mismatchedToken = strings.Replace(normalToken, "ebk-tk-", "ebk-tk.a-", 1)
	_, _, err = service.ParseToken(c, mismatchedToken)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestDeriveUserAvatarToken_SessionExpiration(t *testing.T) {
	service, c, claims, record := initializeScopedTokenServiceTest(t)
	claims.ExpiresAt = time.Now().Unix() + 60
	expectedValue := claims.ExpiresAt
	token, err := service.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, nil, err)

	actualValue, _, err := service.ParseToken(c, token, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue.ExpiresAt)
}

func TestParseToken_UserAvatarExpired(t *testing.T) {
	service, c, _, record := initializeScopedTokenServiceTest(t)
	additionalDataBytes, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: time.Now().Unix() - 1,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	token, err := service.generateEncryptedToken(record, core.USER_AVATAR_DERIVED_NORMAL_TOKEN, additionalDataBytes)
	assert.Equal(t, nil, err)

	actualValue, _, err := service.ParseToken(c, token, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, errs.ErrCurrentTokenExpired, err)
	assert.Nil(t, actualValue)
}

func TestDeriveUserAvatarToken_RevokedSession(t *testing.T) {
	service, c, claims, record := initializeScopedTokenServiceTest(t)
	token, err := service.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, nil, err)

	err = service.DeleteTokenByClaims(c, claims)
	assert.Equal(t, nil, err)

	actualValue, _, err := service.ParseToken(c, token, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
}

func TestDeriveTransactionPictureToken_CurrentSession(t *testing.T) {
	service, c, claims, record := initializeScopedTokenServiceTest(t)
	expectedValue := time.Now().Add(30 * time.Minute).Unix()
	token, err := service.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(token, "ebk-tk.p-"))

	actualValue, actualRecord, err := service.ParseToken(c, token, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, nil, err)
	assert.Equal(t, claims.UserTokenId, actualValue.UserTokenId)
	assert.Equal(t, claims.Uid, actualValue.Uid)
	assert.Equal(t, record, actualRecord)
	assert.InDelta(t, expectedValue, actualValue.ExpiresAt, 2)

	_, _, err = service.ParseToken(c, token)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	_, _, err = service.ParseToken(c, token, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	mismatchedToken := strings.Replace(token, "ebk-tk.p-", "ebk-tk.a-", 1)
	_, _, err = service.ParseToken(c, mismatchedToken, core.NormalTokenScopeUserAvatar, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	missingDerivedTokenType := strings.Replace(token, "ebk-tk.p-", "ebk-tk-", 1)
	actualValue, _, err = service.ParseToken(c, missingDerivedTokenType, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)

	normalToken, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, []byte{byte(core.NormalTokenScopeAll)})
	assert.Equal(t, nil, err)
	mismatchedToken = strings.Replace(normalToken, "ebk-tk-", "ebk-tk.p-", 1)
	_, _, err = service.ParseToken(c, mismatchedToken)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestDeriveTransactionPictureToken_SessionExpiration(t *testing.T) {
	service, c, claims, record := initializeScopedTokenServiceTest(t)
	claims.ExpiresAt = time.Now().Unix() + 60
	expectedValue := claims.ExpiresAt
	token, err := service.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, nil, err)

	actualValue, _, err := service.ParseToken(c, token, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue.ExpiresAt)
}

func TestParseToken_TransactionPictureExpired(t *testing.T) {
	service, c, _, record := initializeScopedTokenServiceTest(t)
	additionalDataBytes, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeTransactionPicture,
		DerivedTokenExpiresAt: time.Now().Unix() - 1,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	token, err := service.generateEncryptedToken(record, core.TRANSACTION_PICTURE_DERIVED_NORMAL_TOKEN, additionalDataBytes)
	assert.Equal(t, nil, err)

	actualValue, _, err := service.ParseToken(c, token, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, errs.ErrCurrentTokenExpired, err)
	assert.Nil(t, actualValue)
}

func TestDeriveTransactionPictureToken_RevokedSession(t *testing.T) {
	service, c, claims, record := initializeScopedTokenServiceTest(t)
	token, err := service.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, nil, err)

	err = service.DeleteTokenByClaims(c, claims)
	assert.Equal(t, nil, err)

	actualValue, _, err := service.ParseToken(c, token, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
}

func TestDeriveUserAvatarToken_InvalidClaims(t *testing.T) {
	service := initializeTokenServiceTest(t)
	expectedValue := ""
	actualValue, err := service.DeriveUserAvatarToken(core.NewNullContext(), nil, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	actualValue, err = service.DeriveUserAvatarToken(core.NewNullContext(), &core.UserTokenClaims{
		Type: core.USER_TOKEN_TYPE_API,
	}, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	actualValue, err = service.DeriveUserAvatarToken(core.NewNullContext(), &core.UserTokenClaims{
		Type:        core.USER_TOKEN_TYPE_NORMAL,
		UserTokenId: "invalid",
	}, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDeriveUserAvatarToken_ExpiredSession(t *testing.T) {
	service, c, claims, record := initializeScopedTokenServiceTest(t)
	claims.ExpiresAt = time.Now().Unix() - 1
	expectedValue := ""
	actualValue, err := service.DeriveUserAvatarToken(c, claims, record)

	assert.Equal(t, errs.ErrTokenExpired, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDeriveUserAvatarToken_InvalidRecordType(t *testing.T) {
	service, c, _, _ := initializeScopedTokenServiceTest(t)
	_, claims, record, err := service.createToken(c, &models.User{
		Uid: 123,
	}, core.USER_TOKEN_TYPE_API, "", "", time.Hour)
	assert.Equal(t, nil, err)

	claims.Type = core.USER_TOKEN_TYPE_NORMAL
	expectedValue := ""
	actualValue, err := service.DeriveUserAvatarToken(c, claims, record)

	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDeriveTransactionPictureToken_InvalidClaims(t *testing.T) {
	service := initializeTokenServiceTest(t)
	expectedValue := ""
	actualValue, err := service.DeriveTransactionPictureToken(core.NewNullContext(), nil, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	actualValue, err = service.DeriveTransactionPictureToken(core.NewNullContext(), &core.UserTokenClaims{
		Type: core.USER_TOKEN_TYPE_API,
	}, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	actualValue, err = service.DeriveTransactionPictureToken(core.NewNullContext(), &core.UserTokenClaims{
		Type:        core.USER_TOKEN_TYPE_NORMAL,
		UserTokenId: "invalid",
	}, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDeriveTransactionPictureToken_ExpiredSession(t *testing.T) {
	service, c, claims, record := initializeScopedTokenServiceTest(t)
	claims.ExpiresAt = time.Now().Unix() - 1
	expectedValue := ""
	actualValue, err := service.DeriveTransactionPictureToken(c, claims, record)

	assert.Equal(t, errs.ErrTokenExpired, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDeriveTransactionPictureToken_InvalidRecordType(t *testing.T) {
	service, c, _, _ := initializeScopedTokenServiceTest(t)
	_, claims, record, err := service.createToken(c, &models.User{
		Uid: 123,
	}, core.USER_TOKEN_TYPE_API, "", "", time.Hour)
	assert.Equal(t, nil, err)

	claims.Type = core.USER_TOKEN_TYPE_NORMAL
	expectedValue := ""
	actualValue, err := service.DeriveTransactionPictureToken(c, claims, record)

	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseNormalTokenAdditionalData_InvalidData(t *testing.T) {
	service := &TokenService{}
	actualValue, err := service.parseNormalTokenAdditionalData(core.NewNullContext(), nil, 2000, nil)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)

	actualValue, err = service.parseNormalTokenAdditionalData(core.NewNullContext(), []byte{1, 0}, 2000, nil)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
}

func TestDeriveUserAvatarToken_InvalidRecord(t *testing.T) {
	service := initializeTokenServiceTest(t)
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
	actualValue, err := service.DeriveUserAvatarToken(c, claims, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	record.Uid++
	actualValue, err = service.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	record.Uid = claims.Uid
	record.UserTokenId++
	actualValue, err = service.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	record.UserTokenId = 456
	record.CreatedUnixTime++
	actualValue, err = service.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDeriveUserAvatarToken_RecordExpiration(t *testing.T) {
	service, c, claims, record := initializeScopedTokenServiceTest(t)
	record.ExpiredUnixTime = time.Now().Unix() + 60
	expectedValue := record.ExpiredUnixTime
	token, err := service.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, nil, err)

	actualValue, _, err := service.ParseToken(c, token, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue.ExpiresAt)

	record.ExpiredUnixTime = time.Now().Unix() - 1
	token, err = service.DeriveUserAvatarToken(c, claims, record)
	assert.Equal(t, errs.ErrTokenExpired, err)
	assert.Equal(t, "", token)
}

func TestDeriveUserCustomIconToken_CurrentSession(t *testing.T) {
	service, c, claims, record := initializeScopedTokenServiceTest(t)
	expectedValue := time.Now().Add(30 * time.Minute).Unix()
	token, expiresAt, err := service.DeriveUserCustomIconToken(c, claims, record)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(token, "ebk-tk.i-"))

	actualValue, actualRecord, err := service.ParseToken(c, token, core.NormalTokenScopeUserCustomIcon)
	assert.Equal(t, nil, err)
	assert.Equal(t, claims.UserTokenId, actualValue.UserTokenId)
	assert.Equal(t, claims.Uid, actualValue.Uid)
	assert.Equal(t, record, actualRecord)
	assert.InDelta(t, expectedValue, actualValue.ExpiresAt, 2)
	assert.Equal(t, actualValue.ExpiresAt, expiresAt)

	_, _, err = service.ParseToken(c, token)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	_, _, err = service.ParseToken(c, token, core.NormalTokenScopeUserAvatar)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	_, _, err = service.ParseToken(c, token, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	mismatchedToken := strings.Replace(token, "ebk-tk.i-", "ebk-tk.a-", 1)
	_, _, err = service.ParseToken(c, mismatchedToken, core.NormalTokenScopeUserAvatar, core.NormalTokenScopeUserCustomIcon)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	missingDerivedTokenType := strings.Replace(token, "ebk-tk.i-", "ebk-tk-", 1)
	actualValue, _, err = service.ParseToken(c, missingDerivedTokenType, core.NormalTokenScopeUserCustomIcon)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)

	normalToken, err := service.generateEncryptedToken(record, core.NONE_DERIVED_TOKEN, []byte{byte(core.NormalTokenScopeAll)})
	assert.Equal(t, nil, err)
	mismatchedToken = strings.Replace(normalToken, "ebk-tk-", "ebk-tk.i-", 1)
	_, _, err = service.ParseToken(c, mismatchedToken)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
}

func TestDeriveUserCustomIconToken_SessionExpiration(t *testing.T) {
	service, c, claims, record := initializeScopedTokenServiceTest(t)
	claims.ExpiresAt = time.Now().Unix() + 60
	expectedValue := claims.ExpiresAt
	token, expiresAt, err := service.DeriveUserCustomIconToken(c, claims, record)
	assert.Equal(t, nil, err)

	actualValue, _, err := service.ParseToken(c, token, core.NormalTokenScopeUserCustomIcon)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue.ExpiresAt)
	assert.Equal(t, expectedValue, expiresAt)
}

func TestDeriveUserCustomIconToken_RecordExpiration(t *testing.T) {
	service, c, claims, record := initializeScopedTokenServiceTest(t)
	record.ExpiredUnixTime = time.Now().Unix() + 60
	expectedValue := record.ExpiredUnixTime
	token, expiresAt, err := service.DeriveUserCustomIconToken(c, claims, record)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, expiresAt)

	actualValue, _, err := service.ParseToken(c, token, core.NormalTokenScopeUserCustomIcon)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue.ExpiresAt)

	record.ExpiredUnixTime = time.Now().Unix() - 1
	token, expiresAt, err = service.DeriveUserCustomIconToken(c, claims, record)
	assert.Equal(t, errs.ErrTokenExpired, err)
	assert.Equal(t, "", token)
	assert.Equal(t, int64(0), expiresAt)
}

func TestParseToken_UserCustomIconExpired(t *testing.T) {
	service, c, _, record := initializeScopedTokenServiceTest(t)
	additionalDataBytes, err := (&core.NormalTokenAdditionalData{
		Scope:                 core.NormalTokenScopeUserCustomIcon,
		DerivedTokenExpiresAt: time.Now().Unix() - 1,
	}).MarshalBinary()
	assert.Equal(t, nil, err)

	token, err := service.generateEncryptedToken(record, core.USER_CUSTOM_ICON_DERIVED_NORMAL_TOKEN, additionalDataBytes)
	assert.Equal(t, nil, err)

	actualValue, _, err := service.ParseToken(c, token, core.NormalTokenScopeUserCustomIcon)
	assert.Equal(t, errs.ErrCurrentTokenExpired, err)
	assert.Nil(t, actualValue)
}

func TestDeriveUserCustomIconToken_RevokedSession(t *testing.T) {
	service, c, claims, record := initializeScopedTokenServiceTest(t)
	token, _, err := service.DeriveUserCustomIconToken(c, claims, record)
	assert.Equal(t, nil, err)

	err = service.DeleteTokenByClaims(c, claims)
	assert.Equal(t, nil, err)

	actualValue, _, err := service.ParseToken(c, token, core.NormalTokenScopeUserCustomIcon)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Nil(t, actualValue)
}

func TestDeriveTransactionPictureToken_InvalidRecord(t *testing.T) {
	service := initializeTokenServiceTest(t)
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
	actualValue, err := service.DeriveTransactionPictureToken(c, claims, nil)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	record.Uid++
	actualValue, err = service.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	record.Uid = claims.Uid
	record.UserTokenId++
	actualValue, err = service.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	record.UserTokenId = 456
	record.CreatedUnixTime++
	actualValue, err = service.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestDeriveTransactionPictureToken_RecordExpiration(t *testing.T) {
	service, c, claims, record := initializeScopedTokenServiceTest(t)
	record.ExpiredUnixTime = time.Now().Unix() + 60
	expectedValue := record.ExpiredUnixTime
	token, err := service.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, nil, err)

	actualValue, _, err := service.ParseToken(c, token, core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue.ExpiresAt)

	record.ExpiredUnixTime = time.Now().Unix() - 1
	token, err = service.DeriveTransactionPictureToken(c, claims, record)
	assert.Equal(t, errs.ErrTokenExpired, err)
	assert.Equal(t, "", token)
}

func TestParseToken_ReturnsTokenRecord(t *testing.T) {
	service, c, _, _ := initializeScopedTokenServiceTest(t)
	token, claims, expectedValue, err := service.createToken(c, &models.User{
		Uid: 123,
	}, core.USER_TOKEN_TYPE_NORMAL, "test-agent", "test-context", time.Hour)
	assert.Equal(t, nil, err)

	actualClaims, actualValue, err := service.ParseToken(c, token)
	assert.Equal(t, nil, err)
	assert.Equal(t, claims.UserTokenId, actualClaims.UserTokenId)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, "test-context", actualValue.GetContext())
}
