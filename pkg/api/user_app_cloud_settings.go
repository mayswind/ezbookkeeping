package api

import (
	"encoding/json"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// UserApplicationCloudSettingsApi represents user application cloud settings api
type UserApplicationCloudSettingsApi struct {
	userAppCloudSettings *services.UserApplicationCloudSettingsService
	users                *services.UserService
}

// Initialize a user application cloud settings api singleton instance
var (
	UserApplicationCloudSettings = &UserApplicationCloudSettingsApi{
		userAppCloudSettings: services.UserApplicationCloudSettings,
		users:                services.Users,
	}
)

// ApplicationSettingsGetHandler returns application cloud settings of current user
func (a *UserApplicationCloudSettingsApi) ApplicationSettingsGetHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()

	userApplicationCloudSettings, err := a.userAppCloudSettings.GetUserApplicationCloudSettingsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[user_app_cloud_settings.ApplicationSettingsGetHandler] failed to get latest user application cloud settings for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if userApplicationCloudSettings == nil {
		return false, nil
	}

	applicationCloudSettingSlice := userApplicationCloudSettings.Settings

	if len(applicationCloudSettingSlice) < 1 {
		return false, nil
	}

	return applicationCloudSettingSlice, nil
}

// ApplicationSettingsUpdateHandler updates user application cloud settings by request parameters for current user
func (a *UserApplicationCloudSettingsApi) ApplicationSettingsUpdateHandler(c *core.WebContext) (any, *errs.Error) {
	var userAppCloudSettingUpdateReq models.UserApplicationCloudSettingsUpdateRequest
	err := c.ShouldBindJSON(&userAppCloudSettingUpdateReq)

	if err != nil {
		log.Warnf(c, "[user_app_cloud_settings.ApplicationSettingsUpdateHandler] parse request failed, because %s", err.Error())
		return false, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Warnf(c, "[user_app_cloud_settings.ApplicationSettingsUpdateHandler] failed to get user for user \"uid:%d\", because %s", uid, err.Error())
		}

		return false, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_SYNC_APPLICATION_SETTINGS) {
		return false, errs.ErrNotPermittedToPerformThisAction
	}

	userApplicationCloudSettings, err := a.userAppCloudSettings.GetUserApplicationCloudSettingsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[user_app_cloud_settings.ApplicationSettingsUpdateHandler] failed to get latest user application cloud settings for user \"uid:%d\", because %s", uid, err.Error())
		return false, errs.Or(err, errs.ErrOperationFailed)
	}

	oldApplicationCloudSettingsMap := make(map[string]models.ApplicationCloudSetting)

	if userApplicationCloudSettings != nil {
		for _, setting := range userApplicationCloudSettings.Settings {
			oldApplicationCloudSettingsMap[setting.SettingKey] = setting
		}
	}

	// Check if the full update settings are the same as the existing settings
	if userAppCloudSettingUpdateReq.FullUpdate {
		if len(userAppCloudSettingUpdateReq.Settings) == len(oldApplicationCloudSettingsMap) {
			needUpdate := false

			for _, setting := range userAppCloudSettingUpdateReq.Settings {
				oldSetting, exists := oldApplicationCloudSettingsMap[setting.SettingKey]

				if !exists || oldSetting.SettingValue != setting.SettingValue {
					needUpdate = true
					break
				}
			}

			if !needUpdate {
				return false, errs.ErrNothingWillBeUpdated
			}
		}
	} else { // Check if the partial update settings are the same as the existing settings or the settings to update are not set to sync
		needUpdate := true

		for _, setting := range userAppCloudSettingUpdateReq.Settings {
			cloudSetting, exists := oldApplicationCloudSettingsMap[setting.SettingKey]

			if !exists {
				needUpdate = false
				log.Infof(c, "[user_app_cloud_settings.ApplicationSettingsUpdateHandler] user application cloud setting key \"%s\" is not set to sync", setting.SettingKey)
			} else if cloudSetting.SettingValue == setting.SettingValue {
				needUpdate = false
				log.Infof(c, "[user_app_cloud_settings.ApplicationSettingsUpdateHandler] user application cloud setting key \"%s\" value \"%s\" is not changed, no need to update", setting.SettingKey, setting.SettingValue)
			}
		}

		if !needUpdate {
			log.Infof(c, "[user_app_cloud_settings.ApplicationSettingsUpdateHandler] no user application cloud settings need to update for user \"uid:%d\"", uid)
			return false, nil
		}
	}

	newApplicationCloudSettingsMap := make(map[string]models.ApplicationCloudSetting)
	var newApplicationCloudSettingSlice models.ApplicationCloudSettingSlice

	if userAppCloudSettingUpdateReq.FullUpdate {
		log.Infof(c, "[user_app_cloud_settings.ApplicationSettingsUpdateHandler] user \"uid:%d\" application cloud settings force update, will overwrite all existing settings", uid)
	} else {
		if len(oldApplicationCloudSettingsMap) > 0 {
			log.Infof(c, "[user_app_cloud_settings.ApplicationSettingsUpdateHandler] user \"uid:%d\" application cloud settings exists, try to merge it with request settings", uid)
			newApplicationCloudSettingsMap = oldApplicationCloudSettingsMap
		}
	}

	for _, setting := range userAppCloudSettingUpdateReq.Settings {
		newApplicationCloudSettingsMap[setting.SettingKey] = setting
	}

	for settingKey, setting := range newApplicationCloudSettingsMap {
		settingType, exists := models.ALL_ALLOWED_CLOUD_SYNC_APP_SETTING_KEY_TYPES[settingKey]

		if !exists {
			log.Warnf(c, "[user_app_cloud_settings.ApplicationSettingsUpdateHandler] user application cloud setting key \"%s\" is not supported to sync", settingKey)
			continue
		}

		if settingType == models.USER_APPLICATION_CLOUD_SETTING_TYPE_STRING {
			// Do Nothing
		} else if settingType == models.USER_APPLICATION_CLOUD_SETTING_TYPE_NUMBER {
			_, err := utils.StringToFloat64(setting.SettingValue)

			if err != nil {
				log.Warnf(c, "[user_app_cloud_settings.ApplicationSettingsUpdateHandler] user application cloud setting key \"%s\" has invalid number value \"%s\"", settingKey, setting.SettingValue)
				continue
			}
		} else if settingType == models.USER_APPLICATION_CLOUD_SETTING_TYPE_BOOLEAN {
			if setting.SettingValue != "true" && setting.SettingValue != "false" {
				log.Warnf(c, "[user_app_cloud_settings.ApplicationSettingsUpdateHandler] user application cloud setting key \"%s\" has invalid boolean value \"%s\"", settingKey, setting.SettingValue)
				continue
			}
		} else if settingType == models.USER_APPLICATION_CLOUD_SETTING_TYPE_STRING_BOOLEAN_MAP {
			var settingValueMap map[string]bool
			err := json.Unmarshal([]byte(setting.SettingValue), &settingValueMap)

			if err != nil {
				log.Warnf(c, "[user_app_cloud_settings.ApplicationSettingsUpdateHandler] user application cloud setting key \"%s\" has invalid map value \"%s\", because %s", settingKey, setting.SettingValue, err.Error())
				continue
			}
		} else {
			log.Warnf(c, "[user_app_cloud_settings.ApplicationSettingsUpdateHandler] user application cloud setting key \"%s\" has unknown type \"%s\"", settingKey, settingType)
			continue
		}

		newApplicationCloudSettingSlice = append(newApplicationCloudSettingSlice, setting)
	}

	err = a.userAppCloudSettings.UpdateUserApplicationCloudSettings(c, uid, newApplicationCloudSettingSlice)

	if err != nil {
		log.Errorf(c, "[user_app_cloud_settings.ApplicationSettingsUpdateHandler] failed to update user application cloud settings for user \"uid:%d\", because %s", uid, err.Error())
		return false, errs.Or(err, errs.ErrOperationFailed)
	}

	return true, nil
}

// ApplicationSettingsDisableHandler disabled user application cloud settings by request parameters for current user
func (a *UserApplicationCloudSettingsApi) ApplicationSettingsDisableHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Warnf(c, "[user_app_cloud_settings.ApplicationSettingsDisableHandler] failed to get user for user \"uid:%d\", because %s", uid, err.Error())
		}

		return false, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_SYNC_APPLICATION_SETTINGS) {
		return false, errs.ErrNotPermittedToPerformThisAction
	}

	err = a.userAppCloudSettings.ClearUserApplicationCloudSettings(c, uid)

	if err != nil {
		log.Errorf(c, "[user_app_cloud_settings.ApplicationSettingsDisableHandler] failed to clear user application cloud settings for user \"uid:%d\", because %s", uid, err.Error())
		return false, errs.Or(err, errs.ErrOperationFailed)
	}

	return true, nil
}
