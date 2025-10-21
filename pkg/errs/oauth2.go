package errs

import (
	"net/http"
)

// Error codes related to oauth 2.0
var (
	ErrOAuth2NotEnabled                 = NewNormalError(NormalSubcategoryOAuth2, 0, http.StatusUnauthorized, "oauth2 not enabled")
	ErrOAuth2AutoRegistrationNotEnabled = NewNormalError(NormalSubcategoryOAuth2, 1, http.StatusUnauthorized, "oauth2 auto registration not enabled")
	ErrInvalidOAuth2LoginRequest        = NewNormalError(NormalSubcategoryOAuth2, 2, http.StatusUnauthorized, "invalid oauth2 login request")
	ErrInvalidOAuth2Callback            = NewNormalError(NormalSubcategoryOAuth2, 3, http.StatusUnauthorized, "invalid oauth2 callback")
	ErrMissingOAuth2State               = NewNormalError(NormalSubcategoryOAuth2, 4, http.StatusUnauthorized, "missing state in oauth2 callback")
	ErrMissingOAuth2Code                = NewNormalError(NormalSubcategoryOAuth2, 5, http.StatusUnauthorized, "missing code in oauth2 callback")
	ErrInvalidOAuth2State               = NewNormalError(NormalSubcategoryOAuth2, 6, http.StatusUnauthorized, "invalid state in oauth2 callback")
	ErrCannotRetrieveOAuth2Token        = NewNormalError(NormalSubcategoryOAuth2, 7, http.StatusUnauthorized, "cannot retrieve oauth2 token")
	ErrInvalidOAuth2Token               = NewNormalError(NormalSubcategoryOAuth2, 8, http.StatusUnauthorized, "invalid oauth2 token")
	ErrCannotRetrieveUserInfo           = NewNormalError(NormalSubcategoryOAuth2, 9, http.StatusUnauthorized, "cannot retrieve user info from oauth2 provider")
)
