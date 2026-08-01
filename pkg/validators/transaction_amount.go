package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/mayswind/ezbookkeeping/pkg/models"

	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// ValidTransactionAmount returns whether the textual amount is within the
// supported single-transaction amount range.
func ValidTransactionAmount(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(int64); ok {
		return value >= models.MinimumTransactionAmount && value <= models.MaximumTransactionAmount
	} else if value, ok := fl.Field().Interface().(string); ok {
		if value == "" {
			return true
		}

		amount, err := utils.StringToInt64(value)

		if err != nil {
			return false
		}

		return amount >= models.MinimumTransactionAmount && amount <= models.MaximumTransactionAmount
	}

	return false
}
