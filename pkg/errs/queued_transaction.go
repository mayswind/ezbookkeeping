package errs

import "net/http"

// Error codes related to queued transactions
var (
	ErrQueuedTransactionIdInvalid      = NewNormalError(NormalSubcategoryQueuedTransaction, 0, http.StatusBadRequest, "queued transaction id is invalid")
	ErrQueuedTransactionNotFound       = NewNormalError(NormalSubcategoryQueuedTransaction, 1, http.StatusBadRequest, "queued transaction not found")
	ErrQueuedTransactionAlreadyPending = NewNormalError(NormalSubcategoryQueuedTransaction, 2, http.StatusBadRequest, "queued transaction is not archived")
	ErrQueuedTransactionHasTooManyTags = NewNormalError(NormalSubcategoryQueuedTransaction, 3, http.StatusBadRequest, "queued transaction has too many tags")
)
