package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserFeatureRestrictionsAdd(t *testing.T) {
	var featureRestrictions UserFeatureRestrictions
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_UPDATE_PASSWORD)
	expectedValue := UserFeatureRestrictions(1)
	assert.Equal(t, expectedValue, featureRestrictions)

	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_UPDATE_EMAIL)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_UPDATE_PROFILE_BASIC_INFO)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_UPDATE_AVATAR)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_REVOKE_OTHER_SESSION)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_ENABLE_2FA)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_DISABLE_2FA)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_FORGET_PASSWORD)
	expectedValue = UserFeatureRestrictions(255)
	assert.Equal(t, expectedValue, featureRestrictions)
}

func TestUserFeatureRestrictionsRemove(t *testing.T) {
	var featureRestrictions UserFeatureRestrictions
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_UPDATE_PASSWORD)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_UPDATE_EMAIL)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_UPDATE_PROFILE_BASIC_INFO)
	featureRestrictions = featureRestrictions.Remove(USER_FEATURE_RESTRICTION_TYPE_UPDATE_EMAIL)
	featureRestrictions = featureRestrictions.Remove(USER_FEATURE_RESTRICTION_TYPE_UPDATE_PROFILE_BASIC_INFO)
	featureRestrictions = featureRestrictions.Remove(USER_FEATURE_RESTRICTION_TYPE_UPDATE_AVATAR)
	featureRestrictions = featureRestrictions.Remove(USER_FEATURE_RESTRICTION_TYPE_REVOKE_OTHER_SESSION)
	expectedValue := UserFeatureRestrictions(1)
	assert.Equal(t, expectedValue, featureRestrictions)

	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_UPDATE_AVATAR)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_REVOKE_OTHER_SESSION)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_ENABLE_2FA)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_DISABLE_2FA)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_FORGET_PASSWORD)
	featureRestrictions = featureRestrictions.Remove(USER_FEATURE_RESTRICTION_TYPE_ENABLE_2FA)
	featureRestrictions = featureRestrictions.Remove(USER_FEATURE_RESTRICTION_TYPE_DISABLE_2FA)
	expectedValue = UserFeatureRestrictions(153)
	assert.Equal(t, expectedValue, featureRestrictions)
}

func TestUserFeatureRestrictionsContains(t *testing.T) {
	var featureRestrictions UserFeatureRestrictions
	assert.False(t, featureRestrictions.Contains(USER_FEATURE_RESTRICTION_TYPE_UPDATE_PASSWORD))

	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_UPDATE_PASSWORD)
	assert.True(t, featureRestrictions.Contains(USER_FEATURE_RESTRICTION_TYPE_UPDATE_PASSWORD))
	assert.False(t, featureRestrictions.Contains(USER_FEATURE_RESTRICTION_TYPE_UPDATE_PROFILE_BASIC_INFO))
}

func TestUserFeatureRestrictionsString(t *testing.T) {
	var featureRestrictions UserFeatureRestrictions
	expectedValue := ""
	actualValue := featureRestrictions.String()
	assert.Equal(t, expectedValue, actualValue)

	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_UPDATE_PASSWORD)
	expectedValue = "Update Password"
	actualValue = featureRestrictions.String()
	assert.Equal(t, expectedValue, actualValue)

	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_FORGET_PASSWORD)
	expectedValue = "Update Password,Forget Password"
	actualValue = featureRestrictions.String()
	assert.Equal(t, expectedValue, actualValue)

	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_UPDATE_EMAIL)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_UPDATE_PROFILE_BASIC_INFO)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_UPDATE_AVATAR)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_REVOKE_OTHER_SESSION)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_ENABLE_2FA)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_DISABLE_2FA)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_IMPORT_TRANSACTION)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_EXPORT_TRANSACTION)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_CLEAR_ALL_DATA)
	featureRestrictions = featureRestrictions.Add(USER_FEATURE_RESTRICTION_TYPE_SYNC_APPLICATION_SETTINGS)
	expectedValue = "Update Password," +
		"Update Email," +
		"Update Profile Basic Info," +
		"Update Avatar," +
		"Logout Other Session," +
		"Enable Two-Factor Authentication," +
		"Disable Enable Two-Factor Authentication," +
		"Forget Password," +
		"Import Transactions," +
		"Export Transactions," +
		"Clear All Data," +
		"Sync Application Settings"
	actualValue = featureRestrictions.String()
	assert.Equal(t, expectedValue, actualValue)
}

func TestParseUserFeatureRestrictions(t *testing.T) {
	expectedValue := UserFeatureRestrictions(0)
	actualValue := ParseUserFeatureRestrictions("")
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = UserFeatureRestrictions(1)
	actualValue = ParseUserFeatureRestrictions("1")
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = UserFeatureRestrictions(1)
	actualValue = ParseUserFeatureRestrictions("1,20")
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = UserFeatureRestrictions(255)
	actualValue = ParseUserFeatureRestrictions("1,2,3,4,5,6,7,8,20,21,22")
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = UserFeatureRestrictions(255)
	actualValue = ParseUserFeatureRestrictions("1,2,3,4,5,6,7,8,a,b,20")
	assert.Equal(t, expectedValue, actualValue)
}
