package services

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const tokenPrefix = "ebk"
const tokenIdentityMinimumSize = 8 + 8 + 8 + 1
const tokenMaxExpiredAtUnixTime = int64(253402300799) // 9999-12-31 23:59:59 UTC

// TokenService represents user token service
type TokenService struct {
	ServiceUsingDB
	ServiceUsingConfig
}

// Initialize a user token service singleton instance
var (
	Tokens = &TokenService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingConfig: ServiceUsingConfig{
			container: settings.Container,
		},
	}
)

// GetAllTokensByUid returns all token models of given user
func (s *TokenService) GetAllTokensByUid(c core.Context, uid int64) ([]*models.TokenRecord, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var tokenRecords []*models.TokenRecord
	err := s.TokenDB(uid).NewSession(c).Cols("uid", "user_token_id", "token_type", "user_agent", "created_unix_time", "expired_unix_time").Where("uid=?", uid).Find(&tokenRecords)

	return tokenRecords, err
}

// GetAllUnexpiredNormalAndMCPTokensByUid returns all available token models of given user
func (s *TokenService) GetAllUnexpiredNormalAndMCPTokensByUid(c core.Context, uid int64) ([]*models.TokenRecord, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	var tokenRecords []*models.TokenRecord
	err := s.TokenDB(uid).NewSession(c).Cols("uid", "user_token_id", "token_type", "user_agent", "created_unix_time", "expired_unix_time", "last_seen_unix_time").Where("uid=? AND (token_type=? OR token_type=? OR token_type=?) AND expired_unix_time>?", uid, core.USER_TOKEN_TYPE_NORMAL, core.USER_TOKEN_TYPE_MCP, core.USER_TOKEN_TYPE_API, now).Find(&tokenRecords)

	return tokenRecords, err
}

// ParseToken returns the token model according to token content
func (s *TokenService) ParseToken(c core.Context, token string) (*core.UserTokenClaims, string, error) {
	parts := strings.Split(token, "-")

	if len(parts) == 3 && parts[0] == tokenPrefix {
		return s.parseEncryptedToken(c, parts)
	}

	return s.parseLegacyJWTToken(c, token)
}

// CreateToken generates a new normal token and saves to database
func (s *TokenService) CreateToken(c *core.WebContext, user *models.User) (string, *core.UserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, core.USER_TOKEN_TYPE_NORMAL, s.getUserAgent(c), "", s.CurrentConfig().TokenExpiredTimeDuration)
	return token, claims, err
}

// CreateRequire2FAToken generates a new token requiring user to verify 2fa passcode and saves to database
func (s *TokenService) CreateRequire2FAToken(c *core.WebContext, user *models.User) (string, *core.UserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, core.USER_TOKEN_TYPE_REQUIRE_2FA, s.getUserAgent(c), "", s.CurrentConfig().TemporaryTokenExpiredTimeDuration)
	return token, claims, err
}

// CreateEmailVerifyToken generates a new email verify token and saves to database
func (s *TokenService) CreateEmailVerifyToken(c *core.WebContext, user *models.User) (string, *core.UserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, core.USER_TOKEN_TYPE_EMAIL_VERIFY, s.getUserAgent(c), "", s.CurrentConfig().EmailVerifyTokenExpiredTimeDuration)
	return token, claims, err
}

// CreateEmailVerifyTokenWithoutUserAgent generates a new email verify token and saves to database
func (s *TokenService) CreateEmailVerifyTokenWithoutUserAgent(c core.Context, user *models.User) (string, *core.UserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, core.USER_TOKEN_TYPE_EMAIL_VERIFY, "", "", s.CurrentConfig().EmailVerifyTokenExpiredTimeDuration)
	return token, claims, err
}

// CreatePasswordResetToken generates a new password reset token and saves to database
func (s *TokenService) CreatePasswordResetToken(c *core.WebContext, user *models.User) (string, *core.UserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, core.USER_TOKEN_TYPE_PASSWORD_RESET, s.getUserAgent(c), "", s.CurrentConfig().PasswordResetTokenExpiredTimeDuration)
	return token, claims, err
}

