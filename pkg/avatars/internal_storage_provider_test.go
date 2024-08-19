package avatars

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

func TestInternalStorageAvatarProvider_GetAvatarUrl(t *testing.T) {
	avatarProvider := NewInternalStorageAvatarProvider(&settings.Config{
		RootUrl: "https://foo.bar/",
	})

	expectedValue := "https://foo.bar/avatar/1234567890.jpg"
	actualValue := avatarProvider.GetAvatarUrl(&models.User{
		Uid:              1234567890,
		CustomAvatarType: "jpg",
	})

	assert.Equal(t, expectedValue, actualValue)
}

func TestInternalStorageAvatarProvider_GetAvatarUrl_EmptyCustomAvatarType(t *testing.T) {
	avatarProvider := NewInternalStorageAvatarProvider(&settings.Config{
		RootUrl: "https://foo.bar/",
	})

	expectedValue := ""
	actualValue := avatarProvider.GetAvatarUrl(&models.User{
		Uid:              1234567890,
		CustomAvatarType: "",
	})

	assert.Equal(t, expectedValue, actualValue)
}
