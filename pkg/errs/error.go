package errs

type ErrorCategory int

const (
	CATEGORY_SYSTEM ErrorCategory = 1
	CATEGORY_NORMAL ErrorCategory = 2

	SYSTEM_SUBCATEGORY_DEFAULT  = 0
	SYSTEM_SUBCATEGORY_SETTING  = 1
	SYSTEM_SUBCATEGORY_DATABASE = 2

	NORMAL_SUBCATEGORY_GLOBAL      = 0
	NORMAL_SUBCATEGORY_USER        = 1
	NORMAL_SUBCATEGORY_TOKEN       = 2
	NORMAL_SUBCATEGORY_TWOFACTOR   = 3
	NORMAL_SUBCATEGORY_ACCOUNT     = 4
	NORMAL_SUBCATEGORY_TRANSACTION = 5
	NORMAL_SUBCATEGORY_CATEGORY    = 6
	NORMAL_SUBCATEGORY_TAG         = 7
)

type Error struct {
	Category       ErrorCategory
	SubCategory    int
	Index          int
	HttpStatusCode int
	Message        string
	BaseError      []error
}

func (err *Error) Error() string {
	return err.Message
}

func (err *Error) Code() int {
	return int(err.Category)*100000 + err.SubCategory*1000 + err.Index
}

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

func NewSystemError(subCategory int, index int, httpStatusCode int, message string) *Error {
	return New(CATEGORY_SYSTEM, subCategory, index, httpStatusCode, message)
}

func NewNormalError(subCategory int, index int, httpStatusCode int, message string) *Error {
	return New(CATEGORY_NORMAL, subCategory, index, httpStatusCode, message)
}

func NewIncompleteOrIncorrectSubmissionError(err error) *Error {
	return New(ErrIncompleteOrIncorrectSubmission.Category,
		ErrIncompleteOrIncorrectSubmission.SubCategory,
		ErrIncompleteOrIncorrectSubmission.Index,
		ErrIncompleteOrIncorrectSubmission.HttpStatusCode,
		ErrIncompleteOrIncorrectSubmission.Message, err)
}

func Or(err error, defaultErr *Error) *Error {
	if finalError, ok := err.(*Error); ok {
		return finalError
	} else {
		return defaultErr
	}
}

func IsCustomError(err error) bool {
	_, ok := err.(*Error)
	return ok
}
