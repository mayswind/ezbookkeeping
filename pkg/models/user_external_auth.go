package models

import "github.com/mayswind/ezbookkeeping/pkg/core"

// UserExternalAuth represents user external auth data stored in database
type UserExternalAuth struct {
	Uid              int64                     `xorm:"PK"`
	ExternalAuthType core.UserExternalAuthType `xorm:"VARCHAR(32) PK UNIQUE(uqe_userexternalauth_authtype_username) UNIQUE(uqe_userexternalauth_authtype_email)"`
	ExternalUsername string                    `xorm:"VARCHAR(32) UNIQUE(uqe_userexternalauth_authtype_username) NOT NULL"`
	ExternalEmail    string                    `xorm:"VARCHAR(100) UNIQUE(uqe_userexternalauth_authtype_email) NOT NULL"`
	CreatedUnixTime  int64
}

// UserExternalAuthRevokeRequest represents all parameters of user external auth revoke request
type UserExternalAuthRevokeRequest struct {
	ExternalAuthType core.UserExternalAuthType `json:"externalAuthType" binding:"required,notBlank"`
}
