package avatars

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// AvatarProviderContainer contains the current user avatar provider
type AvatarProviderContainer struct {
	Current AvatarProvider
}

// Initialize a user avatar provider container singleton instance
var (
	Container = &AvatarProviderContainer{}
)

// InitializeAvatarProvider initializes the current user avatar provider according to the config
func InitializeAvatarProvider(config *settings.Config) error {
	if config.AvatarProvider == core.USER_AVATAR_PROVIDER_INTERNAL {
		Container.Current = NewInternalStorageAvatarProvider(config)
		return nil
	} else if config.AvatarProvider == core.USER_AVATAR_PROVIDER_GRAVATAR {
		Container.Current = NewGravatarAvatarProvider()
		return nil
	} else if config.AvatarProvider == "" {
		Container.Current = NewNullAvatarProvider()
		return nil
	}

	return errs.ErrInvalidAvatarProvider
}

// GetAvatarUrl returns the avatar url by the current user avatar provider
func (p *AvatarProviderContainer) GetAvatarUrl(user *models.User) string {
	return p.Current.GetAvatarUrl(user)
}
