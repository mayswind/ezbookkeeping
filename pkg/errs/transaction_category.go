package errs

import "net/http"

var (
	ErrTransactionCategoryIdInvalid            = NewNormalError(NORMAL_SUBCATEGORY_CATEGORY, 0, http.StatusBadRequest, "transaction category id is invalid")
	ErrTransactionCategoryNotFound             = NewNormalError(NORMAL_SUBCATEGORY_CATEGORY, 1, http.StatusBadRequest, "transaction category not found")
	ErrTransactionCategoryTypeInvalid          = NewNormalError(NORMAL_SUBCATEGORY_CATEGORY, 2, http.StatusBadRequest, "transaction category type is invalid")
	ErrParentTransactionCategoryNotFound       = NewNormalError(NORMAL_SUBCATEGORY_CATEGORY, 3, http.StatusBadRequest, "parent transaction category not found")
	ErrCannotAddToSecondaryTransactionCategory = NewNormalError(NORMAL_SUBCATEGORY_CATEGORY, 4, http.StatusBadRequest, "cannot add to secondary transaction category")
)