// CreatePasswordResetTokenWithoutUserAgent generates a new password reset token and saves to database
func (s *TokenService) CreatePasswordResetTokenWithoutUserAgent(c core.Context, user *models.User) (string, *core.UserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, core.USER_TOKEN_TYPE_PASSWORD_RESET, "", "", s.CurrentConfig().PasswordResetTokenExpiredTimeDuration)
	return token, claims, err
}

// CreateAPIToken generates a new API token and saves to database
func (s *TokenService) CreateAPIToken(c *core.WebContext, user *models.User, expiresInSeconds int64) (string, *core.UserTokenClaims, error) {
	var tokenExpiredTimeDuration time.Duration

	if expiresInSeconds > 0 {
		tokenExpiredTimeDuration = time.Duration(expiresInSeconds) * time.Second
	} else {
		tokenExpiredTimeDuration = time.Unix(tokenMaxExpiredAtUnixTime, 0).Sub(time.Now())
	}

	token, claims, _, err := s.createToken(c, user, core.USER_TOKEN_TYPE_API, s.getUserAgent(c), "", tokenExpiredTimeDuration)
	return token, claims, err
}

// CreateAPITokenViaCli generates a new API token and saves to database
func (s *TokenService) CreateAPITokenViaCli(c *core.CliContext, user *models.User, expiresInSeconds int64) (string, *models.TokenRecord, error) {
	var tokenExpiredTimeDuration time.Duration

	if expiresInSeconds > 0 {
		tokenExpiredTimeDuration = time.Duration(expiresInSeconds) * time.Second
	} else {
		tokenExpiredTimeDuration = time.Unix(tokenMaxExpiredAtUnixTime, 0).Sub(time.Now())
	}

	token, _, tokenRecord, err := s.createToken(c, user, core.USER_TOKEN_TYPE_API, core.TokenUserAgentCreatedViaCli, "", tokenExpiredTimeDuration)
	return token, tokenRecord, err
}

// CreateMCPToken generates a new MCP token and saves to database
func (s *TokenService) CreateMCPToken(c *core.WebContext, user *models.User, expiresInSeconds int64) (string, *core.UserTokenClaims, error) {
	var tokenExpiredTimeDuration time.Duration

	if expiresInSeconds > 0 {
		tokenExpiredTimeDuration = time.Duration(expiresInSeconds) * time.Second
	} else {
		tokenExpiredTimeDuration = time.Unix(tokenMaxExpiredAtUnixTime, 0).Sub(time.Now())
	}

	token, claims, _, err := s.createToken(c, user, core.USER_TOKEN_TYPE_MCP, s.getUserAgent(c), "", tokenExpiredTimeDuration)
	return token, claims, err
}

// CreateMCPTokenViaCli generates a new MCP token and saves to database
func (s *TokenService) CreateMCPTokenViaCli(c *core.CliContext, user *models.User, expiresInSeconds int64) (string, *models.TokenRecord, error) {
	var tokenExpiredTimeDuration time.Duration

	if expiresInSeconds > 0 {
		tokenExpiredTimeDuration = time.Duration(expiresInSeconds) * time.Second
	} else {
		tokenExpiredTimeDuration = time.Unix(tokenMaxExpiredAtUnixTime, 0).Sub(time.Now())
	}

	token, _, tokenRecord, err := s.createToken(c, user, core.USER_TOKEN_TYPE_MCP, core.TokenUserAgentCreatedViaCli, "", tokenExpiredTimeDuration)
	return token, tokenRecord, err
}

// CreateOAuth2CallbackRequireVerifyToken generates a new OAuth 2.0 callback token requiring user to verify and saves to database
func (s *TokenService) CreateOAuth2CallbackRequireVerifyToken(c *core.WebContext, user *models.User, context string) (string, *core.UserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, core.USER_TOKEN_TYPE_OAUTH2_CALLBACK_REQUIRE_VERIFY, s.getUserAgent(c), context, s.CurrentConfig().TemporaryTokenExpiredTimeDuration)
	return token, claims, err
}

