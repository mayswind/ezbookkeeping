package middlewares

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
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

// TokenAuthorizationByHeader verifies whether current request is valid by token in header
func TokenAuthorizationByHeader(config *settings.Config) core.MiddlewareHandlerFunc {
	return tokenAuthorization(config, TOKEN_SOURCE_TYPE_HEADER, true)
}

// TokenAuthorizationByQueryString verifies whether current request is valid by token in query string
func TokenAuthorizationByQueryString(config *settings.Config) core.MiddlewareHandlerFunc {
	return tokenAuthorization(config, TOKEN_SOURCE_TYPE_ARGUMENT, false)
}

// TokenAuthorizationByCookie verifies whether current request is valid by token in cookie
func TokenAuthorizationByCookie(config *settings.Config) core.MiddlewareHandlerFunc {
	return tokenAuthorization(config, TOKEN_SOURCE_TYPE_COOKIE, false)
}

// TokenTwoFactorAuthorization verifies whether current request is valid by 2fa passcode
func TokenTwoFactorAuthorization(config *settings.Config) core.MiddlewareHandlerFunc {
	return func(c *core.WebContext) {
		claims, tokenContext, err := getTokenClaims(c, TOKEN_SOURCE_TYPE_HEADER)

		if err != nil {
			utils.PrintJsonErrorResult(c, err)
			return
		}

		if claims.Type != core.USER_TOKEN_TYPE_REQUIRE_2FA {
			log.Warnf(c, "[authorization.TokenTwoFactorAuthorization] user \"uid:%d\" token is not need two-factor authorization", claims.Uid)
			utils.PrintJsonErrorResult(c, errs.ErrCurrentTokenNotRequire2FA)
			return
		}

		c.SetTokenClaims(claims)
		c.SetTokenContext(tokenContext)
		c.Next()
	}
}

// TokenEmailVerifyAuthorization verifies whether current request is email verification
func TokenEmailVerifyAuthorization(config *settings.Config) core.MiddlewareHandlerFunc {
	return func(c *core.WebContext) {
		claims, tokenContext, err := getTokenClaims(c, TOKEN_SOURCE_TYPE_ARGUMENT)

		if err != nil {
			utils.PrintJsonErrorResult(c, errs.ErrEmailVerifyTokenIsInvalidOrExpired)
			return
		}

		if claims.Type != core.USER_TOKEN_TYPE_EMAIL_VERIFY {
			log.Warnf(c, "[authorization.TokenEmailVerifyAuthorization] user \"uid:%d\" token is not for email verification", claims.Uid)
			utils.PrintJsonErrorResult(c, errs.ErrCurrentInvalidToken)
			return
		}

		c.SetTokenClaims(claims)
		c.SetTokenContext(tokenContext)
		c.Next()
	}
}

// TokenResetPasswordAuthorization verifies whether current request is password reset
func TokenResetPasswordAuthorization(config *settings.Config) core.MiddlewareHandlerFunc {
	return func(c *core.WebContext) {
		claims, tokenContext, err := getTokenClaims(c, TOKEN_SOURCE_TYPE_ARGUMENT)

		if err != nil {
			utils.PrintJsonErrorResult(c, errs.ErrPasswordResetTokenIsInvalidOrExpired)
			return
		}

		if claims.Type != core.USER_TOKEN_TYPE_PASSWORD_RESET {
			log.Warnf(c, "[authorization.TokenResetPasswordAuthorization] user \"uid:%d\" token is not for password request", claims.Uid)
			utils.PrintJsonErrorResult(c, errs.ErrCurrentInvalidToken)
			return
		}

		c.SetTokenClaims(claims)
		c.SetTokenContext(tokenContext)
		c.Next()
	}
}

// TokenMCPAuthorization verifies whether current request is valid by mcp token in header
func TokenMCPAuthorization(config *settings.Config) core.MiddlewareHandlerFunc {
	return func(c *core.WebContext) {
		claims, tokenContext, err := getTokenClaims(c, TOKEN_SOURCE_TYPE_HEADER)

		if err != nil {
			utils.PrintJsonErrorResult(c, err)
			return
		}

		if claims.Type != core.USER_TOKEN_TYPE_MCP {
			log.Warnf(c, "[authorization.tokenAuthorization] user \"uid:%d\" token type (%d) is not mcp token", claims.Uid, claims.Type)
			utils.PrintJsonErrorResult(c, errs.ErrCurrentInvalidTokenType)
			return
		}

		c.SetTokenClaims(claims)
		c.SetTokenContext(tokenContext)
		c.Next()
	}
}

