package services

import (
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
func (s *TokenService) GetAllTokensByUid(c *core.Context, uid int64) ([]*models.TokenRecord, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var tokenRecords []*models.TokenRecord
	err := s.TokenDB(uid).NewSession(c).Cols("uid", "user_token_id", "token_type", "user_agent", "created_unix_time", "expired_unix_time").Where("uid=?", uid).Find(&tokenRecords)

	return tokenRecords, err
}

// GetAllUnexpiredNormalTokensByUid returns all available token models of given user
func (s *TokenService) GetAllUnexpiredNormalTokensByUid(c *core.Context, uid int64) ([]*models.TokenRecord, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	var tokenRecords []*models.TokenRecord
	err := s.TokenDB(uid).NewSession(c).Cols("uid", "user_token_id", "token_type", "user_agent", "created_unix_time", "expired_unix_time", "last_seen_unix_time").Where("uid=? AND token_type=? AND expired_unix_time>?", uid, core.USER_TOKEN_TYPE_NORMAL, now).Find(&tokenRecords)

	return tokenRecords, err
}

// ParseTokenByHeader returns the token model according to request data
func (s *TokenService) ParseTokenByHeader(c *core.Context) (*jwt.Token, *core.UserTokenClaims, error) {
	return s.parseToken(c, request.BearerExtractor{})
}

// ParseTokenByArgument returns the token model according to request data
func (s *TokenService) ParseTokenByArgument(c *core.Context, tokenParameterName string) (*jwt.Token, *core.UserTokenClaims, error) {
	return s.parseToken(c, request.ArgumentExtractor{tokenParameterName})
}

// ParseTokenByCookie returns the token model according to request data
func (s *TokenService) ParseTokenByCookie(c *core.Context, tokenCookieName string) (*jwt.Token, *core.UserTokenClaims, error) {
	return s.parseToken(c, utils.CookieExtractor{tokenCookieName})
}

// CreateToken generates a new normal token and saves to database
func (s *TokenService) CreateToken(c *core.Context, user *models.User) (string, *core.UserTokenClaims, error) {
	return s.createToken(c, user, core.USER_TOKEN_TYPE_NORMAL, s.getUserAgent(c), s.CurrentConfig().TokenExpiredTimeDuration)
}

// CreateRequire2FAToken generates a new token requiring user to verify 2fa passcode and saves to database
func (s *TokenService) CreateRequire2FAToken(c *core.Context, user *models.User) (string, *core.UserTokenClaims, error) {
	return s.createToken(c, user, core.USER_TOKEN_TYPE_REQUIRE_2FA, s.getUserAgent(c), s.CurrentConfig().TemporaryTokenExpiredTimeDuration)
}

// CreateEmailVerifyToken generates a new email verify token and saves to database
func (s *TokenService) CreateEmailVerifyToken(c *core.Context, user *models.User) (string, *core.UserTokenClaims, error) {
	return s.createToken(c, user, core.USER_TOKEN_TYPE_EMAIL_VERIFY, s.getUserAgent(c), s.CurrentConfig().EmailVerifyTokenExpiredTimeDuration)
}

// CreatePasswordResetToken generates a new password reset token and saves to database
func (s *TokenService) CreatePasswordResetToken(c *core.Context, user *models.User) (string, *core.UserTokenClaims, error) {
	return s.createToken(c, user, core.USER_TOKEN_TYPE_PASSWORD_RESET, s.getUserAgent(c), s.CurrentConfig().PasswordResetTokenExpiredTimeDuration)
}

// UpdateTokenLastSeen updates the last seen time of specified token
func (s *TokenService) UpdateTokenLastSeen(c *core.Context, tokenRecord *models.TokenRecord) error {
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
func (s *TokenService) DeleteToken(c *core.Context, tokenRecord *models.TokenRecord) error {
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
func (s *TokenService) DeleteTokens(c *core.Context, uid int64, tokenRecords []*models.TokenRecord) error {
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
func (s *TokenService) DeleteTokenByClaims(c *core.Context, claims *core.UserTokenClaims) error {
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
func (s *TokenService) DeleteTokensBeforeTime(c *core.Context, uid int64, createTime int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.TokenDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Where("uid=? AND created_unix_time<?", uid, createTime).Delete(&models.TokenRecord{})
		return err
	})
}

// DeleteTokensByType deletes specified type tokens
func (s *TokenService) DeleteTokensByType(c *core.Context, uid int64, tokenType core.TokenType) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.TokenDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Where("uid=? AND token_type=?", uid, tokenType).Delete(&models.TokenRecord{})
		return err
	})
}

// DeleteAllExpiredTokens deletes all expired tokens
func (s *TokenService) DeleteAllExpiredTokens(c *core.Context) error {
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
		log.Infof("[tokens.DeleteAllExpiredTokens] %d expired tokens have been deleted", totalCount)
	} else if len(errors) == 0 {
		log.Infof("[tokens.DeleteAllExpiredTokens] no expired tokens have been deleted")
	}

	return errs.NewMultiErrorOrNil(errors...)
}