// CreateOAuth2CallbackToken generates a new OAuth 2.0 callback token and saves to database
func (s *TokenService) CreateOAuth2CallbackToken(c *core.WebContext, user *models.User, context string) (string, *core.UserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, core.USER_TOKEN_TYPE_OAUTH2_CALLBACK, s.getUserAgent(c), context, s.CurrentConfig().TemporaryTokenExpiredTimeDuration)
	return token, claims, err
}

// UpdateTokenLastSeen updates the last seen time of specified token
func (s *TokenService) UpdateTokenLastSeen(c core.Context, tokenRecord *models.TokenRecord) error {
	if tokenRecord.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if tokenRecord.UserTokenId <= 0 {
		return errs.ErrInvalidUserTokenId
	}

	tokenRecord.LastSeenUnixTime = time.Now().Unix()

	return s.TokenDB(tokenRecord.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.Cols("last_seen_unix_time").Where("uid=? AND user_token_id=? AND created_unix_time=?", tokenRecord.Uid, tokenRecord.UserTokenId, tokenRecord.CreatedUnixTime).Update(tokenRecord)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrTokenRecordNotFound
		}

		return nil
	})
}

// DeleteToken deletes given token from database
func (s *TokenService) DeleteToken(c core.Context, tokenRecord *models.TokenRecord) error {
	if tokenRecord.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if tokenRecord.UserTokenId <= 0 {
		return errs.ErrInvalidUserTokenId
	}

	return s.TokenDB(tokenRecord.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		deletedRows, err := sess.Where("uid=? AND user_token_id=? AND created_unix_time=?", tokenRecord.Uid, tokenRecord.UserTokenId, tokenRecord.CreatedUnixTime).Delete(&models.TokenRecord{})

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrTokenRecordNotFound
		}

		return nil
	})
}

// DeleteTokens deletes given tokens from database
func (s *TokenService) DeleteTokens(c core.Context, uid int64, tokenRecords []*models.TokenRecord) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.TokenDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		for i := 0; i < len(tokenRecords); i++ {
			tokenRecord := tokenRecords[i]
			deletedRows, err := sess.Where("uid=? AND user_token_id=? AND created_unix_time=?", uid, tokenRecord.UserTokenId, tokenRecord.CreatedUnixTime).Delete(&models.TokenRecord{})

			if err != nil {
				return err
			} else if deletedRows < 1 {
				return errs.ErrTokenRecordNotFound
			}
		}

		return nil
	})
}

// DeleteTokenByClaims deletes given token from database
func (s *TokenService) DeleteTokenByClaims(c core.Context, claims *core.UserTokenClaims) error {
	userTokenId, err := utils.StringToInt64(claims.UserTokenId)

	if err != nil {
		return errs.ErrInvalidUserTokenId
	}

	return s.DeleteToken(c, &models.TokenRecord{
		Uid:             claims.Uid,
		UserTokenId:     userTokenId,
		CreatedUnixTime: claims.IssuedAt,
	})
}

// DeleteTokensBeforeTime deletes tokens that is created before specific time
func (s *TokenService) DeleteTokensBeforeTime(c core.Context, uid int64, createTime int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.TokenDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Where("uid=? AND created_unix_time<?", uid, createTime).Delete(&models.TokenRecord{})
		return err
	})
}

// DeleteTokensByType deletes specified type tokens
func (s *TokenService) DeleteTokensByType(c core.Context, uid int64, tokenType core.TokenType) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.TokenDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Where("uid=? AND token_type=?", uid, tokenType).Delete(&models.TokenRecord{})
		return err
	})
}

// DeleteAllTokens deletes all tokens of given user
func (s *TokenService) DeleteAllTokens(c core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.TokenDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Where("uid=?", uid).Delete(&models.TokenRecord{})
		return err
	})
}

