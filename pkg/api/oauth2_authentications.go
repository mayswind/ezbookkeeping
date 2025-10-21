package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/locales"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const oauth2CallbackPageUrlSuccessFormat = "%sdesktop/#/oauth2_callback?platform=%s&provider=%s&token=%s"
const oauth2CallbackPageUrlNeedVerifyFormat = "%sdesktop/#/oauth2_callback?platform=%s&provider=%s&userName=%s&token=%s"
const oauth2CallbackPageUrlFailedFormat = "%sdesktop/#/oauth2_callback?errorCode=%d&errorMessage=%s"

// OAuth2AuthenticationApi represents OAuth 2.0 authorization api
type OAuth2AuthenticationApi struct {
	ApiUsingConfig
	ApiUsingDuplicateChecker
	users             *services.UserService
	tokens            *services.TokenService
	userExternalAuths *services.UserExternalAuthService
}

// Initialize a OAuth 2.0 authentication api singleton instance
var (
	OAuth2Authentications = &OAuth2AuthenticationApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		ApiUsingDuplicateChecker: ApiUsingDuplicateChecker{
			ApiUsingConfig: ApiUsingConfig{
				container: settings.Container,
			},
			container: duplicatechecker.Container,
		},
		users:             services.Users,
		tokens:            services.Tokens,
		userExternalAuths: services.UserExternalAuths,
	}
)

// LoginHandler handles user login request via OAuth 2.0
func (a *OAuth2AuthenticationApi) LoginHandler(c *core.WebContext) (string, *errs.Error) {
	var oauth2LoginReq models.OAuth2LoginRequest
	err := c.ShouldBindQuery(&oauth2LoginReq)

	if err != nil {
		log.Warnf(c, "[oauth2_authentications.LoginHandler] parse request failed, because %s", err.Error())
		return "", errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if oauth2LoginReq.Platform != "mobile" && oauth2LoginReq.Platform != "desktop" {
		return "", errs.ErrInvalidOAuth2LoginRequest
	}

	state := fmt.Sprintf("%s|%s", oauth2LoginReq.Platform, oauth2LoginReq.ClientSessionId)
	remark := ""

	if a.CurrentConfig().EnableDuplicateSubmissionsCheck {
		found := false
		found, remark = a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_OAUTH2_REDIRECT, 0, oauth2LoginReq.ClientSessionId)

		if found {
			log.Errorf(c, "[oauth2_authentications.LoginHandler] another oauth 2.0 state \"%s\" has been processing for client session id \"%s\"", remark, oauth2LoginReq.ClientSessionId)
			return "", errs.ErrRepeatedRequest
		}

		randomString, err := utils.GetRandomNumberOrLowercaseLetter(32)

		if err != nil {
			log.Errorf(c, "[oauth2_authentications.LoginHandler] failed to generate random string for oauth 2.0 state, because %s", err.Error())
			return "", errs.ErrSystemError
		}

		remark = fmt.Sprintf("%s|%s|%s", oauth2LoginReq.Platform, oauth2LoginReq.ClientSessionId, randomString)
		state = fmt.Sprintf("%s|%s|%s", oauth2LoginReq.Platform, oauth2LoginReq.ClientSessionId, utils.MD5EncodeToString([]byte(remark)))
	}

	redirectUrl, err := oauth2.GetOAuth2AuthUrl(c, state)

	if err != nil {
		log.Errorf(c, "[oauth2_authentications.LoginHandler] failed to get oauth 2.0 auth url, because %s", err.Error())
		return "", errs.Or(err, errs.ErrSystemError)
	}

	a.SetSubmissionRemarkWithCustomExpirationIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_OAUTH2_REDIRECT, 0, oauth2LoginReq.ClientSessionId, remark, a.CurrentConfig().OAuth2StateExpiredTimeDuration)

	return redirectUrl, nil
}

