package middlewares

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// TokenSourceType represents token source
type TokenSourceType byte

// Token source types
const (
	TOKEN_SOURCE_TYPE_HEADER   TokenSourceType = 1
	TOKEN_SOURCE_TYPE_ARGUMENT TokenSourceType = 2
	TOKEN_SOURCE_TYPE_COOKIE   TokenSourceType = 3
)

// JWTAuthorization verifies whether current request is valid by jwt token in header
func JWTAuthorization(c *core.WebContext) {
	jwtAuthorization(c, TOKEN_SOURCE_TYPE_HEADER)
}

// JWTAuthorizationByQueryString verifies whether current request is valid by jwt token in query string
func JWTAuthorizationByQueryString(c *core.WebContext) {
	jwtAuthorization(c, TOKEN_SOURCE_TYPE_ARGUMENT)
}

// JWTAuthorizationByCookie verifies whether current request is valid by jwt token in cookie
func JWTAuthorizationByCookie(c *core.WebContext) {
	jwtAuthorization(c, TOKEN_SOURCE_TYPE_COOKIE)
}

// JWTTwoFactorAuthorization verifies whether current request is valid by 2fa passcode
func JWTTwoFactorAuthorization(c *core.WebContext) {
	claims, err := getTokenClaims(c, TOKEN_SOURCE_TYPE_HEADER)

	if err != nil {
		utils.PrintJsonErrorResult(c, err)
		return
	}

	if claims.Type != core.USER_TOKEN_TYPE_REQUIRE_2FA {
		log.Warnf(c, "[authorization.JWTTwoFactorAuthorization] user \"uid:%d\" token is not need two-factor authorization", claims.Uid)
		utils.PrintJsonErrorResult(c, errs.ErrCurrentTokenNotRequire2FA)
		return
	}

	c.SetTokenClaims(claims)
	c.Next()
}

// JWTEmailVerifyAuthorization verifies whether current request is email verification
func JWTEmailVerifyAuthorization(c *core.WebContext) {
	claims, err := getTokenClaims(c, TOKEN_SOURCE_TYPE_ARGUMENT)

	if err != nil {
		utils.PrintJsonErrorResult(c, errs.ErrEmailVerifyTokenIsInvalidOrExpired)
		return
	}

	if claims.Type != core.USER_TOKEN_TYPE_EMAIL_VERIFY {
		log.Warnf(c, "[authorization.JWTEmailVerifyAuthorization] user \"uid:%d\" token is not for email verification", claims.Uid)
		utils.PrintJsonErrorResult(c, errs.ErrCurrentInvalidToken)
		return
	}

	c.SetTokenClaims(claims)
	c.Next()
}

// JWTResetPasswordAuthorization verifies whether current request is password reset
func JWTResetPasswordAuthorization(c *core.WebContext) {
	claims, err := getTokenClaims(c, TOKEN_SOURCE_TYPE_ARGUMENT)

	if err != nil {
		utils.PrintJsonErrorResult(c, errs.ErrPasswordResetTokenIsInvalidOrExpired)
		return
	}

	if claims.Type != core.USER_TOKEN_TYPE_PASSWORD_RESET {
		log.Warnf(c, "[authorization.JWTResetPasswordAuthorization] user \"uid:%d\" token is not for password request", claims.Uid)
		utils.PrintJsonErrorResult(c, errs.ErrCurrentInvalidToken)
		return
	}

	c.SetTokenClaims(claims)
	c.Next()
}

// JWTMCPAuthorization verifies whether current request is valid by jwt mcp token in header
func JWTMCPAuthorization(c *core.WebContext) {
	claims, err := getTokenClaims(c, TOKEN_SOURCE_TYPE_HEADER)

	if err != nil {
		utils.PrintJsonErrorResult(c, err)
		return
	}

	if claims.Type != core.USER_TOKEN_TYPE_MCP {
		log.Warnf(c, "[authorization.jwtAuthorization] user \"uid:%d\" token type (%d) is not mcp token", claims.Uid, claims.Type)
		utils.PrintJsonErrorResult(c, errs.ErrCurrentInvalidTokenType)
		return
	}

	c.SetTokenClaims(claims)
	c.Next()
}

func jwtAuthorization(c *core.WebContext, source TokenSourceType) {
	claims, err := getTokenClaims(c, source)

	if err != nil {
		utils.PrintJsonErrorResult(c, err)
		return
	}

	if claims.Type == core.USER_TOKEN_TYPE_REQUIRE_2FA {
		log.Warnf(c, "[authorization.jwtAuthorization] user \"uid:%d\" token requires 2fa", claims.Uid)
		utils.PrintJsonErrorResult(c, errs.ErrCurrentTokenRequire2FA)
		return
	}

	if claims.Type != core.USER_TOKEN_TYPE_NORMAL {
		log.Warnf(c, "[authorization.jwtAuthorization] user \"uid:%d\" token type (%d) is invalid", claims.Uid, claims.Type)
		utils.PrintJsonErrorResult(c, errs.ErrCurrentInvalidTokenType)
		return
	}

	c.SetTokenClaims(claims)
	c.Next()
}

func getTokenClaims(c *core.WebContext, source TokenSourceType) (*core.UserTokenClaims, *errs.Error) {
	token, claims, err := parseToken(c, source)

	if err != nil {
		log.Warnf(c, "[authorization.getTokenClaims] failed to parse token, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrUnauthorizedAccess)
	}

	if !token.Valid {
		log.Warnf(c, "[authorization.getTokenClaims] token is invalid")
		return nil, errs.ErrCurrentInvalidToken
	}

	if claims.Uid <= 0 {
		log.Warnf(c, "[authorization.getTokenClaims] user id in token is invalid")
		return nil, errs.ErrCurrentInvalidToken
	}

	return claims, nil
}

func parseToken(c *core.WebContext, source TokenSourceType) (*jwt.Token, *core.UserTokenClaims, error) {
	tokenString := ""

	if source == TOKEN_SOURCE_TYPE_ARGUMENT {
		tokenString = c.GetTokenStringFromQueryString()
	} else if source == TOKEN_SOURCE_TYPE_COOKIE {
		tokenString = c.GetTokenStringFromCookie()
	} else { // if source == TOKEN_SOURCE_TYPE_HEADER
		tokenString = c.GetTokenStringFromHeader()
	}

	if tokenString == "" {
		return nil, nil, errs.ErrTokenIsEmpty
	}

	return services.Tokens.ParseToken(c, tokenString)
}