// DeleteAllExpiredTokens deletes all expired tokens
func (s *TokenService) DeleteAllExpiredTokens(c core.Context) error {
	var errors []error
	totalCount := int64(0)

	for i := 0; i < s.TokenDBCount(); i++ {
		err := s.TokenDBByIndex(i).DoTransaction(c, func(sess *xorm.Session) error {
			count, err := sess.Where("expired_unix_time<=?", time.Now().Unix()).Delete(&models.TokenRecord{})
			totalCount += count
			return err
		})

		if err != nil {
			errors = append(errors, err)
		}
	}

	if totalCount > 0 {
		log.Infof(c, "[tokens.DeleteAllExpiredTokens] %d expired tokens have been deleted", totalCount)
	} else if len(errors) == 0 {
		log.Infof(c, "[tokens.DeleteAllExpiredTokens] no expired tokens have been deleted")
	}

	return errs.NewMultiErrorOrNil(errors...)
}

// ExistsValidTokenByType returns whether the given token type exists
func (s *TokenService) ExistsValidTokenByType(c core.Context, uid int64, tokenType core.TokenType) (bool, error) {
	if uid <= 0 {
		return false, errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	return s.TokenDB(uid).NewSession(c).Cols("uid", "user_token_id", "expired_unix_time").Where("uid=? AND token_type=? AND expired_unix_time>?", uid, tokenType, now).Exist(&models.TokenRecord{})
}

// ParseFromTokenId returns token model according to token id
func (s *TokenService) ParseFromTokenId(tokenId string) (*models.TokenRecord, error) {
	pairs := strings.Split(tokenId, ":")

	if len(pairs) != 3 {
		return nil, errs.ErrInvalidTokenId
	}

	uid, err := utils.StringToInt64(pairs[0])

	if err != nil {
		return nil, errs.ErrInvalidTokenId
	}

	createdUnixTime, err := utils.StringToInt64(pairs[1])

	if err != nil {
		return nil, errs.ErrInvalidTokenId
	}

	userTokenId, err := utils.StringToInt64(pairs[2])

	if err != nil {
		return nil, errs.ErrInvalidTokenId
	}

	tokenRecord := &models.TokenRecord{
		Uid:             uid,
		UserTokenId:     userTokenId,
		CreatedUnixTime: createdUnixTime,
	}

	return tokenRecord, nil
}

// GenerateTokenId generates token id according to token model
func (s *TokenService) GenerateTokenId(tokenRecord *models.TokenRecord) string {
	return fmt.Sprintf("%d:%d:%d", tokenRecord.Uid, tokenRecord.CreatedUnixTime, tokenRecord.UserTokenId)
}

func (s *TokenService) parseEncryptedToken(c core.Context, parts []string) (*core.UserTokenClaims, string, error) {
	identity, encryptedAdditionalData, err := s.parseEncryptedTokenId(c, parts)

	if err != nil {
		return nil, "", err
	}

	record, err := s.getTokenRecord(c, identity.Uid, identity.UserTokenId, identity.CreatedUnixTime)

	if err != nil {
		log.Warnf(c, "[tokens.parseEncryptedToken] cannot get token \"utid:%d\" of user \"uid:%d\", because %s", identity.UserTokenId, identity.Uid, err.Error())
		return nil, "", err
	}

	if record.TokenType != identity.TokenType {
		log.Warnf(c, "[tokens.parseEncryptedToken] token \"utid:%d\" of user \"uid:%d\" record type is not match, because token type is %d but record type is %d", record.UserTokenId, record.Uid, identity.TokenType, record.TokenType)
		return nil, "", errs.ErrCurrentInvalidToken
	}

	now := time.Now().Unix()

	if record.ExpiredUnixTime < now {
		log.Warnf(c, "[tokens.parseEncryptedToken] token \"utid:%d\" of user \"uid:%d\" record is expired", record.UserTokenId, record.Uid)
		return nil, "", errs.ErrCurrentTokenExpired
	}

	if record.CreatedUnixTime > now {
		log.Warnf(c, "[tokens.parseEncryptedToken] token \"utid:%d\" of user \"uid:%d\" record is invalid, because creation time is later than now", record.UserTokenId, record.Uid)
		return nil, "", errs.ErrCurrentInvalidToken
	}

	var additionalData []byte

	if len(encryptedAdditionalData) > 0 {
		additionalData, err = utils.AESGCMDecrypt(utils.MD5Encode([]byte(record.Secret)), encryptedAdditionalData)

		if err != nil {
			log.Warnf(c, "[tokens.parseEncryptedToken] cannot decrypt additional data of token \"utid:%d\" of user \"uid:%d\", because %s", record.UserTokenId, record.Uid, err.Error())
			return nil, "", errs.ErrCurrentInvalidToken
		}
	}

	if record.TokenType == core.USER_TOKEN_TYPE_NORMAL {
		if len(additionalData) < 1 {
			log.Warnf(c, "[tokens.parseEncryptedToken] token \"utid:%d\" of user \"uid:%d\" record is invalid, because additional data is missing", record.UserTokenId, record.Uid)
			return nil, "", errs.ErrCurrentInvalidToken
		}

		var tokenAdditionalData core.NormalTokenAdditionalData
		err = tokenAdditionalData.UnmarshalBinary(additionalData)

		if err != nil {
			log.Warnf(c, "[tokens.parseEncryptedToken] cannot parse additional data of token \"utid:%d\" of user \"uid:%d\", because %s", record.UserTokenId, record.Uid, err.Error())
			return nil, "", errs.ErrCurrentInvalidToken
		}

		if tokenAdditionalData.Scope != core.NormalTokenScopeAll {
			log.Warnf(c, "[tokens.parseEncryptedToken] token \"utid:%d\" of user \"uid:%d\" record is invalid, because scope \"%d\" is invalid", record.UserTokenId, record.Uid, tokenAdditionalData.Scope)
			return nil, "", errs.ErrCurrentInvalidToken
		}
	}

	return &core.UserTokenClaims{
		UserTokenId: utils.Int64ToString(record.UserTokenId),
		Uid:         record.Uid,
		Type:        record.TokenType,
		IssuedAt:    record.CreatedUnixTime,
		ExpiresAt:   record.ExpiredUnixTime,
	}, record.Context, nil
}

func (s *TokenService) parseEncryptedTokenId(c core.Context, parts []string) (*models.TokenRecord, []byte, error) {
	if len(parts) != 3 || parts[0] != tokenPrefix {
		log.Warnf(c, "[tokens.parseEncryptedTokenId] token is invalid, because token format is invalid")
		return nil, nil, errs.ErrCurrentInvalidToken
	}

	tokenType, exists := core.TextualTokenTypesReverse[parts[1]]

	if !exists {
		log.Warnf(c, "[tokens.parseEncryptedTokenId] token is invalid, because token type \"%s\" is not supported", parts[1])
		return nil, nil, errs.ErrCurrentInvalidTokenType
	}

	decryptedFullTokenId, err := utils.DecryptSecretFromHex(parts[2], s.CurrentConfig().SecretKey)

	if err != nil {
		log.Warnf(c, "[tokens.parseEncryptedTokenId] token is invalid, because decrypting token id failed, %s", err.Error())
		return nil, nil, errs.ErrCurrentInvalidToken
	}

	fullTokenId := []byte(decryptedFullTokenId)

	if len(fullTokenId) < tokenIdentityMinimumSize {
		log.Warnf(c, "[tokens.parseEncryptedTokenId] token is invalid, because token id is too short")
		return nil, nil, errs.ErrCurrentInvalidToken
	}

	buf := bytes.NewBuffer(fullTokenId)

	var uid int64
	var userTokenId int64
	var createdUnixTime int64
	var actualTokenType core.TokenType

	_ = binary.Read(buf, binary.BigEndian, &uid)
	_ = binary.Read(buf, binary.BigEndian, &userTokenId)
	_ = binary.Read(buf, binary.BigEndian, &createdUnixTime)
	_ = binary.Read(buf, binary.BigEndian, &actualTokenType)

	if uid <= 0 || userTokenId <= 0 || createdUnixTime <= 0 {
		log.Warnf(c, "[tokens.parseEncryptedTokenId] token is invalid, because uid, user token id or creation time is invalid")
		return nil, nil, errs.ErrCurrentInvalidToken
	}

	if actualTokenType != tokenType {
		log.Warnf(c, "[tokens.parseEncryptedTokenId] token is invalid, because token type is not match, expected %d but got %d", tokenType, actualTokenType)
		return nil, nil, errs.ErrCurrentInvalidToken
	}

	encryptedAdditionalData := buf.Bytes()

	return &models.TokenRecord{
		Uid:             uid,
		UserTokenId:     userTokenId,
		CreatedUnixTime: createdUnixTime,
		TokenType:       tokenType,
	}, encryptedAdditionalData, nil
}

func (s *TokenService) parseLegacyJWTToken(c core.Context, tokenString string) (*core.UserTokenClaims, string, error) {
	claims := &core.UserTokenClaims{}
	tokenContext := ""

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (any, error) {
			now := time.Now().Unix()
			userTokenId, err := utils.StringToInt64(claims.UserTokenId)

			if err != nil {
				log.Warnf(c, "[tokens.parseToken] token \"utid:%s\" in token of user \"uid:%d\" is invalid, because %s", claims.UserTokenId, claims.Uid, err.Error())
				return nil, errs.ErrInvalidUserTokenId
			}

			tokenRecord, err := s.getTokenRecord(c, claims.Uid, userTokenId, claims.IssuedAt)

			if err != nil {
				log.Warnf(c, "[tokens.parseToken] token \"utid:%s\" of user \"uid:%d\" record not found, because %s", claims.UserTokenId, claims.Uid, err.Error())
				return nil, errs.ErrTokenRecordNotFound
			}

			if tokenRecord.ExpiredUnixTime < now {
				log.Warnf(c, "[tokens.parseToken] token \"utid:%s\" of user \"uid:%d\" record is expired", claims.UserTokenId, claims.Uid)
				return nil, errs.ErrTokenExpired
			}

			tokenContext = tokenRecord.Context
			return []byte(tokenRecord.Secret), nil
		},
		jwt.WithIssuedAt(),
	)

	if err != nil {
		if errors.Is(err, request.ErrNoTokenInRequest) {
			return nil, "", errs.ErrTokenIsEmpty
		}

		if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenUnverifiable) || errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			log.Warnf(c, "[tokens.parseToken] token is invalid, because %s", err.Error())
			return nil, "", errs.ErrCurrentInvalidToken
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, "", errs.ErrCurrentTokenExpired
		}

		if errors.Is(err, jwt.ErrTokenUsedBeforeIssued) {
			log.Warnf(c, "[tokens.parseToken] token is invalid, because issue time is later than now")
			return nil, "", errs.ErrCurrentInvalidToken
		}

		return nil, "", err
	}

	if !token.Valid {
		return nil, "", errs.ErrCurrentInvalidToken
	}

	return claims, tokenContext, nil
}

