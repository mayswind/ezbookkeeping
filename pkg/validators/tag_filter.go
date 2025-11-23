package validators

import (
	"github.com/go-playground/validator/v10"

	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// ValidTagFilter returns whether the given tag filter is valid
func ValidTagFilter(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if value == "" {
			return true
		}

		if value == models.TransactionNoTagFilterValue {
			return true
		}

		_, err := models.ParseTransactionTagFilter(value)

		return err == nil
	}

	return false
}
