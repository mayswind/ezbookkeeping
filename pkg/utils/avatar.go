package utils

import (
	"fmt"
	"strings"
)

const gravatarUrlFormat = "https://www.gravatar.com/avatar/%s"

// GetInternalAvatarUrl returns the internal avatar url
func GetInternalAvatarUrl(uid int64, avatarFileExtesion string, webRootUrl string) string {
	if avatarFileExtesion == "" {
		return ""
	}

	return fmt.Sprintf("%savatar/%d.%s", webRootUrl, uid, avatarFileExtesion)
}

// GetGravatarUrl returns the Gravatar url according to the specified user email address
func GetGravatarUrl(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	emailMd5 := MD5EncodeToString([]byte(email))
	return fmt.Sprintf(gravatarUrlFormat, emailMd5)
}
