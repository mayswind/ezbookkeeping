package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidUsername(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validUsername", ValidUsername)
	assert.Nil(t, err)

	username := "foobar"
	err = validate.Var(username, "validUsername")
	assert.Nil(t, err)

	username = "--foo_bar--"
	err = validate.Var(username, "validUsername")
	assert.Nil(t, err)
}

func TestInvalidUsername(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validUsername", ValidUsername)
	assert.Nil(t, err)

	username := "foo~bar~"
	err = validate.Var(username, "validUsername")
	assert.NotNil(t, err)
}
