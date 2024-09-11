package errs

import "net/http"

// Error codes related to data converters
var (
	ErrNotFoundTransactionDataInFile       = NewNormalError(NormalSubcategoryConverter, 0, http.StatusBadRequest, "not found transaction data")
	ErrMissingRequiredFieldInHeaderRow     = NewNormalError(NormalSubcategoryConverter, 1, http.StatusBadRequest, "missing required field in header row")
	ErrFewerFieldsInDataRowThanInHeaderRow = NewNormalError(NormalSubcategoryConverter, 2, http.StatusBadRequest, "fewer fields in data row than in header row")
	ErrTransactionTimeInvalid              = NewNormalError(NormalSubcategoryConverter, 3, http.StatusBadRequest, "transaction time is invalid")
	ErrTransactionTimeZoneInvalid          = NewNormalError(NormalSubcategoryConverter, 4, http.StatusBadRequest, "transaction time zone is invalid")
	ErrCategoryNameCannotBeBlank           = NewNormalError(NormalSubcategoryConverter, 5, http.StatusBadRequest, "category name cannot be blank")
	ErrSubCategoryNameCannotBeBlank        = NewNormalError(NormalSubcategoryConverter, 6, http.StatusBadRequest, "secondary category name cannot be blank")
	ErrAccountNameCannotBeBlank            = NewNormalError(NormalSubcategoryConverter, 7, http.StatusBadRequest, "account name cannot be blank")
	ErrDestinationAccountNameCannotBeBlank = NewNormalError(NormalSubcategoryConverter, 8, http.StatusBadRequest, "destination account name cannot be blank")
	ErrAmountInvalid                       = NewNormalError(NormalSubcategoryConverter, 9, http.StatusBadRequest, "transaction amount is invalid")
	ErrGeographicLocationInvalid           = NewNormalError(NormalSubcategoryConverter, 10, http.StatusBadRequest, "geographic location is invalid")
	ErrFieldsInMultiTableAreDifferent      = NewNormalError(NormalSubcategoryConverter, 11, http.StatusBadRequest, "fields in multiple table headers are different")
)
