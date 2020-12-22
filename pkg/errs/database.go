package errs

import (
	"net/http"
)

// Error codes related to database
var (
	ErrDatabaseTypeInvalid     = NewSystemError(SYSTEM_SUBCATEGORY_DATABASE, 0, http.StatusInternalServerError, "database type is invalid")
	ErrDatabaseHostInvalid     = NewSystemError(SYSTEM_SUBCATEGORY_DATABASE, 1, http.StatusInternalServerError, "database host is invalid")
	ErrDatabaseIsNull          = NewSystemError(SYSTEM_SUBCATEGORY_DATABASE, 2, http.StatusInternalServerError, "database cannot be null")
	ErrDatabaseOperationFailed = NewSystemError(SYSTEM_SUBCATEGORY_DATABASE, 3, http.StatusInternalServerError, "database operation failed")
)
