package validators

import (
	"github.com/go-playground/validator/v10"

	"github.com/mayswind/lab/pkg/utils"
)

func ValidHexRGBColor(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if utils.IsValidHexRGBColor(value) {
			return true
		}
	}

	return false
}
