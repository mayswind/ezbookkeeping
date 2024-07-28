package errs

import "net/http"

// Error codes related to transaction templates
var (
	ErrTransactionTemplateIdInvalid = NewNormalError(NormalSubcategoryTemplate, 0, http.StatusBadRequest, "transaction template id is invalid")
	ErrTransactionTemplateNotFound  = NewNormalError(NormalSubcategoryTemplate, 1, http.StatusBadRequest, "transaction template not found")
)
