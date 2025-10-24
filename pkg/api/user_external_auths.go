package api

import (
	"sort"

	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
)

// UserExternalAuthsApi represents user external auth api
type UserExternalAuthsApi struct {
	users             *services.UserService
	userExternalAuths *services.UserExternalAuthService
}

// Initialize a user external auth api singleton instance
var (
	UserExternalAuths = &UserExternalAuthsApi{
		users:             services.Users,
		userExternalAuths: services.UserExternalAuths,
	}
)

// ExternalAuthListHanlder returns external authentications list of current user
func (a *UserExternalAuthsApi) ExternalAuthListHanlder(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	userExternalAuths, err := a.userExternalAuths.GetUserAllExternalAuthsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[user_external_auths.ExternalAuthListHanlder] failed to get all external authentications for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	userExternalAuthResps := make(models.UserExternalAuthInfoResponsesSlice, 0, len(userExternalAuths)+1)
	currentExternalAuthType := oauth2.GetExternalUserAuthType()
	hasCurrentExternalAuth := false

	for i := 0; i < len(userExternalAuths); i++ {
		userExternalAuth := userExternalAuths[i]

		if userExternalAuth.ExternalAuthType == currentExternalAuthType {
			hasCurrentExternalAuth = true
		}

		userExternalAuthResps = append(userExternalAuthResps, userExternalAuth.ToUserExternalAuthInfoResponse())
	}

	if !hasCurrentExternalAuth {
		userExternalAuthResps = append(userExternalAuthResps, &models.UserExternalAuthInfoResponse{
			ExternalAuthCategory: currentExternalAuthType.GetCategory(),
			ExternalAuthType:     currentExternalAuthType,
			Linked:               false,
		})
	}

	sort.Sort(userExternalAuthResps)

	return userExternalAuthResps, nil
}

// UnlinkExternalAuthHandler unlinks external authentication for current user
func (a *UserExternalAuthsApi) UnlinkExternalAuthHandler(c *core.WebContext) (any, *errs.Error) {
	var externalAuthLinkReq models.UserExternalAuthUnlinkRequest
	err := c.ShouldBindJSON(&externalAuthLinkReq)

	if err != nil {
		log.Warnf(c, "[user_external_auths.UnlinkExternalAuthHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Warnf(c, "[user_external_auths.UnlinkExternalAuthHandler] failed to get user for user \"uid:%d\", because %s", uid, err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if !a.users.IsPasswordEqualsUserPassword(externalAuthLinkReq.Password, user) {
		return nil, errs.ErrUserPasswordWrong
	}

	externalAuthType := core.UserExternalAuthType(externalAuthLinkReq.ExternalAuthType)

	if !externalAuthType.IsValid() {
		return nil, errs.ErrUserExternalAuthNotFound
	}

	err = a.userExternalAuths.DeleteUserExternalAuth(c, uid, externalAuthType)

	if err != nil {
		log.Errorf(c, "[user_external_auths.UnlinkExternalAuthHandler] failed to unlink external authentication \"%s\" for user \"uid:%d\", because %s", externalAuthLinkReq.ExternalAuthType, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return true, nil
}
