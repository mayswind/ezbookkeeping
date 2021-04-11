package api

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"time"

	"github.com/pquerna/otp/totp"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
)

// TwoFactorAuthorizationsApi represents 2fa api
type TwoFactorAuthorizationsApi struct {
	twoFactorAuthorizations *services.TwoFactorAuthorizationService
	users                   *services.UserService
	tokens                  *services.TokenService
}

// Initialize a 2fa api singleton instance
var (
	TwoFactorAuthorizations = &TwoFactorAuthorizationsApi{
		twoFactorAuthorizations: services.TwoFactorAuthorizations,
		users:                   services.Users,
		tokens:                  services.Tokens,
	}
)

// TwoFactorStatusHandler returns 2fa status of current user
func (a *TwoFactorAuthorizationsApi) TwoFactorStatusHandler(c *core.Context) (interface{}, *errs.Error) {
	uid := c.GetCurrentUid()
	twoFactorSetting, err := a.twoFactorAuthorizations.GetUserTwoFactorSettingByUid(uid)

	if err == errs.ErrTwoFactorIsNotEnabled {
		statusResp := &models.TwoFactorStatusResponse{
			Enable: false,
		}

		return statusResp, nil
	}

	if err != nil {
		log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorStatusHandler] failed to get two factor setting, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	statusResp := &models.TwoFactorStatusResponse{
		Enable:    true,
		CreatedAt: twoFactorSetting.CreatedUnixTime,
	}

	return statusResp, nil
}

// TwoFactorEnableRequestHandler returns a new 2fa secret and qr code for current user to set 2fa and verify passcode next
func (a *TwoFactorAuthorizationsApi) TwoFactorEnableRequestHandler(c *core.Context) (interface{}, *errs.Error) {
	uid := c.GetCurrentUid()
	enabled, err := a.twoFactorAuthorizations.ExistsTwoFactorSetting(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableRequestHandler] failed to check two factor setting, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if enabled {
		return nil, errs.ErrTwoFactorAlreadyEnabled
	}

	user, err := a.users.GetUserById(uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableRequestHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	key, err := a.twoFactorAuthorizations.GenerateTwoFactorSecret(user)

	if err != nil {
		log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableRequestHandler] failed to generate two factor secret, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	img, err := key.Image(240, 240)

	if err != nil {
		log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableRequestHandler] failed to generate two factor qrcode, because %s", err.Error())
		return nil, errs.ErrOperationFailed
	}

	imgData := &bytes.Buffer{}

	if err = png.Encode(imgData, img); err != nil {
		return nil, errs.ErrOperationFailed
	}

	enableResp := &models.TwoFactorEnableResponse{
		Secret: key.Secret(),
		QRCode: "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgData.Bytes()),
	}

	return enableResp, nil
}

// TwoFactorEnableConfirmHandler enables 2fa for current user
func (a *TwoFactorAuthorizationsApi) TwoFactorEnableConfirmHandler(c *core.Context) (interface{}, *errs.Error) {
	var confirmReq models.TwoFactorEnableConfirmRequest
	err := c.ShouldBindJSON(&confirmReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableConfirmHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	exists, err := a.twoFactorAuthorizations.ExistsTwoFactorSetting(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableConfirmHandler] failed to check two factor setting, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if exists {
		return nil, errs.ErrTwoFactorAlreadyEnabled
	}

	user, err := a.users.GetUserById(uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableConfirmHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	twoFactorSetting := &models.TwoFactor{
		Uid:    uid,
		Secret: confirmReq.Secret,
	}

	if !totp.Validate(confirmReq.Passcode, confirmReq.Secret) {
		log.WarnfWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableConfirmHandler] passcode is invalid")
		return nil, errs.ErrPasscodeInvalid
	}

	recoveryCodes, err := a.twoFactorAuthorizations.GenerateTwoFactorRecoveryCodes()

	if err != nil {
		log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableConfirmHandler] failed to generate two factor recovery codes for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	err = a.twoFactorAuthorizations.CreateTwoFactorRecoveryCodes(uid, recoveryCodes, user.Salt)

	if err != nil {
		log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableConfirmHandler] failed to create two factor recovery codes for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	err = a.twoFactorAuthorizations.CreateTwoFactorSetting(twoFactorSetting)

	if err != nil {
		log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableConfirmHandler] failed to create two factor setting for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableConfirmHandler] user \"uid:%d\" has enabled two factor authorization", uid)

	now := time.Now().Unix()
	err = a.tokens.DeleteTokensBeforeTime(uid, now)

	if err == nil {
		log.InfofWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableConfirmHandler] revoke old tokens before unix time \"%d\" for user \"uid:%d\"", now, user.Uid)
	} else {
		log.WarnfWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableConfirmHandler] failed to revoke old tokens for user \"uid:%d\", because %s", user.Uid, err.Error())
	}

	token, claims, err := a.tokens.CreateToken(user, c)

	if err != nil {
		log.WarnfWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableConfirmHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())

		confirmResp := &models.TwoFactorEnableConfirmResponse{
			RecoveryCodes: recoveryCodes,
		}

		return confirmResp, nil
	}

	c.SetTokenClaims(claims)

	log.InfofWithRequestId(c, "[twofactor_authorizations.TwoFactorEnableConfirmHandler] user \"uid:%d\" token refreshed, new token will be expired at %d", user.Uid, claims.ExpiresAt)

	confirmResp := &models.TwoFactorEnableConfirmResponse{
		Token:         token,
		RecoveryCodes: recoveryCodes,
	}

	return confirmResp, nil
}