func (s *TokenService) createToken(c core.Context, user *models.User, tokenType core.TokenType, userAgent string, context string, expiryDate time.Duration) (string, *core.UserTokenClaims, *models.TokenRecord, error) {
	var err error
	now := time.Now()

	tokenRecord := &models.TokenRecord{
		Uid:              user.Uid,
		UserTokenId:      s.getUserTokenId(),
		TokenType:        tokenType,
		UserAgent:        userAgent,
		Context:          context,
		CreatedUnixTime:  now.Unix(),
		ExpiredUnixTime:  now.Add(expiryDate).Unix(),
		LastSeenUnixTime: now.Unix(),
	}

	if tokenRecord.Secret, err = utils.GetRandomString(10); err != nil {
		return "", nil, nil, err
	}

	claims := &core.UserTokenClaims{
		UserTokenId: utils.Int64ToString(tokenRecord.UserTokenId),
		Uid:         tokenRecord.Uid,
		Username:    user.Username,
		Type:        tokenRecord.TokenType,
		IssuedAt:    tokenRecord.CreatedUnixTime,
		ExpiresAt:   tokenRecord.ExpiredUnixTime,
	}

	tokenString, err := s.generateEncryptedToken(tokenRecord)

	if err != nil {
		return "", nil, nil, err
	}

	err = s.createTokenRecord(c, tokenRecord)

	if err != nil {
		return "", nil, nil, err
	}

	return tokenString, claims, tokenRecord, err
}

