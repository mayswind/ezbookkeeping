package errs

import "net/http"

// Error codes related to accounts
var (
	ErrAccountIdInvalid                       = NewNormalError(NormalSubcategoryAccount, 0, http.StatusBadRequest, "account id is invalid")
	ErrAccountNotFound                        = NewNormalError(NormalSubcategoryAccount, 1, http.StatusBadRequest, "account not found")
	ErrAccountTypeInvalid                     = NewNormalError(NormalSubcategoryAccount, 2, http.StatusBadRequest, "account type is invalid")
	ErrAccountCurrencyInvalid                 = NewNormalError(NormalSubcategoryAccount, 3, http.StatusBadRequest, "account currency is invalid")
	ErrAccountHaveNoSubAccount                = NewNormalError(NormalSubcategoryAccount, 4, http.StatusBadRequest, "account must have at least one sub-account")
	ErrAccountCannotHaveSubAccounts           = NewNormalError(NormalSubcategoryAccount, 5, http.StatusBadRequest, "account cannot have sub-accounts")
	ErrParentAccountCannotSetCurrency         = NewNormalError(NormalSubcategoryAccount, 6, http.StatusBadRequest, "parent account cannot set currency")
	ErrParentAccountCannotSetBalance          = NewNormalError(NormalSubcategoryAccount, 7, http.StatusBadRequest, "parent account cannot set balance")
	ErrSubAccountCategoryNotEqualsToParent    = NewNormalError(NormalSubcategoryAccount, 8, http.StatusBadRequest, "sub-account category not equals to parent")
	ErrSubAccountTypeInvalid                  = NewNormalError(NormalSubcategoryAccount, 9, http.StatusBadRequest, "sub-account type invalid")
	ErrCannotAddOrDeleteSubAccountsWhenModify = NewNormalError(NormalSubcategoryAccount, 10, http.StatusBadRequest, "cannot add or delete sub-accounts when modify account")
	ErrSourceAccountNotFound                  = NewNormalError(NormalSubcategoryAccount, 11, http.StatusBadRequest, "source account not found")
	ErrDestinationAccountNotFound             = NewNormalError(NormalSubcategoryAccount, 12, http.StatusBadRequest, "destination account not found")
	ErrAccountInUseCannotBeDeleted            = NewNormalError(NormalSubcategoryAccount, 13, http.StatusBadRequest, "account is in use and cannot be deleted")
	ErrAccountCategoryInvalid                 = NewNormalError(NormalSubcategoryAccount, 14, http.StatusBadRequest, "account category is invalid")
	ErrAccountBalanceTimeNotSet               = NewNormalError(NormalSubcategoryAccount, 15, http.StatusBadRequest, "account balance time is not set")
)