// TokenOAuth2CallbackAuthorization verifies whether current request is OAuth 2.0 callback
func TokenOAuth2CallbackAuthorization(config *settings.Config) core.MiddlewareHandlerFunc {
	return func(c *core.WebContext) {
		claims, tokenContext, err := getTokenClaims(c, TOKEN_SOURCE_TYPE_HEADER)

		if err != nil {
			utils.PrintJsonErrorResult(c, errs.ErrTokenExpired)
			return
		}

		if claims.Type != core.USER_TOKEN_TYPE_OAUTH2_CALLBACK && claims.Type != core.USER_TOKEN_TYPE_OAUTH2_CALLBACK_REQUIRE_VERIFY {
			log.Warnf(c, "[authorization.TokenOAuth2CallbackAuthorization] user \"uid:%d\" token is not for oauth 2.0 callback request", claims.Uid)
			utils.PrintJsonErrorResult(c, errs.ErrCurrentInvalidToken)
			return
		}

		c.SetTokenClaims(claims)
		c.SetTokenContext(tokenContext)
		c.Next()
	}
}

func tokenAuthorization(config *settings.Config, source TokenSourceType, allowedAPIToken bool) core.MiddlewareHandlerFunc {
	return func(c *core.WebContext) {
		claims, tokenContext, err := getTokenClaims(c, source)

		if err != nil {
			utils.PrintJsonErrorResult(c, err)
			return
		}

		if claims.Type == core.USER_TOKEN_TYPE_REQUIRE_2FA {
			log.Warnf(c, "[authorization.tokenAuthorization] user \"uid:%d\" token requires 2fa", claims.Uid)
			utils.PrintJsonErrorResult(c, errs.ErrCurrentTokenRequire2FA)
			return
		}

		if claims.Type != core.USER_TOKEN_TYPE_NORMAL && claims.Type != core.USER_TOKEN_TYPE_API {
			log.Warnf(c, "[authorization.tokenAuthorization] user \"uid:%d\" token type (%d) is invalid", claims.Uid, claims.Type)
			utils.PrintJsonErrorResult(c, errs.ErrCurrentInvalidTokenType)
			return
		}

		if claims.Type == core.USER_TOKEN_TYPE_API && !allowedAPIToken {
			log.Warnf(c, "[authorization.tokenAuthorization] user \"uid:%d\" token type (%d) is not allowed", claims.Uid, claims.Type)
			utils.PrintJsonErrorResult(c, errs.ErrCurrentInvalidTokenType)
			return
		}

		if claims.Type == core.USER_TOKEN_TYPE_API && !config.EnableAPIToken {
			log.Warnf(c, "[authorization.tokenAuthorization] api token is not enabled")
			utils.PrintJsonErrorResult(c, errs.ErrAPITokenNotEnabled)
			return
		}

		c.SetTokenClaims(claims)
		c.SetTokenContext(tokenContext)
		c.Next()
	}
}

func getTokenClaims(c *core.WebContext, source TokenSourceType) (*core.UserTokenClaims, string, *errs.Error) {
	claims, tokenContext, err := parseToken(c, source)

	if err != nil {
		log.Warnf(c, "[authorization.getTokenClaims] failed to parse token, because %s", err.Error())
		return nil, "", errs.Or(err, errs.ErrUnauthorizedAccess)
	}

	if claims.Uid <= 0 {
		log.Warnf(c, "[authorization.getTokenClaims] user id in token is invalid")
		return nil, "", errs.ErrCurrentInvalidToken
	}

	return claims, tokenContext, nil
}

func parseToken(c *core.WebContext, source TokenSourceType) (*core.UserTokenClaims, string, error) {
	tokenString := ""

	if source == TOKEN_SOURCE_TYPE_ARGUMENT {
		tokenString = c.GetTokenStringFromQueryString()
	} else if source == TOKEN_SOURCE_TYPE_COOKIE {
		tokenString = c.GetTokenStringFromCookie()
	} else { // if source == TOKEN_SOURCE_TYPE_HEADER
		tokenString = c.GetTokenStringFromHeader()
	}

	if tokenString == "" {
		return nil, "", errs.ErrTokenIsEmpty
	}

	return services.Tokens.ParseToken(c, tokenString)
}
