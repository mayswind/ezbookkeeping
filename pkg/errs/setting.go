package errs

import "net/http"

// Error codes related to settings
var (
	ErrInvalidProtocol     = NewSystemError(SystemSubcategorySetting, 0, http.StatusInternalServerError, "invalid server protocol")
	ErrInvalidLogMode      = NewSystemError(SystemSubcategorySetting, 1, http.StatusInternalServerError, "invalid log mode")
	ErrGettingLocalAddress = NewSystemError(SystemSubcategorySetting, 2, http.StatusInternalServerError, "failed to get local address")
	ErrInvalidUuidMode     = NewSystemError(SystemSubcategorySetting, 3, http.StatusInternalServerError, "invalid uuid mode")
)
