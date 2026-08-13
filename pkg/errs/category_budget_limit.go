package errs

import "net/http"

// [PLUGIN:budget] Error codes related to category budget limits
var (
	ErrCategoryBudgetLimitIdInvalid       = NewNormalError(NormalSubcategoryCategoryBudgetLimit, 0, http.StatusBadRequest, "category budget limit id is invalid")
	ErrCategoryBudgetLimitNotFound        = NewNormalError(NormalSubcategoryCategoryBudgetLimit, 1, http.StatusBadRequest, "category budget limit not found")
	ErrCategoryBudgetLimitCategoryIdEmpty = NewNormalError(NormalSubcategoryCategoryBudgetLimit, 2, http.StatusBadRequest, "category id is empty for category budget limit")
	ErrCategoryBudgetLimitAmountInvalid   = NewNormalError(NormalSubcategoryCategoryBudgetLimit, 3, http.StatusBadRequest, "category budget limit amount is invalid")
	ErrCategoryBudgetLimitPeriodInvalid   = NewNormalError(NormalSubcategoryCategoryBudgetLimit, 4, http.StatusBadRequest, "category budget limit period is invalid")
	ErrCategoryBudgetLimitDateInvalid     = NewNormalError(NormalSubcategoryCategoryBudgetLimit, 5, http.StatusBadRequest, "category budget limit start date is invalid")
	ErrCategoryBudgetLimitCategoryNotFound = NewNormalError(NormalSubcategoryCategoryBudgetLimit, 6, http.StatusBadRequest, "the category of category budget limit does not exist")
)
