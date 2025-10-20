package errs

import (
	"net/http"
)

// Error codes related to oauth 2.0
var (
	ErrOAuth2NotEnabled                 = NewNormalError(NormalSubcategoryOAuth2, 0, http.StatusUnauthorized, "oauth 2.0 not enabled")
	ErrOAuth2AutoRegistrationNotEnabled = NewNormalError(NormalSubcategoryOAuth2, 1, http.StatusUnauthorized, "oauth 2.0 auto registration not enabled")
	ErrInvalidOAuth2LoginRequest        = NewNormalError(NormalSubcategoryOAuth2, 2, http.StatusUnauthorized, "invalid oauth 2.0 login request")
	ErrInvalidOAuth2Callback            = NewNormalError(NormalSubcategoryOAuth2, 3, http.StatusUnauthorized, "invalid oauth 2.0 callback")
	ErrMissingOAuth2State               = NewNormalError(NormalSubcategoryOAuth2, 4, http.StatusUnauthorized, "missing state in oauth 2.0 callback")
	ErrMissingOAuth2Code                = NewNormalError(NormalSubcategoryOAuth2, 5, http.StatusUnauthorized, "missing code in oauth 2.0 callback")
	ErrInvalidOAuth2State               = NewNormalError(NormalSubcategoryOAuth2, 6, http.StatusUnauthorized, "invalid state in oauth 2.0 callback")
	ErrCannotRetrieveOAuth2Token        = NewNormalError(NormalSubcategoryOAuth2, 7, http.StatusUnauthorized, "cannot retrieve oauth 2.0 token")
	ErrInvalidOAuth2Token               = NewNormalError(NormalSubcategoryOAuth2, 8, http.StatusUnauthorized, "invalid oauth 2.0 token")
	ErrCannotRetrieveUserInfo           = NewNormalError(NormalSubcategoryOAuth2, 9, http.StatusUnauthorized, "cannot retrieve user info from oauth 2.0 provider")
)
