package errs

import (
	"fmt"
	"net/http"
)

// General error codes
var (
	ErrIncompleteOrIncorrectSubmission = NewNormalError(NORMAL_SUBCATEGORY_GLOBAL, 0, http.StatusBadRequest, "incomplete or incorrect submission")
	ErrOperationFailed                 = NewNormalError(NORMAL_SUBCATEGORY_GLOBAL, 1, http.StatusInternalServerError, "operation failed")
	ErrRequestIdInvalid                = NewNormalError(NORMAL_SUBCATEGORY_GLOBAL, 2, http.StatusInternalServerError, "request id is invalid")
	ErrCiphertextInvalid               = NewNormalError(NORMAL_SUBCATEGORY_GLOBAL, 3, http.StatusInternalServerError, "ciphertext is invalid")
	ErrNothingWillBeUpdated            = NewNormalError(NORMAL_SUBCATEGORY_GLOBAL, 4, http.StatusBadRequest, "nothing will be updated")
	ErrFailedToRequestRemoteApi        = NewNormalError(NORMAL_SUBCATEGORY_GLOBAL, 5, http.StatusBadRequest, "failed to request third party api")
	ErrPageIndexInvalid                = NewNormalError(NORMAL_SUBCATEGORY_GLOBAL, 6, http.StatusBadRequest, "page index is invalid")
	ErrPageCountInvalid                = NewNormalError(NORMAL_SUBCATEGORY_GLOBAL, 7, http.StatusBadRequest, "page count is invalid")
)

// GetParameterInvalidMessage returns specific error message for invalid parameter error
func GetParameterInvalidMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" is invalid", field)
}

// GetParameterIsRequiredMessage returns specific error message for missing parameter error
func GetParameterIsRequiredMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" is required", field)
}

// GetParameterMustLessThanMessage returns specific error message for parameter too large error
func GetParameterMustLessThanMessage(field string, param string) string {
	return fmt.Sprintf("parameter \"%s\" must be less than %s", field, param)
}

// GetParameterMustLessThanCharsMessage returns specific error message for parameter too long error
func GetParameterMustLessThanCharsMessage(field string, param string) string {
	return fmt.Sprintf("parameter \"%s\" must be less than %s characters", field, param)
}

// GetParameterMustMoreThanMessage returns specific error message for parameter too small error
func GetParameterMustMoreThanMessage(field string, param string) string {
	return fmt.Sprintf("parameter \"%s\" must be more than %s", field, param)
}

// GetParameterMustMoreThanCharsMessage returns specific error message for parameter too short error
func GetParameterMustMoreThanCharsMessage(field string, param string) string {
	return fmt.Sprintf("parameter \"%s\" must be more than %s characters", field, param)
}

// GetParameterLengthNotEqualMessage returns specific error message for parameter length wrong error
func GetParameterLengthNotEqualMessage(field string, param string) string {
	return fmt.Sprintf("parameter \"%s\" length is not equal to %s", field, param)
}

// GetParameterNotBeBlankMessage returns specific error message for blank parameter error
func GetParameterNotBeBlankMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" cannot be blank", field)
}

// GetParameterInvalidUsernameMessage returns specific error message for invalid username parameter error
func GetParameterInvalidUsernameMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" is invalid username format", field)
}

// GetParameterInvalidEmailMessage returns specific error message for invalid email parameter error
func GetParameterInvalidEmailMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" is invalid email format", field)
}

// GetParameterInvalidCurrencyMessage returns specific error message for invalid currency parameter error
func GetParameterInvalidCurrencyMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" is invalid currency", field)
}

// GetParameterInvalidHexRGBColorMessage returns specific error message for invalid hex rgb color parameter error
func GetParameterInvalidHexRGBColorMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" is invalid color", field)
}
