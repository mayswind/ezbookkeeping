package core

import (
	"fmt"
	"strconv"
	"strings"
)

// UserFeatureRestrictions represents all the restrictions of user features
type UserFeatureRestrictions uint64

// Add returns a new feature restrictions with the specified feature
func (r UserFeatureRestrictions) Add(featureRestrictionType UserFeatureRestrictionType) UserFeatureRestrictions {
	typeValue := uint64(1 << (featureRestrictionType - 1))
	return UserFeatureRestrictions(uint64(r) | typeValue)
}

// Remove returns a new feature restrictions without the specified feature
func (r UserFeatureRestrictions) Remove(featureRestrictionType UserFeatureRestrictionType) UserFeatureRestrictions {
	typeValue := uint64(1 << (featureRestrictionType - 1))
	return UserFeatureRestrictions(uint64(r) & (^typeValue))
}

// Contains returns whether contains the specified feature
func (r UserFeatureRestrictions) Contains(featureRestrictionType UserFeatureRestrictionType) bool {
	typeValue := uint64(1 << (featureRestrictionType - 1))
	return uint64(r)&typeValue == typeValue
}

// String returns a textual representation of all the restrictions of user features
func (r UserFeatureRestrictions) String() string {
	builder := strings.Builder{}

	for restrictionType := userFeatureRestrictionTypeMinValue; restrictionType <= userFeatureRestrictionTypeMaxValue; restrictionType++ {
		if !r.Contains(restrictionType) {
			continue
		}

		if builder.Len() > 0 {
			builder.WriteRune(',')
		}

		builder.WriteString(restrictionType.String())
	}

	return builder.String()
}

// ParseUserFeatureRestrictions returns restrictions of user features according to the textual restrictions of user features  separated by commas
func ParseUserFeatureRestrictions(featureRestrictions string) UserFeatureRestrictions {
	if len(featureRestrictions) < 1 {
		return 0
	}

	restrictions := uint64(0)
	typeValues := strings.Split(featureRestrictions, ",")

	for i := 0; i < len(typeValues); i++ {
		value, err := strconv.ParseInt(typeValues[i], 10, 64)

		if err != nil {
			continue
		}

		if uint64(userFeatureRestrictionTypeMinValue) <= uint64(value) && uint64(value) <= uint64(userFeatureRestrictionTypeMaxValue) {
			typeValue := uint64(1 << (value - 1))
			restrictions = restrictions | typeValue
		}
	}

	return UserFeatureRestrictions(restrictions)
}

// UserFeatureRestrictionType represents the restriction type of user features
type UserFeatureRestrictionType uint64

// User Feature Restriction Type
const (
	USER_FEATURE_RESTRICTION_TYPE_UPDATE_PASSWORD                              UserFeatureRestrictionType = 1
	USER_FEATURE_RESTRICTION_TYPE_UPDATE_EMAIL                                 UserFeatureRestrictionType = 2
	USER_FEATURE_RESTRICTION_TYPE_UPDATE_PROFILE_BASIC_INFO                    UserFeatureRestrictionType = 3
	USER_FEATURE_RESTRICTION_TYPE_UPDATE_AVATAR                                UserFeatureRestrictionType = 4
	USER_FEATURE_RESTRICTION_TYPE_REVOKE_OTHER_SESSION                         UserFeatureRestrictionType = 5
	USER_FEATURE_RESTRICTION_TYPE_ENABLE_2FA                                   UserFeatureRestrictionType = 6
	USER_FEATURE_RESTRICTION_TYPE_DISABLE_2FA                                  UserFeatureRestrictionType = 7
	USER_FEATURE_RESTRICTION_TYPE_FORGET_PASSWORD                              UserFeatureRestrictionType = 8
	USER_FEATURE_RESTRICTION_TYPE_IMPORT_TRANSACTION                           UserFeatureRestrictionType = 9
	USER_FEATURE_RESTRICTION_TYPE_EXPORT_TRANSACTION                           UserFeatureRestrictionType = 10
	USER_FEATURE_RESTRICTION_TYPE_CLEAR_ALL_DATA                               UserFeatureRestrictionType = 11
	USER_FEATURE_RESTRICTION_TYPE_SYNC_APPLICATION_SETTINGS                    UserFeatureRestrictionType = 12
	USER_FEATURE_RESTRICTION_TYPE_MCP_ACCESS                                   UserFeatureRestrictionType = 13
	USER_FEATURE_RESTRICTION_TYPE_CREATE_TRANSACTION_FROM_AI_IMAGE_RECOGNITION UserFeatureRestrictionType = 14
)

const userFeatureRestrictionTypeMinValue UserFeatureRestrictionType = USER_FEATURE_RESTRICTION_TYPE_UPDATE_PASSWORD
const userFeatureRestrictionTypeMaxValue UserFeatureRestrictionType = USER_FEATURE_RESTRICTION_TYPE_CREATE_TRANSACTION_FROM_AI_IMAGE_RECOGNITION

// String returns a textual representation of the restriction type of user features
func (t UserFeatureRestrictionType) String() string {
	switch t {
	case USER_FEATURE_RESTRICTION_TYPE_UPDATE_PASSWORD:
		return "Update Password"
	case USER_FEATURE_RESTRICTION_TYPE_UPDATE_EMAIL:
		return "Update Email"
	case USER_FEATURE_RESTRICTION_TYPE_UPDATE_PROFILE_BASIC_INFO:
		return "Update Profile Basic Info"
	case USER_FEATURE_RESTRICTION_TYPE_UPDATE_AVATAR:
		return "Update Avatar"
	case USER_FEATURE_RESTRICTION_TYPE_REVOKE_OTHER_SESSION:
		return "Logout Other Session"
	case USER_FEATURE_RESTRICTION_TYPE_ENABLE_2FA:
		return "Enable Two-Factor Authentication"
	case USER_FEATURE_RESTRICTION_TYPE_DISABLE_2FA:
		return "Disable Enable Two-Factor Authentication"
	case USER_FEATURE_RESTRICTION_TYPE_FORGET_PASSWORD:
		return "Forget Password"
	case USER_FEATURE_RESTRICTION_TYPE_IMPORT_TRANSACTION:
		return "Import Transactions"
	case USER_FEATURE_RESTRICTION_TYPE_EXPORT_TRANSACTION:
		return "Export Transactions"
	case USER_FEATURE_RESTRICTION_TYPE_CLEAR_ALL_DATA:
		return "Clear All Data"
	case USER_FEATURE_RESTRICTION_TYPE_SYNC_APPLICATION_SETTINGS:
		return "Sync Application Settings"
	case USER_FEATURE_RESTRICTION_TYPE_MCP_ACCESS:
		return "MCP (Model Context Protocol) Access"
	default:
		return fmt.Sprintf("Invalid(%d)", int(t))
	}
}
