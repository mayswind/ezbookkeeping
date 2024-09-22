package errs

import "net/http"

// Error codes related to data converters
var (
	ErrNotFoundTransactionDataInFile       = NewNormalError(NormalSubcategoryConverter, 0, http.StatusBadRequest, "not found transaction data")
	ErrMissingRequiredFieldInHeaderRow     = NewNormalError(NormalSubcategoryConverter, 1, http.StatusBadRequest, "missing required field in header row")
	ErrFewerFieldsInDataRowThanInHeaderRow = NewNormalError(NormalSubcategoryConverter, 2, http.StatusBadRequest, "fewer fields in data row than in header row")
	ErrTransactionTimeInvalid              = NewNormalError(NormalSubcategoryConverter, 3, http.StatusBadRequest, "transaction time is invalid")
	ErrTransactionTimeZoneInvalid          = NewNormalError(NormalSubcategoryConverter, 4, http.StatusBadRequest, "transaction time zone is invalid")
	ErrAmountInvalid                       = NewNormalError(NormalSubcategoryConverter, 5, http.StatusBadRequest, "transaction amount is invalid")
	ErrGeographicLocationInvalid           = NewNormalError(NormalSubcategoryConverter, 6, http.StatusBadRequest, "geographic location is invalid")
	ErrFieldsInMultiTableAreDifferent      = NewNormalError(NormalSubcategoryConverter, 7, http.StatusBadRequest, "fields in multiple table headers are different")
	ErrInvalidFileHeader                   = NewNormalError(NormalSubcategoryConverter, 8, http.StatusBadRequest, "invalid file header")
	ErrInvalidCSVFile                      = NewNormalError(NormalSubcategoryConverter, 9, http.StatusBadRequest, "invalid csv file")
	ErrRelatedIdCannotBeBlank              = NewNormalError(NormalSubcategoryConverter, 10, http.StatusBadRequest, "related id cannot be blank")
	ErrFoundRecordNotHasRelatedRecord      = NewNormalError(NormalSubcategoryConverter, 11, http.StatusBadRequest, "found some transactions without related records")
)
