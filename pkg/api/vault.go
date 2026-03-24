package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/Paxtiny/oscar/pkg/core"
	"github.com/Paxtiny/oscar/pkg/errs"
	"github.com/Paxtiny/oscar/pkg/log"
	"github.com/Paxtiny/oscar/pkg/models"
	"github.com/Paxtiny/oscar/pkg/services"
)

// VaultApi represents vault api
type VaultApi struct {
	users *services.UserService
}

// Initialize a vault api singleton instance
var (
	Vault = &VaultApi{
		users: services.Users,
	}
)

// VaultInitHandler initializes a new vault for the current user
func (a *VaultApi) VaultInitHandler(c *core.WebContext) (any, *errs.Error) {
	var vaultInitReq models.VaultInitRequest
	err := c.ShouldBindJSON(&vaultInitReq)

	if err != nil {
		log.Warnf(c, "[vault.VaultInitHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.Errorf(c, "[vault.VaultInitHandler] failed to get user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if user.VaultVersion > 0 {
		return nil, errs.ErrVaultAlreadyInitialized
	}

	// Decode base64 fields
	salt, err := base64.StdEncoding.DecodeString(vaultInitReq.VaultSalt)
	if err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	encryptedDek, err := base64.StdEncoding.DecodeString(vaultInitReq.EncryptedDek)
	if err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	encryptedX25519, err := base64.StdEncoding.DecodeString(vaultInitReq.EncryptedX25519)
	if err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	x25519Public, err := base64.StdEncoding.DecodeString(vaultInitReq.X25519Public)
	if err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	user.VaultVersion = 1
	user.VaultSalt = salt
	user.VaultArgon2Params = vaultInitReq.Argon2Params
	user.VaultEncryptedDek = encryptedDek
	user.VaultEncryptedX25519 = encryptedX25519
	user.VaultX25519Public = x25519Public
	user.UpdatedUnixTime = time.Now().Unix()

	updateErr := a.users.UpdateUserVault(c, user)

	if updateErr != nil {
		log.Errorf(c, "[vault.VaultInitHandler] failed to init vault for user \"uid:%d\", because %s", uid, updateErr.Error())
		return nil, errs.Or(updateErr, errs.ErrOperationFailed)
	}

	log.Infof(c, "[vault.VaultInitHandler] vault initialized for user \"uid:%d\"", uid)

	return true, nil
}

// VaultGetParamsHandler returns vault params for the current user
func (a *VaultApi) VaultGetParamsHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.Errorf(c, "[vault.VaultGetParamsHandler] failed to get user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if user.VaultVersion == 0 {
		return nil, errs.ErrVaultNotInitialized
	}

	return &models.VaultParamsResponse{
		VaultVersion:    user.VaultVersion,
		VaultSalt:       base64.StdEncoding.EncodeToString(user.VaultSalt),
		Argon2Params:    user.VaultArgon2Params,
		EncryptedDek:    base64.StdEncoding.EncodeToString(user.VaultEncryptedDek),
		EncryptedX25519: base64.StdEncoding.EncodeToString(user.VaultEncryptedX25519),
		X25519Public:    base64.StdEncoding.EncodeToString(user.VaultX25519Public),
	}, nil
}

// VaultUpdateParamsHandler updates vault params for the current user (passphrase change)
func (a *VaultApi) VaultUpdateParamsHandler(c *core.WebContext) (any, *errs.Error) {
	var vaultInitReq models.VaultInitRequest
	err := c.ShouldBindJSON(&vaultInitReq)

	if err != nil {
		log.Warnf(c, "[vault.VaultUpdateParamsHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.Errorf(c, "[vault.VaultUpdateParamsHandler] failed to get user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if user.VaultVersion == 0 {
		return nil, errs.ErrVaultNotInitialized
	}

	salt, err := base64.StdEncoding.DecodeString(vaultInitReq.VaultSalt)
	if err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	encryptedDek, err := base64.StdEncoding.DecodeString(vaultInitReq.EncryptedDek)
	if err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	encryptedX25519, err := base64.StdEncoding.DecodeString(vaultInitReq.EncryptedX25519)
	if err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	x25519Public, err := base64.StdEncoding.DecodeString(vaultInitReq.X25519Public)
	if err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	user.VaultSalt = salt
	user.VaultArgon2Params = vaultInitReq.Argon2Params
	user.VaultEncryptedDek = encryptedDek
	user.VaultEncryptedX25519 = encryptedX25519
	user.VaultX25519Public = x25519Public
	user.UpdatedUnixTime = time.Now().Unix()

	updateErr := a.users.UpdateUserVault(c, user)

	if updateErr != nil {
		log.Errorf(c, "[vault.VaultUpdateParamsHandler] failed to update vault for user \"uid:%d\", because %s", uid, updateErr.Error())
		return nil, errs.Or(updateErr, errs.ErrOperationFailed)
	}

	log.Infof(c, "[vault.VaultUpdateParamsHandler] vault params updated for user \"uid:%d\"", uid)

	return true, nil
}

// VaultShredHandler crypto-shreds vault for the current user
func (a *VaultApi) VaultShredHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.Errorf(c, "[vault.VaultShredHandler] failed to get user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if user.VaultVersion == 0 {
		return nil, errs.ErrVaultNotInitialized
	}

	user.VaultVersion = 0
	user.VaultSalt = nil
	user.VaultArgon2Params = ""
	user.VaultEncryptedDek = nil
	user.VaultEncryptedX25519 = nil
	user.VaultX25519Public = nil
	user.UpdatedUnixTime = time.Now().Unix()

	updateErr := a.users.UpdateUserVault(c, user)

	if updateErr != nil {
		log.Errorf(c, "[vault.VaultShredHandler] failed to shred vault for user \"uid:%d\", because %s", uid, updateErr.Error())
		return nil, errs.Or(updateErr, errs.ErrOperationFailed)
	}

	log.Infof(c, "[vault.VaultShredHandler] vault shredded for user \"uid:%d\"", uid)

	return true, nil
}

// LinkApiKeyHandler links a nicodAImus API key to the current user
func (a *VaultApi) LinkApiKeyHandler(c *core.WebContext) (any, *errs.Error) {
	var linkReq models.LinkApiKeyRequest
	err := c.ShouldBindJSON(&linkReq)

	if err != nil {
		log.Warnf(c, "[vault.LinkApiKeyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	// Stub validation: accept any non-empty key as "alfred" tier
	// TODO: When nicodaimus.api_url is configured, call /v1/validate with Bearer token
	tier := "alfred"

	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.Errorf(c, "[vault.LinkApiKeyHandler] failed to get user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	hash := sha256.Sum256([]byte(linkReq.ApiKey))
	keyHash := hex.EncodeToString(hash[:])
	user.NicodaimusKeyHash = keyHash
	user.NicodaimusTier = tier
	user.UpdatedUnixTime = time.Now().Unix()

	updateErr := a.users.UpdateUserApiKey(c, user)

	if updateErr != nil {
		log.Errorf(c, "[vault.LinkApiKeyHandler] failed to link API key for user \"uid:%d\", because %s", uid, updateErr.Error())
		return nil, errs.Or(updateErr, errs.ErrOperationFailed)
	}

	log.Infof(c, "[vault.LinkApiKeyHandler] API key linked for user \"uid:%d\", tier=%s", uid, tier)

	return &models.LinkApiKeyResponse{
		Tier: tier,
	}, nil
}
