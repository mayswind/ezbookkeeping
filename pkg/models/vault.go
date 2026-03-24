package models

// VaultInitRequest represents all parameters of vault initialization request
type VaultInitRequest struct {
	VaultSalt       string `json:"vaultSalt" binding:"required"`
	Argon2Params    string `json:"argon2Params" binding:"required"`
	EncryptedDek    string `json:"encryptedDek" binding:"required"`
	EncryptedX25519 string `json:"encryptedX25519" binding:"required"`
	X25519Public    string `json:"x25519Public" binding:"required"`
}

// VaultParamsResponse represents the vault parameters returned to the client
type VaultParamsResponse struct {
	VaultVersion    int    `json:"vaultVersion"`
	VaultSalt       string `json:"vaultSalt"`
	Argon2Params    string `json:"argon2Params"`
	EncryptedDek    string `json:"encryptedDek"`
	EncryptedX25519 string `json:"encryptedX25519"`
	X25519Public    string `json:"x25519Public"`
}

// LinkApiKeyRequest represents all parameters of API key linking request
type LinkApiKeyRequest struct {
	ApiKey string `json:"apiKey" binding:"required,notBlank"`
}

// LinkApiKeyResponse represents the response after linking an API key
type LinkApiKeyResponse struct {
	Tier string `json:"tier"`
}
