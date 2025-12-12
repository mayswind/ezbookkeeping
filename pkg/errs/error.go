package errs

import (
	"strings"
)

// ErrorCategory represents error category
type ErrorCategory int32

// Error categories
const (
	CATEGORY_SYSTEM ErrorCategory = 1
	CATEGORY_NORMAL ErrorCategory = 2
)

// Sub categories of system error
const (
	SystemSubcategoryDefault  = 0
	SystemSubcategorySetting  = 1
	SystemSubcategoryDatabase = 2
	SystemSubcategoryMail     = 3
	SystemSubcategoryLogging  = 4
	SystemSubcategoryCron     = 5
)

// Sub categories of normal error
const (
	NormalSubcategoryGlobal                 = 0
	NormalSubcategoryUser                   = 1
	NormalSubcategoryToken                  = 2
	NormalSubcategoryTwofactor              = 3
	NormalSubcategoryAccount                = 4
	NormalSubcategoryTransaction            = 5
	NormalSubcategoryCategory               = 6
	NormalSubcategoryTag                    = 7
	NormalSubcategoryDataManagement         = 8
	NormalSubcategoryMapProxy               = 9
	NormalSubcategoryTemplate               = 10
	NormalSubcategoryPicture                = 11
	NormalSubcategoryConverter              = 12
	NormalSubcategoryUserCustomExchangeRate = 13
	NormalSubcategoryModelContextProtocol   = 14
	NormalSubcategoryLargeLanguageModel     = 15
	NormalSubcategoryUserExternalAuth       = 16
	NormalSubcategoryOAuth2                 = 17
	NormalSubcategoryProject                = 18
)

// Error represents the specific error returned to user
type Error struct {
	Category       ErrorCategory
	SubCategory    int32
	Index          int32
	HttpStatusCode int
	Message        string
	BaseError      []error
	Context        any
}

type MultiErrors struct {
	errors []error
}

// Error returns the error message
func (err *Error) Error() string {
	return err.Message
}

// Code returns the error code
func (err *Error) Code() int32 {
	return int32(err.Category)*100000 + err.SubCategory*1000 + err.Index
}

// New returns a new error instance
func New(category ErrorCategory, subCategory int32, index int32, httpStatusCode int, message string, baseError ...error) *Error {
	return &Error{
		Category:       category,
		SubCategory:    subCategory,
		Index:          index,
		HttpStatusCode: httpStatusCode,
		Message:        message,
		BaseError:      baseError,
	}
}

// Error returns the error message
func (err *MultiErrors) Error() string {
	if len(err.errors) == 1 {
		return err.errors[0].Error()
	}

	var ret strings.Builder
	var lastErrorChar byte

	ret.WriteString("multi errors: ")

	for i := 0; i < len(err.errors); i++ {
		if i > 0 {
			if lastErrorChar == '.' {
				ret.WriteString(" ")
			} else {
				ret.WriteString(", ")
			}
		}

		errorContent := err.errors[i].Error()
		lastErrorChar = errorContent[len(errorContent)-1]
		ret.WriteString(errorContent)
	}

	return ret.String()
}

// NewSystemError returns a new system error instance
func NewSystemError(subCategory int32, index int32, httpStatusCode int, message string) *Error {
	return New(CATEGORY_SYSTEM, subCategory, index, httpStatusCode, message)
}

// NewNormalError returns a new normal error instance
func NewNormalError(subCategory int32, index int32, httpStatusCode int, message string) *Error {
	return New(CATEGORY_NORMAL, subCategory, index, httpStatusCode, message)
}

// NewLoggingError returns a new logging error instance
func NewLoggingError(message string, err ...error) *Error {
	return New(ErrLoggingError.Category,
		ErrLoggingError.SubCategory,
		ErrLoggingError.Index,
		ErrLoggingError.HttpStatusCode,
		message, err...)
}

// NewIncompleteOrIncorrectSubmissionError returns a new incomplete or incorrect submission error instance
func NewIncompleteOrIncorrectSubmissionError(err error) *Error {
	return New(ErrIncompleteOrIncorrectSubmission.Category,
		ErrIncompleteOrIncorrectSubmission.SubCategory,
		ErrIncompleteOrIncorrectSubmission.Index,
		ErrIncompleteOrIncorrectSubmission.HttpStatusCode,
		ErrIncompleteOrIncorrectSubmission.Message, err)
}

// NewErrorWithContext returns a new error instance with specified context
func NewErrorWithContext(baseError *Error, context any) *Error {
	return &Error{
		Category:       baseError.Category,
		SubCategory:    baseError.SubCategory,
		Index:          baseError.Index,
		HttpStatusCode: baseError.HttpStatusCode,
		Message:        baseError.Message,
		BaseError:      baseError.BaseError,
		Context:        context,
	}
}

// NewMultiErrorOrNil returns a new multi error instance
func NewMultiErrorOrNil(errors ...error) error {
	count := len(errors)

	if count < 1 {
		return nil
	} else if count == 1 {
		return errors[0]
	}

	return &MultiErrors{
		errors: errors,
	}
}

// Or would return the error from err parameter if the this error is defined in this project,
// or return the default error
func Or(err error, defaultErr *Error) *Error {
	if finalError, ok := err.(*Error); ok {
		return finalError
	} else {
		return defaultErr
	}
}

// IsCustomError returns whether this error is defined in this project
func IsCustomError(err error) bool {
	_, ok := err.(*Error)
	return ok
}
