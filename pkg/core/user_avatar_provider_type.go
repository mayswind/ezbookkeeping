package core

// UserAvatarProviderType represents type of the user avatar provider
type UserAvatarProviderType string

// User avatar provider types
const (
	USER_AVATAR_PROVIDER_INTERNAL UserAvatarProviderType = "internal"
	USER_AVATAR_PROVIDER_GRAVATAR UserAvatarProviderType = "gravatar"
)
