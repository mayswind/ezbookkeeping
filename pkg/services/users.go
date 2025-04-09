package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/locales"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/mail"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/storage"
	"github.com/mayswind/ezbookkeeping/pkg/templates"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

const verifyEmailUrlFormat = "%sdesktop/#/verify_email?token=%s"

// UserService represents user service
type UserService struct {
	ServiceUsingDB
	ServiceUsingConfig
	ServiceUsingMailer
	ServiceUsingUuid
	ServiceUsingStorage
}

// Initialize a user service singleton instance
var (
	Users = &UserService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingConfig: ServiceUsingConfig{
			container: settings.Container,
		},
		ServiceUsingMailer: ServiceUsingMailer{
			container: mail.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
		ServiceUsingStorage: ServiceUsingStorage{
			container: storage.Container,
		},
	}
)

// GetUserByUsernameOrEmailAndPassword returns the user model according to login name and password
func (s *UserService) GetUserByUsernameOrEmailAndPassword(c core.Context, loginname string, password string) (*models.User, int64, error) {
	var user *models.User
	var err error

	if utils.IsValidUsername(loginname) {
		user, err = s.GetUserByUsername(c, loginname)
	} else if utils.IsValidEmail(loginname) {
		user, err = s.GetUserByEmail(c, loginname)
	} else {
		err = errs.ErrLoginNameInvalid
	}

	if err != nil {
		return nil, 0, err
	}

	if user == nil {
		return nil, 0, errs.ErrUserNotFound
	}

	if !s.IsPasswordEqualsUserPassword(password, user) {
		return nil, user.Uid, errs.ErrUserPasswordWrong
	}

	return user, user.Uid, nil
}

// GetUserById returns the user model according to user uid
func (s *UserService) GetUserById(c core.Context, uid int64) (*models.User, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	user := &models.User{}
	has, err := s.UserDB().NewSession(c).ID(uid).Where("deleted=?", false).Get(user)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrUserNotFound
	}

	return user, nil
}

// GetUserByUsername returns the user model according to user name
func (s *UserService) GetUserByUsername(c core.Context, username string) (*models.User, error) {
	if username == "" {
		return nil, errs.ErrUsernameIsEmpty
	}

	user := &models.User{}
	has, err := s.UserDB().NewSession(c).Where("username=? AND deleted=?", username, false).Get(user)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrUserNotFound
	}

	return user, nil
}

// GetUserByEmail returns the user model according to user email
func (s *UserService) GetUserByEmail(c core.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errs.ErrEmailIsEmpty
	}

	user := &models.User{}
	has, err := s.UserDB().NewSession(c).Where("email=? AND deleted=?", email, false).Get(user)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrUserNotFound
	}

	return user, nil
}

// GetUserAvatar returns the user avatar image data according to user uid
func (s *UserService) GetUserAvatar(c core.Context, uid int64, fileExtension string) ([]byte, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	user := &models.User{}
	has, err := s.UserDB().NewSession(c).ID(uid).Cols("uid", "deleted", "custom_avatar_type").Where("deleted=?", false).Get(user)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrUserNotFound
	}

	if user.CustomAvatarType == "" {
		return nil, errs.ErrUserAvatarNotSet
	}

	if user.CustomAvatarType != fileExtension {
		return nil, errs.ErrUserAvatarExtensionInvalid
	}

	avatarFile, err := s.ReadAvatar(user.Uid, user.CustomAvatarType)

	if os.IsNotExist(err) {
		return nil, errs.ErrUserAvatarNoExists
	}

	if err != nil {
		return nil, err
	}

	defer avatarFile.Close()

	avatarData, err := io.ReadAll(avatarFile)

	if err != nil {
		return nil, err
	}

	return avatarData, nil
}

