package services

import (
	"time"
	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// UserCustomExchangeRatesService represents user custom exchange rate data service
type UserCustomExchangeRatesService struct {
	ServiceUsingDB
}

// Initialize a user custom exchange rate data service singleton instance
var (
	UserCustomExchangeRates = &UserCustomExchangeRatesService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
	}
)

// GetAllCustomExchangeRatesByUid returns all user exchange rate data models of user
func (s *UserCustomExchangeRatesService) GetAllCustomExchangeRatesByUid(c core.Context, uid int64) ([]*models.UserCustomExchangeRate, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var customExchangeRates []*models.UserCustomExchangeRate
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted_unix_time=?", uid, 0).Find(&customExchangeRates)

	return customExchangeRates, err
}

// UpdateCustomExchangeRate updates user exchange rate data model to database
func (s *UserCustomExchangeRatesService) UpdateCustomExchangeRate(c core.Context, uid int64, currency string, rate string, defaultCurrency string) (*models.UserCustomExchangeRate, *models.UserCustomExchangeRate, error) {
	if uid <= 0 {
		return nil, nil, errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()
	newCustomExchangeRate := &models.UserCustomExchangeRate{}
	defaultCurrencyExchangeRate := &models.UserCustomExchangeRate{}

	err := s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		oldCustomExchangeRate := &models.UserCustomExchangeRate{}
		has, err := sess.Where("uid=? AND deleted_unix_time=? AND currency=?", uid, 0, currency).Get(oldCustomExchangeRate)

		if err != nil {
			return err
		}

		if has {
			updateOldExchangeRateModel := &models.UserCustomExchangeRate{
				DeletedUnixTime: now,
			}

			_, err = sess.Cols("deleted_unix_time").Where("uid=? AND deleted_unix_time=? AND currency=?", uid, 0, currency).Update(updateOldExchangeRateModel)

			if err != nil {
				return err
			}
		}

		if currency != defaultCurrency {
			has, err := sess.Where("uid=? AND deleted_unix_time=? AND currency=?", uid, 0, defaultCurrency).Get(defaultCurrencyExchangeRate)

			if err != nil {
				return err
			}

			if !has {
				defaultCurrencyExchangeRate, _ = models.CreateUserCustomExchangeRate(uid, defaultCurrency, "1", 0)
				defaultCurrencyExchangeRate.CreatedUnixTime = now
				defaultCurrencyExchangeRate.UpdatedUnixTime = now
				defaultCurrencyExchangeRate.DeletedUnixTime = 0
				_, err = sess.Insert(defaultCurrencyExchangeRate)

				if err != nil {
					return err
				}
			}
		} else {
			defaultCurrencyExchangeRate = oldCustomExchangeRate
		}

		newCustomExchangeRate, err = models.CreateUserCustomExchangeRate(uid, currency, rate, defaultCurrencyExchangeRate.Rate)
		newCustomExchangeRate.CreatedUnixTime = now
		newCustomExchangeRate.UpdatedUnixTime = now
		newCustomExchangeRate.DeletedUnixTime = 0
		_, err = sess.Insert(newCustomExchangeRate)

		return err
	})

	if err != nil {
		return nil, nil, err
	}

	return newCustomExchangeRate, defaultCurrencyExchangeRate, err
}

// DeleteCustomExchangeRate deletes an existed user exchange rate data from database
func (s *UserCustomExchangeRatesService) DeleteCustomExchangeRate(c core.Context, uid int64, currency string) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.UserCustomExchangeRate{
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		deletedRows, err := sess.Cols("deleted_unix_time").Where("uid=? AND deleted_unix_time=? AND currency=?", uid, 0, currency).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrUserCustomExchangeRateNotFound
		}

		return err
	})
}

// DeleteAllCustomExchangeRates deletes all existed user exchange rate data from database
func (s *UserCustomExchangeRatesService) DeleteAllCustomExchangeRates(c core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.UserCustomExchangeRate{
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Cols("deleted_unix_time").Where("uid=? AND deleted_unix_time=?", uid, 0).Update(updateModel)

		if err != nil {
			return err
		}

		return nil
	})
}
