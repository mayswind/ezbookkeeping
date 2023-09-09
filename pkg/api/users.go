package api

import (
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

// UsersApi represents user api
type UsersApi struct {
	users    *services.UserService
	tokens   *services.TokenService
	accounts *services.AccountService
}

// Initialize a user api singleton instance
var (
	Users = &UsersApi{
		users:    services.Users,
		tokens:   services.Tokens,
		accounts: services.Accounts,
	}
)

// UserRegisterHandler saves a new user by request parameters
func (a *UsersApi) UserRegisterHandler(c *core.Context) (interface{}, *errs.Error) {
	if !settings.Container.Current.EnableUserRegister {
		return nil, errs.ErrUserRegistrationNotAllowed
	}

	var userRegisterReq models.UserRegisterRequest
	err := c.ShouldBindJSON(&userRegisterReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[users.UserRegisterHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if userRegisterReq.DefaultCurrency == validators.ParentAccountCurrencyPlaceholder {
		log.WarnfWithRequestId(c, "[users.UserRegisterHandler] user default currency is invalid")
		return nil, errs.ErrUserDefaultCurrencyIsInvalid
	}

	userRegisterReq.Username = strings.TrimSpace(userRegisterReq.Username)
	userRegisterReq.Email = strings.TrimSpace(userRegisterReq.Email)
	userRegisterReq.Nickname = strings.TrimSpace(userRegisterReq.Nickname)

	user := &models.User{
		Username:             userRegisterReq.Username,
		Email:                userRegisterReq.Email,
		Nickname:             userRegisterReq.Nickname,
		Password:             userRegisterReq.Password,
		Language:             userRegisterReq.Language,
		DefaultCurrency:      userRegisterReq.DefaultCurrency,
		FirstDayOfWeek:       userRegisterReq.FirstDayOfWeek,
		TransactionEditScope: models.TRANSACTION_EDIT_SCOPE_ALL,
	}

	err = a.users.CreateUser(c, user)

	if err != nil {
		log.ErrorfWithRequestId(c, "[users.UserRegisterHandler] failed to create user \"%s\", because %s", user.Username, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[users.UserRegisterHandler] user \"%s\" has registered successfully, uid is %d", user.Username, user.Uid)

	authResp := &models.AuthResponse{
		Need2FA:         false,
		NeedVerifyEmail: settings.Container.Current.EnableUserForceVerifyEmail,
		User:            user.ToUserBasicInfo(),
	}

	if authResp.NeedVerifyEmail {
		return authResp, nil
	}

	token, claims, err := a.tokens.CreateToken(c, user)

	if err != nil {
		log.WarnfWithRequestId(c, "[users.UserRegisterHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return authResp, nil
	}

	authResp.Token = token
	c.SetTextualToken(token)
	c.SetTokenClaims(claims)

	log.InfofWithRequestId(c, "[users.UserRegisterHandler] user \"uid:%d\" has logined, token will be expired at %d", user.Uid, claims.ExpiresAt)

	return authResp, nil
}

// UserEmailVerifyHandler sets user email address verified
func (a *UsersApi) UserEmailVerifyHandler(c *core.Context) (interface{}, *errs.Error) {
	var userVerifyEmailReq models.UserVerifyEmailRequest
	err := c.ShouldBindJSON(&userVerifyEmailReq)

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.ErrorfWithRequestId(c, "[users.UserEmailVerifyHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if user.Disabled {
		log.WarnfWithRequestId(c, "[users.UserEmailVerifyHandler] user \"uid:%d\" is disabled", user.Uid)
		return nil, errs.ErrUserIsDisabled
	}

	if user.EmailVerified {
		log.WarnfWithRequestId(c, "[users.UserEmailVerifyHandler] user \"uid:%d\" email has been verified", user.Uid)
		return nil, errs.ErrEmailIsVerified
	}

	err = a.users.SetUserEmailVerified(c, user.Username)

	if err != nil {
		log.ErrorfWithRequestId(c, "[users.UserEmailVerifyHandler] failed to update user \"uid:%d\" email address verified, because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	now := time.Now().Unix()
	err = a.tokens.DeleteTokensByTypeBeforeTime(c, uid, core.USER_TOKEN_TYPE_EMAIL_VERIFY, now)

	if err == nil {
		log.InfofWithRequestId(c, "[users.UserEmailVerifyHandler] revoke old email verify tokens before unix time \"%d\" for user \"uid:%d\"", now, user.Uid)
	} else {
		log.WarnfWithRequestId(c, "[users.UserEmailVerifyHandler] failed to revoke old email verify tokens for user \"uid:%d\", because %s", user.Uid, err.Error())
	}

	resp := &models.UserVerifyEmailResponse{}

	if userVerifyEmailReq.RequestNewToken {
		token, claims, err := a.tokens.CreateToken(c, user)

		if err != nil {
			log.WarnfWithRequestId(c, "[users.UserEmailVerifyHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
			return resp, nil
		}

		resp.NewToken = token
		resp.User = user.ToUserBasicInfo()
		c.SetTextualToken(token)
		c.SetTokenClaims(claims)

		log.InfofWithRequestId(c, "[users.UserEmailVerifyHandler] user \"uid:%d\" token created, new token will be expired at %d", user.Uid, claims.ExpiresAt)
	}

	return resp, nil
}

// UserProfileHandler returns user profile of current user
func (a *UsersApi) UserProfileHandler(c *core.Context) (interface{}, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.ErrorfWithRequestId(c, "[users.UserRegisterHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	userResp := user.ToUserProfileResponse()
	return userResp, nil
}

// UserUpdateProfileHandler saves user profile by request parameters for current user
func (a *UsersApi) UserUpdateProfileHandler(c *core.Context) (interface{}, *errs.Error) {
	var userUpdateReq models.UserProfileUpdateRequest
	err := c.ShouldBindJSON(&userUpdateReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[users.UserUpdateProfileHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.ErrorfWithRequestId(c, "[users.UserUpdateProfileHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	userUpdateReq.Email = strings.TrimSpace(userUpdateReq.Email)
	userUpdateReq.Nickname = strings.TrimSpace(userUpdateReq.Nickname)

	anythingUpdate := false
	userNew := &models.User{
		Uid:  user.Uid,
		Salt: user.Salt,
	}

	if userUpdateReq.Email != "" && userUpdateReq.Email != user.Email {
		user.Email = userUpdateReq.Email
		userNew.Email = userUpdateReq.Email
		anythingUpdate = true
	}

	if userUpdateReq.Password != "" {
		if !a.users.IsPasswordEqualsUserPassword(userUpdateReq.OldPassword, user) {
			return nil, errs.ErrUserPasswordWrong
		}

		if !a.users.IsPasswordEqualsUserPassword(userUpdateReq.Password, user) {
			userNew.Password = userUpdateReq.Password
			anythingUpdate = true
		}
	}

	if userUpdateReq.Nickname != "" && userUpdateReq.Nickname != user.Nickname {
		user.Nickname = userUpdateReq.Nickname
		userNew.Nickname = userUpdateReq.Nickname
		anythingUpdate = true
	}

	if userUpdateReq.DefaultAccountId > 0 && userUpdateReq.DefaultAccountId != user.DefaultAccountId {
		accounts, err := a.accounts.GetAccountsByAccountIds(c, uid, []int64{userUpdateReq.DefaultAccountId})

		if err != nil || len(accounts) < 1 {
			return nil, errs.Or(err, errs.ErrUserDefaultAccountIsInvalid)
		}

		user.DefaultAccountId = userUpdateReq.DefaultAccountId
		userNew.DefaultAccountId = userUpdateReq.DefaultAccountId
		anythingUpdate = true
	}

	if userUpdateReq.TransactionEditScope != nil && *userUpdateReq.TransactionEditScope != user.TransactionEditScope {
		user.TransactionEditScope = *userUpdateReq.TransactionEditScope
		userNew.TransactionEditScope = *userUpdateReq.TransactionEditScope
		anythingUpdate = true
	} else {
		userNew.TransactionEditScope = models.TRANSACTION_EDIT_SCOPE_INVALID
	}

	modifyUserLanguage := false

	if userUpdateReq.Language != user.Language {
		user.Language = userUpdateReq.Language
		userNew.Language = userUpdateReq.Language
		modifyUserLanguage = true
		anythingUpdate = true
	}

	if userUpdateReq.DefaultCurrency != "" && userUpdateReq.DefaultCurrency != user.DefaultCurrency {
		user.DefaultCurrency = userUpdateReq.DefaultCurrency
		userNew.DefaultCurrency = userUpdateReq.DefaultCurrency
		anythingUpdate = true
	}

	if userUpdateReq.FirstDayOfWeek != nil && *userUpdateReq.FirstDayOfWeek != user.FirstDayOfWeek {
		user.FirstDayOfWeek = *userUpdateReq.FirstDayOfWeek
		userNew.FirstDayOfWeek = *userUpdateReq.FirstDayOfWeek
		anythingUpdate = true
	} else {
		userNew.FirstDayOfWeek = models.WEEKDAY_INVALID
	}

	if userUpdateReq.LongDateFormat != nil && *userUpdateReq.LongDateFormat != user.LongDateFormat {
		user.LongDateFormat = *userUpdateReq.LongDateFormat
		userNew.LongDateFormat = *userUpdateReq.LongDateFormat
		anythingUpdate = true
	} else {
		userNew.LongDateFormat = models.LONG_DATE_FORMAT_INVALID
	}

	if userUpdateReq.ShortDateFormat != nil && *userUpdateReq.ShortDateFormat != user.ShortDateFormat {
		user.ShortDateFormat = *userUpdateReq.ShortDateFormat
		userNew.ShortDateFormat = *userUpdateReq.ShortDateFormat
		anythingUpdate = true
	} else {
		userNew.ShortDateFormat = models.SHORT_DATE_FORMAT_INVALID
	}

	if userUpdateReq.LongTimeFormat != nil && *userUpdateReq.LongTimeFormat != user.LongTimeFormat {
		user.LongTimeFormat = *userUpdateReq.LongTimeFormat
		userNew.LongTimeFormat = *userUpdateReq.LongTimeFormat
		anythingUpdate = true
	} else {
		userNew.LongTimeFormat = models.LONG_TIME_FORMAT_INVALID
	}

	if userUpdateReq.ShortTimeFormat != nil && *userUpdateReq.ShortTimeFormat != user.ShortTimeFormat {
		user.ShortTimeFormat = *userUpdateReq.ShortTimeFormat
		userNew.ShortTimeFormat = *userUpdateReq.ShortTimeFormat
		anythingUpdate = true
	} else {
		userNew.ShortTimeFormat = models.SHORT_TIME_FORMAT_INVALID
	}

	if !anythingUpdate {
		return nil, errs.ErrNothingWillBeUpdated
	}

	keyProfileUpdated, err := a.users.UpdateUser(c, userNew, modifyUserLanguage)

	if err != nil {
		log.ErrorfWithRequestId(c, "[users.UserUpdateProfileHandler] failed to update user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[users.UserUpdateProfileHandler] user \"uid:%d\" has updated successfully", user.Uid)

	resp := &models.UserProfileUpdateResponse{
		User: user.ToUserBasicInfo(),
	}

	if keyProfileUpdated {
		now := time.Now().Unix()
		err = a.tokens.DeleteTokensBeforeTime(c, uid, now)

		if err == nil {
			log.InfofWithRequestId(c, "[users.UserUpdateProfileHandler] revoke old tokens before unix time \"%d\" for user \"uid:%d\"", now, user.Uid)
		} else {
			log.WarnfWithRequestId(c, "[users.UserUpdateProfileHandler] failed to revoke old tokens for user \"uid:%d\", because %s", user.Uid, err.Error())
		}

		token, claims, err := a.tokens.CreateToken(c, user)

		if err != nil {
			log.WarnfWithRequestId(c, "[users.UserUpdateProfileHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
			return resp, nil
		}

		resp.NewToken = token
		c.SetTextualToken(token)
		c.SetTokenClaims(claims)

		log.InfofWithRequestId(c, "[users.UserUpdateProfileHandler] user \"uid:%d\" token refreshed, new token will be expired at %d", user.Uid, claims.ExpiresAt)

		return resp, nil
	}

	return resp, nil
}

// UserSendVerifyEmailByUnloginUserHandler sends unlogin user verify email
func (a *UsersApi) UserSendVerifyEmailByUnloginUserHandler(c *core.Context) (interface{}, *errs.Error) {
	var userResendVerifyEmailReq models.UserResendVerifyEmailRequest
	err := c.ShouldBindJSON(&userResendVerifyEmailReq)

	user, err := a.users.GetUserByEmail(c, userResendVerifyEmailReq.Email)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.ErrorfWithRequestId(c, "[users.UserSendVerifyEmailByUnloginUserHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if !a.users.IsPasswordEqualsUserPassword(userResendVerifyEmailReq.Password, user) {
		log.WarnfWithRequestId(c, "[users.UserSendVerifyEmailByUnloginUserHandler] request password not equals to the user password")
		return nil, errs.ErrUserPasswordWrong
	}

	if user.Disabled {
		log.WarnfWithRequestId(c, "[users.UserSendVerifyEmailByUnloginUserHandler] user \"uid:%d\" is disabled", user.Uid)
		return nil, errs.ErrUserIsDisabled
	}

	if user.EmailVerified {
		log.WarnfWithRequestId(c, "[users.UserSendVerifyEmailByUnloginUserHandler] user \"uid:%d\" email has been verified", user.Uid)
		return nil, errs.ErrEmailIsVerified
	}

	if !settings.Container.Current.EnableSMTP {
		return nil, errs.ErrSMTPServerNotEnabled
	}

	token, _, err := a.tokens.CreateEmailVerifyToken(c, user)

	if err != nil {
		log.ErrorfWithRequestId(c, "[users.UserSendVerifyEmailByUnloginUserHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrTokenGenerating
	}

	err = a.users.SendVerifyEmail(user, token, c.GetClientLocale())

	if err != nil {
		log.WarnfWithRequestId(c, "[users.UserSendVerifyEmailByUnloginUserHandler] cannot send email to \"%s\", because %s", user.Email, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return true, nil
}

// UserSendVerifyEmailByLoginedUserHandler sends logined user verify email
func (a *UsersApi) UserSendVerifyEmailByLoginedUserHandler(c *core.Context) (interface{}, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.ErrorfWithRequestId(c, "[users.UserSendVerifyEmailByLoginedUserHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if user.EmailVerified {
		log.WarnfWithRequestId(c, "[users.UserSendVerifyEmailByLoginedUserHandler] user \"uid:%d\" email has been verified", user.Uid)
		return nil, errs.ErrEmailIsVerified
	}

	if !settings.Container.Current.EnableSMTP {
		return nil, errs.ErrSMTPServerNotEnabled
	}

	token, _, err := a.tokens.CreateEmailVerifyToken(c, user)

	if err != nil {
		log.ErrorfWithRequestId(c, "[users.UserSendVerifyEmailByLoginedUserHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrTokenGenerating
	}

	err = a.users.SendVerifyEmail(user, token, c.GetClientLocale())

	if err != nil {
		log.WarnfWithRequestId(c, "[users.UserSendVerifyEmailByLoginedUserHandler] cannot send email to \"%s\", because %s", user.Email, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return true, nil
}
