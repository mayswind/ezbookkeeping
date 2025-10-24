package errs

import (
	"net/http"
)

// Error codes related to user external authentication
var (
	ErrUserExternalAuthNotFound      = NewNormalError(NormalSubcategoryUserExternalAuth, 0, http.StatusBadRequest, "user external auth is not found")
	ErrUserExternalAuthAlreadyExists = NewNormalError(NormalSubcategoryUserExternalAuth, 1, http.StatusBadRequest, "user external auth already exists")
	ErrUserExternalAuthTypeInvalid   = NewNormalError(NormalSubcategoryUserExternalAuth, 2, http.StatusBadRequest, "user external auth type invalid")
)
