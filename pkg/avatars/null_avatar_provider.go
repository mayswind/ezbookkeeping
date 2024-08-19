package avatars

import (
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// NullAvatarProvider represents the null avatar provider
type NullAvatarProvider struct {
}

// NewNullAvatarProvider returns a new null avatar provider
func NewNullAvatarProvider() *NullAvatarProvider {
	return &NullAvatarProvider{}
}

// GetAvatarUrl returns an empty url
func (p *NullAvatarProvider) GetAvatarUrl(user *models.User) string {
	return ""
}
