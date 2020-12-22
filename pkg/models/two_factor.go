package models

// TwoFactor represents user 2fa data stored in database
type TwoFactor struct {
	Uid             int64  `xorm:"PK"`
	Secret          string `xorm:"VARCHAR(80) NOT NULL"`
	CreatedUnixTime int64
}

// TwoFactorLoginRequest represents all parameters of 2fa login request
type TwoFactorLoginRequest struct {
	Passcode string `json:"passcode" binding:"required,notBlank,len=6"`
}

// TwoFactorEnableConfirmRequest represents all parameters of 2fa confirm request
type TwoFactorEnableConfirmRequest struct {
	Secret   string `json:"secret" binding:"required,notBlank,len=32"`
	Passcode string `json:"passcode" binding:"required,notBlank,len=6"`
}

// TwoFactorEnableResponse represents all response parameters when user requests to enable 2fa
type TwoFactorEnableResponse struct {
	Secret string `json:"secret"`
	QRCode string `json:"qrcode"`
}

// TwoFactorEnableConfirmResponse represents all response parameters after user have enabled 2fa
type TwoFactorEnableConfirmResponse struct {
	Token         string   `json:"token,omitempty"`
	RecoveryCodes []string `json:"recoveryCodes"`
}

// TwoFactorDisableRequest represents all parameters of 2fa disabling request
type TwoFactorDisableRequest struct {
	Password string `json:"password" binding:"omitempty,min=6,max=128"`
}

// TwoFactorRegenerateRecoveryCodeRequest represents all parameters of 2fa regenerating recovery codes request
type TwoFactorRegenerateRecoveryCodeRequest struct {
	Password string `json:"password" binding:"omitempty,min=6,max=128"`
}

// TwoFactorStatusResponse represents a view-object of 2fa status
type TwoFactorStatusResponse struct {
	Enable    bool  `json:"enable"`
	CreatedAt int64 `json:"createdAt,omitempty"`
}