// CreateUser saves a new user model to database
func (s *UserService) CreateUser(c core.Context, user *models.User) error {
	exists, err := s.ExistsUsername(c, user.Username)

	if err != nil {
		return err
	} else if exists {
		return errs.ErrUsernameAlreadyExists
	}

	exists, err = s.ExistsEmail(c, user.Email)

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

	if user.Uid < 1 {
		return errs.ErrSystemIsBusy
	}

	user.Password = utils.EncodePassword(user.Password, user.Salt)

	user.Deleted = false

	user.CreatedUnixTime = time.Now().Unix()
	user.UpdatedUnixTime = time.Now().Unix()
	user.LastLoginUnixTime = time.Now().Unix()

	return s.UserDB().DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(user)
		return err
	})
}

// UpdateUser saves an existed user model to database
func (s *UserService) UpdateUser(c core.Context, user *models.User, modifyUserLanguage bool) (keyProfileUpdated bool, emailSetToUnverified bool, err error) {
	if user.Uid <= 0 {
		return false, false, errs.ErrUserIdInvalid
	}

	updateCols := make([]string, 0, 8)

	now := time.Now().Unix()
	keyProfileUpdated = false
	emailSetToUnverified = false

	if user.Email != "" {
		exists, err := s.ExistsEmail(c, user.Email)

		if err != nil {
			return false, false, err
		} else if exists {
			return false, false, errs.ErrUserEmailAlreadyExists
		}

		user.EmailVerified = false

		updateCols = append(updateCols, "email")
		updateCols = append(updateCols, "email_verified")

		emailSetToUnverified = true
	}

	if user.Password != "" {
		user.Password = utils.EncodePassword(user.Password, user.Salt)

		keyProfileUpdated = true
		updateCols = append(updateCols, "password")
	}

	if user.Nickname != "" {
		updateCols = append(updateCols, "nickname")
	}

	if user.DefaultAccountId > 0 {
		updateCols = append(updateCols, "default_account_id")
	}

	if models.TRANSACTION_EDIT_SCOPE_NONE <= user.TransactionEditScope && user.TransactionEditScope <= models.TRANSACTION_EDIT_SCOPE_THIS_YEAR_OR_LATER {
		updateCols = append(updateCols, "transaction_edit_scope")
	}

	if modifyUserLanguage || user.Language != "" {
		updateCols = append(updateCols, "language")
	}

	if user.DefaultCurrency != "" {
		updateCols = append(updateCols, "default_currency")
	}

	if core.WEEKDAY_SUNDAY <= user.FirstDayOfWeek && user.FirstDayOfWeek <= core.WEEKDAY_SATURDAY {
		updateCols = append(updateCols, "first_day_of_week")
	}

	if core.FISCAL_YEAR_START_DEFAULT <= user.FiscalYearStart && core.FISCAL_YEAR_START_MIN <= user.FiscalYearStart && core.FISCAL_YEAR_START_MAX >= user.FiscalYearStart {
		updateCols = append(updateCols, "fiscal_year_start")
	}

	if core.LONG_DATE_FORMAT_DEFAULT <= user.LongDateFormat && user.LongDateFormat <= core.LONG_DATE_FORMAT_D_M_YYYY {
		updateCols = append(updateCols, "long_date_format")
	}

	if core.SHORT_DATE_FORMAT_DEFAULT <= user.ShortDateFormat && user.ShortDateFormat <= core.SHORT_DATE_FORMAT_D_M_YYYY {
		updateCols = append(updateCols, "short_date_format")
	}

	if core.LONG_TIME_FORMAT_DEFAULT <= user.LongTimeFormat && user.LongTimeFormat <= core.LONG_TIME_FORMAT_HH_MM_SS_A {
		updateCols = append(updateCols, "long_time_format")
	}

	if core.SHORT_TIME_FORMAT_DEFAULT <= user.ShortTimeFormat && user.ShortTimeFormat <= core.SHORT_TIME_FORMAT_HH_MM_A {
		updateCols = append(updateCols, "short_time_format")
	}

	if core.DECIMAL_SEPARATOR_DEFAULT <= user.DecimalSeparator && user.DecimalSeparator <= core.DECIMAL_SEPARATOR_COMMA {
		updateCols = append(updateCols, "decimal_separator")
	}

	if core.DIGIT_GROUPING_SYMBOL_DEFAULT <= user.DigitGroupingSymbol && user.DigitGroupingSymbol <= core.DIGIT_GROUPING_SYMBOL_APOSTROPHE {
		updateCols = append(updateCols, "digit_grouping_symbol")
	}

	if core.DIGIT_GROUPING_TYPE_DEFAULT <= user.DigitGrouping && user.DigitGrouping <= core.DIGIT_GROUPING_TYPE_THOUSANDS_SEPARATOR {
		updateCols = append(updateCols, "digit_grouping")
	}

	if core.CURRENCY_DISPLAY_TYPE_DEFAULT <= user.CurrencyDisplayType && user.CurrencyDisplayType <= core.CURRENCY_DISPLAY_TYPE_NAME_AFTER_AMOUNT {
		updateCols = append(updateCols, "currency_display_type")
	}

	if models.AMOUNT_COLOR_TYPE_DEFAULT <= user.ExpenseAmountColor && user.ExpenseAmountColor <= models.AMOUNT_COLOR_TYPE_BLACK_OR_WHITE {
		updateCols = append(updateCols, "expense_amount_color")
	}

	if models.AMOUNT_COLOR_TYPE_DEFAULT <= user.IncomeAmountColor && user.IncomeAmountColor <= models.AMOUNT_COLOR_TYPE_BLACK_OR_WHITE {
		updateCols = append(updateCols, "income_amount_color")
	}

	user.UpdatedUnixTime = now
	updateCols = append(updateCols, "updated_unix_time")

	err = s.UserDB().DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(user.Uid).Cols(updateCols...).Where("deleted=?", false).Update(user)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrUserNotFound
		}

		return nil
	})

	if err != nil {
		return false, false, err
	}

	return keyProfileUpdated, emailSetToUnverified, nil
}

