package extservices

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// UserSettingService stores the ext preferences of a person (not of a business)
type UserSettingService struct{}

// UserSettings is the user setting service singleton
var UserSettings = &UserSettingService{}

// Get returns a person's settings and whether they were ever saved. Somebody who never chose gets the defaults:
// everything off, so business features stay hidden until they opt in.
func (s *UserSettingService) Get(c core.Context, uid int64) (setting *extmodels.UserSetting, configured bool, err error) {
	if uid <= 0 {
		return nil, false, errs.ErrUserIdInvalid
	}

	setting = &extmodels.UserSetting{}
	has, err := globalDB().NewSession(c).Where("uid=?", uid).Get(setting)

	if err != nil {
		return nil, false, err
	} else if !has {
		return &extmodels.UserSetting{Uid: uid}, false, nil
	}

	return setting, true, nil
}

// SetBusinessFeatures turns the business features on or off for a person
func (s *UserSettingService) SetBusinessFeatures(c core.Context, uid int64, enabled bool) (*extmodels.UserSetting, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	now := nowUnix()

	err := globalDB().DoTransaction(c, func(sess *xormSession) error {
		updated, err := sess.Where("uid=?", uid).Cols("business_features", "updated_unix_time").Update(&extmodels.UserSetting{BusinessFeatures: enabled, UpdatedUnixTime: now})

		if err != nil || updated > 0 {
			return err
		}

		_, err = sess.Insert(&extmodels.UserSetting{Uid: uid, BusinessFeatures: enabled, CreatedUnixTime: now, UpdatedUnixTime: now})

		return err
	})

	if err != nil {
		// two first-time requests can race on the insert; the loser finds the row now and only has to update it
		if _, configured, getErr := s.Get(c, uid); getErr == nil && configured {
			_, err = globalDB().NewSession(c).Where("uid=?", uid).Cols("business_features", "updated_unix_time").Update(&extmodels.UserSetting{BusinessFeatures: enabled, UpdatedUnixTime: now})
		}

		if err != nil {
			return nil, err
		}
	}

	setting, _, err := s.Get(c, uid)

	return setting, err
}
