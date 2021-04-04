package errs

import (
	"net/http"
)

// Error codes related to users
var (
	ErrLoginNameInvalid             = NewNormalError(NormalSubcategoryUser, 0, http.StatusUnauthorized, "login name is invalid")
	ErrLoginNameOrPasswordInvalid   = NewNormalError(NormalSubcategoryUser, 1, http.StatusUnauthorized, "login name or password is invalid")
	ErrLoginNameOrPasswordWrong     = NewNormalError(NormalSubcategoryUser, 2, http.StatusUnauthorized, "login name or password is wrong")
	ErrUserIdInvalid                = NewNormalError(NormalSubcategoryUser, 3, http.StatusBadRequest, "user id is invalid")
	ErrUsernameIsEmpty              = NewNormalError(NormalSubcategoryUser, 4, http.StatusBadRequest, "username is empty")
	ErrEmailIsEmpty                 = NewNormalError(NormalSubcategoryUser, 5, http.StatusBadRequest, "email is empty")
	ErrNicknameIsEmpty              = NewNormalError(NormalSubcategoryUser, 6, http.StatusBadRequest, "nickname is empty")
	ErrPasswordIsEmpty              = NewNormalError(NormalSubcategoryUser, 7, http.StatusBadRequest, "password is empty")
	ErrUserDefaultCurrencyIsEmpty   = NewNormalError(NormalSubcategoryUser, 8, http.StatusBadRequest, "user default currency is empty")
	ErrUserDefaultCurrencyIsInvalid = NewNormalError(NormalSubcategoryUser, 9, http.StatusBadRequest, "user default currency is invalid")
	ErrUserNotFound                 = NewNormalError(NormalSubcategoryUser, 10, http.StatusBadRequest, "user not found")
	ErrUserPasswordWrong            = NewNormalError(NormalSubcategoryUser, 11, http.StatusBadRequest, "password is wrong")
	ErrUsernameAlreadyExists        = NewNormalError(NormalSubcategoryUser, 12, http.StatusBadRequest, "username already exists")
	ErrUserEmailAlreadyExists       = NewNormalError(NormalSubcategoryUser, 13, http.StatusBadRequest, "email already exists")
	ErrUserRegistrationNotAllowed   = NewNormalError(NormalSubcategoryUser, 14, http.StatusBadRequest, "user registration not allowed")
)
