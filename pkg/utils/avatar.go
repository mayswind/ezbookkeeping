package utils

import (
	"fmt"
	"strings"
)

const gravatarUrlFormat = "https://www.gravatar.com/avatar/%s"

// GetGravatarUrl returns the Gravatar url according to the specified user email address
func GetGravatarUrl(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	emailMd5 := MD5EncodeToString([]byte(email))
	return fmt.Sprintf(gravatarUrlFormat, emailMd5)
}
