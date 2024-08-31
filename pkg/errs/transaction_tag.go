package errs

import "net/http"

// Error codes related to transaction tags
var (
	ErrTransactionTagIdInvalid            = NewNormalError(NormalSubcategoryTag, 0, http.StatusBadRequest, "transaction tag id is invalid")
	ErrTransactionTagNotFound             = NewNormalError(NormalSubcategoryTag, 1, http.StatusBadRequest, "transaction tag not found")
	ErrTransactionTagNameIsEmpty          = NewNormalError(NormalSubcategoryTag, 2, http.StatusBadRequest, "transaction tag name is empty")
	ErrTransactionTagNameAlreadyExists    = NewNormalError(NormalSubcategoryTag, 3, http.StatusBadRequest, "transaction tag name already exists")
	ErrTransactionTagInUseCannotBeDeleted = NewNormalError(NormalSubcategoryTag, 4, http.StatusBadRequest, "transaction tag is in use and cannot be deleted")
	ErrTransactionTagIndexNotFound        = NewNormalError(NormalSubcategoryTag, 5, http.StatusBadRequest, "transaction tag index not found")
)
