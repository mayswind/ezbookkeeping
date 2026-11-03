package models

import "github.com/mayswind/ezbookkeeping/pkg/core"

// TokenMaxUserAgentLength represents the maximum size of user agent stored in database
const TokenMaxUserAgentLength = 255

// TokenRecord represents token data stored in database
type TokenRecord struct {
	core.TokenRecord
	Uid              int64          `xorm:"PK INDEX(IDX_token_record_uid_type_expired_time) INDEX(IDX_token_record_expired_time)"`
	UserTokenId      int64          `xorm:"PK"`
	TokenType        core.TokenType `xorm:"INDEX(IDX_token_record_uid_type_expired_time) TINYINT NOT NULL"`
	Secret           string         `xorm:"VARCHAR(10) NOT NULL"`
	UserAgent        string         `xorm:"VARCHAR(255)"`
	Context          string         `xorm:"BLOB"`
	CreatedUnixTime  int64          `xorm:"PK"`
	ExpiredUnixTime  int64          `xorm:"INDEX(IDX_token_record_uid_type_expired_time) INDEX(IDX_token_record_expired_time)"`
	LastSeenUnixTime int64
}

// MCPTokenContext represents the context data of an MCP token
type MCPTokenContext struct {
	EnabledTools int64 `json:"enabledTools,string"`
}

// OAuth2CallbackTokenContext represents the context data of oauth 2.0 callback token
type OAuth2CallbackTokenContext struct {
	ExternalAuthType core.UserExternalAuthType `json:"externalAuthType"`
	ExternalUsername string                    `json:"externalUsername"`
	ExternalEmail    string                    `json:"externalEmail"`
}

// TokenGenerateAPIRequest represents all parameters of api token generation request
type TokenGenerateAPIRequest struct {
	ExpiredInSeconds int64  `json:"expiresInSeconds" binding:"omitempty,min=0,max=4294967295"`
	Password         string `json:"password" binding:"omitempty,min=6,max=128"`
}

// TokenGenerateMCPRequest represents all parameters of mcp token generation request
type TokenGenerateMCPRequest struct {
	EnabledTools     []string `json:"enabledTools"`
	ExpiredInSeconds int64    `json:"expiresInSeconds" binding:"omitempty,min=0,max=4294967295"`
	Password         string   `json:"password" binding:"omitempty,min=6,max=128"`
}

// TokenRevokeRequest represents all parameters of token revoking request
type TokenRevokeRequest struct {
	TokenId string `json:"tokenId" binding:"required,notBlank"`
}

// TokenGenerateAPIResponse represents all response parameters of generated api token
type TokenGenerateAPIResponse struct {
	Token      string `json:"token"`
	APIBaseUrl string `json:"apiBaseUrl"`
}

// TokenGenerateMCPResponse represents all response parameters of generated mcp token
type TokenGenerateMCPResponse struct {
	Token  string `json:"token"`
	MCPUrl string `json:"mcpUrl"`
}

// TokenRefreshResponse represents all response parameters of token refreshing
type TokenRefreshResponse struct {
	NewToken                 string                        `json:"newToken,omitempty"`
	OldTokenId               string                        `json:"oldTokenId,omitempty"`
	User                     *UserBasicInfo                `json:"user"`
	ApplicationCloudSettings *ApplicationCloudSettingSlice `json:"applicationCloudSettings,omitempty"`
	NotificationContent      string                        `json:"notificationContent,omitempty"`
}

// TokenInfoResponse represents a view-object of token
type TokenInfoResponse struct {
	TokenId   string         `json:"tokenId"`
	TokenType core.TokenType `json:"tokenType"`
	UserAgent string         `json:"userAgent"`
	LastSeen  int64          `json:"lastSeen"`
	IsCurrent bool           `json:"isCurrent"`
}

// GetUid returns the uid of the token record
func (r *TokenRecord) GetUid() int64 {
	return r.Uid
}

// GetUserTokenId returns the user token id of the token record
func (r *TokenRecord) GetUserTokenId() int64 {
	return r.UserTokenId
}

// GetTokenType returns the token type of the token record
func (r *TokenRecord) GetTokenType() core.TokenType {
	return r.TokenType
}

// GetSecret returns the secret of the token record
func (r *TokenRecord) GetSecret() string {
	return r.Secret
}

// GetContext returns the context of the token record
func (r *TokenRecord) GetContext() string {
	return r.Context
}

// GetCreatedUnixTime returns the created unix time of the token record
func (r *TokenRecord) GetCreatedUnixTime() int64 {
	return r.CreatedUnixTime
}

// GetExpiredUnixTime returns the expired unix time of the token record
func (r *TokenRecord) GetExpiredUnixTime() int64 {
	return r.ExpiredUnixTime
}

// TokenInfoResponseSlice represents the slice data structure of TokenInfoResponse
type TokenInfoResponseSlice []*TokenInfoResponse

// GetEnabledToolsMark returns the enabled tool mask
func (c *MCPTokenContext) GetEnabledToolsMark() int64 {
	if c == nil {
		return core.MCPAllTools
	}

	return c.EnabledTools
}

// Len returns the count of items
func (a TokenInfoResponseSlice) Len() int {
	return len(a)
}

// Swap swaps two items
func (a TokenInfoResponseSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less reports whether the first item is less than the second one
func (a TokenInfoResponseSlice) Less(i, j int) bool {
	return a[i].LastSeen > a[j].LastSeen
}
