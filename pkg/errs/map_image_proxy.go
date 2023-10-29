package errs

import "net/http"

// Error codes related to map image proxy
var (
	ErrMapProviderNotCurrent = NewNormalError(NormalSubcategoryMapProxy, 0, http.StatusBadRequest, "specified map provider is not set")
)
