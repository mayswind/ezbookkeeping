package models

import "github.com/mayswind/lab/pkg/core"

const TOKEN_USER_AGENT_MAX_LENGTH = 255

type TokenRecord struct {
	Uid             int64          `xorm:"PK INDEX(IDX_token_record_uid_type_expired_time)"`
	UserTokenId     int64          `xorm:"PK"`
	TokenType       core.TokenType `xorm:"INDEX(IDX_token_record_uid_type_expired_time) TINYINT NOT NULL"`
	Secret          string         `xorm:"VARCHAR(10) NOT NULL"`
	UserAgent       string         `xorm:"VARCHAR(255)"`
	CreatedUnixTime int64          `xorm:"PK"`
	ExpiredUnixTime int64          `xorm:"INDEX(IDX_token_record_uid_type_expired_time)"`
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

type TokenInfoResponseSlice []*TokenInfoResponse

func (a TokenInfoResponseSlice) Len() int {
	return len(a)
}

func (a TokenInfoResponseSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a TokenInfoResponseSlice) Less(i, j int) bool {
	return a[i].ExpiredAt > a[j].ExpiredAt
}
