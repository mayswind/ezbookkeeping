package errs

import "net/http"

// Error codes related to settings
var (
	ErrInvalidServerMode                              = NewSystemError(SystemSubcategorySetting, 0, http.StatusInternalServerError, "invalid server mode")
	ErrInvalidProtocol                                = NewSystemError(SystemSubcategorySetting, 1, http.StatusInternalServerError, "invalid server protocol")
	ErrInvalidLogMode                                 = NewSystemError(SystemSubcategorySetting, 2, http.StatusInternalServerError, "invalid log mode")
	ErrInvalidLogLevel                                = NewSystemError(SystemSubcategorySetting, 3, http.StatusInternalServerError, "invalid log level")
	ErrGettingLocalAddress                            = NewSystemError(SystemSubcategorySetting, 4, http.StatusInternalServerError, "failed to get local address")
	ErrInvalidUuidMode                                = NewSystemError(SystemSubcategorySetting, 5, http.StatusInternalServerError, "invalid uuid mode")
	ErrInvalidDuplicateCheckerType                    = NewSystemError(SystemSubcategorySetting, 6, http.StatusInternalServerError, "invalid duplicate checker type")
	ErrInvalidInMemoryDuplicateCheckerCleanupInterval = NewSystemError(SystemSubcategorySetting, 7, http.StatusInternalServerError, "invalid in-memory duplicate checker cleanup interval")
	ErrInvalidTokenExpiredTime                        = NewSystemError(SystemSubcategorySetting, 8, http.StatusInternalServerError, "invalid token expired time")
	ErrInvalidTokenMinRefreshInterval                 = NewSystemError(SystemSubcategorySetting, 9, http.StatusInternalServerError, "invalid token min refresh interval")
	ErrInvalidTemporaryTokenExpiredTime               = NewSystemError(SystemSubcategorySetting, 10, http.StatusInternalServerError, "invalid temporary token expired time")
	ErrInvalidEmailVerifyTokenExpiredTime             = NewSystemError(SystemSubcategorySetting, 11, http.StatusInternalServerError, "invalid email verify token expired time")
	ErrInvalidAvatarProvider                          = NewSystemError(SystemSubcategorySetting, 12, http.StatusInternalServerError, "invalid avatar provider")
	ErrInvalidMapProvider                             = NewSystemError(SystemSubcategorySetting, 13, http.StatusInternalServerError, "invalid map provider")
	ErrInvalidAmapSecurityVerificationMethod          = NewSystemError(SystemSubcategorySetting, 14, http.StatusInternalServerError, "invalid amap security verification method")
	ErrInvalidPasswordResetTokenExpiredTime           = NewSystemError(SystemSubcategorySetting, 15, http.StatusInternalServerError, "invalid password reset token expired time")
	ErrInvalidExchangeRatesDataSource                 = NewSystemError(SystemSubcategorySetting, 16, http.StatusInternalServerError, "invalid exchange rates data source")
)
