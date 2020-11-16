package validators

import (
	"github.com/go-playground/validator/v10"

	"github.com/mayswind/lab/pkg/utils"
)

func ValidRGBColor(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if utils.IsValidRGBColor(value) {
			return true
		}
	}

	return false
}