// UpdateUserAvatar updates the custom avatar type of specified user
func (s *UserService) UpdateUserAvatar(c core.Context, uid int64, avatarFile multipart.File, fileExtension string, oldFileExtension string) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	defer avatarFile.Close()

	err := s.SaveAvatar(uid, avatarFile, fileExtension)

	if err != nil {
		return err
	}

	now := time.Now().Unix()

	updateModel := &models.User{
		CustomAvatarType: fileExtension,
		UpdatedUnixTime:  now,
	}

	err = s.UserDB().DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.ID(uid).Cols("custom_avatar_type", "updated_unix_time").Where("deleted=?", false).Update(updateModel)
		return err
	})

	if err != nil {
		return err
	}

	if fileExtension != oldFileExtension && oldFileExtension != "" {
		err = s.DeleteAvatar(uid, oldFileExtension)

		if err != nil {
			log.Warnf(c, "[users.UpdateUserAvatar] failed to delete old avatar with extension \"%s\" for user \"uid:%d\", because %s", oldFileExtension, uid, err.Error())
		}
	}

	return nil
}

// RemoveUserAvatar removes the custom avatar type of specified user
func (s *UserService) RemoveUserAvatar(c core.Context, uid int64, fileExtension string) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	err := s.DeleteAvatar(uid, fileExtension)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	now := time.Now().Unix()

	updateModel := &models.User{
		CustomAvatarType: "",
		UpdatedUnixTime:  now,
	}

	return s.UserDB().DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.ID(uid).Cols("custom_avatar_type", "updated_unix_time").Where("deleted=?", false).Update(updateModel)
		return err
	})
}

// UpdateUserLastLoginTime updates the last login time field
func (s *UserService) UpdateUserLastLoginTime(c core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.UserDB().DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.ID(uid).Cols("last_login_unix_time").Where("deleted=?", false).Update(&models.User{LastLoginUnixTime: time.Now().Unix()})
		return err
	})
}

// EnableUser sets user enabled
func (s *UserService) EnableUser(c core.Context, username string) error {
	if username == "" {
		return errs.ErrUsernameIsEmpty
	}

	now := time.Now().Unix()

	updateModel := &models.User{
		Disabled:        false,
		UpdatedUnixTime: now,
	}

	updatedRows, err := s.UserDB().NewSession(c).Cols("disabled", "updated_unix_time").Where("username=? AND deleted=?", username, false).Update(updateModel)

	if err != nil {
		return err
	} else if updatedRows < 1 {
		return errs.ErrUserNotFound
	}
	return nil
}

