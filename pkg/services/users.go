package services

import (
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

// UserService represents user service
type UserService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a user service singleton instance
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

// GetUserByUsernameOrEmailAndPassword returns the user model according to login name and password
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

// GetUserById returns the user model according to user uid
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

// GetUserByUsername returns the user model according to user name
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

// GetUserByEmail returns the user model according to user email
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

// CreateUser saves a new user model to database
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

	user.Uid = s.GenerateUuid(uuid.UUID_TYPE_USER)
	user.Password = utils.EncodePassword(user.Password, user.Salt)

	user.Deleted = false

	user.CreatedUnixTime = time.Now().Unix()
	user.UpdatedUnixTime = time.Now().Unix()
	user.LastLoginUnixTime = time.Now().Unix()

	return s.UserDB().DoTransaction(func(sess *xorm.Session) error {
		_, err := sess.Insert(user)
		return err
	})
}

// UpdateUser saves an existed user model to database
func (s *UserService) UpdateUser(user *models.User) (keyProfileUpdated bool, err error) {
	if user.Uid <= 0 {
		return false, errs.ErrUserIdInvalid
	}

	updateCols := make([]string, 0, 8)

	now := time.Now().Unix()
	keyProfileUpdated = false

	if user.Email != "" {
		exists, err := s.ExistsEmail(user.Email)

		if err != nil {
			return false, err
		} else if exists {
			return false, errs.ErrUserEmailAlreadyExists
		}

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

	if user.DefaultAccountId > 0 {
		updateCols = append(updateCols, "default_account_id")
	}

	if models.WEEKDAY_SUNDAY <= user.FirstDayOfWeek && user.FirstDayOfWeek <= models.WEEKDAY_SATURDAY {
		updateCols = append(updateCols, "first_day_of_week")
	}

	if models.TRANSACTION_EDIT_SCOPE_NONE <= user.TransactionEditScope && user.TransactionEditScope <= models.TRANSACTION_EDIT_SCOPE_THIS_YEAR_OR_LATER {
		updateCols = append(updateCols, "transaction_edit_scope")
	}

	user.UpdatedUnixTime = now
	updateCols = append(updateCols, "updated_unix_time")

	err = s.UserDB().DoTransaction(func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(user.Uid).Cols(updateCols...).Where("deleted=?", false).Update(user)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrUserNotFound
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return keyProfileUpdated, nil
}

// UpdateUserLastLoginTime updates the last login time field
func (s *UserService) UpdateUserLastLoginTime(uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.UserDB().DoTransaction(func(sess *xorm.Session) error {
		_, err := sess.ID(uid).Cols("last_login_unix_time").Where("deleted=?", false).Update(&models.User{LastLoginUnixTime: time.Now().Unix()})
		return err
	})
}

// DeleteUser deletes an existed user from database
func (s *UserService) DeleteUser(username string) error {
	if username == "" {
		return errs.ErrUsernameIsEmpty
	}

	now := time.Now().Unix()

	updateModel := &models.User{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	deletedRows, err := s.UserDB().Cols("deleted", "deleted_unix_time").Where("username=? AND deleted=?", username, false).Update(updateModel)

	if err != nil {
		return err
	} else if deletedRows < 1 {
		return errs.ErrUserNotFound
	}
	return nil
}

// ExistsUsername returns whether the given user name exists
func (s *UserService) ExistsUsername(username string) (bool, error) {
	if username == "" {
		return false, errs.ErrUsernameIsEmpty
	}

	return s.UserDB().Cols("username").Where("username=? AND deleted=?", username, false).Exist(&models.User{})
}

// ExistsEmail returns whether the given user email exists
func (s *UserService) ExistsEmail(email string) (bool, error) {
	if email == "" {
		return false, errs.ErrEmailIsEmpty
	}

	return s.UserDB().Cols("email").Where("email=? AND deleted=?", email, false).Exist(&models.User{})
}

// IsPasswordEqualsUserPassword returns whether the given password is correct
func (s *UserService) IsPasswordEqualsUserPassword(password string, user *models.User) bool {
	return user.Password == utils.EncodePassword(password, user.Salt)
}
