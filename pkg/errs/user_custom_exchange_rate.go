package errs

import "net/http"

// Error codes related to user custom exchange rates
var (
	ErrUserCustomExchangeRateNotFound             = NewNormalError(NormalSubcategoryUserCustomExchangeRate, 0, http.StatusBadRequest, "user custom exchange rate data not found")
	ErrCannotUpdateExchangeRateForDefaultCurrency = NewNormalError(NormalSubcategoryUserCustomExchangeRate, 1, http.StatusBadRequest, "cannot update exchange rate data for base currency")
	ErrCannotDeleteExchangeRateForDefaultCurrency = NewNormalError(NormalSubcategoryUserCustomExchangeRate, 2, http.StatusBadRequest, "cannot delete exchange rate data for base currency")
)
