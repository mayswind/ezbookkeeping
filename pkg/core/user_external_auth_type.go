package core

// UserExternalAuthType represents the type of user external authentication
type UserExternalAuthType string

// User External Auth Type
const (
	USER_EXTERNAL_AUTH_TYPE_OAUTH2_NEXTCLOUD UserExternalAuthType = "nextcloud"
)

// IsValid checks if the UserExternalAuthType is valid
func (t UserExternalAuthType) IsValid() bool {
	switch t {
	case USER_EXTERNAL_AUTH_TYPE_OAUTH2_NEXTCLOUD:
		return true
	}
	return false
}
