package errs

import (
	"net/http"
)

var (
	ErrTokenGenerating           = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 0, http.StatusInternalServerError, "failed to generate token")
	ErrUnauthorizedAccess        = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 1, http.StatusUnauthorized, "unauthorized access")
	ErrCurrentInvalidToken       = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 2, http.StatusUnauthorized, "current token is invalid")
	ErrCurrentTokenExpired       = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 3, http.StatusUnauthorized, "current token is expired")
	ErrCurrentInvalidTokenType   = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 4, http.StatusUnauthorized, "current token type is invalid")
	ErrCurrentTokenRequire2FA    = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 5, http.StatusUnauthorized, "current token requires two factor authorization")
	ErrCurrentTokenNotRequire2FA = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 6, http.StatusUnauthorized, "current token does not require two factor authorization")
	ErrInvalidToken              = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 7, http.StatusBadRequest, "token is invalid")
	ErrInvalidTokenId            = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 8, http.StatusBadRequest, "token id is invalid")
	ErrInvalidUserTokenId        = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 9, http.StatusBadRequest, "user token id is invalid")
	ErrTokenRecordNotFound       = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 10, http.StatusBadRequest, "token is not found")
	ErrTokenExpired              = NewNormalError(NORMAL_SUBCATEGORY_TOKEN, 11, http.StatusBadRequest, "token is expired")
)
