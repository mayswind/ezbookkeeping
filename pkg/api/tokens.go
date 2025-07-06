package api

import (
	"sort"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/avatars"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// TokensApi represents token api
type TokensApi struct {
	ApiUsingConfig
	ApiWithUserInfo
	tokens               *services.TokenService
	users                *services.UserService
	userAppCloudSettings *services.UserApplicationCloudSettingsService
}

// Initialize a token api singleton instance
var (
	Tokens = &TokensApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		ApiWithUserInfo: ApiWithUserInfo{
			ApiUsingConfig: ApiUsingConfig{
				container: settings.Container,
			},
			ApiUsingAvatarProvider: ApiUsingAvatarProvider{
				container: avatars.Container,
			},
		},
		tokens:               services.Tokens,
		users:                services.Users,
		userAppCloudSettings: services.UserApplicationCloudSettings,
	}
)

// TokenListHandler returns available token list of current user
func (a *TokensApi) TokenListHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	tokens, err := a.tokens.GetAllUnexpiredNormalAndMCPTokensByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[tokens.TokenListHandler] failed to get all tokens for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	tokenResps := make(models.TokenInfoResponseSlice, len(tokens))
	claims := c.GetTokenClaims()

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tokenResp := &models.TokenInfoResponse{
			TokenId:   a.tokens.GenerateTokenId(token),
			TokenType: token.TokenType,
			UserAgent: token.UserAgent,
			LastSeen:  token.LastSeenUnixTime,
		}

		if token.Uid == claims.Uid && utils.Int64ToString(token.UserTokenId) == claims.UserTokenId && token.CreatedUnixTime == claims.IssuedAt {
			tokenResp.IsCurrent = true
		}

		if token.TokenType == core.USER_TOKEN_TYPE_MCP && token.UserAgent != services.TokenUserAgentCreatedViaCli {
			tokenResp.UserAgent = services.TokenUserAgentForMCP
		}

		tokenResps[i] = tokenResp
	}

	sort.Sort(tokenResps)

	return tokenResps, nil
}

