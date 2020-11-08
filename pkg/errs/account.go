package errs

import "net/http"

var (
	ErrAccountIdInvalid             = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 0, http.StatusBadRequest, "account id is invalid")
	ErrAccountNotFound              = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 1, http.StatusBadRequest, "account not found")
	ErrAccountTypeInvalid           = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 2, http.StatusBadRequest, "account type is invalid")
	ErrAccountHaveNoSubAccount      = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 3, http.StatusBadRequest, "account must have at least one sub account")
	ErrAccountCannotHaveSubAccounts = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 4, http.StatusBadRequest, "account cannot have sub accounts")
)
