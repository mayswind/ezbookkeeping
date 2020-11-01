package models

type TwoFactor struct {
	Uid             int64  `xorm:"PK"`
	Secret          string `xorm:"VARCHAR(80) NOT NULL"`
	CreatedUnixTime int64
}

type TwoFactorLoginRequest struct {
	Passcode string `json:"passcode" binding:"required,notBlank,len=6"`
}

type TwoFactorEnableConfirmRequest struct {
	Secret   string `json:"secret" binding:"required,notBlank,len=32"`
	Passcode string `json:"passcode" binding:"required,notBlank,len=6"`
}

type TwoFactorEnableResponse struct {
	Secret string `json:"secret"`
	QRCode string `json:"qrcode"`
}

type TwoFactorEnableConfirmResponse struct {
	Token         string   `json:"token,omitempty"`
	RecoveryCodes []string `json:"recoveryCodes"`
}

type TwoFactorDisableRequest struct {
	Password string `json:"password" binding:"omitempty,min=6,max=128"`
}

type TwoFactorRegenerateRecoveryCodeRequest struct {
	Password string `json:"password" binding:"omitempty,min=6,max=128"`
}

type TwoFactorStatusResponse struct {
	Enable    bool  `json:"enable"`
	CreatedAt int64 `json:"createdAt,omitempty"`
}
