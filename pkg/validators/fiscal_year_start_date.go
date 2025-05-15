package validators

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

// ValidateFiscalYearStart validates if a fiscal year start date is valid
func ValidateFiscalYearStart(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(core.FiscalYearStart)
	if !ok {
		return false
	}

	// Use the core functionality to validate
	_, _, err := date.GetMonthDay()
	return err == nil
}

// RegisterFiscalYearStartValidator registers the fiscal year start date validator
func RegisterFiscalYearStartValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validFiscalYearStart", ValidateFiscalYearStart)
	}
}
