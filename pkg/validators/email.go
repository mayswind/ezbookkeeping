package validators

import (
	"github.com/go-playground/validator/v10"

	"github.com/mayswind/lab/pkg/utils"
)

func ValidEmail(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if utils.IsValidEmail(value) {
			return true
		}
	}

	return false
}
