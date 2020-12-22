package errs

import "net/http"

// Error codes related to two factor authorization
var (
	ErrPasscodeInvalid               = NewNormalError(NORMAL_SUBCATEGORY_TWOFACTOR, 0, http.StatusUnauthorized, "passcode is invalid")
	ErrTwoFactorRecoveryCodeInvalid  = NewNormalError(NORMAL_SUBCATEGORY_TWOFACTOR, 1, http.StatusUnauthorized, "two factor backup code is invalid")
	ErrTwoFactorRecoveryCodeNotExist = NewNormalError(NORMAL_SUBCATEGORY_TWOFACTOR, 2, http.StatusUnauthorized, "two factor backup code does not exist")
	ErrTwoFactorIsNotEnabled         = NewNormalError(NORMAL_SUBCATEGORY_TWOFACTOR, 3, http.StatusBadRequest, "two factor is not enabled")
	ErrTwoFactorAlreadyEnabled       = NewNormalError(NORMAL_SUBCATEGORY_TWOFACTOR, 4, http.StatusBadRequest, "two factor has already been enabled")
)
