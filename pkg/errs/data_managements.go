package errs

import (
	"net/http"
)

// Error codes related to data management
var (
	ErrDataExportNotAllowed = NewNormalError(NormalSubcategoryDataManagement, 1, http.StatusBadRequest, "data export not allowed")
)