// TwoFactorDisableHandler disables 2fa for current user
func (a *TwoFactorAuthorizationsApi) TwoFactorDisableHandler(c *core.Context) (interface{}, *errs.Error) {
	var disableReq models.TwoFactorDisableRequest
	err := c.ShouldBindJSON(&disableReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[twofactor_authorizations.TwoFactorDisableHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.WarnfWithRequestId(c, "[twofactor_authorizations.TwoFactorDisableHandler] failed to get user for user \"uid:%d\", because %s", uid, err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if !a.users.IsPasswordEqualsUserPassword(disableReq.Password, user) {
		return nil, errs.ErrUserPasswordWrong
	}

	enableTwoFactor, err := a.twoFactorAuthorizations.ExistsTwoFactorSetting(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorDisableHandler] failed to check two factor setting, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if !enableTwoFactor {
		return nil, errs.ErrTwoFactorIsNotEnabled
	}

	err = a.twoFactorAuthorizations.DeleteTwoFactorRecoveryCodes(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorDisableHandler] failed to delete two factor recovery codes for user \"uid:%d\"", uid)
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	err = a.twoFactorAuthorizations.DeleteTwoFactorSetting(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorDisableHandler] failed to delete two factor setting for user \"uid:%d\"", uid)
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[twofactor_authorizations.TwoFactorDisableHandler] user \"uid:%d\" has disabled two factor authorization", uid)

	return true, nil
}

// TwoFactorRecoveryCodeRegenerateHandler returns new 2fa recovery codes and revokes old recovery codes for current user
func (a *TwoFactorAuthorizationsApi) TwoFactorRecoveryCodeRegenerateHandler(c *core.Context) (interface{}, *errs.Error) {
	var regenerateReq models.TwoFactorRegenerateRecoveryCodeRequest
	err := c.ShouldBindJSON(&regenerateReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[twofactor_authorizations.TwoFactorRecoveryCodeRegenerateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.WarnfWithRequestId(c, "[twofactor_authorizations.TwoFactorRecoveryCodeRegenerateHandler] failed to get user for user \"uid:%d\", because %s", uid, err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if !a.users.IsPasswordEqualsUserPassword(regenerateReq.Password, user) {
		return nil, errs.ErrUserPasswordWrong
	}

	enableTwoFactor, err := a.twoFactorAuthorizations.ExistsTwoFactorSetting(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorRecoveryCodeRegenerateHandler] failed to check two factor setting, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if !enableTwoFactor {
		return nil, errs.ErrTwoFactorIsNotEnabled
	}

	recoveryCodes, err := a.twoFactorAuthorizations.GenerateTwoFactorRecoveryCodes()

	if err != nil {
		log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorRecoveryCodeRegenerateHandler] failed to generate two factor recovery codes for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	err = a.twoFactorAuthorizations.CreateTwoFactorRecoveryCodes(uid, recoveryCodes, user.Salt)

	if err != nil {
		log.ErrorfWithRequestId(c, "[twofactor_authorizations.TwoFactorRecoveryCodeRegenerateHandler] failed to create two factor recovery codes for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	recoveryCodesResp := &models.TwoFactorEnableConfirmResponse{
		RecoveryCodes: recoveryCodes,
	}

	log.InfofWithRequestId(c, "[twofactor_authorizations.TwoFactorRecoveryCodeRegenerateHandler] user \"uid:%d\" has regenerated two factor recovery codes", uid)

	return recoveryCodesResp, nil
}
