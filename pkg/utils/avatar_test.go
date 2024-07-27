package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInternalAvatarUrl(t *testing.T) {
	expectedValue := "https://demo.ezbookkeeping.mayswind.net/avatar/1234567890.jpg"
	actualValue := GetInternalAvatarUrl(1234567890, "jpg", "https://demo.ezbookkeeping.mayswind.net/")
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = ""
	actualValue = GetInternalAvatarUrl(1234567890, "", "https://demo.ezbookkeeping.mayswind.net/")
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetGravatarUrl(t *testing.T) {
	// Reference: https://en.gravatar.com/site/implement/hash/
	expectedValue := "https://www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346"
	actualValue := GetGravatarUrl("MyEmailAddress@example.com")
	assert.Equal(t, expectedValue, actualValue)
}
