package errs

import (
	"net/http"
)

// Error codes related to data management
var (
	ErrDataExportNotAllowed     = NewNormalError(NormalSubcategoryDataManagement, 1, http.StatusBadRequest, "data export not allowed")
	ErrDataImportNotAllowed     = NewNormalError(NormalSubcategoryDataManagement, 2, http.StatusBadRequest, "data import not allowed")
	ErrImportTooManyTransaction = NewNormalError(NormalSubcategoryDataManagement, 3, http.StatusBadRequest, "import too many transactions")
)