func (s *TokenService) generateEncryptedToken(record *models.TokenRecord) (string, error) {
	typeName, exists := core.TextualTokenTypes[record.TokenType]

	if !exists {
		return "", errs.ErrCurrentInvalidTokenType
	}

	var encryptedAdditionalData []byte

	if record.TokenType == core.USER_TOKEN_TYPE_NORMAL {
		additionalData := &core.NormalTokenAdditionalData{
			Scope: core.NormalTokenScopeAll,
		}
		context, err := additionalData.MarshalBinary()

		if err != nil {
			return "", err
		}

		encryptedAdditionalData, err = utils.AESGCMEncrypt(utils.MD5Encode([]byte(record.Secret)), context)

		if err != nil {
			return "", err
		}
	}

	buf := &bytes.Buffer{}
	buf.Grow(tokenIdentityMinimumSize + len(encryptedAdditionalData))
	_ = binary.Write(buf, binary.BigEndian, record.Uid)
	_ = binary.Write(buf, binary.BigEndian, record.UserTokenId)
	_ = binary.Write(buf, binary.BigEndian, record.CreatedUnixTime)
	_ = binary.Write(buf, binary.BigEndian, record.TokenType)

	if len(encryptedAdditionalData) > 0 {
		_, _ = buf.Write(encryptedAdditionalData)
	}

	encryptedFullTokenId, err := utils.EncryptSecretToHex(buf.String(), s.CurrentConfig().SecretKey)

	if err != nil {
		return "", err
	}

	return tokenPrefix + "-" + typeName + "-" + encryptedFullTokenId, nil
}

