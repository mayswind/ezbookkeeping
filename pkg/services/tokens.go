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
func (s *TokenService) GetAllTokensByUid(uid int64) ([]*models.TokenRecord, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var tokenRecords []*models.TokenRecord
	err := s.TokenDB(uid).Cols("uid", "user_token_id", "token_type", "user_agent", "created_unix_time", "expired_unix_time").Where("uid=?", uid).Find(&tokenRecords)

	return tokenRecords, err
}

// GetAllUnexpiredNormalTokensByUid returns all available token models of given user
func (s *TokenService) GetAllUnexpiredNormalTokensByUid(uid int64) ([]*models.TokenRecord, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	var tokenRecords []*models.TokenRecord
	err := s.TokenDB(uid).Cols("uid", "user_token_id", "token_type", "user_agent", "created_unix_time", "expired_unix_time").Where("uid=? AND token_type=? AND expired_unix_time>?", uid, core.USER_TOKEN_TYPE_NORMAL, now).Find(&tokenRecords)

	return tokenRecords, err
}

// ParseToken returns the token model according to request data
func (s *TokenService) ParseToken(c *core.Context) (*jwt.Token, *core.UserTokenClaims, error) {
	claims := &core.UserTokenClaims{}

	token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			now := time.Now().Unix()
			userTokenId, err := utils.StringToInt64(claims.UserTokenId)

			if err != nil {
				log.WarnfWithRequestId(c, "[tokens.ParseToken] token \"utid:%s\" in token of user \"uid:%d\" is invalid, because %s", claims.UserTokenId, claims.Uid, err.Error())
				return nil, errs.ErrInvalidUserTokenId
			}

			tokenRecord, err := s.getTokenRecord(claims.Uid, userTokenId, claims.IssuedAt)

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

// CreateToken generates a new normal token and saves to database
func (s *TokenService) CreateToken(user *models.User, ctx *core.Context) (string, *core.UserTokenClaims, error) {
	return s.createToken(user, core.USER_TOKEN_TYPE_NORMAL, s.getUserAgent(ctx), s.CurrentConfig().TokenExpiredTimeDuration)
}

// CreateRequire2FAToken generates a new token requiring user to verify 2fa passcode and saves to database
func (s *TokenService) CreateRequire2FAToken(user *models.User, ctx *core.Context) (string, *core.UserTokenClaims, error) {
	return s.createToken(user, core.USER_TOKEN_TYPE_REQUIRE_2FA, s.getUserAgent(ctx), s.CurrentConfig().TemporaryTokenExpiredTimeDuration)
}

// DeleteToken deletes given token from database
func (s *TokenService) DeleteToken(tokenRecord *models.TokenRecord) error {
	if tokenRecord.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if tokenRecord.UserTokenId <= 0 {
		return errs.ErrInvalidUserTokenId
	}

	return s.TokenDB(tokenRecord.Uid).DoTransaction(func(sess *xorm.Session) error {
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
func (s *TokenService) DeleteTokens(uid int64, tokenRecords []*models.TokenRecord) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.TokenDB(uid).DoTransaction(func(sess *xorm.Session) error {
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
func (s *TokenService) DeleteTokenByClaims(claims *core.UserTokenClaims) error {
	userTokenId, err := utils.StringToInt64(claims.UserTokenId)

	if err != nil {
		return errs.ErrInvalidUserTokenId
	}

	return s.DeleteToken(&models.TokenRecord{
		Uid:             claims.Uid,
		UserTokenId:     userTokenId,
		CreatedUnixTime: claims.IssuedAt,
	})
}

// DeleteTokensBeforeTime deletes tokens that is created before specific tim
func (s *TokenService) DeleteTokensBeforeTime(uid int64, expireTime int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.TokenDB(uid).DoTransaction(func(sess *xorm.Session) error {
		_, err := sess.Where("uid=? AND created_unix_time<?", uid, expireTime).Delete(&models.TokenRecord{})
		return err
	})
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

func (s *TokenService) createToken(user *models.User, tokenType core.TokenType, userAgent string, expiryDate time.Duration) (string, *core.UserTokenClaims, error) {
	var err error
	now := time.Now()

	tokenRecord := &models.TokenRecord{
		Uid:             user.Uid,
		UserTokenId:     s.getUserTokenId(),
		TokenType:       tokenType,
		UserAgent:       userAgent,
		CreatedUnixTime: now.Unix(),
		ExpiredUnixTime: now.Add(expiryDate).Unix(),
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

	err = s.createTokenRecord(tokenRecord)

	if err != nil {
		return "", nil, err
	}

	return tokenString, claims, err
}

func (s *TokenService) getTokenRecord(uid int64, userTokenId int64, createUnixTime int64) (*models.TokenRecord, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if userTokenId <= 0 {
		return nil, errs.ErrInvalidUserTokenId
	}

	tokenRecord := &models.TokenRecord{}
	has, err := s.TokenDB(uid).Where("uid=? AND user_token_id=? AND created_unix_time=?", uid, userTokenId, createUnixTime).Limit(1).Get(tokenRecord)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errs.ErrTokenRecordNotFound
	}

	return tokenRecord, nil
}

func (s *TokenService) createTokenRecord(tokenRecord *models.TokenRecord) error {
	if tokenRecord.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if tokenRecord.UserTokenId <= 0 {
		return errs.ErrInvalidUserTokenId
	}

	return s.TokenDB(tokenRecord.Uid).DoTransaction(func(sess *xorm.Session) error {
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
