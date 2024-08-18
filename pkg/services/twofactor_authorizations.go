package services

import (
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

const (
	twoFactorPeriod             uint = 30 // seconds
	twoFactorSecretSize         uint = 20 // bytes
	twoFactorRecoveryCodeCount  int  = 10
	twoFactorRecoveryCodeLength int  = 10 // bytes
)

// TwoFactorAuthorizationService represents 2fa service
type TwoFactorAuthorizationService struct {
	ServiceUsingDB
	ServiceUsingConfig
	ServiceUsingUuid
}

// Initialize a 2fa service singleton instance
var (
	TwoFactorAuthorizations = &TwoFactorAuthorizationService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingConfig: ServiceUsingConfig{
			container: settings.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetUserTwoFactorSettingByUid returns the 2fa setting model according to user uid
func (s *TwoFactorAuthorizationService) GetUserTwoFactorSettingByUid(c core.Context, uid int64) (*models.TwoFactor, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	twoFactor := &models.TwoFactor{}
	has, err := s.UserDB().NewSession(c).Where("uid=?", uid).Get(twoFactor)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrTwoFactorIsNotEnabled
	}

	twoFactor.Secret, err = utils.DecryptSecret(twoFactor.Secret, s.CurrentConfig().SecretKey)

	if err != nil {
		return nil, err
	}

	return twoFactor, nil
}

// GenerateTwoFactorSecret generates a new 2fa secret
func (s *TwoFactorAuthorizationService) GenerateTwoFactorSecret(c core.Context, user *models.User) (*otp.Key, error) {
	if user == nil {
		return nil, errs.ErrUserNotFound
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.CurrentConfig().AppName,
		AccountName: user.Username,
		Period:      twoFactorPeriod,
		SecretSize:  twoFactorSecretSize,
	})

	return key, err
}

// CreateTwoFactorSetting saves a new 2fa setting to database
func (s *TwoFactorAuthorizationService) CreateTwoFactorSetting(c core.Context, twoFactor *models.TwoFactor) error {
	if twoFactor.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	var err error
	twoFactor.Secret, err = utils.EncryptSecret(twoFactor.Secret, s.CurrentConfig().SecretKey)

	if err != nil {
		return err
	}

	twoFactor.CreatedUnixTime = time.Now().Unix()

	return s.UserDB().DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(twoFactor)
		return err
	})
}

// DeleteTwoFactorSetting deletes an existed 2fa setting from database
func (s *TwoFactorAuthorizationService) DeleteTwoFactorSetting(c core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.UserDB().DoTransaction(c, func(sess *xorm.Session) error {
		deletedRows, err := sess.Where("uid=?", uid).Delete(&models.TwoFactor{})

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrTwoFactorIsNotEnabled
		}

		return nil
	})
}

// ExistsTwoFactorSetting returns whether the given user has existed 2fa setting
func (s *TwoFactorAuthorizationService) ExistsTwoFactorSetting(c core.Context, uid int64) (bool, error) {
	if uid <= 0 {
		return false, errs.ErrUserIdInvalid
	}

	return s.UserDB().NewSession(c).Cols("uid").Where("uid=?", uid).Exist(&models.TwoFactor{})
}

// GetAndUseUserTwoFactorRecoveryCode checks whether the given 2fa recovery code exists and marks it used
func (s *TwoFactorAuthorizationService) GetAndUseUserTwoFactorRecoveryCode(c core.Context, uid int64, recoveryCode string, salt string) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	recoveryCode = utils.EncodePassword(recoveryCode, salt)
	exists, err := s.UserDB().NewSession(c).Cols("uid", "recovery_code").Where("uid=? AND recovery_code=? AND used=?", uid, recoveryCode, false).Exist(&models.TwoFactorRecoveryCode{})

	if err != nil {
		return err
	} else if !exists {
		return errs.ErrTwoFactorRecoveryCodeNotExist
	}

	return s.UserDB().DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Cols("used", "used_unix_time").Where("uid=? AND recovery_code=?", uid, recoveryCode).Update(&models.TwoFactorRecoveryCode{Used: true, UsedUnixTime: time.Now().Unix()})
		return err
	})
}

// GenerateTwoFactorRecoveryCodes generates new 2fa recovery codes
func (s *TwoFactorAuthorizationService) GenerateTwoFactorRecoveryCodes() ([]string, error) {
	recoveryCodes := make([]string, twoFactorRecoveryCodeCount)

	for i := 0; i < twoFactorRecoveryCodeCount; i++ {
		recoveryCode, err := utils.GetRandomNumberOrLowercaseLetter(twoFactorRecoveryCodeLength)

		if err != nil {
			return nil, err
		}

		recoveryCodes[i] = recoveryCode[:5] + "-" + recoveryCode[5:]
	}

	return recoveryCodes, nil
}

// CreateTwoFactorRecoveryCodes saves new 2fa recovery codes to database
func (s *TwoFactorAuthorizationService) CreateTwoFactorRecoveryCodes(c core.Context, uid int64, recoveryCodes []string, salt string) error {
	twoFactorRecoveryCodes := make([]*models.TwoFactorRecoveryCode, len(recoveryCodes))

	for i := 0; i < len(recoveryCodes); i++ {
		twoFactorRecoveryCodes[i] = &models.TwoFactorRecoveryCode{
			Uid:             uid,
			Used:            false,
			RecoveryCode:    utils.EncodePassword(recoveryCodes[i], salt),
			CreatedUnixTime: time.Now().Unix(),
		}
	}

	return s.UserDB().DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Where("uid=?", uid).Delete(&models.TwoFactorRecoveryCode{})

		if err != nil {
			return err
		}

		for i := 0; i < len(twoFactorRecoveryCodes); i++ {
			twoFactorRecoveryCode := twoFactorRecoveryCodes[i]
			_, err := sess.Insert(twoFactorRecoveryCode)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteTwoFactorRecoveryCodes deletes existed 2fa recovery codes from database
func (s *TwoFactorAuthorizationService) DeleteTwoFactorRecoveryCodes(c core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.UserDB().DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Where("uid=?", uid).Delete(&models.TwoFactorRecoveryCode{})
		return err
	})
}
