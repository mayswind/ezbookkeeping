package utils

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// PrintJsonSuccessResult writes success response in json format to current http context
func PrintJsonSuccessResult(c *core.WebContext, result any) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}

// PrintDataSuccessResult writes success response in custom content type to current http context
func PrintDataSuccessResult(c *core.WebContext, contentType string, fileName string, result []byte) {
	if fileName != "" {
		c.Header("Content-Disposition", "attachment;filename="+fileName)
	}

	c.Data(http.StatusOK, contentType, result)
}

// PrintJsonErrorResult writes error response in json format to current http context
func PrintJsonErrorResult(c *core.WebContext, err *errs.Error) {
	c.SetResponseError(err)

	errorMessage := err.Error()

	if err.Code() == errs.ErrIncompleteOrIncorrectSubmission.Code() && len(err.BaseError) > 0 {
		validationErrors, ok := err.BaseError[0].(validator.ValidationErrors)

		if ok {
			for _, err := range validationErrors {
				errorMessage = getValidationErrorText(err)
				break
			}
		}
	}

	result := gin.H{
		"success":      false,
		"errorCode":    err.Code(),
		"errorMessage": errorMessage,
		"path":         c.Request.URL.Path,
	}

	if err.Context != nil {
		result["context"] = err.Context
	}

	c.AbortWithStatusJSON(err.HttpStatusCode, result)
}

// PrintDataErrorResult writes error response in custom content type to current http context
func PrintDataErrorResult(c *core.WebContext, contentType string, err *errs.Error) {
	c.SetResponseError(err)

	errorMessage := err.Error()

	if err.Code() == errs.ErrIncompleteOrIncorrectSubmission.Code() && len(err.BaseError) > 0 {
		validationErrors, ok := err.BaseError[0].(validator.ValidationErrors)

		if ok {
			for _, err := range validationErrors {
				errorMessage = getValidationErrorText(err)
				break
			}
		}
	}

	c.Data(err.HttpStatusCode, contentType, []byte(errorMessage))
	c.Abort()
}

// SetEventStreamHeader sets the headers for event stream response
func SetEventStreamHeader(c *core.WebContext) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
}

func WriteEventStreamJsonSuccessResult(c *core.WebContext, result any) {
	data, err := json.Marshal(result)

	if err != nil {
		c.Abort()
		return
	}

	_, err = c.Writer.WriteString("data: " + string(data) + "\n\n")

	if err != nil {
		c.Abort()
		return
	}

	c.Writer.Flush()
}

func WriteEventStreamJsonErrorResult(c *core.WebContext, originalErr *errs.Error) {
	c.SetResponseError(originalErr)

	errorMessage := originalErr.Error()

	if originalErr.Code() == errs.ErrIncompleteOrIncorrectSubmission.Code() && len(originalErr.BaseError) > 0 {
		validationErrors, ok := originalErr.BaseError[0].(validator.ValidationErrors)

		if ok {
			for _, err := range validationErrors {
				errorMessage = getValidationErrorText(err)
				break
			}
		}
	}

	result := gin.H{
		"success":      false,
		"errorCode":    originalErr.Code(),
		"errorMessage": errorMessage,
		"path":         c.Request.URL.Path,
	}

	if originalErr.Context != nil {
		result["context"] = originalErr.Context
	}

	data, err := json.Marshal(result)

	if err != nil {
		c.Abort()
		return
	}

	_, err = c.Writer.WriteString("data: " + string(data) + "\n\n")

	if err != nil {
		c.Abort()
		return
	}

	c.Writer.Flush()
}

func getValidationErrorText(err validator.FieldError) string {
	fieldName := GetFirstLowerCharString(err.Field())

	switch err.Tag() {
	case "required":
		return errs.GetParameterIsRequiredMessage(fieldName)
	case "max":
		if isIntegerParameter(err.Kind()) {
			return errs.GetParameterMustLessThanMessage(fieldName, err.Param())
		} else if isStringParameter(err.Kind()) {
			return errs.GetParameterMustLessThanCharsMessage(fieldName, err.Param())
		}
	case "min":
		if isIntegerParameter(err.Kind()) {
			return errs.GetParameterMustMoreThanMessage(fieldName, err.Param())
		} else if isStringParameter(err.Kind()) {
			return errs.GetParameterMustMoreThanCharsMessage(fieldName, err.Param())
		}
	case "len":
		return errs.GetParameterLengthNotEqualMessage(fieldName, err.Param())
	case "notBlank":
		return errs.GetParameterNotBeBlankMessage(fieldName)
	case "validUsername":
		return errs.GetParameterInvalidUsernameMessage(fieldName)
	case "validEmail":
		return errs.GetParameterInvalidEmailMessage(fieldName)
	case "validCurrency":
		return errs.GetParameterInvalidCurrencyMessage(fieldName)
	case "validHexRGBColor":
		return errs.GetParameterInvalidHexRGBColorMessage(fieldName)
	case "validAmountFilter":
		return errs.GetParameterInvalidAmountFilterMessage(fieldName)
	}

	return errs.GetParameterInvalidMessage(fieldName)
}

func isIntegerParameter(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func isStringParameter(kind reflect.Kind) bool {
	return kind == reflect.String
}