// DisableUser sets user disabled
func (s *UserService) DisableUser(c core.Context, username string) error {
	if username == "" {
		return errs.ErrUsernameIsEmpty
	}

	now := time.Now().Unix()

	updateModel := &models.User{
		Disabled:        true,
		UpdatedUnixTime: now,
	}

	updatedRows, err := s.UserDB().NewSession(c).Cols("disabled", "updated_unix_time").Where("username=? AND deleted=?", username, false).Update(updateModel)

	if err != nil {
		return err
	} else if updatedRows < 1 {
		return errs.ErrUserNotFound
	}
	return nil
}

// UpdateUserFeatureRestriction sets user user feature restrictions
func (s *UserService) UpdateUserFeatureRestriction(c core.Context, username string, featureRestriction core.UserFeatureRestrictions) error {
	if username == "" {
		return errs.ErrUsernameIsEmpty
	}

	now := time.Now().Unix()

	updateModel := &models.User{
		FeatureRestriction: featureRestriction,
		UpdatedUnixTime:    now,
	}

	updatedRows, err := s.UserDB().NewSession(c).Cols("feature_restriction", "updated_unix_time").Where("username=? AND deleted=?", username, false).Update(updateModel)

	if err != nil {
		return err
	} else if updatedRows < 1 {
		return errs.ErrUserNotFound
	}
	return nil
}

// AddUserFeatureRestriction adds user user feature restrictions
func (s *UserService) AddUserFeatureRestriction(c core.Context, username string, featureRestriction core.UserFeatureRestrictions) error {
	if username == "" {
		return errs.ErrUsernameIsEmpty
	}

	now := time.Now().Unix()

	updateModel := &models.User{
		FeatureRestriction: featureRestriction,
		UpdatedUnixTime:    now,
	}

	updatedRows, err := s.UserDB().NewSession(c).SetExpr("feature_restriction", fmt.Sprintf("feature_restriction|(%d)", updateModel.FeatureRestriction)).Cols("updated_unix_time").Where("username=? AND deleted=?", username, false).Update(updateModel)

	if err != nil {
		return err
	} else if updatedRows < 1 {
		return errs.ErrUserNotFound
	}
	return nil
}

// RemoveUserFeatureRestriction removes user user feature restrictions
func (s *UserService) RemoveUserFeatureRestriction(c core.Context, username string, featureRestriction core.UserFeatureRestrictions) error {
	if username == "" {
		return errs.ErrUsernameIsEmpty
	}

	now := time.Now().Unix()

	updateModel := &models.User{
		FeatureRestriction: ^featureRestriction,
		UpdatedUnixTime:    now,
	}

	updatedRows, err := s.UserDB().NewSession(c).SetExpr("feature_restriction", fmt.Sprintf("feature_restriction&(%d)", updateModel.FeatureRestriction)).Cols("updated_unix_time").Where("username=? AND deleted=?", username, false).Update(updateModel)

	if err != nil {
		return err
	} else if updatedRows < 1 {
		return errs.ErrUserNotFound
	}
	return nil
}

// SetUserEmailVerified sets user email address verified
func (s *UserService) SetUserEmailVerified(c core.Context, username string) error {
	if username == "" {
		return errs.ErrUsernameIsEmpty
	}

	now := time.Now().Unix()

	updateModel := &models.User{
		EmailVerified:   true,
		UpdatedUnixTime: now,
	}

	updatedRows, err := s.UserDB().NewSession(c).Cols("email_verified", "updated_unix_time").Where("username=? AND deleted=?", username, false).Update(updateModel)

	if err != nil {
		return err
	} else if updatedRows < 1 {
		return errs.ErrUserNotFound
	}
	return nil
}

