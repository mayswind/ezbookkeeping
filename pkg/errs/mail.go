package errs

import "net/http"

// Error codes related to mail
var (
	ErrSMTPServerNotEnabled  = NewSystemError(SystemSubcategoryMail, 0, http.StatusInternalServerError, "SMTP server is not enabled")
	ErrSMTPServerHostInvalid = NewSystemError(SystemSubcategoryMail, 1, http.StatusInternalServerError, "SMTP server host is invalid")
)
