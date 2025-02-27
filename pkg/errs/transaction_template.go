package errs

import "net/http"

// Error codes related to transaction templates
var (
	ErrTransactionTemplateIdInvalid                          = NewNormalError(NormalSubcategoryTemplate, 0, http.StatusBadRequest, "transaction template id is invalid")
	ErrTransactionTemplateNotFound                           = NewNormalError(NormalSubcategoryTemplate, 1, http.StatusBadRequest, "transaction template not found")
	ErrTransactionTemplateTypeInvalid                        = NewNormalError(NormalSubcategoryTemplate, 2, http.StatusBadRequest, "transaction template type is invalid")
	ErrScheduledTransactionNotEnabled                        = NewNormalError(NormalSubcategoryTemplate, 3, http.StatusBadRequest, "scheduled transaction is not enabled")
	ErrScheduledTransactionFrequencyInvalid                  = NewNormalError(NormalSubcategoryTemplate, 4, http.StatusBadRequest, "scheduled transaction frequency is invalid")
	ErrTransactionTemplateHasTooManyTags                     = NewNormalError(NormalSubcategoryTemplate, 5, http.StatusBadRequest, "transaction template has too many tags")
	ErrScheduledTransactionTemplateStartDataLaterThanEndDate = NewNormalError(NormalSubcategoryTemplate, 6, http.StatusBadRequest, "scheduled transaction start date is later than end time")
)
