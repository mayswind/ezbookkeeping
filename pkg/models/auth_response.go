package models

// AuthResponse returns a view-object of user authorization
type AuthResponse struct {
	Token                    string                        `json:"token"`
	MapProxyToken            string                        `json:"mapProxyToken,omitempty"`
	Need2FA                  bool                          `json:"need2FA"`
	User                     *UserBasicInfo                `json:"user,omitempty"`
	ApplicationCloudSettings *ApplicationCloudSettingSlice `json:"applicationCloudSettings,omitempty"`
	NotificationContent      string                        `json:"notificationContent,omitempty"`
}

// RegisterResponse returns a view-object of user register response
type RegisterResponse struct {
	AuthResponse
	NeedVerifyEmail       bool `json:"needVerifyEmail"`
	PresetCategoriesSaved bool `json:"presetCategoriesSaved"`
}
