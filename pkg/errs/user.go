package errs

import (
	"net/http"
)

var (
	ErrUserIdInvalid              = NewNormalError(NORMAL_SUBCATEGORY_USER, 0, http.StatusBadRequest, "user id is invalid")
	ErrUsernameIsEmpty            = NewNormalError(NORMAL_SUBCATEGORY_USER, 1, http.StatusBadRequest, "username is empty")
	ErrEmailIsEmpty               = NewNormalError(NORMAL_SUBCATEGORY_USER, 2, http.StatusBadRequest, "email is empty")
	ErrPasswordIsEmpty            = NewNormalError(NORMAL_SUBCATEGORY_USER, 3, http.StatusBadRequest, "password is empty")
	ErrUserNotFound               = NewNormalError(NORMAL_SUBCATEGORY_USER, 4, http.StatusBadRequest, "user not found")
	ErrUserPasswordWrong          = NewNormalError(NORMAL_SUBCATEGORY_USER, 5, http.StatusBadRequest, "password is wrong")
	ErrUsernameAlreadyExists      = NewNormalError(NORMAL_SUBCATEGORY_USER, 6, http.StatusBadRequest, "username already exists")
	ErrUserEmailAlreadyExists     = NewNormalError(NORMAL_SUBCATEGORY_USER, 7, http.StatusBadRequest, "email already exists")
	ErrLoginNameInvalid           = NewNormalError(NORMAL_SUBCATEGORY_USER, 8, http.StatusUnauthorized, "login name is invalid")
	ErrLoginNameOrPasswordInvalid = NewNormalError(NORMAL_SUBCATEGORY_USER, 9, http.StatusUnauthorized, "login name or password is invalid")
	ErrLoginNameOrPasswordWrong   = NewNormalError(NORMAL_SUBCATEGORY_USER, 10, http.StatusUnauthorized, "login name or password is wrong")
)
