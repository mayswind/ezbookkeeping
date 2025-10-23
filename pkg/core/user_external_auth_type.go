package core

// UserExternalAuthType represents the type of user external authentication
type UserExternalAuthType string

// User External Auth Type
const (
	USER_EXTERNAL_AUTH_TYPE_OAUTH2_OIDC      UserExternalAuthType = "oidc"
	USER_EXTERNAL_AUTH_TYPE_OAUTH2_NEXTCLOUD UserExternalAuthType = "nextcloud"
	USER_EXTERNAL_AUTH_TYPE_OAUTH2_GITEA     UserExternalAuthType = "gitea"
	USER_EXTERNAL_AUTH_TYPE_OAUTH2_GITHUB    UserExternalAuthType = "github"
)

// IsValid checks if the UserExternalAuthType is valid
func (t UserExternalAuthType) IsValid() bool {
	switch t {
	case USER_EXTERNAL_AUTH_TYPE_OAUTH2_OIDC,
		USER_EXTERNAL_AUTH_TYPE_OAUTH2_NEXTCLOUD,
		USER_EXTERNAL_AUTH_TYPE_OAUTH2_GITEA,
		USER_EXTERNAL_AUTH_TYPE_OAUTH2_GITHUB:
		return true
	}
	return false
}
