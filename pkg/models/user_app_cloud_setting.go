package models

import (
	"encoding/json"
)

type UserApplicationCloudSettingType string

const (
	USER_APPLICATION_CLOUD_SETTING_TYPE_STRING             UserApplicationCloudSettingType = "string"
	USER_APPLICATION_CLOUD_SETTING_TYPE_NUMBER             UserApplicationCloudSettingType = "number"
	USER_APPLICATION_CLOUD_SETTING_TYPE_BOOLEAN            UserApplicationCloudSettingType = "boolean"
	USER_APPLICATION_CLOUD_SETTING_TYPE_STRING_BOOLEAN_MAP UserApplicationCloudSettingType = "string_boolean_map"
)

var ALL_ALLOWED_CLOUD_SYNC_APP_SETTING_KEY_TYPES = map[string]UserApplicationCloudSettingType{
	// Basic Settings
	"showAccountBalance": USER_APPLICATION_CLOUD_SETTING_TYPE_BOOLEAN,
	// Overview Page
	"showAmountInHomePage":                        USER_APPLICATION_CLOUD_SETTING_TYPE_BOOLEAN,
	"timezoneUsedForStatisticsInHomePage":         USER_APPLICATION_CLOUD_SETTING_TYPE_NUMBER,
	"overviewAccountFilterInHomePage":             USER_APPLICATION_CLOUD_SETTING_TYPE_STRING_BOOLEAN_MAP,
	"overviewTransactionCategoryFilterInHomePage": USER_APPLICATION_CLOUD_SETTING_TYPE_STRING_BOOLEAN_MAP,
	// Transaction List Page
	"itemsCountInTransactionListPage":      USER_APPLICATION_CLOUD_SETTING_TYPE_NUMBER,
	"showTotalAmountInTransactionListPage": USER_APPLICATION_CLOUD_SETTING_TYPE_BOOLEAN,
	"showTagInTransactionListPage":         USER_APPLICATION_CLOUD_SETTING_TYPE_BOOLEAN,
	// Transaction Edit Page
	"autoSaveTransactionDraft":                                 USER_APPLICATION_CLOUD_SETTING_TYPE_STRING,
	"autoGetCurrentGeoLocation":                                USER_APPLICATION_CLOUD_SETTING_TYPE_BOOLEAN,
	"alwaysShowTransactionPicturesInMobileTransactionEditPage": USER_APPLICATION_CLOUD_SETTING_TYPE_BOOLEAN,
	// Account List Page
	"totalAmountExcludeAccountIds": USER_APPLICATION_CLOUD_SETTING_TYPE_STRING_BOOLEAN_MAP,
	// Exchange Rates Data Page
	"currencySortByInExchangeRatesPage": USER_APPLICATION_CLOUD_SETTING_TYPE_NUMBER,
	// Statistics Settings
	"statistics.defaultChartDataType":                 USER_APPLICATION_CLOUD_SETTING_TYPE_NUMBER,
	"statistics.defaultTimezoneType":                  USER_APPLICATION_CLOUD_SETTING_TYPE_NUMBER,
	"statistics.defaultAccountFilter":                 USER_APPLICATION_CLOUD_SETTING_TYPE_STRING_BOOLEAN_MAP,
	"statistics.defaultTransactionCategoryFilter":     USER_APPLICATION_CLOUD_SETTING_TYPE_STRING_BOOLEAN_MAP,
	"statistics.defaultSortingType":                   USER_APPLICATION_CLOUD_SETTING_TYPE_NUMBER,
	"statistics.defaultCategoricalChartType":          USER_APPLICATION_CLOUD_SETTING_TYPE_NUMBER,
	"statistics.defaultCategoricalChartDataRangeType": USER_APPLICATION_CLOUD_SETTING_TYPE_NUMBER,
	"statistics.defaultTrendChartType":                USER_APPLICATION_CLOUD_SETTING_TYPE_NUMBER,
	"statistics.defaultTrendChartDataRangeType":       USER_APPLICATION_CLOUD_SETTING_TYPE_NUMBER,
}

// UserApplicationCloudSetting represents user application cloud setting stored in database
type UserApplicationCloudSetting struct {
	Uid             int64                        `xorm:"PK"`
	Settings        ApplicationCloudSettingSlice `xorm:"BLOB"`
	UpdatedUnixTime int64
}

// UserApplicationCloudSettingsUpdateRequest represents all parameters of application cloud settings update request
type UserApplicationCloudSettingsUpdateRequest struct {
	Settings   ApplicationCloudSettingSlice `json:"settings"`
	FullUpdate bool                         `json:"fullUpdate"`
}

// ApplicationCloudSettingSlice represents the slice data structure of ApplicationCloudSetting
type ApplicationCloudSettingSlice []ApplicationCloudSetting

// ApplicationCloudSetting represents one application cloud setting
type ApplicationCloudSetting struct {
	SettingKey   string `json:"settingKey"`
	SettingValue string `json:"settingValue"`
}

// FromDB fills the fields from the data stored in database
func (s *ApplicationCloudSettingSlice) FromDB(data []byte) error {
	return json.Unmarshal(data, s)
}

// ToDB returns the actual stored data in database
func (s *ApplicationCloudSettingSlice) ToDB() ([]byte, error) {
	return json.Marshal(s)
}

// Len returns the count of items
func (s ApplicationCloudSettingSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s ApplicationCloudSettingSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s ApplicationCloudSettingSlice) Less(i, j int) bool {
	return s[i].SettingKey < s[j].SettingKey
}
