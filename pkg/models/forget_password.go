package models

// ForgetPasswordRequest represents all parameters of forget password request
type ForgetPasswordRequest struct {
	Email string `json:"email" binding:"required,notBlank,max=100,validEmail"`
}

// PasswordResetRequest represents all parameters of reset password request
type PasswordResetRequest struct {
	Email    string `json:"email" binding:"required,notBlank,max=100,validEmail"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}
