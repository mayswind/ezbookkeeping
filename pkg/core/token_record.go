package core

// TokenRecord represents a user token record
type TokenRecord interface {
	// GetUid returns the uid of the token record
	GetUid() int64

	// GetUserTokenId returns the user token id of the token record
	GetUserTokenId() int64

	// GetTokenType returns the token type of the token record
	GetTokenType() TokenType

	// GetSecret returns the secret of the token record
	GetSecret() string

	// GetContext returns the context of the token record
	GetContext() string

	// GetCreatedUnixTime returns the created unix time of the token record
	GetCreatedUnixTime() int64

	// GetExpiredUnixTime returns the expired unix time of the token record
	GetExpiredUnixTime() int64
}
