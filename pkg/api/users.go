package api

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"

	"github.com/mayswind/ezbookkeeping/pkg/avatars"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/locales"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

// UsersApi represents user api
type UsersApi struct {
	ApiUsingConfig
	ApiWithUserInfo
	users    *services.UserService
	tokens   *services.TokenService
	accounts *services.AccountService
}

// Initialize a user api singleton instance
var (
	Users = &UsersApi{
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
		users:    services.Users,
		tokens:   services.Tokens,
		accounts: services.Accounts,
	}
)

// UserRegisterHandler saves a new user by request parameters
func (a *UsersApi) UserRegisterHandler(c *core.WebContext) (any, *errs.Error) {
	if !a.CurrentConfig().EnableUserRegister {
		return nil, errs.ErrUserRegistrationNotAllowed
	}

	var userRegisterReq models.UserRegisterRequest
	err := c.ShouldBindBodyWith(&userRegisterReq, binding.JSON)

	if err != nil {
		log.Warnf(c, "[users.UserRegisterHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if userRegisterReq.DefaultCurrency == validators.ParentAccountCurrencyPlaceholder {
		log.Warnf(c, "[users.UserRegisterHandler] user default currency is invalid")
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
		FeatureRestriction:   a.CurrentConfig().DefaultFeatureRestrictions,
	}

	err = a.users.CreateUser(c, user)

	if err != nil {
		log.Errorf(c, "[users.UserRegisterHandler] failed to create user \"%s\", because %s", user.Username, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[users.UserRegisterHandler] user \"%s\" has registered successfully, uid is %d", user.Username, user.Uid)

	presetCategoriesSaved := false

	if len(userRegisterReq.Categories) > 0 {
		_, err = TransactionCategories.createBatchCategories(c, user.Uid, &userRegisterReq.TransactionCategoryCreateBatchRequest)

		if err == nil {
			presetCategoriesSaved = true
		}
	}

	authResp := &models.RegisterResponse{
		AuthResponse: models.AuthResponse{
			Need2FA:             false,
			User:                a.GetUserBasicInfo(user),
			NotificationContent: a.GetAfterRegisterNotificationContent(user.Language, c.GetClientLocale()),
		},
		NeedVerifyEmail:       a.CurrentConfig().EnableUserVerifyEmail && a.CurrentConfig().EnableUserForceVerifyEmail,
		PresetCategoriesSaved: presetCategoriesSaved,
	}

	if a.CurrentConfig().EnableUserVerifyEmail && a.CurrentConfig().EnableSMTP {
		token, _, err := a.tokens.CreateEmailVerifyToken(c, user)

		if err != nil {
			log.Errorf(c, "[users.UserRegisterHandler] failed to create email verify token for user \"uid:%d\", because %s", user.Uid, err.Error())
		} else {
			go func() {
				err = a.users.SendVerifyEmail(user, token, c.GetClientLocale())

				if err != nil {
					log.Warnf(c, "[users.UserRegisterHandler] cannot send verify email to \"%s\", because %s", user.Email, err.Error())
				}
			}()
		}
	}

	if a.CurrentConfig().EnableUserForceVerifyEmail {
		return authResp, nil
	}

	token, claims, err := a.tokens.CreateToken(c, user)

	if err != nil {
		log.Warnf(c, "[users.UserRegisterHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return authResp, nil
	}

	authResp.Token = token
	c.SetTextualToken(token)
	c.SetTokenClaims(claims)

	log.Infof(c, "[users.UserRegisterHandler] user \"uid:%d\" has logined, token will be expired at %d", user.Uid, claims.ExpiresAt)

	return authResp, nil
}

// UserEmailVerifyHandler sets user email address verified
func (a *UsersApi) UserEmailVerifyHandler(c *core.WebContext) (any, *errs.Error) {
	var userVerifyEmailReq models.UserVerifyEmailRequest
	err := c.ShouldBindJSON(&userVerifyEmailReq)

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[users.UserEmailVerifyHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if user.Disabled {
		log.Warnf(c, "[users.UserEmailVerifyHandler] user \"uid:%d\" is disabled", user.Uid)
		return nil, errs.ErrUserIsDisabled
	}

	if user.EmailVerified {
		log.Warnf(c, "[users.UserEmailVerifyHandler] user \"uid:%d\" email has been verified", user.Uid)
		return nil, errs.ErrEmailIsVerified
	}

	err = a.users.SetUserEmailVerified(c, user.Username)

	if err != nil {
		log.Errorf(c, "[users.UserEmailVerifyHandler] failed to update user \"uid:%d\" email address verified, because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	err = a.tokens.DeleteTokensByType(c, uid, core.USER_TOKEN_TYPE_EMAIL_VERIFY)

	if err == nil {
		log.Infof(c, "[users.UserEmailVerifyHandler] revoke old email verify tokens for user \"uid:%d\"", user.Uid)
	} else {
		log.Warnf(c, "[users.UserEmailVerifyHandler] failed to revoke old email verify tokens for user \"uid:%d\", because %s", user.Uid, err.Error())
	}

	resp := &models.UserVerifyEmailResponse{}

	if userVerifyEmailReq.RequestNewToken {
		token, claims, err := a.tokens.CreateToken(c, user)

		if err != nil {
			log.Warnf(c, "[users.UserEmailVerifyHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
			return resp, nil
		}

		resp.NewToken = token
		resp.User = a.GetUserBasicInfo(user)
		resp.NotificationContent = a.GetAfterLoginNotificationContent(user.Language, c.GetClientLocale())

		c.SetTextualToken(token)
		c.SetTokenClaims(claims)

		log.Infof(c, "[users.UserEmailVerifyHandler] user \"uid:%d\" token created, new token will be expired at %d", user.Uid, claims.ExpiresAt)
	}

	return resp, nil
}

// UserProfileHandler returns user profile of current user
func (a *UsersApi) UserProfileHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[users.UserRegisterHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	userResp := a.getUserProfileResponse(user)
	return userResp, nil
}

// UserUpdateProfileHandler saves user profile by request parameters for current user
func (a *UsersApi) UserUpdateProfileHandler(c *core.WebContext) (any, *errs.Error) {
	var userUpdateReq models.UserProfileUpdateRequest
	err := c.ShouldBindJSON(&userUpdateReq)

	if err != nil {
		log.Warnf(c, "[users.UserUpdateProfileHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[users.UserUpdateProfileHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	userUpdateReq.Email = strings.TrimSpace(userUpdateReq.Email)
	userUpdateReq.Nickname = strings.TrimSpace(userUpdateReq.Nickname)

	modifyProfileBasicInfo := false
	anythingUpdate := false
	userNew := &models.User{
		Uid:  user.Uid,
		Salt: user.Salt,
	}

	if userUpdateReq.Email != "" && userUpdateReq.Email != user.Email {
		if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_UPDATE_EMAIL) {
			return nil, errs.ErrNotPermittedToPerformThisAction
		}

		user.Email = userUpdateReq.Email
		userNew.Email = userUpdateReq.Email
		anythingUpdate = true
	}

	if userUpdateReq.Password != "" {
		if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_UPDATE_PASSWORD) {
			return nil, errs.ErrNotPermittedToPerformThisAction
		}

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
		modifyProfileBasicInfo = true
		anythingUpdate = true
	}

	if userUpdateReq.DefaultAccountId > 0 && userUpdateReq.DefaultAccountId != user.DefaultAccountId {
		accountMap, err := a.accounts.GetAccountsByAccountIds(c, uid, []int64{userUpdateReq.DefaultAccountId})

		if err != nil || len(accountMap) < 1 {
			return nil, errs.Or(err, errs.ErrUserDefaultAccountIsInvalid)
		}

		if _, exists := accountMap[userUpdateReq.DefaultAccountId]; !exists {
			log.Warnf(c, "[users.UserUpdateProfileHandler] account \"id:%d\" does not exist for user \"uid:%d\"", userUpdateReq.DefaultAccountId, uid)
			return nil, errs.ErrUserDefaultAccountIsInvalid
		}

		if accountMap[userUpdateReq.DefaultAccountId].Hidden {
			log.Warnf(c, "[users.UserUpdateProfileHandler] account \"id:%d\" is hidden of user \"uid:%d\"", userUpdateReq.DefaultAccountId, uid)
			return nil, errs.ErrUserDefaultAccountIsHidden
		}

		user.DefaultAccountId = userUpdateReq.DefaultAccountId
		userNew.DefaultAccountId = userUpdateReq.DefaultAccountId
		modifyProfileBasicInfo = true
		anythingUpdate = true
	}

	if userUpdateReq.TransactionEditScope != nil && *userUpdateReq.TransactionEditScope != user.TransactionEditScope {
		user.TransactionEditScope = *userUpdateReq.TransactionEditScope
		userNew.TransactionEditScope = *userUpdateReq.TransactionEditScope
		modifyProfileBasicInfo = true
		anythingUpdate = true
	} else {
		userNew.TransactionEditScope = models.TRANSACTION_EDIT_SCOPE_INVALID
	}

	modifyUserLanguage := false

	if userUpdateReq.Language != user.Language {
		user.Language = userUpdateReq.Language
		userNew.Language = userUpdateReq.Language
		modifyUserLanguage = true
		modifyProfileBasicInfo = true
		anythingUpdate = true
	}

	if userUpdateReq.DefaultCurrency != "" && userUpdateReq.DefaultCurrency != user.DefaultCurrency {
		user.DefaultCurrency = userUpdateReq.DefaultCurrency
		userNew.DefaultCurrency = userUpdateReq.DefaultCurrency
		modifyProfileBasicInfo = true
		anythingUpdate = true
	}

	if userUpdateReq.FirstDayOfWeek != nil && *userUpdateReq.FirstDayOfWeek != user.FirstDayOfWeek {
		user.FirstDayOfWeek = *userUpdateReq.FirstDayOfWeek
		userNew.FirstDayOfWeek = *userUpdateReq.FirstDayOfWeek
		modifyProfileBasicInfo = true
		anythingUpdate = true
	} else {
		userNew.FirstDayOfWeek = core.WEEKDAY_INVALID
	}

	if userUpdateReq.FiscalYearStart != nil && *userUpdateReq.FiscalYearStart != user.FiscalYearStart {
		user.FiscalYearStart = *userUpdateReq.FiscalYearStart
		userNew.FiscalYearStart = *userUpdateReq.FiscalYearStart
		modifyProfileBasicInfo = true
		anythingUpdate = true
	} else {
		userNew.FiscalYearStart = core.FISCAL_YEAR_START_INVALID
	}

	if userUpdateReq.LongDateFormat != nil && *userUpdateReq.LongDateFormat != user.LongDateFormat {
		user.LongDateFormat = *userUpdateReq.LongDateFormat
		userNew.LongDateFormat = *userUpdateReq.LongDateFormat
		modifyProfileBasicInfo = true
		anythingUpdate = true
	} else {
		userNew.LongDateFormat = core.LONG_DATE_FORMAT_INVALID
	}

	if userUpdateReq.ShortDateFormat != nil && *userUpdateReq.ShortDateFormat != user.ShortDateFormat {
		user.ShortDateFormat = *userUpdateReq.ShortDateFormat
		userNew.ShortDateFormat = *userUpdateReq.ShortDateFormat
		modifyProfileBasicInfo = true
		anythingUpdate = true
	} else {
		userNew.ShortDateFormat = core.SHORT_DATE_FORMAT_INVALID
	}

	if userUpdateReq.LongTimeFormat != nil && *userUpdateReq.LongTimeFormat != user.LongTimeFormat {
		user.LongTimeFormat = *userUpdateReq.LongTimeFormat
		userNew.LongTimeFormat = *userUpdateReq.LongTimeFormat
		modifyProfileBasicInfo = true
		anythingUpdate = true
	} else {
		userNew.LongTimeFormat = core.LONG_TIME_FORMAT_INVALID
	}

	if userUpdateReq.ShortTimeFormat != nil && *userUpdateReq.ShortTimeFormat != user.ShortTimeFormat {
		user.ShortTimeFormat = *userUpdateReq.ShortTimeFormat
		userNew.ShortTimeFormat = *userUpdateReq.ShortTimeFormat
		modifyProfileBasicInfo = true
		anythingUpdate = true
	} else {
		userNew.ShortTimeFormat = core.SHORT_TIME_FORMAT_INVALID
	}

	if userUpdateReq.DecimalSeparator != nil && *userUpdateReq.DecimalSeparator != user.DecimalSeparator {
		user.DecimalSeparator = *userUpdateReq.DecimalSeparator
		userNew.DecimalSeparator = *userUpdateReq.DecimalSeparator
		modifyProfileBasicInfo = true
		anythingUpdate = true
	} else {
		userNew.DecimalSeparator = core.DECIMAL_SEPARATOR_INVALID
	}

	if userUpdateReq.DigitGroupingSymbol != nil && *userUpdateReq.DigitGroupingSymbol != user.DigitGroupingSymbol {
		user.DigitGroupingSymbol = *userUpdateReq.DigitGroupingSymbol
		userNew.DigitGroupingSymbol = *userUpdateReq.DigitGroupingSymbol
		modifyProfileBasicInfo = true
		anythingUpdate = true
	} else {
		userNew.DigitGroupingSymbol = core.DIGIT_GROUPING_SYMBOL_INVALID
	}

	if userUpdateReq.DigitGrouping != nil && *userUpdateReq.DigitGrouping != user.DigitGrouping {
		user.DigitGrouping = *userUpdateReq.DigitGrouping
		userNew.DigitGrouping = *userUpdateReq.DigitGrouping
		modifyProfileBasicInfo = true
		anythingUpdate = true
	} else {
		userNew.DigitGrouping = core.DIGIT_GROUPING_TYPE_INVALID
	}

	if userUpdateReq.CurrencyDisplayType != nil && *userUpdateReq.CurrencyDisplayType != user.CurrencyDisplayType {
		user.CurrencyDisplayType = *userUpdateReq.CurrencyDisplayType
		userNew.CurrencyDisplayType = *userUpdateReq.CurrencyDisplayType
		modifyProfileBasicInfo = true
		anythingUpdate = true
	} else {
		userNew.CurrencyDisplayType = core.CURRENCY_DISPLAY_TYPE_INVALID
	}

	if userUpdateReq.ExpenseAmountColor != nil && *userUpdateReq.ExpenseAmountColor != user.ExpenseAmountColor {
		user.ExpenseAmountColor = *userUpdateReq.ExpenseAmountColor
		userNew.ExpenseAmountColor = *userUpdateReq.ExpenseAmountColor
		modifyProfileBasicInfo = true
		anythingUpdate = true
	} else {
		userNew.ExpenseAmountColor = models.AMOUNT_COLOR_TYPE_INVALID
	}

	if userUpdateReq.IncomeAmountColor != nil && *userUpdateReq.IncomeAmountColor != user.IncomeAmountColor {
		user.IncomeAmountColor = *userUpdateReq.IncomeAmountColor
		userNew.IncomeAmountColor = *userUpdateReq.IncomeAmountColor
		modifyProfileBasicInfo = true
		anythingUpdate = true
	} else {
		userNew.IncomeAmountColor = models.AMOUNT_COLOR_TYPE_INVALID
	}

	if modifyProfileBasicInfo && user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_UPDATE_PROFILE_BASIC_INFO) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	if modifyUserLanguage || userNew.DecimalSeparator != core.DECIMAL_SEPARATOR_INVALID || userNew.DigitGroupingSymbol != core.DIGIT_GROUPING_SYMBOL_INVALID {
		decimalSeparator := userNew.DecimalSeparator
		digitGroupingSymbol := userNew.DigitGroupingSymbol

		if userNew.DecimalSeparator == core.DECIMAL_SEPARATOR_INVALID {
			decimalSeparator = user.DecimalSeparator
		}

		if userNew.DigitGroupingSymbol == core.DIGIT_GROUPING_SYMBOL_INVALID {
			digitGroupingSymbol = user.DigitGroupingSymbol
		}

		locale := user.Language

		if modifyUserLanguage {
			locale = userNew.Language
		}

		if locale == "" {
			locale = c.GetClientLocale()
		}

		if locales.IsDecimalSeparatorEqualsDigitGroupingSymbol(decimalSeparator, digitGroupingSymbol, locale) {
			return nil, errs.ErrDecimalSeparatorAndDigitGroupingSymbolCannotBeEqual
		}
	}

	if !anythingUpdate {
		return nil, errs.ErrNothingWillBeUpdated
	}

	keyProfileUpdated, emailSetToUnverified, err := a.users.UpdateUser(c, userNew, modifyUserLanguage)

	if err != nil {
		log.Errorf(c, "[users.UserUpdateProfileHandler] failed to update user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if emailSetToUnverified {
		user.EmailVerified = false
	}

	log.Infof(c, "[users.UserUpdateProfileHandler] user \"uid:%d\" has updated successfully", user.Uid)

	resp := &models.UserProfileUpdateResponse{
		User: a.GetUserBasicInfo(user),
	}

	if emailSetToUnverified && a.CurrentConfig().EnableUserVerifyEmail && a.CurrentConfig().EnableSMTP {
		err = a.tokens.DeleteTokensByType(c, uid, core.USER_TOKEN_TYPE_EMAIL_VERIFY)

		if err != nil {
			log.Errorf(c, "[users.UserUpdateProfileHandler] failed to revoke old email verify tokens for user \"uid:%d\", because %s", user.Uid, err.Error())
		} else {
			token, _, err := a.tokens.CreateEmailVerifyToken(c, user)

			if err != nil {
				log.Errorf(c, "[users.UserUpdateProfileHandler] failed to create email verify token for user \"uid:%d\", because %s", user.Uid, err.Error())
			} else {
				go func() {
					err = a.users.SendVerifyEmail(user, token, c.GetClientLocale())

					if err != nil {
						log.Warnf(c, "[users.UserUpdateProfileHandler] cannot send verify email to \"%s\", because %s", user.Email, err.Error())
					}
				}()
			}
		}
	}

	if keyProfileUpdated {
		now := time.Now().Unix()
		err = a.tokens.DeleteTokensBeforeTime(c, uid, now)

		if err == nil {
			log.Infof(c, "[users.UserUpdateProfileHandler] revoke old tokens before unix time \"%d\" for user \"uid:%d\"", now, user.Uid)
		} else {
			log.Warnf(c, "[users.UserUpdateProfileHandler] failed to revoke old tokens for user \"uid:%d\", because %s", user.Uid, err.Error())
		}

		token, claims, err := a.tokens.CreateToken(c, user)

		if err != nil {
			log.Warnf(c, "[users.UserUpdateProfileHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
			return resp, nil
		}

		resp.NewToken = token
		c.SetTextualToken(token)
		c.SetTokenClaims(claims)

		log.Infof(c, "[users.UserUpdateProfileHandler] user \"uid:%d\" token refreshed, new token will be expired at %d", user.Uid, claims.ExpiresAt)

		return resp, nil
	}

	return resp, nil
}

// UserUpdateAvatarHandler saves user avatar by request parameters for current user
func (a *UsersApi) UserUpdateAvatarHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[users.UserUpdateAvatarHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_UPDATE_AVATAR) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	form, err := c.MultipartForm()

	if err != nil {
		log.Errorf(c, "[users.UserUpdateAvatarHandler] failed to get multi-part form data for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrParameterInvalid
	}

	avatarFiles := form.File["avatar"]

	if len(avatarFiles) < 1 {
		log.Warnf(c, "[users.UserUpdateAvatarHandler] there is no user avatar in request for user \"uid:%d\"", user.Uid)
		return nil, errs.ErrNoUserAvatar
	}

	if avatarFiles[0].Size < 1 {
		log.Warnf(c, "[users.UserUpdateAvatarHandler] the size of user avatar in request is zero for user \"uid:%d\"", user.Uid)
		return nil, errs.ErrUserAvatarIsEmpty
	}

	if avatarFiles[0].Size > int64(a.CurrentConfig().MaxAvatarFileSize) {
		log.Warnf(c, "[users.UserUpdateAvatarHandler] the upload file size \"%d\" exceeds the maximum size \"%d\" of user avatar for user \"uid:%d\"", avatarFiles[0].Size, a.CurrentConfig().MaxAvatarFileSize, uid)
		return nil, errs.ErrExceedMaxUserAvatarFileSize
	}

	fileExtension := utils.GetFileNameExtension(avatarFiles[0].Filename)

	if utils.GetImageContentType(fileExtension) == "" {
		log.Warnf(c, "[users.UserUpdateAvatarHandler] the file extension \"%s\" of user avatar in request is not supported for user \"uid:%d\"", fileExtension, user.Uid)
		return nil, errs.ErrImageTypeNotSupported
	}

	avatarFile, err := avatarFiles[0].Open()

	if err != nil {
		log.Errorf(c, "[users.UserUpdateAvatarHandler] failed to get avatar file from request for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	err = a.users.UpdateUserAvatar(c, user.Uid, avatarFile, fileExtension, user.CustomAvatarType)

	if err != nil {
		log.Errorf(c, "[users.UserUpdateAvatarHandler] failed to update avatar for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	user.CustomAvatarType = fileExtension
	userResp := a.getUserProfileResponse(user)
	return userResp, nil
}

// UserRemoveAvatarHandler removes user avatar by request parameters for current user
func (a *UsersApi) UserRemoveAvatarHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[users.UserRemoveAvatarHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_UPDATE_AVATAR) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	if user.CustomAvatarType == "" {
		return nil, errs.ErrNothingWillBeUpdated
	}

	err = a.users.RemoveUserAvatar(c, user.Uid, user.CustomAvatarType)

	if err != nil {
		log.Errorf(c, "[users.UserRemoveAvatarHandler] failed to remove avatar for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	user.CustomAvatarType = ""
	userResp := a.getUserProfileResponse(user)
	return userResp, nil
}

// UserSendVerifyEmailByUnloginUserHandler sends unlogin user verify email
func (a *UsersApi) UserSendVerifyEmailByUnloginUserHandler(c *core.WebContext) (any, *errs.Error) {
	if !a.CurrentConfig().EnableUserVerifyEmail {
		return nil, errs.ErrEmailValidationNotAllowed
	}

	var userResendVerifyEmailReq models.UserResendVerifyEmailRequest
	err := c.ShouldBindJSON(&userResendVerifyEmailReq)

	user, err := a.users.GetUserByEmail(c, userResendVerifyEmailReq.Email)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[users.UserSendVerifyEmailByUnloginUserHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if !a.users.IsPasswordEqualsUserPassword(userResendVerifyEmailReq.Password, user) {
		log.Warnf(c, "[users.UserSendVerifyEmailByUnloginUserHandler] request password not equals to the user password")
		return nil, errs.ErrUserPasswordWrong
	}

	if user.Disabled {
		log.Warnf(c, "[users.UserSendVerifyEmailByUnloginUserHandler] user \"uid:%d\" is disabled", user.Uid)
		return nil, errs.ErrUserIsDisabled
	}

	if user.EmailVerified {
		log.Warnf(c, "[users.UserSendVerifyEmailByUnloginUserHandler] user \"uid:%d\" email has been verified", user.Uid)
		return nil, errs.ErrEmailIsVerified
	}

	if !a.CurrentConfig().EnableSMTP {
		return nil, errs.ErrSMTPServerNotEnabled
	}

	token, _, err := a.tokens.CreateEmailVerifyToken(c, user)

	if err != nil {
		log.Errorf(c, "[users.UserSendVerifyEmailByUnloginUserHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrTokenGenerating
	}

	go func() {
		err = a.users.SendVerifyEmail(user, token, c.GetClientLocale())

		if err != nil {
			log.Warnf(c, "[users.UserSendVerifyEmailByUnloginUserHandler] cannot send email to \"%s\", because %s", user.Email, err.Error())
		}
	}()

	return true, nil
}

// UserSendVerifyEmailByLoginedUserHandler sends logined user verify email
func (a *UsersApi) UserSendVerifyEmailByLoginedUserHandler(c *core.WebContext) (any, *errs.Error) {
	if !a.CurrentConfig().EnableUserVerifyEmail {
		return nil, errs.ErrEmailValidationNotAllowed
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[users.UserSendVerifyEmailByLoginedUserHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if user.EmailVerified {
		log.Warnf(c, "[users.UserSendVerifyEmailByLoginedUserHandler] user \"uid:%d\" email has been verified", user.Uid)
		return nil, errs.ErrEmailIsVerified
	}

	if !a.CurrentConfig().EnableSMTP {
		return nil, errs.ErrSMTPServerNotEnabled
	}

	token, _, err := a.tokens.CreateEmailVerifyToken(c, user)

	if err != nil {
		log.Errorf(c, "[users.UserSendVerifyEmailByLoginedUserHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrTokenGenerating
	}

	go func() {
		err = a.users.SendVerifyEmail(user, token, c.GetClientLocale())

		if err != nil {
			log.Warnf(c, "[users.UserSendVerifyEmailByLoginedUserHandler] cannot send email to \"%s\", because %s", user.Email, err.Error())
		}
	}()

	return true, nil
}

// UserGetAvatarHandler returns user avatar data for current user
func (a *UsersApi) UserGetAvatarHandler(c *core.WebContext) ([]byte, string, *errs.Error) {
	fileName := c.Param("fileName")
	fileExtension := utils.GetFileNameExtension(fileName)
	contentType := utils.GetImageContentType(fileExtension)

	if contentType == "" {
		return nil, "", errs.ErrImageTypeNotSupported
	}

	uid := c.GetCurrentUid()
	fileBaseName := utils.GetFileNameWithoutExtension(fileName)

	if utils.Int64ToString(uid) != fileBaseName {
		log.Warnf(c, "[users.UserGetAvatarHandler] cannot get other user avatar \"uid:%s\" for user \"uid:%d\"", fileBaseName, uid)
		return nil, "", errs.ErrUserIdInvalid
	}

	avatarData, err := a.users.GetUserAvatar(c, uid, fileExtension)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[users.UserGetAvatarHandler] failed to get user avatar, because %s", err.Error())
		}

		return nil, "", errs.Or(err, errs.ErrOperationFailed)
	}

	return avatarData, contentType, nil
}

func (a *UsersApi) getUserProfileResponse(user *models.User) *models.UserProfileResponse {
	return user.ToUserProfileResponse(a.GetUserBasicInfo(user))
}
