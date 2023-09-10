package api

import (
	"github.com/pquerna/otp/totp"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// AuthorizationsApi represents authorization api
type AuthorizationsApi struct {
	users                   *services.UserService
	tokens                  *services.TokenService
	twoFactorAuthorizations *services.TwoFactorAuthorizationService
}

// Initialize a authorization api singleton instance
var (
	Authorizations = &AuthorizationsApi{
		users:                   services.Users,
		tokens:                  services.Tokens,
		twoFactorAuthorizations: services.TwoFactorAuthorizations,
	}
)

// AuthorizeHandler verifies and authorizes current login request
func (a *AuthorizationsApi) AuthorizeHandler(c *core.Context) (interface{}, *errs.Error) {
	var credential models.UserLoginRequest
	err := c.ShouldBindJSON(&credential)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.AuthorizeHandler] parse request failed, because %s", err.Error())
		return nil, errs.ErrLoginNameOrPasswordInvalid
	}

	user, err := a.users.GetUserByUsernameOrEmailAndPassword(c, credential.LoginName, credential.Password)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.AuthorizeHandler] login failed for user \"%s\", because %s", credential.LoginName, err.Error())
		return nil, errs.ErrLoginNameOrPasswordWrong
	}

	if user.Disabled {
		log.WarnfWithRequestId(c, "[authorizations.AuthorizeHandler] login failed for user \"%s\", because user is disabled", credential.LoginName)
		return nil, errs.ErrUserIsDisabled
	}

	if settings.Container.Current.EnableUserForceVerifyEmail && !user.EmailVerified {
		log.WarnfWithRequestId(c, "[authorizations.AuthorizeHandler] login failed for user \"%s\", because user has not verified email", credential.LoginName)
		return nil, errs.NewErrorWithContext(errs.ErrEmailIsNotVerified, map[string]string{
			"email": user.Email,
		})
	}

	err = a.users.UpdateUserLastLoginTime(c, user.Uid)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.AuthorizeHandler] failed to update last login time for user \"uid:%d\", because %s", user.Uid, err.Error())
	}

	twoFactorEnable := a.tokens.CurrentConfig().EnableTwoFactor

	if twoFactorEnable {
		twoFactorEnable, err = a.twoFactorAuthorizations.ExistsTwoFactorSetting(c, user.Uid)

		if err != nil {
			log.ErrorfWithRequestId(c, "[authorizations.AuthorizeHandler] failed to check two factor setting for user \"uid:%d\", because %s", user.Uid, err.Error())
			return nil, errs.Or(err, errs.ErrSystemError)
		}
	}

	var token string
	var claims *core.UserTokenClaims

	if twoFactorEnable {
		token, claims, err = a.tokens.CreateRequire2FAToken(c, user)
	} else {
		token, claims, err = a.tokens.CreateToken(c, user)
	}

	if err != nil {
		log.ErrorfWithRequestId(c, "[authorizations.AuthorizeHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrTokenGenerating
	}

	if !twoFactorEnable {
		c.SetTextualToken(token)
	}

	c.SetTokenClaims(claims)

	log.InfofWithRequestId(c, "[authorizations.AuthorizeHandler] user \"uid:%d\" has logined, token type is %d, token will be expired at %d", user.Uid, claims.Type, claims.ExpiresAt)

	authResp := a.getAuthResponse(token, twoFactorEnable, user)
	return authResp, nil
}

