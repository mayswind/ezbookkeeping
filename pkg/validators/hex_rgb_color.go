package validators

import (
	"github.com/go-playground/validator/v10"

	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// ValidHexRGBColor returns whether the given hex reb color is valid
func ValidHexRGBColor(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if utils.IsValidHexRGBColor(value) {
			return true
		}
	}

	return false
}