func (s *TokenService) getTokenRecord(c core.Context, uid int64, userTokenId int64, createUnixTime int64) (*models.TokenRecord, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if userTokenId <= 0 {
		return nil, errs.ErrInvalidUserTokenId
	}

	tokenRecord := &models.TokenRecord{}
	has, err := s.TokenDB(uid).NewSession(c).Where("uid=? AND user_token_id=? AND created_unix_time=?", uid, userTokenId, createUnixTime).Limit(1).Get(tokenRecord)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errs.ErrTokenRecordNotFound
	}

	return tokenRecord, nil
}

func (s *TokenService) createTokenRecord(c core.Context, tokenRecord *models.TokenRecord) error {
	if tokenRecord.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if tokenRecord.UserTokenId <= 0 {
		return errs.ErrInvalidUserTokenId
	}

	return s.TokenDB(tokenRecord.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(tokenRecord)
		return err
	})
}

func (s *TokenService) getUserTokenId() int64 {
	nanoSeconds := time.Now().Nanosecond()
	randomNumber, _ := utils.GetRandomInteger(math.MaxInt32)
	userTokenId := (int64(nanoSeconds) << 32) | int64(randomNumber)

	return userTokenId
}

func (s *TokenService) getUserAgent(ctx *core.WebContext) string {
	userAgent := ""

	if ctx != nil && ctx.Request != nil {
		userAgent = ctx.Request.UserAgent()
	}

	if len(userAgent) > models.TokenMaxUserAgentLength {
		userAgent = utils.SubString(userAgent, 0, models.TokenMaxUserAgentLength)
	}

	return userAgent
}
