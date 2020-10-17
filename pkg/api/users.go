package api

import (
	"strings"
	"time"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/services"
	"github.com/mayswind/lab/pkg/utils"
)

type UsersApi struct {
	users *services.UserService
	tokens *services.TokenService
}

var (
	Users = &UsersApi{
		users: services.Users,
		tokens: services.Tokens,
	}
)

func (a *UsersApi) UserRegisterHandler(c *core.Context) (interface{}, *errs.Error) {
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
		Username: userRegisterReq.Username,
		Email:    userRegisterReq.Email,
		Nickname: userRegisterReq.Nickname,
		Password: userRegisterReq.Password,
	}

	err = a.users.CreateUser(user)

	if err != nil {
		log.ErrorfWithRequestId(c, "[users.UserRegisterHandler] failed to create user \"%s\", because %s", user.Username, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[users.UserRegisterHandler] user \"%s\" has registered successfully, uid is %d", user.Username, user.Uid)

	token, claims, err := a.tokens.CreateToken(user, c)

	if err != nil {
		log.WarnfWithRequestId(c, "[users.UserRegisterHandler] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return true, nil
	}

	c.SetTokenClaims(claims)

	log.InfofWithRequestId(c, "[users.UserRegisterHandler] user \"uid:%d\" has logined, token will be expired at %d", user.Uid, claims.ExpiresAt)
	return token, nil
}

func (a *UsersApi) UserProfileHandler(c *core.Context) (interface{}, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.ErrorfWithRequestId(c, "[users.UserRegisterHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	userResp := &models.UserProfileResponse{
		Uid : utils.Int64ToString(user.Uid),
		Username: user.Username,
		Email: user.Email,
		Nickname: user.Nickname,
		Type: user.Type,
		CreatedAt: user.CreatedUnixTime,
		UpdatedAt: user.UpdatedUnixTime,
		LastLoginAt: user.LastLoginUnixTime,
	}

	return userResp, nil
}

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

	if userUpdateReq.Email != "" && userUpdateReq.Email != user.Email {
		anythingUpdate = true
	} else {
		userUpdateReq.Email = ""
	}

	if userUpdateReq.Password != "" && !a.users.IsPasswordEqualsUserPassword(userUpdateReq.Password, user) {
		anythingUpdate = true
	} else {
		userUpdateReq.Password = ""
	}

	if userUpdateReq.Nickname != "" && userUpdateReq.Nickname != user.Nickname {
		anythingUpdate = true
	} else {
		userUpdateReq.Nickname = ""
	}

	if !anythingUpdate {
		return nil, errs.ErrNothingWillBeUpdated
	}

	user.Email = userUpdateReq.Email
	user.Password = userUpdateReq.Password
	user.Nickname = userUpdateReq.Nickname

	keyProfileUpdated, err := a.users.UpdateUser(user)

	if err != nil {
		log.ErrorfWithRequestId(c, "[users.UserUpdateProfileHandler] failed to update user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[users.UserUpdateProfileHandler] user \"uid:%d\" has updated successfully", user.Uid)

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
			return true, nil
		}

		c.SetTokenClaims(claims)

		log.InfofWithRequestId(c, "[users.UserUpdateProfileHandler] user \"uid:%d\" token refreshed, new token will be expired at %d", user.Uid, claims.ExpiresAt)
		return token, nil
	} else {
		return true, nil
	}
}