// TwoFactorAuthorizeHandler verifies and authorizes current 2fa login by passcode
func (a *AuthorizationsApi) TwoFactorAuthorizeHandler(c *core.Context) (interface{}, *errs.Error) {
	var credential models.TwoFactorLoginRequest
	err := c.ShouldBindJSON(&credential)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] parse request failed, because %s", err.Error())
		return nil, errs.ErrPasscodeInvalid
	}

	uid := c.GetCurrentUid()
	twoFactorSetting, err := a.twoFactorAuthorizations.GetUserTwoFactorSettingByUid(c, uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] failed to get two factor setting for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrSystemError)
	}

	if !totp.Validate(credential.Passcode, twoFactorSetting.Secret) {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] passcode is invalid for user \"uid:%d\"", uid)
		return nil, errs.ErrPasscodeInvalid
	}

	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] failed to get user \"uid:%d\" info, because %s", user.Uid, err.Error())
		return nil, errs.ErrUserNotFound
	}

	if user.Disabled {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] user \"uid:%d\" is disabled", user.Uid)
		return nil, errs.ErrUserIsDisabled
	}

	if settings.Container.Current.EnableUserForceVerifyEmail && !user.EmailVerified {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] user \"uid:%d\" has not verified email", user.Uid)
		return nil, errs.ErrEmailIsNotVerified
	}

	oldTokenClaims := c.GetTokenClaims()
	err = a.tokens.DeleteTokenByClaims(c, oldTokenClaims)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] failed to revoke temporary token \"utid:%s\" for user \"uid:%d\", because %s", oldTokenClaims.UserTokenId, user.Uid, err.Error())
	}

	token, claims, err := a.tokens.CreateToken(c, user)

	if err != nil {
		log.ErrorfWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrTokenGenerating
	}

	c.SetTextualToken(token)
	c.SetTokenClaims(claims)

	log.InfofWithRequestId(c, "[authorizations.TwoFactorAuthorizeHandler] user \"uid:%d\" has authorized two factor via passcode, token will be expired at %d", user.Uid, claims.ExpiresAt)

	authResp := a.getAuthResponse(token, false, user)
	return authResp, nil
}

// TwoFactorAuthorizeByRecoveryCodeHandler verifies and authorizes current 2fa login by recovery code
func (a *AuthorizationsApi) TwoFactorAuthorizeByRecoveryCodeHandler(c *core.Context) (interface{}, *errs.Error) {
	var credential models.TwoFactorRecoveryCodeLoginRequest
	err := c.ShouldBindJSON(&credential)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] parse request failed, because %s", err.Error())
		return nil, errs.ErrTwoFactorRecoveryCodeInvalid
	}

	uid := c.GetCurrentUid()
	enableTwoFactor, err := a.twoFactorAuthorizations.ExistsTwoFactorSetting(c, uid)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] failed to get two factor setting for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrSystemError)
	}

	if !enableTwoFactor {
		return nil, errs.ErrTwoFactorIsNotEnabled
	}

	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] failed to get user \"uid:%d\" info, because %s", user.Uid, err.Error())
		return nil, errs.ErrUserNotFound
	}

	if user.Disabled {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] user \"uid:%d\" is disabled", user.Uid)
		return nil, errs.ErrUserIsDisabled
	}

	if settings.Container.Current.EnableUserForceVerifyEmail && !user.EmailVerified {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] user \"uid:%d\" has not verified email", user.Uid)
		return nil, errs.ErrEmailIsNotVerified
	}

	err = a.twoFactorAuthorizations.GetAndUseUserTwoFactorRecoveryCode(c, uid, credential.RecoveryCode, user.Salt)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] failed to get two factor recovery code for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrTwoFactorRecoveryCodeNotExist)
	}

	oldTokenClaims := c.GetTokenClaims()
	err = a.tokens.DeleteTokenByClaims(c, oldTokenClaims)

	if err != nil {
		log.WarnfWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] failed to revoke temporary token \"utid:%s\" for user \"uid:%d\", because %s", oldTokenClaims.UserTokenId, user.Uid, err.Error())
	}

	token, claims, err := a.tokens.CreateToken(c, user)

	if err != nil {
		log.ErrorfWithRequestId(c, "[authorizations.TwoFactorAuthorizeByRecoveryCodeHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrTokenGenerating
	}

	c.SetTextualToken(token)
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
