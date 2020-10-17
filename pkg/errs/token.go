package errs

import (
	"net/http"
)

var (
	ErrTokenGenerating     = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 0, http.StatusInternalServerError, "failed to generate token")
	ErrUnauthorizedAccess  = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 1, http.StatusUnauthorized, "unauthorized access")
	ErrTokenExpired        = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 2, http.StatusUnauthorized, "token is expired")
	ErrInvalidToken        = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 3, http.StatusUnauthorized, "token is invalid")
	ErrInvalidUserTokenId  = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 4, http.StatusUnauthorized, "user token id is invalid")
	ErrInvalidTokenId      = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 5, http.StatusUnauthorized, "token id is invalid")
	ErrTokenRecordNotFound = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 6, http.StatusUnauthorized, "token is not found")
	ErrInvalidTokenType    = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 7, http.StatusUnauthorized, "token type is invalid")
	ErrTokenRequire2FA     = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 8, http.StatusUnauthorized, "token requires two factor authorization")
	ErrTokenNotRequire2FA  = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 9, http.StatusUnauthorized, "token does not require two factor authorization")
)
