package errs

import "net/http"

// Error codes related to transaction tags
var (
	ErrTransactionTagIdInvalid            = NewNormalError(NORMAL_SUBCATEGORY_TAG, 0, http.StatusBadRequest, "transaction tag id is invalid")
	ErrTransactionTagNotFound             = NewNormalError(NORMAL_SUBCATEGORY_TAG, 1, http.StatusBadRequest, "transaction tag not found")
	ErrTransactionTagNameIsEmpty          = NewNormalError(NORMAL_SUBCATEGORY_TAG, 2, http.StatusBadRequest, "transaction tag name is empty")
	ErrTransactionTagNameAlreadyExists    = NewNormalError(NORMAL_SUBCATEGORY_TAG, 3, http.StatusBadRequest, "transaction tag name already exists")
	ErrTransactionTagInUseCannotBeDeleted = NewNormalError(NORMAL_SUBCATEGORY_TAG, 4, http.StatusBadRequest, "transaction tag is in use and cannot be deleted")
)
