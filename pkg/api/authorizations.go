package api

import (
	"github.com/pquerna/otp/totp"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/services"
)

type AuthorizationsApi struct {
	users                   *services.UserService
	tokens                  *services.TokenService
	twoFactorAuthorizations *services.TwoFactorAuthorizationService
}

var (
	Authorizations = &AuthorizationsApi{
		users:                   services.Users,
		tokens:                  services.Tokens,
		twoFactorAuthorizations: services.TwoFactorAuthorizations,
	}
)

func (a *AuthorizationsApi) AuthorizeHandler(c *core.Context) (interface{}, *errs.Error) {
	var credential models.UserLoginRequest
	err := c.ShouldBindJSON(&credential)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.AuthorizeHandler] parse request failed, because %s", err.Error())
		return nil, errs.ErrLoginNameOrPasswordInvalid
	}

	user, err := a.users.GetUserByUsernameOrEmailAndPassword(credential.LoginName, credential.Password)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.AuthorizeHandler] login failed for user \"%s\", because %s", credential.LoginName, err.Error())
		return nil, errs.ErrLoginNameOrPasswordWrong
	}

	err = a.users.UpdateUserLastLoginTime(user.Uid)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.AuthorizeHandler] failed to update last login time for user \"uid:%d\", because %s", user.Uid, err.Error())
	}

	twoFactorEnable := a.tokens.CurrentConfig().EnableTwoFactor

	if twoFactorEnable {
		twoFactorEnable, err = a.twoFactorAuthorizations.ExistsTwoFactorSetting(user.Uid)

		if err != nil {
			log.ErrorfWithRequestId(c, "[authorizations.AuthorizeHandler] failed to check two factor setting for user \"uid:%d\", because %s", user.Uid, err.Error())
			return nil, errs.Or(err, errs.ErrSystemError)
		}
	}

	var token string
	var claims *core.UserTokenClaims

	if twoFactorEnable {
		token, claims, err = a.tokens.CreateRequire2FAToken(user, c)
	} else {
		token, claims, err = a.tokens.CreateToken(user, c)
	}

	if err != nil {
		log.ErrorfWithRequestId(c, "[authorizations.AuthorizeHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrTokenGenerating
	}

	c.SetTokenClaims(claims)

	log.InfofWithRequestId(c, "[authorizations.AuthorizeHandler] user \"uid:%d\" has logined, token type is %d, token will be expired at %d", user.Uid, claims.Type, claims.ExpiresAt)

	authResp := a.getAuthResponse(token, twoFactorEnable, user)
	return authResp, nil
}

func (a *AuthorizationsApi) TwoFactorAuthorizeHandler(c *core.Context) (interface{}, *errs.Error) {
	var credential models.TwoFactorLoginRequest
	err := c.ShouldBindJSON(&credential)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] parse request failed, because %s", err.Error())
		return nil, errs.ErrPasscodeInvalid
	}

	uid := c.GetCurrentUid()
	twoFactorSetting, err := a.twoFactorAuthorizations.GetUserTwoFactorSettingByUid(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] failed to get two factor setting for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrSystemError)
	}

	if !totp.Validate(credential.Passcode, twoFactorSetting.Secret) {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] passcode is invalid for user \"uid:%d\"", uid)
		return nil, errs.ErrPasscodeInvalid
	}

	user, err := a.users.GetUserById(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] failed to get user \"uid:%d\" info, because %s", user.Uid, err.Error())
		return nil, errs.ErrUserNotFound
	}

	oldTokenClaims := c.GetTokenClaims()
	err = a.tokens.DeleteTokenByClaims(oldTokenClaims)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] failed to revoke temporary token \"utid:%s\" for user \"uid:%d\", because %s", oldTokenClaims.UserTokenId, user.Uid, err.Error())
	}

	token, claims, err := a.tokens.CreateToken(user, c)

	if err != nil {
		log.ErrorfWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrTokenGenerating
	}

	c.SetTokenClaims(claims)

	log.InfofWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] user \"uid:%d\" has authorized two factor via passcode, token will be expired at %d", user.Uid, claims.ExpiresAt)

	authResp := a.getAuthResponse(token, false, user)
	return authResp, nil
}

func (a *AuthorizationsApi) TwoFactorAuthorizeByRecoveryCodeHandler(c *core.Context) (interface{}, *errs.Error) {
	var credential models.TwoFactorRecoveryCodeLoginRequest
	err := c.ShouldBindJSON(&credential)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] parse request failed, because %s", err.Error())
		return nil, errs.ErrTwoFactorRecoveryCodeInvalid
	}

	uid := c.GetCurrentUid()
	enableTwoFactor, err := a.twoFactorAuthorizations.ExistsTwoFactorSetting(uid)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] failed to get two factor setting for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrSystemError)
	}

	if !enableTwoFactor {
		return nil, errs.ErrTwoFactorIsNotEnabled
	}

	user, err := a.users.GetUserById(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] failed to get user \"uid:%d\" info, because %s", user.Uid, err.Error())
		return nil, errs.ErrUserNotFound
	}

	err = a.twoFactorAuthorizations.GetAndUseUserTwoFactorRecoveryCode(uid, credential.RecoveryCode, user.Salt)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] failed to get two factor recovery code for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrTwoFactorRecoveryCodeNotExist)
	}

	oldTokenClaims := c.GetTokenClaims()
	err = a.tokens.DeleteTokenByClaims(oldTokenClaims)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] failed to revoke temporary token \"utid:%s\" for user \"uid:%d\", because %s", oldTokenClaims.UserTokenId, user.Uid, err.Error())
	}

	token, claims, err := a.tokens.CreateToken(user, c)

	if err != nil {
		log.ErrorfWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrTokenGenerating
	}

	c.SetTokenClaims(claims)

	log.InfofWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] user \"uid:%d\" has authorized two factor via recovery code \"%s\", token will be expired at %d", user.Uid, credential.RecoveryCode, claims.ExpiresAt)

	authResp := a.getAuthResponse(token, false, user)
	return authResp, nil
}

func (a *AuthorizationsApi) getAuthResponse(token string, need2FA bool, user *models.User) *models.AuthResponse {
	return &models.AuthResponse{
		Token:   token,
		Need2FA: need2FA,
		User:    user.ToUserBasicInfo(),
	}
}