// CallbackHandler handles OAuth 2.0 callback request
func (a *OAuth2AuthenticationApi) CallbackHandler(c *core.WebContext) (string, *errs.Error) {
	var oauth2CallbackReq models.OAuth2CallbackRequest
	err := c.ShouldBindQuery(&oauth2CallbackReq)

	if err != nil {
		log.Warnf(c, "[oauth2_authentications.CallbackHandler] parse request failed, because %s", err.Error())
		return a.redirectToFailedCallbackPage(c, errs.NewIncompleteOrIncorrectSubmissionError(err))
	}

	if oauth2CallbackReq.State == "" {
		return a.redirectToFailedCallbackPage(c, errs.ErrMissingOAuth2State)
	}

	if oauth2CallbackReq.Code == "" {
		return a.redirectToFailedCallbackPage(c, errs.ErrMissingOAuth2Code)
	}

	platform := ""
	clientSessionId := ""

	stateParts := strings.Split(oauth2CallbackReq.State, "|")

	if len(stateParts) >= 2 {
		platform = stateParts[0]
		clientSessionId = stateParts[1]
	} else {
		return a.redirectToFailedCallbackPage(c, errs.ErrInvalidOAuth2State)
	}

	if platform != "mobile" && platform != "desktop" {
		return a.redirectToFailedCallbackPage(c, errs.ErrInvalidOAuth2LoginRequest)
	}

	if a.CurrentConfig().EnableDuplicateSubmissionsCheck {
		found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_OAUTH2_REDIRECT, 0, clientSessionId)

		if !found {
			log.Errorf(c, "[oauth2_authentications.CallbackHandler] cannot find oauth 2.0 state in duplicate checker for client session id \"%s\"", clientSessionId)
			return a.redirectToFailedCallbackPage(c, errs.ErrInvalidOAuth2Callback)
		}

		remarkParts := strings.Split(remark, "|")

		if len(remarkParts) != 3 || remarkParts[0] != platform || remarkParts[1] != clientSessionId {
			log.Errorf(c, "[oauth2_authentications.CallbackHandler] invalid oauth 2.0 state \"%s\" in duplicate checker for client session id \"%s\"", remark, clientSessionId)
			return a.redirectToFailedCallbackPage(c, errs.ErrInvalidOAuth2State)
		}

		expectedState := fmt.Sprintf("%s|%s|%s", platform, clientSessionId, remarkParts[2])
		expectedState = fmt.Sprintf("%s|%s|%s", platform, clientSessionId, utils.MD5EncodeToString([]byte(expectedState)))

		if oauth2CallbackReq.State != expectedState {
			log.Errorf(c, "[oauth2_authentications.CallbackHandler] mismatched random string in oauth 2.0 state, expected \"%s\", got \"%s\"", expectedState, oauth2CallbackReq.State)
			return a.redirectToFailedCallbackPage(c, errs.ErrInvalidOAuth2State)
		}

		a.RemoveSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_OAUTH2_REDIRECT, 0, clientSessionId)
	}

	oauth2Token, err := oauth2.GetOAuth2Token(c, oauth2CallbackReq.Code)

	if err != nil {
		log.Errorf(c, "[oauth2_authentications.CallbackHandler] failed to retrieve oauth 2.0 token, because %s", err.Error())
		return a.redirectToFailedCallbackPage(c, errs.Or(err, errs.ErrCannotRetrieveOAuth2Token))
	}

	oauth2UserInfo, err := oauth2.GetOAuth2UserInfo(c, oauth2Token)

	if err != nil {
		log.Errorf(c, "[oauth2_authentications.CallbackHandler] failed to retrieve oauth 2.0 user info, because %s", err.Error())
		return a.redirectToFailedCallbackPage(c, errs.Or(err, errs.ErrInvalidOAuth2Token))
	}

	if oauth2UserInfo == nil {
		log.Errorf(c, "[oauth2_authentications.CallbackHandler] failed to retrieve oauth 2.0 user info, because user info is nil")
		return a.redirectToFailedCallbackPage(c, errs.ErrCannotRetrieveUserInfo)
	}

	if oauth2UserInfo.UserName == "" || oauth2UserInfo.Email == "" {
		log.Errorf(c, "[oauth2_authentications.CallbackHandler] invalid oauth 2.0 user info, userName: %s, email: %s", oauth2UserInfo.UserName, oauth2UserInfo.Email)
		return a.redirectToFailedCallbackPage(c, errs.ErrCannotRetrieveUserInfo)
	}

	userExternalAuthType := oauth2.GetExternalUserAuthType()
	var userExternalAuth *models.UserExternalAuth

	if a.CurrentConfig().OAuth2UserIdentifier == settings.OAuth2UserIdentifierEmail {
		userExternalAuth, err = a.userExternalAuths.GetUserExternalAuthByExternalEmail(c, oauth2UserInfo.Email, userExternalAuthType)
	} else if a.CurrentConfig().OAuth2UserIdentifier == settings.OAuth2UserIdentifierUsername {
		userExternalAuth, err = a.userExternalAuths.GetUserExternalAuthByExternalUserName(c, oauth2UserInfo.UserName, userExternalAuthType)
	} else {
		userExternalAuth, err = a.userExternalAuths.GetUserExternalAuthByExternalEmail(c, oauth2UserInfo.Email, userExternalAuthType)
	}

	if err != nil && !errors.Is(err, errs.ErrUserExternalAuthNotFound) {
		log.Errorf(c, "[oauth2_authentications.CallbackHandler] failed to get user external auth, because %s", err.Error())
		return a.redirectToFailedCallbackPage(c, errs.Or(err, errs.ErrOperationFailed))
	}

	var user *models.User

	if err == nil { // user already bound to external auth, redirect to success page
		user, err = a.users.GetUserById(c, userExternalAuth.Uid)

		if err != nil {
			log.Errorf(c, "[oauth2_authentications.CallbackHandler] failed to get user by id %d, because %s", userExternalAuth.Uid, err.Error())
			return a.redirectToFailedCallbackPage(c, errs.Or(err, errs.ErrOperationFailed))
		}
	} else { // errors.Is(err, errs.ErrUserExternalAuthNotFound) // user not bound to external auth, try to bind or register new user
		if a.CurrentConfig().OAuth2UserIdentifier == settings.OAuth2UserIdentifierEmail {
			user, err = a.users.GetUserByEmail(c, oauth2UserInfo.Email)
		} else if a.CurrentConfig().OAuth2UserIdentifier == settings.OAuth2UserIdentifierUsername {
			user, err = a.users.GetUserByUsername(c, oauth2UserInfo.UserName)
		} else {
			user, err = a.users.GetUserByEmail(c, oauth2UserInfo.Email)
		}

		if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
			log.Errorf(c, "[oauth2_authentications.CallbackHandler] failed to get user, because %s", err.Error())
			return a.redirectToFailedCallbackPage(c, errs.Or(err, errs.ErrOperationFailed))
		}

		if user == nil && a.CurrentConfig().EnableUserRegister && a.CurrentConfig().OAuth2AutoRegister {
			userName := strings.TrimSpace(oauth2UserInfo.UserName)
			email := strings.TrimSpace(oauth2UserInfo.Email)
			nickName := strings.TrimSpace(oauth2UserInfo.NickName)
			languageCode := ""
			currencyCode := "USD"

			if _, exists := locales.AllLanguages[oauth2UserInfo.LanguageCode]; exists {
				languageCode = oauth2UserInfo.LanguageCode
			}

			if _, exists := validators.AllCurrencyNames[oauth2UserInfo.CurrencyCode]; exists {
				currencyCode = oauth2UserInfo.CurrencyCode
			}

			user = &models.User{
				Username:             userName,
				Email:                email,
				Nickname:             nickName,
				Password:             "",
				Language:             languageCode,
				DefaultCurrency:      currencyCode,
				FirstDayOfWeek:       oauth2UserInfo.FirstDayOfWeek,
				FiscalYearStart:      core.FISCAL_YEAR_START_DEFAULT,
				TransactionEditScope: models.TRANSACTION_EDIT_SCOPE_ALL,
				FeatureRestriction:   a.CurrentConfig().DefaultFeatureRestrictions,
			}

			err = a.users.CreateUser(c, user)

			if err != nil {
				log.Errorf(c, "[oauth2_authentications.CallbackHandler] failed to create user \"%s\", because %s", user.Username, err.Error())
				return a.redirectToFailedCallbackPage(c, errs.Or(err, errs.ErrOperationFailed))
			}

			log.Infof(c, "[oauth2_authentications.CallbackHandler] user \"%s\" has registered successfully, uid is %d", user.Username, user.Uid)

			userExternalAuth := &models.UserExternalAuth{
				Uid:              user.Uid,
				ExternalAuthType: userExternalAuthType,
				ExternalUsername: oauth2UserInfo.UserName,
				ExternalEmail:    oauth2UserInfo.Email,
			}

			err = a.userExternalAuths.CreateUserExternalAuth(c, userExternalAuth)

			if err != nil {
				log.Errorf(c, "[oauth2_authentications.CallbackHandler] failed to create user external auth for user \"uid:%d\", because %s", user.Uid, err.Error())
				return a.redirectToFailedCallbackPage(c, errs.Or(err, errs.ErrOperationFailed))
			}

			log.Infof(c, "[oauth2_authentications.CallbackHandler] user external auth has been created for user \"uid:%d\"", user.Uid)
		} else if user == nil {
			return a.redirectToFailedCallbackPage(c, errs.ErrOAuth2AutoRegistrationNotEnabled)
		}
	}

	if userExternalAuth == nil {
		tokenContext, err := json.Marshal(&models.OAuth2CallbackTokenContext{
			ExternalAuthType: userExternalAuthType,
			ExternalUsername: oauth2UserInfo.UserName,
			ExternalEmail:    oauth2UserInfo.Email,
		})

		if err != nil {
			log.Errorf(c, "[oauth2_authentications.CallbackHandler] failed to marshal oauth 2.0 callback verify token context, because %s", err.Error())
			return a.redirectToFailedCallbackPage(c, errs.ErrOperationFailed)
		}

		token, _, err := a.tokens.CreateOAuth2CallbackRequireVerifyToken(c, user, string(tokenContext))

		if err != nil {
			log.Errorf(c, "[oauth2_authentications.CallbackHandler] failed to create oauth 2.0 callback verify token, because %s", err.Error())
			return a.redirectToFailedCallbackPage(c, errs.ErrTokenGenerating)
		}

		return a.redirectToVerifyCallbackPage(c, platform, userExternalAuthType, user.Username, token)
	} else {
		tokenContext, err := json.Marshal(&models.OAuth2CallbackTokenContext{
			ExternalAuthType: userExternalAuthType,
		})

		if err != nil {
			log.Errorf(c, "[oauth2_authentications.CallbackHandler] failed to marshal oauth 2.0 callback token context, because %s", err.Error())
			return a.redirectToFailedCallbackPage(c, errs.ErrOperationFailed)
		}

		token, _, err := a.tokens.CreateOAuth2CallbackToken(c, user, string(tokenContext))

		if err != nil {
			log.Errorf(c, "[oauth2_authentications.CallbackHandler] failed to create oauth 2.0 callback token, because %s", err.Error())
			return a.redirectToFailedCallbackPage(c, errs.ErrTokenGenerating)
		}

		return a.redirectToSuccessCallbackPage(c, platform, userExternalAuthType, token)
	}
}

func (a *OAuth2AuthenticationApi) redirectToSuccessCallbackPage(c *core.WebContext, platform string, externalAuthType core.UserExternalAuthType, token string) (string, *errs.Error) {
	return fmt.Sprintf(oauth2CallbackPageUrlSuccessFormat, a.CurrentConfig().RootUrl, platform, externalAuthType, url.QueryEscape(token)), nil
}

func (a *OAuth2AuthenticationApi) redirectToVerifyCallbackPage(c *core.WebContext, platform string, externalAuthType core.UserExternalAuthType, userName string, token string) (string, *errs.Error) {
	return fmt.Sprintf(oauth2CallbackPageUrlNeedVerifyFormat, a.CurrentConfig().RootUrl, platform, externalAuthType, userName, url.QueryEscape(token)), nil
}

func (a *OAuth2AuthenticationApi) redirectToFailedCallbackPage(c *core.WebContext, err *errs.Error) (string, *errs.Error) {
	return fmt.Sprintf(oauth2CallbackPageUrlFailedFormat, a.CurrentConfig().RootUrl, err.Code(), url.QueryEscape(utils.GetDisplayErrorMessage(err))), nil
}
