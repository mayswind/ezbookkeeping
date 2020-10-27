package models

type AuthResponse struct {
	Token    string `json:"token"`
	Username string `json:"username,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Need2FA  bool   `json:"need2FA"`
}
