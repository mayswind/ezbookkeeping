package models

import "github.com/mayswind/lab/pkg/core"

const TOKEN_USER_AGENT_MAX_LENGTH = 255

type TokenRecord struct {
	Uid             int64          `xorm:"PK"`
	UserTokenId     int64          `xorm:"PK"`
	TokenType       core.TokenType `xorm:"TINYINT NOT NULL"`
	Secret          string         `xorm:"VARCHAR(10) NOT NULL"`
	UserAgent       string         `xorm:"VARCHAR(255)"`
	CreatedUnixTime int64          `xorm:"PK"`
	ExpiredUnixTime int64
}

type TokenRevokeRequest struct {
	TokenId string `json:"tokenId" binding:"required,notBlank"`
}

type TokenRefreshResponse struct {
	NewToken   string         `json:"newToken"`
	OldTokenId string         `json:"oldTokenId"`
	User       *UserBasicInfo `json:"user"`
}

type TokenInfoResponse struct {
	TokenId   string         `json:"tokenId"`
	TokenType core.TokenType `json:"tokenType"`
	UserAgent string         `json:"userAgent"`
	CreatedAt int64          `json:"createdAt"`
	ExpiredAt int64          `json:"expiredAt"`
	IsCurrent bool           `json:"isCurrent"`
}
