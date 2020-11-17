package errs

import (
	"fmt"
	"net/http"
)

var (
	ErrIncompleteOrIncorrectSubmission = NewNormalError(NORMAL_SUBCATEGORY_GLOBAL, 0, http.StatusBadRequest, "incomplete or incorrect submission")
	ErrOperationFailed                 = NewNormalError(NORMAL_SUBCATEGORY_GLOBAL, 1, http.StatusInternalServerError, "operation failed")
	ErrRequestIdInvalid                = NewNormalError(NORMAL_SUBCATEGORY_GLOBAL, 2, http.StatusInternalServerError, "request id is invalid")
	ErrCiphertextInvalid               = NewNormalError(NORMAL_SUBCATEGORY_GLOBAL, 3, http.StatusInternalServerError, "ciphertext is invalid")
	ErrNothingWillBeUpdated            = NewNormalError(NORMAL_SUBCATEGORY_GLOBAL, 4, http.StatusBadRequest, "nothing will be updated")
)

func GetParameterInvalidMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" is invalid", field)
}

func GetParameterIsRequiredMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" is required", field)
}

func GetParameterMustLessThanMessage(field string, param string) string {
	return fmt.Sprintf("parameter \"%s\" must be less than %s", field, param)
}

func GetParameterMustMoreThanMessage(field string, param string) string {
	return fmt.Sprintf("parameter \"%s\" must be more than %s", field, param)
}

func GetParameterLengthNotEqualMessage(field string, param string) string {
	return fmt.Sprintf("parameter \"%s\" length is not equal to %s", field, param)
}

func GetParameterNotBeBlankMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" cannot be blank", field)
}

func GetParameterInvalidUsernameMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" is invalid username format", field)
}

func GetParameterInvalidEmailMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" is invalid email format", field)
}

func GetParameterInvalidCurrencyMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" is invalid currency", field)
}

func GetParameterInvalidHexRGBColorMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" is invalid color", field)
}
