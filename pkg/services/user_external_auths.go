package services

import (
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// UserExternalAuthService represents user external auth service
type UserExternalAuthService struct {
	ServiceUsingDB
}

// Initialize a user external auth service singleton instance
var (
	UserExternalAuths = &UserExternalAuthService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
	}
)

// GetUserAllExternalAuthsByUid returns the user all external auth list according to user uid
func (s *UserExternalAuthService) GetUserAllExternalAuthsByUid(c core.Context, uid int64) ([]*models.UserExternalAuth, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var userExternalAuths []*models.UserExternalAuth
	err := s.UserDB().NewSession(c).Where("uid=?", uid).Find(&userExternalAuths)

	return userExternalAuths, err
}

// GetUserExternalAuthByUid returns the user external auth record by uid
func (s *UserExternalAuthService) GetUserExternalAuthByUid(c core.Context, uid int64, externalAuthType core.UserExternalAuthType) (*models.UserExternalAuth, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	userExternalAuth := &models.UserExternalAuth{}
	has, err := s.UserDB().NewSession(c).Where("uid=? AND external_auth_type=?", uid, externalAuthType).Get(userExternalAuth)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrUserExternalAuthNotFound
	}

	return userExternalAuth, err
}

// GetUserExternalAuthByExternalUserName returns the user external auth record by external username
func (s *UserExternalAuthService) GetUserExternalAuthByExternalUserName(c core.Context, externalUserName string, externalAuthType core.UserExternalAuthType) (*models.UserExternalAuth, error) {
	userExternalAuth := &models.UserExternalAuth{}
	has, err := s.UserDB().NewSession(c).Where("external_auth_type=? AND external_username=?", externalAuthType, externalUserName).Get(userExternalAuth)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrUserExternalAuthNotFound
	}

	return userExternalAuth, err
}

// GetUserExternalAuthByExternalEmail returns the user external auth record by external email
func (s *UserExternalAuthService) GetUserExternalAuthByExternalEmail(c core.Context, externalEmail string, externalAuthType core.UserExternalAuthType) (*models.UserExternalAuth, error) {
	userExternalAuth := &models.UserExternalAuth{}
	has, err := s.UserDB().NewSession(c).Where("external_auth_type=? AND external_email=?", externalAuthType, externalEmail).Get(userExternalAuth)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrUserExternalAuthNotFound
	}

	return userExternalAuth, err
}

// CreateUserExternalAuth creates a new user external auth record in database
func (s *UserExternalAuthService) CreateUserExternalAuth(c core.Context, userExternalAuth *models.UserExternalAuth) error {
	if userExternalAuth.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	userExternalAuth.CreatedUnixTime = time.Now().Unix()

	return s.UserDB().DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(userExternalAuth)
		return err
	})
}

// DeleteUserExternalAuth deletes given user external auth record from database
func (s *UserExternalAuthService) DeleteUserExternalAuth(c core.Context, uid int64, externalAuthType core.UserExternalAuthType) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.UserDB().DoTransaction(c, func(sess *xorm.Session) error {
		deletedRows, err := sess.Where("uid=? AND external_auth_type=?", uid, externalAuthType).Delete(&models.UserExternalAuth{})

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrUserExternalAuthNotFound
		}

		return nil
	})
}
