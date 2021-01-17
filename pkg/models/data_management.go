package models

// ClearDataRequest represents all parameters of clear user data request
type ClearDataRequest struct {
	Password string `json:"password" binding:"omitempty,min=6,max=128"`
}
