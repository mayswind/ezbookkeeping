package models

// OAuth2LoginRequest represents all parameters of OAuth 2.0 login request
type OAuth2LoginRequest struct {
	Platform        string `form:"platform" binding:"required"`
	ClientSessionId string `form:"client_session_id" binding:"required"`
}

// OAuth2CallbackRequest represents all parameters of OAuth 2.0 callback request
type OAuth2CallbackRequest struct {
	State string `form:"state"`
	Code  string `form:"code"`
}

// OAuth2CallbackLoginRequest represents all parameters of OAuth 2.0 callback login request
type OAuth2CallbackLoginRequest struct {
	Provider string `json:"provider" binding:"required,notBlank"`
	Password string `json:"password" binding:"omitempty,min=6,max=128"`
}
