package middlewares

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const tokenQueryStringParam = "token"

// JWTAuthorization verifies whether current request is valid by jwt token
func JWTAuthorization(c *core.Context) {
	claims, err := getTokenClaims(c)

	if err != nil {
		utils.PrintJsonErrorResult(c, err)
		return
	}

	if claims.Type == core.USER_TOKEN_TYPE_REQUIRE_2FA {
		log.WarnfWithRequestId(c, "[authorization.JWTAuthorization] user \"uid:%d\" token requires 2fa", claims.Uid)
		utils.PrintJsonErrorResult(c, errs.ErrCurrentTokenRequire2FA)
		return
	}

	if claims.Type != core.USER_TOKEN_TYPE_NORMAL {
		log.WarnfWithRequestId(c, "[authorization.JWTAuthorization] user \"uid:%d\" token type is invalid", claims.Uid)
		utils.PrintJsonErrorResult(c, errs.ErrCurrentInvalidTokenType)
		return
	}

	c.SetTokenClaims(claims)
	c.Next()
}

// JWTAuthorizationByQueryString verifies whether current request is valid by jwt token
func JWTAuthorizationByQueryString(c *core.Context) {
	token, exists := c.GetQuery(tokenQueryStringParam)

	if !exists {
		log.WarnfWithRequestId(c, "[authorization.JWTAuthorizationByQueryString] no token provided")
		utils.PrintJsonErrorResult(c, errs.ErrUnauthorizedAccess)
		return
	}

	c.Request.Header.Set("Authorization", token)

	JWTAuthorization(c)
}

// JWTTwoFactorAuthorization verifies whether current request is valid by 2fa passcode
func JWTTwoFactorAuthorization(c *core.Context) {
	claims, err := getTokenClaims(c)

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

func getTokenClaims(c *core.Context) (*core.UserTokenClaims, *errs.Error) {
	token, claims, err := services.Tokens.ParseToken(c)

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
