package avatars

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/models"
)

func TestGravatarAvatarProvider_GetGravatarUrl(t *testing.T) {
	avatarProvider := NewGravatarAvatarProvider()

	expectedValue := "https://www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346"
	actualValue := avatarProvider.GetAvatarUrl(&models.User{
		Email: "MyEmailAddress@example.com",
	})

	assert.Equal(t, expectedValue, actualValue)
}
