package middlewares

import (
	"time"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/services"
	"github.com/mayswind/lab/pkg/utils"
)

func JWTAuthorization(c *core.Context) {
	claims, err := getTokenClaims(c)

	if err != nil {
		utils.PrintErrorResult(c, err)
		return
	}

	if claims.Type == core.USER_TOKEN_TYPE_REQUIRE_2FA {
		log.WarnfWithRequestId(c, "[authorization.JWTAuthorization] user \"uid:%s\" token requires 2fa", claims.Id)
		utils.PrintErrorResult(c, errs.ErrCurrentTokenRequire2FA)
		return
	}

	if claims.Type != core.USER_TOKEN_TYPE_NORMAL {
		log.WarnfWithRequestId(c, "[authorization.JWTAuthorization] user \"uid:%s\" token type is invalid", claims.Id)
		utils.PrintErrorResult(c, errs.ErrCurrentInvalidTokenType)
		return
	}

	c.SetTokenClaims(claims)
	c.Next()
}

func JWTTwoFactorAuthorization(c *core.Context) {
	claims, err := getTokenClaims(c)

	if err != nil {
		utils.PrintErrorResult(c, err)
		return
	}

	if claims.Type != core.USER_TOKEN_TYPE_REQUIRE_2FA {
		log.WarnfWithRequestId(c, "[authorization.JWTTwoFactorAuthorization] user \"uid:%s\" token is not need two factor authorization", claims.Id)
		utils.PrintErrorResult(c, errs.ErrCurrentTokenNotRequire2FA)
		return
	}

	c.SetTokenClaims(claims)
	c.Next()
}

func getTokenClaims(c *core.Context) (*core.UserTokenClaims, *errs.Error) {
	token, claims, err := services.Tokens.ParseToken(c)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorization.getTokenClaims] failed to parse token, because %s", err.Error())
		return nil, errs.ErrUnauthorizedAccess
	}

	if !token.Valid {
		log.WarnfWithRequestId(c, "[authorization.getTokenClaims] token is invalid")
		return nil, errs.ErrCurrentInvalidToken
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		log.WarnfWithRequestId(c, "[authorization.getTokenClaims] token is expired")
		return nil, errs.ErrCurrentTokenExpired
	}

	if claims.Id == "" {
		log.WarnfWithRequestId(c, "[authorization.getTokenClaims] user id in token is empty")
		return nil, errs.ErrCurrentInvalidToken
	}

	return claims, nil
}
