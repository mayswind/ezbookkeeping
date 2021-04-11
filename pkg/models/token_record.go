package models

import "github.com/mayswind/ezbookkeeping/pkg/core"

// TokenMaxUserAgentLength represents the maximum size of user agent stored in database
const TokenMaxUserAgentLength = 255

// TokenRecord represents token data stored in database
type TokenRecord struct {
	Uid             int64          `xorm:"PK INDEX(IDX_token_record_uid_type_expired_time)"`
	UserTokenId     int64          `xorm:"PK"`
	TokenType       core.TokenType `xorm:"INDEX(IDX_token_record_uid_type_expired_time) TINYINT NOT NULL"`
	Secret          string         `xorm:"VARCHAR(10) NOT NULL"`
	UserAgent       string         `xorm:"VARCHAR(255)"`
	CreatedUnixTime int64          `xorm:"PK"`
	ExpiredUnixTime int64          `xorm:"INDEX(IDX_token_record_uid_type_expired_time)"`
}

// TokenRevokeRequest represents all parameters of token revoking request
type TokenRevokeRequest struct {
	TokenId string `json:"tokenId" binding:"required,notBlank"`
}

// TokenRefreshResponse represents all parameters of token refreshing request
type TokenRefreshResponse struct {
	NewToken   string         `json:"newToken"`
	OldTokenId string         `json:"oldTokenId"`
	User       *UserBasicInfo `json:"user"`
}

// TokenInfoResponse represents a view-object of token
type TokenInfoResponse struct {
	TokenId   string         `json:"tokenId"`
	TokenType core.TokenType `json:"tokenType"`
	UserAgent string         `json:"userAgent"`
	CreatedAt int64          `json:"createdAt"`
	ExpiredAt int64          `json:"expiredAt"`
	IsCurrent bool           `json:"isCurrent"`
}

// TokenInfoResponseSlice represents the slice data structure of TokenInfoResponse
type TokenInfoResponseSlice []*TokenInfoResponse

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
	return a[i].ExpiredAt > a[j].ExpiredAt
}
