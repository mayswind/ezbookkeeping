package errs

import "net/http"

// Error codes related to overview
var (
	ErrQueryItemsEmpty   = NewNormalError(NormalSubcategoryOverview, 0, http.StatusBadRequest, "query items cannot be empty")
	ErrQueryItemsTooMuch = NewNormalError(NormalSubcategoryOverview, 1, http.StatusBadRequest, "query items too much")
)
