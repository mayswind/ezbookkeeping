package services

import (
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/lab/pkg/datastore"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/utils"
	"github.com/mayswind/lab/pkg/uuid"
)

type UserService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

var (
	Users = &UserService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

func (s *UserService) GetUserByUsernameOrEmailAndPassword(loginname string, password string) (*models.User, error) {
	var user *models.User
	var err error

	if utils.IsValidUsername(loginname) {
		user, err = s.GetUserByUsername(loginname)
	} else if utils.IsValidEmail(loginname) {
		user, err = s.GetUserByEmail(loginname)
	} else {
		err = errs.ErrLoginNameInvalid
	}

	if err != nil {
		return nil, err
	}

	if !s.IsPasswordEqualsUserPassword(password, user) {
		return nil, errs.ErrUserPasswordWrong
	}

	return user, nil
}

func (s *UserService) GetUserById(uid int64) (*models.User, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	user := &models.User{}
	has, err := s.UserDB().ID(uid).Where("deleted=?", false).Get(user)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	if username == "" {
		return nil, errs.ErrUsernameIsEmpty
	}

	user := &models.User{}
	has, err := s.UserDB().Where("username=? AND deleted=?", username, false).Get(user)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, errs.ErrEmailIsEmpty
	}

	user := &models.User{}
	has, err := s.UserDB().Where("email=? AND deleted=?", email, false).Get(user)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) CreateUser(user *models.User) error {
	exists, err := s.ExistsUsername(user.Username)

	if err != nil {
		return err
	} else if exists {
		return errs.ErrUsernameAlreadyExists
	}

	exists, err = s.ExistsEmail(user.Email)

	if err != nil {
		return err
	} else if exists {
		return errs.ErrUserEmailAlreadyExists
	}

	if user.Password == "" {
		return errs.ErrPasswordIsEmpty
	}

	if user.Salt, err = utils.GetRandomString(10); err != nil {
		return err
	}

	if user.Rands, err = utils.GetRandomString(10); err != nil {
		return err
	}

	user.Uid = s.GenerateUuid(uuid.UUID_TYPE_USER)
	user.Password = utils.EncodePassword(user.Password, user.Salt)

	user.Deleted = false

	user.CreatedUnixTime = time.Now().Unix()
	user.UpdatedUnixTime = time.Now().Unix()
	user.LastLoginUnixTime = time.Now().Unix()

	return s.UserDB().DoTranscation(func(sess *xorm.Session) error {
		_, err := sess.Insert(user)
		return err
	})
}

func (s *UserService) UpdateUser(user *models.User) (keyProfileUpdated bool, err error) {
	if user.Uid <= 0 {
		return false, errs.ErrUserIdInvalid
	}

	var updateCols []string

	now := time.Now().Unix()
	keyProfileUpdated = false

	if user.Email != "" {
		user.EmailVerified = false

		updateCols = append(updateCols, "email")
		updateCols = append(updateCols, "email_verified")
	}

	if user.Password != "" {
		user.Password = utils.EncodePassword(user.Password, user.Salt)

		keyProfileUpdated = true
		updateCols = append(updateCols, "password")
	}

	if user.Nickname != "" {
		updateCols = append(updateCols, "nickname")
	}

	if user.DefaultCurrency != "" {
		updateCols = append(updateCols, "default_currency")
	}

	user.UpdatedUnixTime = now
	updateCols = append(updateCols, "updated_unix_time")

	err = s.UserDB().DoTranscation(func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(user.Uid).Where("deleted=?", false).Cols(updateCols...).Update(user)

		if updatedRows < 1 {
			return errs.ErrUserNotFound
		}

		return err
	})

	if err != nil {
		return false, err
	}

	return keyProfileUpdated, nil
}

func (s *UserService) UpdateUserLastLoginTime(uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.UserDB().DoTranscation(func(sess *xorm.Session) error {
		_, err := sess.ID(uid).Where("deleted=?", false).Cols("last_login_unix_time").Update(&models.User{LastLoginUnixTime: time.Now().Unix()})
		return err
	})
}

func (s *UserService) ExistsUsername(username string) (bool, error) {
	if username == "" {
		return false, errs.ErrUsernameIsEmpty
	}

	return s.UserDB().Cols("username").Where("username=? AND deleted=?", username, false).Exist(&models.User{})
}

func (s *UserService) ExistsEmail(email string) (bool, error) {
	if email == "" {
		return false, errs.ErrEmailIsEmpty
	}

	return s.UserDB().Cols("email").Where("email=? AND deleted=?", email, false).Exist(&models.User{})
}

func (s *UserService) IsPasswordEqualsUserPassword(password string, user *models.User) bool {
	return user.Password == utils.EncodePassword(password, user.Salt)
}
