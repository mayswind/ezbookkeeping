package api

import (
	"fmt"
	"sort"

	"github.com/mayswind/ezbookkeeping/pkg/avatars"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
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
func (a *ApiUsingConfig) GetTransactionPictureInfoResponse(pictureInfo *models.TransactionPictureInfo) *models.TransactionPictureInfoBasicResponse {
	originalUrl := fmt.Sprintf(internalTransactionPictureUrlFormat, a.CurrentConfig().RootUrl, pictureInfo.PictureId, pictureInfo.PictureExtension)
	return pictureInfo.ToTransactionPictureInfoBasicResponse(originalUrl)
}

// GetTransactionPictureInfoResponseList returns the view-object list of transaction picture basic info according to the transaction picture model
func (a *ApiUsingConfig) GetTransactionPictureInfoResponseList(pictureInfos []*models.TransactionPictureInfo) models.TransactionPictureInfoBasicResponseSlice {
	pictureInfoResps := make(models.TransactionPictureInfoBasicResponseSlice, len(pictureInfos))

	for i := 0; i < len(pictureInfos); i++ {
		pictureInfoResps[i] = a.GetTransactionPictureInfoResponse(pictureInfos[i])
	}

	sort.Sort(pictureInfoResps)

	return pictureInfoResps
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
	ApiUsingConfig
	container *duplicatechecker.DuplicateCheckerContainer
}

// GetSubmissionRemark returns whether the same submission has been processed and related remark by the current duplicate checker
func (a *ApiUsingDuplicateChecker) GetSubmissionRemark(checkerType duplicatechecker.DuplicateCheckerType, uid int64, identification string) (bool, string) {
	return a.container.GetSubmissionRemark(checkerType, uid, identification)
}

// SetSubmissionRemarkIfEnable saves the identification and remark by the current duplicate checker if the duplicate submission check is enabled
func (a *ApiUsingDuplicateChecker) SetSubmissionRemarkIfEnable(checkerType duplicatechecker.DuplicateCheckerType, uid int64, identification string, remark string) {
	if a.CurrentConfig().EnableDuplicateSubmissionsCheck {
		a.container.SetSubmissionRemark(checkerType, uid, identification, remark)
	}
}

// RemoveSubmissionRemarkIfEnable removes the identification and remark by the current duplicate checker if the duplicate submission check is enabled
func (a *ApiUsingDuplicateChecker) RemoveSubmissionRemarkIfEnable(checkerType duplicatechecker.DuplicateCheckerType, uid int64, identification string) {
	if a.CurrentConfig().EnableDuplicateSubmissionsCheck {
		a.container.RemoveSubmissionRemark(checkerType, uid, identification)
	}
}

// CheckFailureCount returns whether the failure count of the specified IP and user has reached the limit and increases the failure count
func (a *ApiUsingDuplicateChecker) CheckFailureCount(c *core.WebContext, uid int64) error {
	if a.CurrentConfig().MaxFailuresPerIpPerMinute > 0 {
		clientIp := c.ClientIP()
		ipFailureCount := a.container.GetFailureCount(clientIp)

		if ipFailureCount >= a.CurrentConfig().MaxFailuresPerIpPerMinute {
			log.Warnf(c, "[base.CheckFailureCount] operation failure via IP \"%s\", current failure count: %d reached the limit", clientIp, ipFailureCount)
			return errs.ErrFailureCountLimitReached
		}
	}

	if a.CurrentConfig().MaxFailuresPerUserPerMinute > 0 && uid > 0 {
		uidFailureCount := a.container.GetFailureCount(utils.Int64ToString(uid))

		if uidFailureCount >= a.CurrentConfig().MaxFailuresPerUserPerMinute {
			log.Warnf(c, "[base.CheckFailureCount] operation failure via uid \"%d\", current failure count: %d reached the limit", uid, uidFailureCount)
			return errs.ErrFailureCountLimitReached
		}
	}

	return nil
}

// CheckAndIncreaseFailureCount returns whether the failure count of the specified IP and user has reached the limit and increases the failure count
func (a *ApiUsingDuplicateChecker) CheckAndIncreaseFailureCount(c *core.WebContext, uid int64) error {
	clientIp := c.ClientIP()
	ipFailureCount := uint32(0)
	uidFailureCount := uint32(0)

	if a.CurrentConfig().MaxFailuresPerIpPerMinute > 0 {
		ipFailureCount = a.container.GetFailureCount(clientIp)
	}

	if a.CurrentConfig().MaxFailuresPerUserPerMinute > 0 && uid > 0 {
		uidFailureCount = a.container.GetFailureCount(utils.Int64ToString(uid))
	}

	if a.CurrentConfig().MaxFailuresPerIpPerMinute > 0 && ipFailureCount < a.CurrentConfig().MaxFailuresPerIpPerMinute {
		log.Warnf(c, "[base.CheckAndIncreaseFailureCount] operation failure via IP \"%s\", previous failure count: %d", clientIp, ipFailureCount)
		a.container.IncreaseFailureCount(clientIp)
	}

	if a.CurrentConfig().MaxFailuresPerUserPerMinute > 0 && uid > 0 && uidFailureCount < a.CurrentConfig().MaxFailuresPerUserPerMinute {
		log.Warnf(c, "[base.CheckAndIncreaseFailureCount] operation failure via uid \"%d\", previous failure count: %d", uid, uidFailureCount)
		a.container.IncreaseFailureCount(utils.Int64ToString(uid))
	}

	if a.CurrentConfig().MaxFailuresPerIpPerMinute > 0 && ipFailureCount >= a.CurrentConfig().MaxFailuresPerIpPerMinute {
		log.Warnf(c, "[base.CheckAndIncreaseFailureCount] operation failure via IP \"%s\", current failure count: %d reached the limit", clientIp, ipFailureCount)
		return errs.ErrFailureCountLimitReached
	}

	if a.CurrentConfig().MaxFailuresPerUserPerMinute > 0 && uid > 0 && uidFailureCount >= a.CurrentConfig().MaxFailuresPerUserPerMinute {
		log.Warnf(c, "[base.CheckAndIncreaseFailureCount] operation failure via uid \"%d\", current failure count: %d reached the limit", uid, uidFailureCount)
		return errs.ErrFailureCountLimitReached
	}

	return nil
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
