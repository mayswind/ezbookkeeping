package errs

import "net/http"

// Error codes related to mail
var (
	ErrSmtpServerNotEnabled  = NewSystemError(SystemSubcategoryMail, 0, http.StatusInternalServerError, "smtp server is not enabled")
	ErrSmtpServerHostInvalid = NewSystemError(SystemSubcategoryMail, 1, http.StatusInternalServerError, "smtp server host is invalid")
)
