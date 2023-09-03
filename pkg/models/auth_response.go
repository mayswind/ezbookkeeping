package models

// AuthResponse returns a view-object of user authorization
type AuthResponse struct {
	Token           string         `json:"token"`
	Need2FA         bool           `json:"need2FA"`
	NeedVerifyEmail bool           `json:"needVerifyEmail"`
	User            *UserBasicInfo `json:"user"`
}
