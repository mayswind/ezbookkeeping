package api

import (
	"strings"
	"time"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/services"
	"github.com/mayswind/lab/pkg/settings"
)

// UsersApi represents user api
type UsersApi struct {
	users  *services.UserService
	tokens *services.TokenService
}

// Initialize a user api singleton instance
var (
	Users = &UsersApi{
		users:  services.Users,
		tokens: services.Tokens,
	}
)

// UserRegisterHandler saves a new user by request parameters
func (a *UsersApi) UserRegisterHandler(c *core.Context) (interface{}, *errs.Error) {
	if !settings.Container.Current.EnableUserRegister {
		return nil, errs.ErrUserRegistrationNotAllowed
	}

	var userRegisterReq models.UserRegisterRequest
	err := c.ShouldBindJSON(&userRegisterReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[users.UserRegisterHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	userRegisterReq.Username = strings.TrimSpace(userRegisterReq.Username)
	userRegisterReq.Email = strings.TrimSpace(userRegisterReq.Email)
	userRegisterReq.Nickname = strings.TrimSpace(userRegisterReq.Nickname)

	user := &models.User{
		Username:        userRegisterReq.Username,
		Email:           userRegisterReq.Email,
		Nickname:        userRegisterReq.Nickname,
		Password:        userRegisterReq.Password,
		DefaultCurrency: userRegisterReq.DefaultCurrency,
	}

	err = a.users.CreateUser(user)

	if err != nil {
		log.ErrorfWithRequestId(c, "[users.UserRegisterHandler] failed to create user \"%s\", because %s", user.Username, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[users.UserRegisterHandler] user \"%s\" has registered successfully, uid is %d", user.Username, user.Uid)

	authResp := &models.AuthResponse{
		Need2FA: false,
		User:    user.ToUserBasicInfo(),
	}

	token, claims, err := a.tokens.CreateToken(user, c)

	if err != nil {
		log.WarnfWithRequestId(c, "[users.UserRegisterHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return authResp, nil
	}

	authResp.Token = token
	c.SetTokenClaims(claims)

	log.InfofWithRequestId(c, "[users.UserRegisterHandler] user \"uid:%d\" has logined, token will be expired at %d", user.Uid, claims.ExpiresAt)

	return authResp, nil
}

// UserProfileHandler returns user profile of current user
func (a *UsersApi) UserProfileHandler(c *core.Context) (interface{}, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.ErrorfWithRequestId(c, "[users.UserRegisterHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	userResp := user.ToUserProfileResponse()
	return userResp, nil
}

// UserUpdateProfileHandler saves user profile by request parameters for current user
func (a *UsersApi) UserUpdateProfileHandler(c *core.Context) (interface{}, *errs.Error) {
	var userUpdateReq models.UserProfileUpdateRequest
	err := c.ShouldBindJSON(&userUpdateReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[users.UserUpdateProfileHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.ErrorfWithRequestId(c, "[users.UserUpdateProfileHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	userUpdateReq.Email = strings.TrimSpace(userUpdateReq.Email)
	userUpdateReq.Nickname = strings.TrimSpace(userUpdateReq.Nickname)

	anythingUpdate := false
	userNew := &models.User{
		Uid:   user.Uid,
		Salt:  user.Salt,
		Rands: user.Rands,
	}

	if userUpdateReq.Email != "" && userUpdateReq.Email != user.Email {
		user.Email = userUpdateReq.Email
		userNew.Email = userUpdateReq.Email
		anythingUpdate = true
	}

	if userUpdateReq.Password != "" {
		if !a.users.IsPasswordEqualsUserPassword(userUpdateReq.OldPassword, user) {
			return nil, errs.ErrUserPasswordWrong
		}

		if !a.users.IsPasswordEqualsUserPassword(userUpdateReq.Password, user) {
			userNew.Password = userUpdateReq.Password
			anythingUpdate = true
		}
	}

	if userUpdateReq.Nickname != "" && userUpdateReq.Nickname != user.Nickname {
		user.Nickname = userUpdateReq.Nickname
		userNew.Nickname = userUpdateReq.Nickname
		anythingUpdate = true
	}

	if userUpdateReq.DefaultCurrency != "" && userUpdateReq.DefaultCurrency != user.DefaultCurrency {
		user.DefaultCurrency = userUpdateReq.DefaultCurrency
		userNew.DefaultCurrency = userUpdateReq.DefaultCurrency
		anythingUpdate = true
	}

	if !anythingUpdate {
		return nil, errs.ErrNothingWillBeUpdated
	}

	keyProfileUpdated, err := a.users.UpdateUser(userNew)

	if err != nil {
		log.ErrorfWithRequestId(c, "[users.UserUpdateProfileHandler] failed to update user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[users.UserUpdateProfileHandler] user \"uid:%d\" has updated successfully", user.Uid)

	resp := &models.UserProfileUpdateResponse{
		User: user.ToUserBasicInfo(),
	}

	if keyProfileUpdated {
		now := time.Now().Unix()
		err = a.tokens.DeleteTokensBeforeTime(uid, now)

		if err == nil {
			log.InfofWithRequestId(c, "[users.UserUpdateProfileHandler] revoke old tokens before unix time \"%d\" for user \"uid:%d\"", now, user.Uid)
		} else {
			log.WarnfWithRequestId(c, "[users.UserUpdateProfileHandler] failed to revoke old tokens for user \"uid:%d\", because %s", user.Uid, err.Error())
		}

		token, claims, err := a.tokens.CreateToken(user, c)

		if err != nil {
			log.WarnfWithRequestId(c, "[users.UserUpdateProfileHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
			return resp, nil
		}

		resp.NewToken = token
		c.SetTokenClaims(claims)

		log.InfofWithRequestId(c, "[users.UserUpdateProfileHandler] user \"uid:%d\" token refreshed, new token will be expired at %d", user.Uid, claims.ExpiresAt)

		return resp, nil
	}

	return resp, nil
}