// ExistsValidTokenByType returns whether the given token type exists
func (s *TokenService) ExistsValidTokenByType(c *core.Context, uid int64, tokenType core.TokenType) (bool, error) {
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

func (s *TokenService) parseToken(c *core.Context, extractor request.Extractor) (*jwt.Token, *core.UserTokenClaims, error) {
	claims := &core.UserTokenClaims{}

	token, err := request.ParseFromRequest(c.Request, extractor,
		func(token *jwt.Token) (any, error) {
			now := time.Now().Unix()
			userTokenId, err := utils.StringToInt64(claims.UserTokenId)

			if err != nil {
				log.WarnfWithRequestId(c, "[tokens.ParseToken] token \"utid:%s\" in token of user \"uid:%d\" is invalid, because %s", claims.UserTokenId, claims.Uid, err.Error())
				return nil, errs.ErrInvalidUserTokenId
			}

			tokenRecord, err := s.getTokenRecord(c, claims.Uid, userTokenId, claims.IssuedAt)

			if err != nil {
				log.WarnfWithRequestId(c, "[tokens.ParseToken] token \"utid:%s\" of user \"uid:%d\" record not found, because %s", claims.UserTokenId, claims.Uid, err.Error())
				return nil, errs.ErrTokenRecordNotFound
			}

			if tokenRecord.ExpiredUnixTime < now {
				log.WarnfWithRequestId(c, "[tokens.ParseToken] token \"utid:%s\" of user \"uid:%d\" record is expired", claims.UserTokenId, claims.Uid)
				return nil, errs.ErrTokenExpired
			}

			return []byte(tokenRecord.Secret), nil
		},
		request.WithClaims(claims),
		request.WithParser(jwt.NewParser(jwt.WithIssuedAt())),
	)

	if err != nil {
		if err == request.ErrNoTokenInRequest {
			return nil, nil, errs.ErrTokenIsEmpty
		}

		if err == jwt.ErrTokenMalformed || err == jwt.ErrTokenUnverifiable || err == jwt.ErrTokenSignatureInvalid {
			log.WarnfWithRequestId(c, "[tokens.ParseToken] token is invalid, because %s", err.Error())
			return nil, nil, errs.ErrCurrentInvalidToken
		}

		if err == jwt.ErrTokenExpired {
			return nil, nil, errs.ErrCurrentTokenExpired
		}

		if err == jwt.ErrTokenUsedBeforeIssued {
			log.WarnfWithRequestId(c, "[tokens.ParseToken] token is invalid, because issue time is later than now")
			return nil, nil, errs.ErrCurrentInvalidToken
		}

		return nil, nil, err
	}

	return token, claims, err
}

func (s *TokenService) createToken(c *core.Context, user *models.User, tokenType core.TokenType, userAgent string, expiryDate time.Duration) (string, *core.UserTokenClaims, error) {
	var err error
	now := time.Now()

	tokenRecord := &models.TokenRecord{
		Uid:              user.Uid,
		UserTokenId:      s.getUserTokenId(),
		TokenType:        tokenType,
		UserAgent:        userAgent,
		CreatedUnixTime:  now.Unix(),
		ExpiredUnixTime:  now.Add(expiryDate).Unix(),
		LastSeenUnixTime: now.Unix(),
	}

	if tokenRecord.Secret, err = utils.GetRandomString(10); err != nil {
		return "", nil, err
	}

	claims := &core.UserTokenClaims{
		UserTokenId: utils.Int64ToString(tokenRecord.UserTokenId),
		Uid:         tokenRecord.Uid,
		Username:    user.Username,
		Type:        tokenRecord.TokenType,
		IssuedAt:    tokenRecord.CreatedUnixTime,
		ExpiresAt:   tokenRecord.ExpiredUnixTime,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString([]byte(tokenRecord.Secret))

	if err != nil {
		return "", nil, err
	}

	err = s.createTokenRecord(c, tokenRecord)

	if err != nil {
		return "", nil, err
	}

	return tokenString, claims, err
}

func (s *TokenService) getTokenRecord(c *core.Context, uid int64, userTokenId int64, createUnixTime int64) (*models.TokenRecord, error) {
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

func (s *TokenService) createTokenRecord(c *core.Context, tokenRecord *models.TokenRecord) error {
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

func (s *TokenService) getUserAgent(ctx *core.Context) string {
	userAgent := ""

	if ctx != nil && ctx.Request != nil {
		userAgent = ctx.Request.UserAgent()
	}

	if len(userAgent) > models.TokenMaxUserAgentLength {
		userAgent = utils.SubString(userAgent, 0, models.TokenMaxUserAgentLength)
	}

	return userAgent
}
