package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// TokenUserAgentCreatedViaCli is the user agent of token created via cli
const TokenUserAgentCreatedViaCli = ApplicationName + " Cli"

// TokenUserAgentForAPI is the user agent for API token
const TokenUserAgentForAPI = ApplicationName + " API"

// TokenUserAgentForMCP is the user agent for MCP token
const TokenUserAgentForMCP = ApplicationName + " MCP"

// TokenType represents token type
type TokenType byte

// Token types
const (
	USER_TOKEN_TYPE_NORMAL                         TokenType = 1
	USER_TOKEN_TYPE_REQUIRE_2FA                    TokenType = 2
	USER_TOKEN_TYPE_EMAIL_VERIFY                   TokenType = 3
	USER_TOKEN_TYPE_PASSWORD_RESET                 TokenType = 4
	USER_TOKEN_TYPE_MCP                            TokenType = 5
	USER_TOKEN_TYPE_OAUTH2_CALLBACK_REQUIRE_VERIFY TokenType = 6
	USER_TOKEN_TYPE_OAUTH2_CALLBACK                TokenType = 7
	USER_TOKEN_TYPE_API                            TokenType = 8
)

// TextualTokenTypes maps token types to their textual token prefixes.
var TextualTokenTypes = map[TokenType]string{
	USER_TOKEN_TYPE_NORMAL:                         "tk",
	USER_TOKEN_TYPE_REQUIRE_2FA:                    "2fa",
	USER_TOKEN_TYPE_EMAIL_VERIFY:                   "vfy",
	USER_TOKEN_TYPE_PASSWORD_RESET:                 "rst",
	USER_TOKEN_TYPE_MCP:                            "mcp",
	USER_TOKEN_TYPE_OAUTH2_CALLBACK_REQUIRE_VERIFY: "olk",
	USER_TOKEN_TYPE_OAUTH2_CALLBACK:                "oauth2",
	USER_TOKEN_TYPE_API:                            "ak",
}

// TextualTokenTypesReverse maps textual token prefixes to their token types
var TextualTokenTypesReverse = func() map[string]TokenType {
	reverseMap := make(map[string]TokenType)

	for k, v := range TextualTokenTypes {
		reverseMap[v] = k
	}

	return reverseMap
}()

// TextualDerivedTokenType represents the textual type of a derived token
type TextualDerivedTokenType string

// Textual derived token types
const (
	NONE_DERIVED_TOKEN                       TextualDerivedTokenType = ""
	USER_AVATAR_DERIVED_NORMAL_TOKEN         TextualDerivedTokenType = "a"
	TRANSACTION_PICTURE_DERIVED_NORMAL_TOKEN TextualDerivedTokenType = "p"
)

// ValidTextualDerivedTokenTypes contains all supported textual derived token types
var ValidTextualDerivedTokenTypes = map[TextualDerivedTokenType]bool{
	USER_AVATAR_DERIVED_NORMAL_TOKEN:         true,
	TRANSACTION_PICTURE_DERIVED_NORMAL_TOKEN: true,
}

// UserTokenClaims represents user token
type UserTokenClaims struct {
	UserTokenId string    `json:"userTokenId"`
	Uid         int64     `json:"jti,string"`
	Username    string    `json:"username,omitempty"`
	Type        TokenType `json:"type"`
	IssuedAt    int64     `json:"iat"`
	ExpiresAt   int64     `json:"exp"`
}

// NormalTokenAdditionalData represents the additional data of a normal token
type NormalTokenAdditionalData struct {
	Scope                 NormalTokenScope
	DerivedTokenExpiresAt int64
}

// NormalTokenScope represents the scope of a normal token
type NormalTokenScope byte

// Normal token scopes
const (
	NormalTokenScopeInvalid            NormalTokenScope = 0
	NormalTokenScopeAll                NormalTokenScope = 1
	NormalTokenScopeUserAvatar         NormalTokenScope = 2
	NormalTokenScopeTransactionPicture NormalTokenScope = 3

	NormalTokenScopeMinValue NormalTokenScope = NormalTokenScopeAll
	NormalTokenScopeMaxValue NormalTokenScope = NormalTokenScopeTransactionPicture
)

// NormalTokenScopeDerivedTokenTypeMap maps normal token scopes to their corresponding derived token types
var NormalTokenScopeDerivedTokenTypeMap = map[NormalTokenScope]TextualDerivedTokenType{
	NormalTokenScopeUserAvatar:         USER_AVATAR_DERIVED_NORMAL_TOKEN,
	NormalTokenScopeTransactionPicture: TRANSACTION_PICTURE_DERIVED_NORMAL_TOKEN,
}

// String returns a textual representation of the token type enum
func (t TokenType) String() string {
	switch t {
	case USER_TOKEN_TYPE_NORMAL:
		return "Normal Token"
	case USER_TOKEN_TYPE_REQUIRE_2FA:
		return "Require 2FA Token"
	case USER_TOKEN_TYPE_EMAIL_VERIFY:
		return "Email Verify Token"
	case USER_TOKEN_TYPE_PASSWORD_RESET:
		return "Password Reset Token"
	case USER_TOKEN_TYPE_MCP:
		return "MCP Token"
	case USER_TOKEN_TYPE_OAUTH2_CALLBACK_REQUIRE_VERIFY:
		return "OAuth2 Callback Require Verify Token"
	case USER_TOKEN_TYPE_OAUTH2_CALLBACK:
		return "OAuth2 Callback Token"
	case USER_TOKEN_TYPE_API:
		return "API Token"
	default:
		return fmt.Sprintf("Invalid(%d)", int(t))
	}
}

// MarshalBinary returns the binary representation of the normal token additional data.
func (d *NormalTokenAdditionalData) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, d.Scope)

	if d.Scope != NormalTokenScopeAll {
		_ = binary.Write(buf, binary.BigEndian, d.DerivedTokenExpiresAt)
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary decodes the binary normal token additional data.
func (d *NormalTokenAdditionalData) UnmarshalBinary(data []byte) error {
	if len(data) < 1 {
		return errs.ErrCurrentInvalidToken
	}

	if NormalTokenScope(data[0]) == NormalTokenScopeAll && len(data) != 1 {
		return errs.ErrCurrentInvalidToken
	} else if NormalTokenScope(data[0]) != NormalTokenScopeAll && len(data) != 9 {
		return errs.ErrCurrentInvalidToken
	}

	d.Scope = NormalTokenScopeInvalid
	d.DerivedTokenExpiresAt = 0

	buf := bytes.NewBuffer(data)
	_ = binary.Read(buf, binary.BigEndian, &d.Scope)

	if d.Scope != NormalTokenScopeAll {
		_ = binary.Read(buf, binary.BigEndian, &d.DerivedTokenExpiresAt)
	}

	return nil
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
