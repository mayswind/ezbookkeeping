package services

import (
	"time"
	"xorm.io/xorm"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/mayswind/lab/pkg/datastore"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/settings"
	"github.com/mayswind/lab/pkg/utils"
	"github.com/mayswind/lab/pkg/uuid"
)

const (
	TWOFACTOR_PERIOD               uint = 30 // seconds
	TWOFACTOR_SECRET_SIZE          uint = 20 // bytes
	TWOFACTOR_RECOVERY_CODE_COUNT  int  = 10
	TWOFACTOR_RECOVERY_CODE_LENGTH int  = 10 // bytes
)

type TwoFactorAuthorizationService struct {
	ServiceUsingDB
	ServiceUsingConfig
	ServiceUsingUuid
}

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

func (s *TwoFactorAuthorizationService) GetUserTwoFactorSettingByUid(uid int64) (*models.TwoFactor, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	twoFactor := &models.TwoFactor{}
	has, err := s.UserDB().Where("uid=?", uid).Get(twoFactor)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrTwoFactorKeyIsNotEnabled
	}

	twoFactor.Secret, err = utils.DecryptSecret(twoFactor.Secret, s.CurrentConfig().SecretKey)

	if err != nil {
		return nil, err
	}

	return twoFactor, nil
}

func (s *TwoFactorAuthorizationService) GenerateTwoFactorSecret(user *models.User) (*otp.Key, error) {
	if user == nil {
		return nil, errs.ErrUserNotFound
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.CurrentConfig().AppName,
		AccountName: user.Username,
		Period:      TWOFACTOR_PERIOD,
		SecretSize:  TWOFACTOR_SECRET_SIZE,
	})

	return key, err
}

func (s *TwoFactorAuthorizationService) CreateTwoFactorSetting(twoFactor *models.TwoFactor) error {
	if twoFactor.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	var err error
	twoFactor.Secret, err = utils.EncyptSecret(twoFactor.Secret, s.CurrentConfig().SecretKey)

	if err != nil {
		return err
	}

	twoFactor.CreatedUnixTime = time.Now().Unix()

	return s.UserDB().DoTranscation(func(sess *xorm.Session) error {
		_, err := sess.Insert(twoFactor)
		return err
	})
}

func (s *TwoFactorAuthorizationService) DeleteTwoFactorSetting(uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.UserDB().DoTranscation(func(sess *xorm.Session) error {
		deletedRows, err := sess.Where("uid=?", uid).Delete(&models.TwoFactor{})

		if deletedRows < 1 {
			return errs.ErrTwoFactorKeyIsNotEnabled
		}

		return err
	})
}

func (s *TwoFactorAuthorizationService) ExistsTwoFactorSetting(uid int64) (bool, error) {
	if uid <= 0 {
		return false, errs.ErrUserIdInvalid
	}

	return s.UserDB().Cols("uid").Where("uid=?", uid).Exist(&models.TwoFactor{})
}

func (s *TwoFactorAuthorizationService) GetAndUseUserTwoFactorRecoveryCode(uid int64, recoveryCode string, salt string) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	recoveryCode = utils.EncodePassword(recoveryCode, salt)
	exists, err := s.UserDB().Cols("uid", "recovery_code").Where("uid=? AND recovery_code=? AND used=?", uid, recoveryCode, false).Exist(&models.TwoFactorRecoveryCode{})

	if err != nil {
		return err
	} else if !exists {
		return errs.ErrTwoFactorRecoveryCodeNotExist
	}

	return s.UserDB().DoTranscation(func(sess *xorm.Session) error {
		_, err := sess.Cols("used", "used_unix_time").Where("uid=? AND recovery_code=?", uid, recoveryCode).Update(&models.TwoFactorRecoveryCode{Used: true, UsedUnixTime: time.Now().Unix()})
		return err
	})
}

func (s *TwoFactorAuthorizationService) GenerateTwoFactorRecoveryCodes() ([]string, error) {
	recoveryCodes := make([]string, TWOFACTOR_RECOVERY_CODE_COUNT)

	for i := 0; i < TWOFACTOR_RECOVERY_CODE_COUNT; i++ {
		recoveryCode, err := utils.GetRandomNumberOrLetter(TWOFACTOR_RECOVERY_CODE_LENGTH)

		if err != nil {
			return nil, err
		}

		recoveryCodes[i] = recoveryCode[:5] + "-" + recoveryCode[5:]
	}

	return recoveryCodes, nil
}

func (s *TwoFactorAuthorizationService) CreateTwoFactorRecoveryCodes(uid int64, recoveryCodes []string, salt string) error {
	twoFactorRecoveryCodes := make([]*models.TwoFactorRecoveryCode, len(recoveryCodes))

	for i := 0; i < len(recoveryCodes); i++ {
		twoFactorRecoveryCodes[i] = &models.TwoFactorRecoveryCode{
			Uid:             uid,
			Used:            false,
			RecoveryCode:    utils.EncodePassword(recoveryCodes[i], salt),
			CreatedUnixTime: time.Now().Unix(),
		}
	}

	return s.UserDB().DoTranscation(func(sess *xorm.Session) error {
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

func (s *TwoFactorAuthorizationService) DeleteTwoFactorRecoveryCodes(uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.UserDB().DoTranscation(func(sess *xorm.Session) error {
		_, err := sess.Where("uid=?", uid).Delete(&models.TwoFactorRecoveryCode{})
		return err
	})
}
