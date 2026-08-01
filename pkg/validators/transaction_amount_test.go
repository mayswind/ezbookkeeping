package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestValidTransactionAmount(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validTransactionAmount", ValidTransactionAmount)
	assert.NoError(t, err)

	amount := ""
	err = validate.Var(amount, "validTransactionAmount")
	assert.Nil(t, err)

	amount = "-1"
	err = validate.Var(amount, "validTransactionAmount")
	assert.Nil(t, err)

	amount = "0"
	err = validate.Var(amount, "validTransactionAmount")
	assert.Nil(t, err)

	amount = "1"
	err = validate.Var(amount, "validTransactionAmount")
	assert.Nil(t, err)

	amount = utils.Int64ToString(models.MinimumTransactionAmount)
	err = validate.Var(amount, "validTransactionAmount")
	assert.Nil(t, err)

	amount = utils.Int64ToString(models.MaximumTransactionAmount)
	err = validate.Var(amount, "validTransactionAmount")
	assert.Nil(t, err)
}

func TestInvalidTransactionAmount(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validTransactionAmount", ValidTransactionAmount)
	assert.NoError(t, err)

	amount := "1.23"
	err = validate.Var(amount, "validTransactionAmount")
	assert.NotNil(t, err)

	amount = "invalid"
	err = validate.Var(amount, "validTransactionAmount")
	assert.NotNil(t, err)

	amount = utils.Int64ToString(models.MinimumTransactionAmount - 1)
	err = validate.Var(amount, "validTransactionAmount")
	assert.NotNil(t, err)

	amount = utils.Int64ToString(models.MaximumTransactionAmount + 1)
	err = validate.Var(amount, "validTransactionAmount")
	assert.NotNil(t, err)
}
