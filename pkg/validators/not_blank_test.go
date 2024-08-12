package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidNotBlank(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("notBlank", NotBlank)
	assert.Nil(t, err)

	err = validate.Var("1", "notBlank")
	assert.Nil(t, err)

	err = validate.Var(" 1 ", "notBlank")
	assert.Nil(t, err)
}

func TestInvalidNotBlank(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("notBlank", NotBlank)
	assert.Nil(t, err)

	err = validate.Var("", "notBlank")
	assert.NotNil(t, err)

	err = validate.Var("   ", "notBlank")
	assert.NotNil(t, err)
}
