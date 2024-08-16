package api

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// ForgetPasswordsApi represents user forget password api
type ForgetPasswordsApi struct {
	ApiUsingConfig
	users           *services.UserService
	tokens          *services.TokenService
	forgetPasswords *services.ForgetPasswordService
}

// Initialize a user api singleton instance
var (
	ForgetPasswords = &ForgetPasswordsApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		users:           services.Users,
		tokens:          services.Tokens,
		forgetPasswords: services.ForgetPasswords,
	}
)

// UserForgetPasswordRequestHandler generates password reset link and send user an email with this link
func (a *ForgetPasswordsApi) UserForgetPasswordRequestHandler(c *core.Context) (any, *errs.Error) {
	var request models.ForgetPasswordRequest
	err := c.ShouldBindJSON(&request)

	if err != nil {
		log.WarnfWithRequestId(c, "[forget_passwords.UserForgetPasswordRequestHandler] parse request failed, because %s", err.Error())
		return nil, errs.ErrEmailIsEmptyOrInvalid
	}

	user, err := a.users.GetUserByEmail(c, request.Email)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.ErrorfWithRequestId(c, "[forget_passwords.UserForgetPasswordRequestHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if user.Disabled {
		log.WarnfWithRequestId(c, "[forget_passwords.UserForgetPasswordRequestHandler] user \"uid:%d\" is disabled", user.Uid)
		return nil, errs.ErrUserIsDisabled
	}

	if a.CurrentConfig().ForgetPasswordRequireVerifyEmail && !user.EmailVerified {
		log.WarnfWithRequestId(c, "[forget_passwords.UserForgetPasswordRequestHandler] user \"uid:%d\" has not verified email", user.Uid)
		return nil, errs.ErrEmailIsNotVerified
	}

	if !a.CurrentConfig().EnableSMTP {
		return nil, errs.ErrSMTPServerNotEnabled
	}

	token, _, err := a.tokens.CreatePasswordResetToken(c, user)

	if err != nil {
		log.ErrorfWithRequestId(c, "[forget_passwords.UserForgetPasswordRequestHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrTokenGenerating
	}

	go func() {
		err = a.forgetPasswords.SendPasswordResetEmail(c, user, token, c.GetClientLocale())

		if err != nil {
			log.WarnfWithRequestId(c, "[forget_passwords.UserForgetPasswordRequestHandler] cannot send email to \"%s\", because %s", user.Email, err.Error())
		}
	}()

	return true, nil
}

// UserResetPasswordHandler resets user password by request parameters
func (a *ForgetPasswordsApi) UserResetPasswordHandler(c *core.Context) (any, *errs.Error) {
	var request models.PasswordResetRequest
	err := c.ShouldBindJSON(&request)

	if err != nil {
		log.WarnfWithRequestId(c, "[forget_passwords.UserResetPasswordHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.ErrorfWithRequestId(c, "[forget_passwords.UserResetPasswordHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if user.Disabled {
		log.WarnfWithRequestId(c, "[forget_passwords.UserResetPasswordHandler] user \"uid:%d\" is disabled", user.Uid)
		return nil, errs.ErrUserIsDisabled
	}

	if a.CurrentConfig().ForgetPasswordRequireVerifyEmail && !user.EmailVerified {
		log.WarnfWithRequestId(c, "[forget_passwords.UserResetPasswordHandler] user \"uid:%d\" has not verified email", user.Uid)
		return nil, errs.ErrEmailIsNotVerified
	}

	if user.Email != request.Email {
		log.WarnfWithRequestId(c, "[forget_passwords.UserResetPasswordHandler] request email not equals the user email")
		return nil, errs.ErrEmptyIsInvalid
	}

	if a.users.IsPasswordEqualsUserPassword(request.Password, user) {
		oldTokenClaims := c.GetTokenClaims()
		err = a.tokens.DeleteTokenByClaims(c, oldTokenClaims)

		if err != nil {
			log.WarnfWithRequestId(c, "[forget_passwords.UserResetPasswordHandler] failed to revoke password reset token \"utid:%s\" for user \"uid:%d\", because %s", oldTokenClaims.UserTokenId, user.Uid, err.Error())
		}

		return nil, errs.ErrNewPasswordEqualsOldInvalid
	}

	userNew := &models.User{
		Uid:      user.Uid,
		Salt:     user.Salt,
		Password: request.Password,
	}

	_, _, err = a.users.UpdateUser(c, userNew, false)

	if err != nil {
		log.ErrorfWithRequestId(c, "[forget_passwords.UserResetPasswordHandler] failed to update user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	now := time.Now().Unix()
	err = a.tokens.DeleteTokensBeforeTime(c, uid, now)

	if err == nil {
		log.InfofWithRequestId(c, "[forget_passwords.UserResetPasswordHandler] revoke old tokens before unix time \"%d\" for user \"uid:%d\"", now, user.Uid)
	} else {
		log.WarnfWithRequestId(c, "[forget_passwords.UserResetPasswordHandler] failed to revoke old tokens for user \"uid:%d\", because %s", user.Uid, err.Error())
	}

	return true, nil
}
