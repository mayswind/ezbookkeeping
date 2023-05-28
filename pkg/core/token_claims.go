package core

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType represents token type
type TokenType byte

// Token types
const (
	USER_TOKEN_TYPE_NORMAL      TokenType = 1
	USER_TOKEN_TYPE_REQUIRE_2FA TokenType = 2
)

// UserTokenClaims represents user token
type UserTokenClaims struct {
	UserTokenId string    `json:"userTokenId"`
	Uid         int64     `json:"jti,string"`
	Username    string    `json:"username,omitempty"`
	Type        TokenType `json:"type"`
	IssuedAt    int64     `json:"iat"`
	ExpiresAt   int64     `json:"exp"`
}

// GetExpirationTime returns the expiration time of this token
func (c *UserTokenClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: time.Unix(c.ExpiresAt, 0),
	}, nil
}

// GetIssuedAt returns the issue time of this token
func (c *UserTokenClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: time.Unix(c.IssuedAt, 0),
	}, nil
}

// GetNotBefore returns the earliest valid time of this token
func (c *UserTokenClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{}, nil
}

// GetIssuer returns the issuer of this token
func (c *UserTokenClaims) GetIssuer() (string, error) {
	return "", nil
}

// GetSubject returns the subject of this token
func (c *UserTokenClaims) GetSubject() (string, error) {
	return "", nil
}

// GetAudience returns the audience of this token
func (c *UserTokenClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}
