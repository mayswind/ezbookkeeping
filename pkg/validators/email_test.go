package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidEmail(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validEmail", ValidEmail)
	assert.Nil(t, err)

	email := "foo@bar.com"
	err = validate.Var(email, "validEmail")
	assert.Nil(t, err)

	email = "foo@1.2.3.4"
	err = validate.Var(email, "validEmail")
	assert.Nil(t, err)

	email = "foo_bar@foo.bar"
	err = validate.Var(email, "validEmail")
	assert.Nil(t, err)
}

func TestInvalidEmail(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validEmail", ValidEmail)
	assert.Nil(t, err)

	email := "foo"
	err = validate.Var(email, "validEmail")
	assert.NotNil(t, err)

	email = "@bar"
	err = validate.Var(email, "validEmail")
	assert.NotNil(t, err)

	email = "foo@bar"
	err = validate.Var(email, "validEmail")
	assert.NotNil(t, err)

	email = "foo@bar."
	err = validate.Var(email, "validEmail")
	assert.NotNil(t, err)
}
