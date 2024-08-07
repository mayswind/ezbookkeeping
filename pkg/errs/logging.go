package errs

import (
	"net/http"
)

// Error codes related to logging
var (
	ErrLoggingError = NewSystemError(SystemSubcategoryLogging, 0, http.StatusInternalServerError, "logging error")
)