// SetUserEmailUnverified sets user email address unverified
func (s *UserService) SetUserEmailUnverified(c core.Context, username string) error {
	if username == "" {
		return errs.ErrUsernameIsEmpty
	}

	now := time.Now().Unix()

	updateModel := &models.User{
		EmailVerified:   false,
		UpdatedUnixTime: now,
	}

	updatedRows, err := s.UserDB().NewSession(c).Cols("email_verified", "updated_unix_time").Where("username=? AND deleted=?", username, false).Update(updateModel)

	if err != nil {
		return err
	} else if updatedRows < 1 {
		return errs.ErrUserNotFound
	}
	return nil
}

// DeleteUser deletes an existed user from database
func (s *UserService) DeleteUser(c core.Context, username string) error {
	if username == "" {
		return errs.ErrUsernameIsEmpty
	}

	now := time.Now().Unix()

	updateModel := &models.User{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	deletedRows, err := s.UserDB().NewSession(c).Cols("deleted", "deleted_unix_time").Where("username=? AND deleted=?", username, false).Update(updateModel)

	if err != nil {
		return err
	} else if deletedRows < 1 {
		return errs.ErrUserNotFound
	}
	return nil
}

// ExistsUsername returns whether the given user name exists
func (s *UserService) ExistsUsername(c core.Context, username string) (bool, error) {
	if username == "" {
		return false, errs.ErrUsernameIsEmpty
	}

	return s.UserDB().NewSession(c).Cols("username").Where("username=? AND deleted=?", username, false).Exist(&models.User{})
}

// ExistsEmail returns whether the given user email exists
func (s *UserService) ExistsEmail(c core.Context, email string) (bool, error) {
	if email == "" {
		return false, errs.ErrEmailIsEmpty
	}

	return s.UserDB().NewSession(c).Cols("email").Where("email=? AND deleted=?", email, false).Exist(&models.User{})
}

// SendVerifyEmail sends verify email according to specified parameters
func (s *UserService) SendVerifyEmail(user *models.User, verifyEmailToken string, backupLocale string) error {
	if !s.CurrentConfig().EnableSMTP {
		return errs.ErrSMTPServerNotEnabled
	}

	locale := user.Language

	if locale == "" {
		locale = backupLocale
	}

	localeTextItems := locales.GetLocaleTextItems(locale)
	verifyEmailTextItems := localeTextItems.VerifyEmailTextItems

	expireTimeInMinutes := s.CurrentConfig().EmailVerifyTokenExpiredTimeDuration.Minutes()
	verifyEmailUrl := fmt.Sprintf(verifyEmailUrlFormat, s.CurrentConfig().RootUrl, url.QueryEscape(verifyEmailToken))

	tmpl, err := templates.GetTemplate(templates.TEMPLATE_VERIFY_EMAIL)

	if err != nil {
		return err
	}

	templateParams := map[string]any{
		"AppName": s.CurrentConfig().AppName,
		"VerifyEmail": map[string]any{
			"Title":               verifyEmailTextItems.Title,
			"Salutation":          fmt.Sprintf(verifyEmailTextItems.SalutationFormat, user.Nickname),
			"DescriptionAboveBtn": verifyEmailTextItems.DescriptionAboveBtn,
			"VerifyEmailUrl":      verifyEmailUrl,
			"VerifyEmail":         verifyEmailTextItems.VerifyEmail,
			"DescriptionBelowBtn": fmt.Sprintf(verifyEmailTextItems.DescriptionBelowBtnFormat, s.CurrentConfig().AppName, expireTimeInMinutes),
		},
	}

	var bodyBuffer bytes.Buffer
	err = tmpl.Execute(&bodyBuffer, templateParams)

	if err != nil {
		return err
	}

	message := &mail.MailMessage{
		To:      user.Email,
		Subject: verifyEmailTextItems.Title,
		Body:    bodyBuffer.String(),
	}

	err = s.SendMail(message)

	return err
}

// IsPasswordEqualsUserPassword returns whether the given password is correct
func (s *UserService) IsPasswordEqualsUserPassword(password string, user *models.User) bool {
	return user.Password == utils.EncodePassword(password, user.Salt)
}
