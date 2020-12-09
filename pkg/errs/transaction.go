package errs

import "net/http"

var (
	ErrTransactionIdInvalid                                = NewNormalError(NORMAL_SUBCATEGORY_TRANSACTION, 0, http.StatusBadRequest, "transaction id is invalid")
	ErrTransactionNotFound                                 = NewNormalError(NORMAL_SUBCATEGORY_TRANSACTION, 1, http.StatusBadRequest, "transaction not found")
	ErrTransactionTypeInvalid                              = NewNormalError(NORMAL_SUBCATEGORY_TRANSACTION, 2, http.StatusBadRequest, "transaction type is invalid")
	ErrTransactionSourceAndDestinationIdNotEqual           = NewNormalError(NORMAL_SUBCATEGORY_TRANSACTION, 3, http.StatusBadRequest, "transaction source and destination account id not equal")
	ErrTransactionSourceAndDestinationIdCannotBeEqual      = NewNormalError(NORMAL_SUBCATEGORY_TRANSACTION, 4, http.StatusBadRequest, "transaction source and destination account id cannot be equal")
	ErrTransactionSourceAndDestinationAmountNotEqual       = NewNormalError(NORMAL_SUBCATEGORY_TRANSACTION, 5, http.StatusBadRequest, "transaction source and destination amount not equal")
	ErrTooMuchTransactionInOneSecond                       = NewNormalError(NORMAL_SUBCATEGORY_TRANSACTION, 6, http.StatusBadRequest, "too much transaction in one second")
	ErrBalanceModificationTransactionCannotSetCategory     = NewNormalError(NORMAL_SUBCATEGORY_TRANSACTION, 7, http.StatusBadRequest, "balance modification transaction cannot set category")
	ErrBalanceModificationTransactionCannotChangeAccountId = NewNormalError(NORMAL_SUBCATEGORY_TRANSACTION, 8, http.StatusBadRequest, "balance modification transaction cannot change account id")
	ErrBalanceModificationTransactionCannotAddWhenNotEmpty = NewNormalError(NORMAL_SUBCATEGORY_TRANSACTION, 9, http.StatusBadRequest, "balance modification transaction cannot add when other transaction exists")
	ErrCannotAddTransactionToHiddenAccount                 = NewNormalError(NORMAL_SUBCATEGORY_TRANSACTION, 10, http.StatusBadRequest, "cannot add transaction to hidden account")
	ErrCannotModifyTransactionInHiddenAccount              = NewNormalError(NORMAL_SUBCATEGORY_TRANSACTION, 11, http.StatusBadRequest, "cannot modify transaction of hidden account")
	ErrCannotDeleteTransactionInHiddenAccount              = NewNormalError(NORMAL_SUBCATEGORY_TRANSACTION, 12, http.StatusBadRequest, "cannot delete transaction in hidden account")
)
