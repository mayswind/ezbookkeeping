package core

import (
	"github.com/dgrijalva/jwt-go"
)

type TokenType byte

const (
	USER_TOKEN_TYPE_NORMAL      TokenType = 1
	USER_TOKEN_TYPE_REQUIRE_2FA TokenType = 2
)

type UserTokenClaims struct {
	UserTokenId string    `json:"userTokenId"`
	Username    string    `json:"username,omitempty"`
	Type        TokenType `json:"type"`
	jwt.StandardClaims
}
