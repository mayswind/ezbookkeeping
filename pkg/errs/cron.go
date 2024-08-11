package errs

import "net/http"

// Error codes related to cron jobs
var (
	ErrCronJobNameIsEmpty           = NewSystemError(SystemSubcategoryCron, 0, http.StatusInternalServerError, "cron job name is empty")
	ErrCronJobNotExistsOrNotEnabled = NewSystemError(SystemSubcategoryCron, 1, http.StatusInternalServerError, "cron job not exists or not enabled")
)
