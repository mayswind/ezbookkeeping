package avatars

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/models"
)

type tokenRequiredAvatarProvider struct{}

func (p *tokenRequiredAvatarProvider) GetAvatarUrl(user *models.User) string {
	return ""
}

func (p *tokenRequiredAvatarProvider) NeedAvatarToken() bool {
	return true
}

func TestAvatarProviderContainer_NeedAvatarToken_Internal(t *testing.T) {
	container := &AvatarProviderContainer{current: &InternalStorageAvatarProvider{}}

	actualValue := container.NeedAvatarToken()
	assert.True(t, actualValue)
}

func TestAvatarProviderContainer_NeedAvatarToken_Gravatar(t *testing.T) {
	container := &AvatarProviderContainer{current: NewGravatarAvatarProvider()}

	actualValue := container.NeedAvatarToken()
	assert.False(t, actualValue)
}

func TestAvatarProviderContainer_NeedAvatarToken_Null(t *testing.T) {
	container := &AvatarProviderContainer{current: NewNullAvatarProvider()}

	actualValue := container.NeedAvatarToken()
	assert.False(t, actualValue)
}

func TestAvatarProviderContainer_NeedAvatarToken_Uninitialized(t *testing.T) {
	container := &AvatarProviderContainer{}

	actualValue := container.NeedAvatarToken()
	assert.False(t, actualValue)
}

func TestAvatarProviderContainer_NeedAvatarToken_CustomProvider(t *testing.T) {
	container := &AvatarProviderContainer{current: &tokenRequiredAvatarProvider{}}

	actualValue := container.NeedAvatarToken()
	assert.True(t, actualValue)
}
