package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestEmptyAmount(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validAmountFilter", ValidAmountFilter)
	assert.Nil(t, err)

	err = validate.Var("", "validAmountFilter")
	assert.Nil(t, err)
}

func TestValidOneParameterAmount(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validAmountFilter", ValidAmountFilter)
	assert.Nil(t, err)

	for _, filterType := range []string{"gt", "lt", "eq", "ne"} {
		err = validate.Var(filterType+":0", "validAmountFilter")
		assert.Nil(t, err)

		err = validate.Var(filterType+":1", "validAmountFilter")
		assert.Nil(t, err)

		err = validate.Var(filterType+":-1", "validAmountFilter")
		assert.Nil(t, err)

		err = validate.Var(filterType+":1073741824", "validAmountFilter")
		assert.Nil(t, err)

		err = validate.Var(filterType+":-2147483648", "validAmountFilter")
		assert.Nil(t, err)
	}
}

func TestInvalidOneParameterAmount(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validAmountFilter", ValidAmountFilter)
	assert.Nil(t, err)

	for _, filterType := range []string{"gt", "lt", "eq", "ne"} {
		err = validate.Var(filterType, "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+": ", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":-", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":a", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":-a", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":1.1", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":-1.1", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":1:2", "validAmountFilter")
		assert.NotNil(t, err)
	}
}

func TestValidTwoParameterAmount(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validAmountFilter", ValidAmountFilter)
	assert.Nil(t, err)

	for _, filterType := range []string{"bt", "nb"} {
		err = validate.Var(filterType+":0:0", "validAmountFilter")
		assert.Nil(t, err)

		err = validate.Var(filterType+":1:1", "validAmountFilter")
		assert.Nil(t, err)

		err = validate.Var(filterType+":0:1", "validAmountFilter")
		assert.Nil(t, err)

		err = validate.Var(filterType+":-1:-1", "validAmountFilter")
		assert.Nil(t, err)

		err = validate.Var(filterType+":-1:0", "validAmountFilter")
		assert.Nil(t, err)

		err = validate.Var(filterType+":-2147483648:1073741824", "validAmountFilter")
		assert.Nil(t, err)
	}
}

func TestInvalidTwoParameterAmount(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validAmountFilter", ValidAmountFilter)
	assert.Nil(t, err)

	for _, filterType := range []string{"bt", "nb"} {
		err = validate.Var(filterType+":", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+"::", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":1", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":1:", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+"::1", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":-:-", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":a:b", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":-a:-b", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":1.1:1.2", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":-1.2:-1.1", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":1:0", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":0:-1", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":1073741824:-2147483648", "validAmountFilter")
		assert.NotNil(t, err)

		err = validate.Var(filterType+":1:2:3", "validAmountFilter")
		assert.NotNil(t, err)
	}
}
