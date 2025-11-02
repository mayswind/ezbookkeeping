package errs

import (
	"net/http"
)

// Error codes related to tokens
var (
	ErrTokenGenerating                      = NewNormalError(NormalSubcategoryToken, 0, http.StatusInternalServerError, "failed to generate token")
	ErrUnauthorizedAccess                   = NewNormalError(NormalSubcategoryToken, 1, http.StatusUnauthorized, "unauthorized access")
	ErrCurrentInvalidToken                  = NewNormalError(NormalSubcategoryToken, 2, http.StatusUnauthorized, "current token is invalid")
	ErrCurrentTokenExpired                  = NewNormalError(NormalSubcategoryToken, 3, http.StatusUnauthorized, "current token is expired")
	ErrCurrentInvalidTokenType              = NewNormalError(NormalSubcategoryToken, 4, http.StatusUnauthorized, "current token type is invalid")
	ErrCurrentTokenRequire2FA               = NewNormalError(NormalSubcategoryToken, 5, http.StatusUnauthorized, "current token requires two-factor authorization")
	ErrCurrentTokenNotRequire2FA            = NewNormalError(NormalSubcategoryToken, 6, http.StatusUnauthorized, "current token does not require two-factor authorization")
	ErrInvalidToken                         = NewNormalError(NormalSubcategoryToken, 7, http.StatusBadRequest, "token is invalid")
	ErrInvalidTokenId                       = NewNormalError(NormalSubcategoryToken, 8, http.StatusBadRequest, "token id is invalid")
	ErrInvalidUserTokenId                   = NewNormalError(NormalSubcategoryToken, 9, http.StatusBadRequest, "user token id is invalid")
	ErrTokenRecordNotFound                  = NewNormalError(NormalSubcategoryToken, 10, http.StatusBadRequest, "token is not found")
	ErrTokenExpired                         = NewNormalError(NormalSubcategoryToken, 11, http.StatusBadRequest, "token is expired")
	ErrTokenIsEmpty                         = NewNormalError(NormalSubcategoryToken, 12, http.StatusBadRequest, "token is empty")
	ErrEmailVerifyTokenIsInvalidOrExpired   = NewNormalError(NormalSubcategoryToken, 13, http.StatusBadRequest, "email verify token is invalid or expired")
	ErrPasswordResetTokenIsInvalidOrExpired = NewNormalError(NormalSubcategoryToken, 14, http.StatusBadRequest, "password reset token is invalid or expired")
	ErrNotAllowedToGenerateAPIToken         = NewNormalError(NormalSubcategoryToken, 15, http.StatusForbidden, "not allowed to generate api token")
)
