package errs

import "net/http"

// Error codes related to two-factor authorization
var (
	ErrPasscodeInvalid               = NewNormalError(NormalSubcategoryTwofactor, 0, http.StatusUnauthorized, "passcode is invalid")
	ErrTwoFactorRecoveryCodeInvalid  = NewNormalError(NormalSubcategoryTwofactor, 1, http.StatusUnauthorized, "two-factor backup code is invalid")
	ErrTwoFactorRecoveryCodeNotExist = NewNormalError(NormalSubcategoryTwofactor, 2, http.StatusUnauthorized, "two-factor backup code does not exist")
	ErrTwoFactorIsNotEnabled         = NewNormalError(NormalSubcategoryTwofactor, 3, http.StatusBadRequest, "two-factor is not enabled")
	ErrTwoFactorAlreadyEnabled       = NewNormalError(NormalSubcategoryTwofactor, 4, http.StatusBadRequest, "two-factor has already been enabled")
)
