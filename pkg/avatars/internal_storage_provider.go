package avatars

import (
	"fmt"

	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const internalAvatarUrlFormat = "%savatar/%d.%s"

// InternalStorageAvatarProvider represents the internal storage avatar provider
type InternalStorageAvatarProvider struct {
	webRootUrl string
}

// NewInternalStorageAvatarProvider returns a new internal storage avatar provider
func NewInternalStorageAvatarProvider(config *settings.Config) *InternalStorageAvatarProvider {
	return &InternalStorageAvatarProvider{
		webRootUrl: config.RootUrl,
	}
}

// GetAvatarUrl returns the built-in avatar url
func (p *InternalStorageAvatarProvider) GetAvatarUrl(user *models.User) string {
	if user.CustomAvatarType == "" {
		return ""
	}

	return fmt.Sprintf(internalAvatarUrlFormat, p.webRootUrl, user.Uid, user.CustomAvatarType)
}
