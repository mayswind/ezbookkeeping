package api

import (
	"fmt"

	"github.com/mayswind/ezbookkeeping/pkg/avatars"
	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const internalTransactionPictureUrlFormat = "%spictures/%d.%s"

// ApiUsingConfig represents an api that need to use config
type ApiUsingConfig struct {
	container *settings.ConfigContainer
}

// CurrentConfig returns the current config
func (a *ApiUsingConfig) CurrentConfig() *settings.Config {
	return a.container.Current
}

// GetTransactionPictureInfoResponse returns the view-object of transaction picture basic info according to the transaction picture model
func (a *ApiUsingConfig) GetTransactionPictureInfoResponse(picture *models.TransactionPictureInfo) *models.TransactionPictureInfoBasicResponse {
	originalUrl := fmt.Sprintf(internalTransactionPictureUrlFormat, a.CurrentConfig().RootUrl, picture.PictureId, picture.PictureExtension)
	return picture.ToTransactionPictureInfoBasicResponse(originalUrl)
}

// GetAfterRegisterNotificationContent returns the notification content displayed each time users register
func (a *ApiUsingConfig) GetAfterRegisterNotificationContent(userLanguage string, clientLanguage string) string {
	language := userLanguage

	if language == "" {
		language = clientLanguage
	}

	if !a.container.Current.AfterRegisterNotification.Enabled {
		return ""
	}

	if multiLanguageContent, exists := a.container.Current.AfterRegisterNotification.MultiLanguageContent[language]; exists {
		return multiLanguageContent
	}

	return a.container.Current.AfterRegisterNotification.DefaultContent
}

// GetAfterLoginNotificationContent returns the notification content displayed each time users log in
func (a *ApiUsingConfig) GetAfterLoginNotificationContent(userLanguage string, clientLanguage string) string {
	language := userLanguage

	if language == "" {
		language = clientLanguage
	}

	if !a.container.Current.AfterLoginNotification.Enabled {
		return ""
	}

	if multiLanguageContent, exists := a.container.Current.AfterLoginNotification.MultiLanguageContent[language]; exists {
		return multiLanguageContent
	}

	return a.container.Current.AfterLoginNotification.DefaultContent
}

// GetAfterOpenNotificationContent returns the notification content displayed each time users open the app
func (a *ApiUsingConfig) GetAfterOpenNotificationContent(userLanguage string, clientLanguage string) string {
	language := userLanguage

	if language == "" {
		language = clientLanguage
	}

	if !a.container.Current.AfterOpenNotification.Enabled {
		return ""
	}

	if multiLanguageContent, exists := a.container.Current.AfterOpenNotification.MultiLanguageContent[language]; exists {
		return multiLanguageContent
	}

	return a.container.Current.AfterOpenNotification.DefaultContent
}

// ApiUsingDuplicateChecker represents an api that need to use duplicate checker
type ApiUsingDuplicateChecker struct {
	container *duplicatechecker.DuplicateCheckerContainer
}

// GetSubmissionRemark returns whether the same submission has been processed and related remark by the current duplicate checker
func (a *ApiUsingDuplicateChecker) GetSubmissionRemark(checkerType duplicatechecker.DuplicateCheckerType, uid int64, identification string) (bool, string) {
	return a.container.GetSubmissionRemark(checkerType, uid, identification)
}

// SetSubmissionRemark saves the identification and remark to in-memory cache by the current duplicate checker
func (a *ApiUsingDuplicateChecker) SetSubmissionRemark(checkerType duplicatechecker.DuplicateCheckerType, uid int64, identification string, remark string) {
	a.container.SetSubmissionRemark(checkerType, uid, identification, remark)
}

// ApiUsingAvatarProvider represents an api that need to use avatar provider
type ApiUsingAvatarProvider struct {
	container *avatars.AvatarProviderContainer
}

// GetAvatarUrl returns the avatar url by the current user avatar provider
func (a *ApiUsingAvatarProvider) GetAvatarUrl(user *models.User) string {
	return a.container.GetAvatarUrl(user)
}

// ApiWithUserInfo represents an api that can returns user info
type ApiWithUserInfo struct {
	ApiUsingConfig
	ApiUsingAvatarProvider
}

// GetUserBasicInfo returns the view-object of user basic info according to the user model
func (a *ApiWithUserInfo) GetUserBasicInfo(user *models.User) *models.UserBasicInfo {
	return user.ToUserBasicInfo(a.CurrentConfig().AvatarProvider, a.GetAvatarUrl(user))
}
