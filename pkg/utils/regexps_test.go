package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidUsername_ValidUserName(t *testing.T) {
	username := "foobar"
	expectedValue := true
	actualValue := IsValidUsername(username)
	assert.Equal(t, expectedValue, actualValue)

	username = "--foo_bar--"
	expectedValue = true
	actualValue = IsValidUsername(username)
	assert.Equal(t, expectedValue, actualValue)
}

func TestIsValidUsername_InvalidUserName(t *testing.T) {
	username := "foo~bar~"
	expectedValue := false
	actualValue := IsValidUsername(username)
	assert.Equal(t, expectedValue, actualValue)
}

func TestIsValidEmail_ValidEmail(t *testing.T) {
	email := "foo@bar.com"
	expectedValue := true
	actualValue := IsValidEmail(email)
	assert.Equal(t, expectedValue, actualValue)

	email = "foo@1.2.3.4"
	expectedValue = true
	actualValue = IsValidEmail(email)
	assert.Equal(t, expectedValue, actualValue)

	email = "foo_bar@foo.bar"
	expectedValue = true
	actualValue = IsValidEmail(email)
	assert.Equal(t, expectedValue, actualValue)
}

func TestIsValidEmail_InvalidEmail(t *testing.T) {
	email := "foo"
	expectedValue := false
	actualValue := IsValidEmail(email)
	assert.Equal(t, expectedValue, actualValue)

	email = "@bar"
	expectedValue = false
	actualValue = IsValidEmail(email)
	assert.Equal(t, expectedValue, actualValue)

	email = "foo@bar"
	expectedValue = false
	actualValue = IsValidEmail(email)
	assert.Equal(t, expectedValue, actualValue)

	email = "foo@bar."
	expectedValue = false
	actualValue = IsValidEmail(email)
	assert.Equal(t, expectedValue, actualValue)
}
