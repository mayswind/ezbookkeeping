package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGravatarUrl(t *testing.T) {
	// Reference: https://en.gravatar.com/site/implement/hash/
	expectedValue := "https://www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346"
	actualValue := GetGravatarUrl("MyEmailAddress@example.com")
	assert.Equal(t, expectedValue, actualValue)
}
