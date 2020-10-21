package api

import (
	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/services"
	"github.com/mayswind/lab/pkg/utils"
)

type TokensApi struct {
	tokens *services.TokenService
	users  *services.UserService
}

var (
	Tokens = &TokensApi{
		tokens: services.Tokens,
		users:  services.Users,
	}
)

func (a *TokensApi) TokenListHandler(c *core.Context) (interface{}, *errs.Error) {
	uid := c.GetCurrentUid()
	tokens, err := a.tokens.GetAllTokensByUid(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[tokens.TokenListHandler] failed to get all tokens for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	tokenResps := make([]*models.TokenInfoResponse, len(tokens))
	claims := c.GetTokenClaims()

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tokenResp := &models.TokenInfoResponse{
			TokenId:   a.tokens.GenerateTokenId(token),
			TokenType: token.TokenType,
			UserAgent: token.UserAgent,
			CreatedAt: token.CreatedUnixTime,
			ExpiredAt: token.ExpiredUnixTime,
		}

		if utils.Int64ToString(token.Uid) == claims.Id && utils.Int64ToString(token.UserTokenId) == claims.UserTokenId && token.CreatedUnixTime == claims.IssuedAt {
			tokenResp.IsCurrent = true
		}

		tokenResps[i] = tokenResp
	}

	return tokenResps, nil
}

func (a *TokensApi) TokenRevokeCurrentHandler(c *core.Context) (interface{}, *errs.Error) {
	claims := c.GetTokenClaims()
	uid := c.GetCurrentUid()
	userTokenId, err := utils.StringToInt64(claims.UserTokenId)

	if err != nil {
		log.WarnfWithRequestId(c, "[tokens.TokenRevokeCurrentHandler] parse user token id failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	tokenRecord := &models.TokenRecord{
		Uid: uid,
		UserTokenId: userTokenId,
		CreatedUnixTime: claims.IssuedAt,
	}

	tokenId := a.tokens.GenerateTokenId(tokenRecord)
	err = a.tokens.DeleteToken(tokenRecord)

	if err != nil {
		log.ErrorfWithRequestId(c, "[token.TokenRevokeCurrentHandler] failed to revoke token \"id:%s\" for user \"uid:%d\", because %s", tokenId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[token.TokenRevokeCurrentHandler] user \"uid:%d\" has revoked token \"id:%s\"", uid, tokenId)
	return true, nil
}

func (a *TokensApi) TokenRevokeHandler(c *core.Context) (interface{}, *errs.Error) {
	var tokenRevokeReq models.TokenRevokeRequest
	err := c.ShouldBindJSON(&tokenRevokeReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[tokens.TokenRevokeHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	tokenRecord, err := a.tokens.ParseFromTokenId(tokenRevokeReq.TokenId)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.ErrorfWithRequestId(c, "[token.TokenRevokeHandler] failed to parse token \"id:%s\", because %s", tokenRevokeReq.TokenId, err.Error())
		}

		return nil, errs.Or(err, errs.ErrInvalidTokenId)
	}

	uid := c.GetCurrentUid()

	if tokenRecord.Uid != uid {
		log.WarnfWithRequestId(c, "[token.TokenRevokeHandler] token \"id:%s\" is not owned by user \"uid:%d\"", tokenRevokeReq.TokenId, uid)
		return nil, errs.ErrInvalidTokenId
	}

	err = a.tokens.DeleteToken(tokenRecord)

	if err != nil {
		log.ErrorfWithRequestId(c, "[token.TokenRevokeHandler] failed to revoke token \"id:%s\" for user \"uid:%d\", because %s", tokenRevokeReq.TokenId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[token.TokenRevokeHandler] user \"uid:%d\" has revoked token \"id:%s\"", uid, tokenRevokeReq.TokenId)
	return true, nil
}

func (a *TokensApi) TokenRefreshHandler(c *core.Context) (interface{}, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(uid)

	if err != nil {
		log.WarnfWithRequestId(c, "[token.TokenRefreshHandler] failed to get user \"uid:%d\" info, because %s", uid, err.Error())
		return nil, errs.ErrUserNotFound
	}

	token, claims, err := a.tokens.CreateToken(user, c)

	if err != nil {
		log.ErrorfWithRequestId(c, "[token.TokenRefreshHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrTokenGenerating)
	}

	oldTokenClaims := c.GetTokenClaims()
	err = a.tokens.DeleteTokenByClaims(oldTokenClaims)

	if err != nil {
		log.WarnfWithRequestId(c, "[token.TokenRefreshHandler] failed to revoke token \"id:%s\" for user \"uid:%d\", because %s", oldTokenClaims.UserTokenId, user.Uid, err.Error())
	}

	c.SetTokenClaims(claims)

	log.InfofWithRequestId(c, "[token.TokenRefreshHandler] user \"uid:%d\" token refreshed, new token will be expired at %d", user.Uid, claims.ExpiresAt)
	return token, nil
}
