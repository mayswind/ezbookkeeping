package errs

import "net/http"

var (
	ErrAccountIdInvalid                       = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 0, http.StatusBadRequest, "account id is invalid")
	ErrAccountNotFound                        = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 1, http.StatusBadRequest, "account not found")
	ErrAccountTypeInvalid                     = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 2, http.StatusBadRequest, "account type is invalid")
	ErrAccountHaveNoSubAccount                = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 3, http.StatusBadRequest, "account must have at least one sub account")
	ErrAccountCannotHaveSubAccounts           = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 4, http.StatusBadRequest, "account cannot have sub accounts")
	ErrParentAccountCannotSetCurrency         = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 5, http.StatusBadRequest, "parent account cannot set currency")
	ErrSubAccountCategoryNotEqualsToParent    = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 6, http.StatusBadRequest, "sub account category not equals to parent")
	ErrSubAccountTypeInvalid                  = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 7, http.StatusBadRequest, "sub account type invalid")
	ErrCannotAddOrDeleteSubAccountsWhenModify = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 8, http.StatusBadRequest, "cannot add or delete sub accounts when modify account")
	ErrSourceAccountNotFound                  = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 9, http.StatusBadRequest, "source account not found")
	ErrDestinationAccountNotFound             = NewNormalError(NORMAL_SUBCATEGORY_ACCOUNT, 10, http.StatusBadRequest, "destination account not found")
)
