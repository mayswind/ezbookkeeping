package errs

import "net/http"

var (
	ErrInvalidProtocol     = NewSystemError(SYSTEM_SUBCATEGORY_SETTING, 0, http.StatusInternalServerError, "invalid server protocol")
	ErrInvalidLogMode      = NewSystemError(SYSTEM_SUBCATEGORY_SETTING, 1, http.StatusInternalServerError, "invalid log mode")
	ErrGettingLocalAddress = NewSystemError(SYSTEM_SUBCATEGORY_SETTING, 2, http.StatusInternalServerError, "failed to get local address")
	ErrInvalidUuidMode     = NewSystemError(SYSTEM_SUBCATEGORY_SETTING, 3, http.StatusInternalServerError, "invalid uuid mode")
)
