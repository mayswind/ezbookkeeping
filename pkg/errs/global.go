package errs

import (
	"fmt"
	"net/http"
)

// General error codes
var (
	ErrIncompleteOrIncorrectSubmission = NewNormalError(NormalSubcategoryGlobal, 0, http.StatusBadRequest, "incomplete or incorrect submission")
	ErrOperationFailed                 = NewNormalError(NormalSubcategoryGlobal, 1, http.StatusInternalServerError, "operation failed")
	ErrRequestIdInvalid                = NewNormalError(NormalSubcategoryGlobal, 2, http.StatusInternalServerError, "request id is invalid")
	ErrCiphertextInvalid               = NewNormalError(NormalSubcategoryGlobal, 3, http.StatusInternalServerError, "ciphertext is invalid")
	ErrNothingWillBeUpdated            = NewNormalError(NormalSubcategoryGlobal, 4, http.StatusBadRequest, "nothing will be updated")
	ErrFailedToRequestRemoteApi        = NewNormalError(NormalSubcategoryGlobal, 5, http.StatusBadRequest, "failed to request third party api")
	ErrPageIndexInvalid                = NewNormalError(NormalSubcategoryGlobal, 6, http.StatusBadRequest, "page index is invalid")
	ErrPageCountInvalid                = NewNormalError(NormalSubcategoryGlobal, 7, http.StatusBadRequest, "page count is invalid")
	ErrClientTimezoneOffsetInvalid     = NewNormalError(NormalSubcategoryGlobal, 8, http.StatusBadRequest, "client timezone offset is invalid")
	ErrQueryItemsEmpty                 = NewNormalError(NormalSubcategoryGlobal, 9, http.StatusBadRequest, "query items cannot be blank")
	ErrQueryItemsTooMuch               = NewNormalError(NormalSubcategoryGlobal, 10, http.StatusBadRequest, "query items too much")
	ErrQueryItemsInvalid               = NewNormalError(NormalSubcategoryGlobal, 11, http.StatusBadRequest, "query items have invalid item")
	ErrParameterInvalid                = NewNormalError(NormalSubcategoryGlobal, 12, http.StatusBadRequest, "parameter invalid")
	ErrFormatInvalid                   = NewNormalError(NormalSubcategoryGlobal, 13, http.StatusBadRequest, "format invalid")
	ErrNumberInvalid                   = NewNormalError(NormalSubcategoryGlobal, 14, http.StatusBadRequest, "number invalid")
	ErrNoFilesUpload                   = NewNormalError(NormalSubcategoryGlobal, 15, http.StatusBadRequest, "no files uploaded")
	ErrUploadedFileEmpty               = NewNormalError(NormalSubcategoryGlobal, 16, http.StatusBadRequest, "uploaded file is empty")
	ErrExceedMaxUploadFileSize         = NewNormalError(NormalSubcategoryGlobal, 17, http.StatusBadRequest, "uploaded file size exceeds the maximum allowed size")
	ErrFailureCountLimitReached        = NewNormalError(NormalSubcategoryGlobal, 18, http.StatusBadRequest, "failure count exceeded maximum limit")
	ErrRepeatedRequest                 = NewNormalError(NormalSubcategoryGlobal, 19, http.StatusBadRequest, "repeated request")
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

// GetParameterInvalidAmountFilterMessage returns specific error message for invalid amount filter parameter error
func GetParameterInvalidAmountFilterMessage(field string) string {
	return fmt.Sprintf("parameter \"%s\" is invalid amount filter", field)
}
