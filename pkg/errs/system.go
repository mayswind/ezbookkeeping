package errs

import "net/http"

var (
	ErrSystemError      = NewSystemError(SYSTEM_SUBCATEGORY_DEFAULT, 0, http.StatusInternalServerError, "system error")
	ErrApiNotFound      = NewSystemError(SYSTEM_SUBCATEGORY_DEFAULT, 1, http.StatusNotFound, "api not found")
	ErrMethodNotAllowed = NewSystemError(SYSTEM_SUBCATEGORY_DEFAULT, 2, http.StatusMethodNotAllowed, "method not allowed")
	ErrNotImplemented   = NewSystemError(SYSTEM_SUBCATEGORY_DEFAULT, 3, http.StatusNotImplemented, "not implemented")
)
