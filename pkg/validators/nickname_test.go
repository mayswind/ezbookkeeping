package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidNickname(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validNickname", ValidNickname)
	assert.Nil(t, err)

	nickname := "0123456789012345678901234567890123456789012345678901234567890123"
	err = validate.Var(nickname, "validNickname")
	assert.Nil(t, err)
}

func TestInvalidNickname(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validNickname", ValidNickname)
	assert.Nil(t, err)

	nickname := "01234567890123456789012345678901234567890123456789012345678901234"
	err = validate.Var(nickname, "validNickname")
	assert.NotNil(t, err)
}
