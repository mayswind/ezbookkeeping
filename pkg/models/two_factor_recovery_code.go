package models

type TwoFactorRecoveryCode struct {
	Uid             int64  `xorm:"PK"`
	RecoveryCode    string `xorm:"VARCHAR(64) PK"`
	Used            bool   `xorm:"NOT NULL"`
	CreatedUnixTime int64
	UsedUnixTime    int64
}

type TwoFactorRecoveryCodeLoginRequest struct {
	RecoveryCode string `json:"recoveryCode" binding:"required,notBlank,len=11"`
}
