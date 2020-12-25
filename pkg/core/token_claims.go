package core

import (
	"github.com/dgrijalva/jwt-go"
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
	Username    string    `json:"username,omitempty"`
	Type        TokenType `json:"type"`
	jwt.StandardClaims
}
