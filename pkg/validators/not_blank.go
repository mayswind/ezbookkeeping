package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func NotBlank(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if value != "" && strings.Trim(value, " ") != "" {
			return true
		}
	}

	return false
}