// TokenGenerateMCPHandler generates a new MCP token for current user
func (a *TokensApi) TokenGenerateMCPHandler(c *core.WebContext) (any, *errs.Error) {
	var generateMCPTokenReq models.TokenGenerateMCPRequest
	err := c.ShouldBindJSON(&generateMCPTokenReq)

	if err != nil {
		log.Warnf(c, "[tokens.TokenGenerateMCPHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.Warnf(c, "[tokens.TokenGenerateMCPHandler] failed to get user \"uid:%d\" info, because %s", uid, err.Error())
		return nil, errs.ErrUserNotFound
	}

	if !a.users.IsPasswordEqualsUserPassword(generateMCPTokenReq.Password, user) {
		return nil, errs.ErrUserPasswordWrong
	}

	token, claims, err := a.tokens.CreateMCPToken(c, user)

	if err != nil {
		log.Errorf(c, "[tokens.TokenGenerateMCPHandler] failed to create mcp token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrTokenGenerating)
	}

	log.Infof(c, "[tokens.TokenGenerateMCPHandler] user \"uid:%d\" has generated mcp token, new token will be expired at %d", user.Uid, claims.ExpiresAt)

	generateMCPTokenResp := &models.TokenGenerateMCPResponse{
		Token:  token,
		MCPUrl: a.CurrentConfig().RootUrl + "mcp",
	}

	return generateMCPTokenResp, nil
}

// TokenRevokeCurrentHandler revokes current token of current user
func (a *TokensApi) TokenRevokeCurrentHandler(c *core.WebContext) (any, *errs.Error) {
	_, claims, err := a.tokens.ParseTokenByHeader(c)

	if err != nil {
		return nil, errs.Or(err, errs.NewIncompleteOrIncorrectSubmissionError(err))
	}

	userTokenId, err := utils.StringToInt64(claims.UserTokenId)

	if err != nil {
		log.Warnf(c, "[tokens.TokenRevokeCurrentHandler] parse user token id failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	tokenRecord := &models.TokenRecord{
		Uid:             claims.Uid,
		UserTokenId:     userTokenId,
		CreatedUnixTime: claims.IssuedAt,
	}

	tokenId := a.tokens.GenerateTokenId(tokenRecord)
	err = a.tokens.DeleteToken(c, tokenRecord)

	if err != nil {
		log.Errorf(c, "[tokens.TokenRevokeCurrentHandler] failed to revoke token \"id:%s\" for user \"uid:%d\", because %s", tokenId, claims.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[tokens.TokenRevokeCurrentHandler] user \"uid:%d\" has revoked token \"id:%s\"", claims.Uid, tokenId)
	return true, nil
}

// TokenRevokeHandler revokes specific token of current user
func (a *TokensApi) TokenRevokeHandler(c *core.WebContext) (any, *errs.Error) {
	var tokenRevokeReq models.TokenRevokeRequest
	err := c.ShouldBindJSON(&tokenRevokeReq)

	if err != nil {
		log.Warnf(c, "[tokens.TokenRevokeHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	tokenRecord, err := a.tokens.ParseFromTokenId(tokenRevokeReq.TokenId)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[tokens.TokenRevokeHandler] failed to parse token \"id:%s\", because %s", tokenRevokeReq.TokenId, err.Error())
		}

		return nil, errs.Or(err, errs.ErrInvalidTokenId)
	}

	uid := c.GetCurrentUid()

	if tokenRecord.Uid != uid {
		log.Warnf(c, "[tokens.TokenRevokeHandler] token \"id:%s\" is not owned by user \"uid:%d\"", tokenRevokeReq.TokenId, uid)
		return nil, errs.ErrInvalidTokenId
	}

	if utils.Int64ToString(tokenRecord.UserTokenId) != c.GetTokenClaims().UserTokenId || tokenRecord.CreatedUnixTime != c.GetTokenClaims().IssuedAt {
		user, err := a.users.GetUserById(c, uid)

		if err != nil {
			if !errs.IsCustomError(err) {
				log.Errorf(c, "[tokens.TokenRevokeHandler] failed to get user, because %s", err.Error())
			}

			return nil, errs.ErrUserNotFound
		}

		if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_REVOKE_OTHER_SESSION) {
			return nil, errs.ErrNotPermittedToPerformThisAction
		}
	}

	err = a.tokens.DeleteToken(c, tokenRecord)

	if err != nil {
		log.Errorf(c, "[tokens.TokenRevokeHandler] failed to revoke token \"id:%s\" for user \"uid:%d\", because %s", tokenRevokeReq.TokenId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[tokens.TokenRevokeHandler] user \"uid:%d\" has revoked token \"id:%s\"", uid, tokenRevokeReq.TokenId)
	return true, nil
}

// TokenRevokeAllHandler revokes all tokens of current user except current token
func (a *TokensApi) TokenRevokeAllHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	tokens, err := a.tokens.GetAllTokensByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[tokens.TokenRevokeAllHandler] failed to get all tokens for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	claims := c.GetTokenClaims()
	currentTokenIndex := 0

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if token.Uid == claims.Uid && utils.Int64ToString(token.UserTokenId) == claims.UserTokenId && token.CreatedUnixTime == claims.IssuedAt {
			currentTokenIndex = i
			break
		}
	}

	tokens = append(tokens[:currentTokenIndex], tokens[currentTokenIndex+1:]...)

	if len(tokens) < 1 {
		return nil, errs.ErrTokenRecordNotFound
	}

	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[tokens.TokenRevokeAllHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_REVOKE_OTHER_SESSION) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	err = a.tokens.DeleteTokens(c, uid, tokens)

	if err != nil {
		log.Errorf(c, "[tokens.TokenRevokeAllHandler] failed to revoke all tokens for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[tokens.TokenRevokeAllHandler] user \"uid:%d\" has revoked all tokens", uid)
	return true, nil
}

// TokenRefreshHandler refresh current token of current user
func (a *TokensApi) TokenRefreshHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.Warnf(c, "[tokens.TokenRefreshHandler] failed to get user \"uid:%d\" info, because %s", uid, err.Error())
		return nil, errs.ErrUserNotFound
	}

	now := time.Now().Unix()
	oldTokenClaims := c.GetTokenClaims()

	if now-oldTokenClaims.IssuedAt < int64(a.CurrentConfig().TokenMinRefreshInterval) {
		log.Infof(c, "[tokens.TokenRefreshHandler] token of user \"uid:%d\" does not need to be refreshed", uid)

		userTokenId, err := utils.StringToInt64(oldTokenClaims.UserTokenId)

		if err != nil {
			log.Warnf(c, "[tokens.TokenRefreshHandler] parse user token id failed, because %s", err.Error())
		} else {
			tokenRecord := &models.TokenRecord{
				Uid:             oldTokenClaims.Uid,
				UserTokenId:     userTokenId,
				CreatedUnixTime: oldTokenClaims.IssuedAt,
			}

			tokenId := a.tokens.GenerateTokenId(tokenRecord)

			err = a.tokens.UpdateTokenLastSeen(c, tokenRecord)

			if err != nil {
				log.Warnf(c, "[tokens.TokenRefreshHandler] failed to update last seen of token \"id:%s\" for user \"uid:%d\", because %s", tokenId, uid, err.Error())
			}
		}

		userApplicationCloudSettings, err := a.userAppCloudSettings.GetUserApplicationCloudSettingsByUid(c, user.Uid)
		var applicationCloudSettingSlice *models.ApplicationCloudSettingSlice = nil

		if err != nil {
			log.Warnf(c, "[tokens.TokenRefreshHandler] failed to get latest user application cloud settings for user \"uid:%d\", because %s", user.Uid, err.Error())
		} else if userApplicationCloudSettings != nil && len(userApplicationCloudSettings.Settings) > 0 {
			applicationCloudSettingSlice = &userApplicationCloudSettings.Settings
		}

		refreshResp := &models.TokenRefreshResponse{
			User:                     a.GetUserBasicInfo(user),
			ApplicationCloudSettings: applicationCloudSettingSlice,
			NotificationContent:      a.GetAfterOpenNotificationContent(user.Language, c.GetClientLocale()),
		}

		return refreshResp, nil
	}

	token, claims, err := a.tokens.CreateToken(c, user)

	if err != nil {
		log.Errorf(c, "[tokens.TokenRefreshHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrTokenGenerating)
	}

	oldUserTokenId, _ := utils.StringToInt64(oldTokenClaims.UserTokenId)
	oldTokenRecord := &models.TokenRecord{
		Uid:             uid,
		UserTokenId:     oldUserTokenId,
		CreatedUnixTime: oldTokenClaims.IssuedAt,
	}

	c.SetTextualToken(token)
	c.SetTokenClaims(claims)

	userApplicationCloudSettings, err := a.userAppCloudSettings.GetUserApplicationCloudSettingsByUid(c, user.Uid)
	var applicationCloudSettingSlice *models.ApplicationCloudSettingSlice = nil

	if err != nil {
		log.Warnf(c, "[tokens.TokenRefreshHandler] failed to get latest user application cloud settings for user \"uid:%d\", because %s", user.Uid, err.Error())
	} else if userApplicationCloudSettings != nil && len(userApplicationCloudSettings.Settings) > 0 {
		applicationCloudSettingSlice = &userApplicationCloudSettings.Settings
	}

	log.Infof(c, "[tokens.TokenRefreshHandler] user \"uid:%d\" token refreshed, new token will be expired at %d", user.Uid, claims.ExpiresAt)

	refreshResp := &models.TokenRefreshResponse{
		NewToken:                 token,
		OldTokenId:               a.tokens.GenerateTokenId(oldTokenRecord),
		User:                     a.GetUserBasicInfo(user),
		ApplicationCloudSettings: applicationCloudSettingSlice,
		NotificationContent:      a.GetAfterOpenNotificationContent(user.Language, c.GetClientLocale()),
	}

	return refreshResp, nil
}
