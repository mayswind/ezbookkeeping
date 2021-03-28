package errs

// ErrorCategory represents error category
type ErrorCategory int

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
)

// Sub categories of normal error
const (
	NormalSubcategoryGlobal         = 0
	NormalSubcategoryUser           = 1
	NormalSubcategoryToken          = 2
	NormalSubcategoryTwofactor      = 3
	NormalSubcategoryAccount        = 4
	NormalSubcategoryTransaction    = 5
	NormalSubcategoryCategory       = 6
	NormalSubcategoryTag            = 7
	NormalSubcategoryDataManagement = 8
)

// Error represents the specific error returned to user
type Error struct {
	Category       ErrorCategory
	SubCategory    int
	Index          int
	HttpStatusCode int
	Message        string
	BaseError      []error
}

// Error returns the error message
func (err *Error) Error() string {
	return err.Message
}

// Code returns the error code
func (err *Error) Code() int {
	return int(err.Category)*100000 + err.SubCategory*1000 + err.Index
}

// New returns a new error instance
func New(category ErrorCategory, subCategory int, index int, httpStatusCode int, message string, baseError ...error) *Error {
	return &Error{
		Category:       category,
		SubCategory:    subCategory,
		Index:          index,
		HttpStatusCode: httpStatusCode,
		Message:        message,
		BaseError:      baseError,
	}
}

// NewSystemError returns a new system error instance
func NewSystemError(subCategory int, index int, httpStatusCode int, message string) *Error {
	return New(CATEGORY_SYSTEM, subCategory, index, httpStatusCode, message)
}

// NewNormalError returns a new normal error instance
func NewNormalError(subCategory int, index int, httpStatusCode int, message string) *Error {
	return New(CATEGORY_NORMAL, subCategory, index, httpStatusCode, message)
}

// NewIncompleteOrIncorrectSubmissionError returns a new incomplete or incorrect submission error instance
func NewIncompleteOrIncorrectSubmissionError(err error) *Error {
	return New(ErrIncompleteOrIncorrectSubmission.Category,
		ErrIncompleteOrIncorrectSubmission.SubCategory,
		ErrIncompleteOrIncorrectSubmission.Index,
		ErrIncompleteOrIncorrectSubmission.HttpStatusCode,
		ErrIncompleteOrIncorrectSubmission.Message, err)
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
