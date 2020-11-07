package models

type AuthResponse struct {
	Token   string         `json:"token"`
	Need2FA bool           `json:"need2FA"`
	User    *UserBasicInfo `json:"user"`
}
