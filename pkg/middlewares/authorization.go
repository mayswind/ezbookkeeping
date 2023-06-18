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

const tokenQueryStringParam = "token"

// JWTAuthorization verifies whether current request is valid by jwt token in header
func JWTAuthorization(c *core.Context) {
	jwtAuthorization(c, TOKEN_SOURCE_TYPE_HEADER)
}

// JWTAuthorizationByQueryString verifies whether current request is valid by jwt token in query string
func JWTAuthorizationByQueryString(c *core.Context) {
	jwtAuthorization(c, TOKEN_SOURCE_TYPE_ARGUMENT)
}

// JWTAuthorizationByCookie verifies whether current request is valid by jwt token in cookie
func JWTAuthorizationByCookie(c *core.Context) {
	jwtAuthorization(c, TOKEN_SOURCE_TYPE_COOKIE)
}

// JWTTwoFactorAuthorization verifies whether current request is valid by 2fa passcode
func JWTTwoFactorAuthorization(c *core.Context) {
	claims, err := getTokenClaims(c, TOKEN_SOURCE_TYPE_HEADER)

	if err != nil {
		utils.PrintJsonErrorResult(c, err)
		return
	}

	if claims.Type != core.USER_TOKEN_TYPE_REQUIRE_2FA {
		log.WarnfWithRequestId(c, "[authorization.JWTTwoFactorAuthorization] user \"uid:%d\" token is not need two factor authorization", claims.Uid)
		utils.PrintJsonErrorResult(c, errs.ErrCurrentTokenNotRequire2FA)
		return
	}

	c.SetTokenClaims(claims)
	c.Next()
}

func jwtAuthorization(c *core.Context, source TokenSourceType) {
	claims, err := getTokenClaims(c, source)

	if err != nil {
		utils.PrintJsonErrorResult(c, err)
		return
	}

	if claims.Type == core.USER_TOKEN_TYPE_REQUIRE_2FA {
		log.WarnfWithRequestId(c, "[authorization.jwtAuthorization] user \"uid:%d\" token requires 2fa", claims.Uid)
		utils.PrintJsonErrorResult(c, errs.ErrCurrentTokenRequire2FA)
		return
	}

	if claims.Type != core.USER_TOKEN_TYPE_NORMAL {
		log.WarnfWithRequestId(c, "[authorization.jwtAuthorization] user \"uid:%d\" token type is invalid", claims.Uid)
		utils.PrintJsonErrorResult(c, errs.ErrCurrentInvalidTokenType)
		return
	}

	c.SetTokenClaims(claims)
	c.Next()
}

func getTokenClaims(c *core.Context, source TokenSourceType) (*core.UserTokenClaims, *errs.Error) {
	token, claims, err := parseToken(c, source)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorization.getTokenClaims] failed to parse token, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrUnauthorizedAccess)
	}

	if !token.Valid {
		log.WarnfWithRequestId(c, "[authorization.getTokenClaims] token is invalid")
		return nil, errs.ErrCurrentInvalidToken
	}

	if claims.Uid <= 0 {
		log.WarnfWithRequestId(c, "[authorization.getTokenClaims] user id in token is invalid")
		return nil, errs.ErrCurrentInvalidToken
	}

	return claims, nil
}

func parseToken(c *core.Context, source TokenSourceType) (*jwt.Token, *core.UserTokenClaims, error) {
	if source == TOKEN_SOURCE_TYPE_ARGUMENT {
		return services.Tokens.ParseTokenByArgument(c, tokenQueryStringParam)
	} else if source == TOKEN_SOURCE_TYPE_COOKIE {
		return services.Tokens.ParseTokenByCookie(c, tokenCookieParam)
	}

	return services.Tokens.ParseTokenByHeader(c)
}
