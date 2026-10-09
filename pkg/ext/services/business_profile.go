package extservices

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// BusinessProfileService stores what a business prints on its receipts
type BusinessProfileService struct{}

// BusinessProfiles is the business profile service singleton
var BusinessProfiles = &BusinessProfileService{}

// ProfileInput is the editable part of the profile
type ProfileInput struct {
	ReceiptName string
	Address     string
	Phone       string
	Footer      string
}

// Get returns a business's profile; a business that never set one gets an empty profile
func (s *BusinessProfileService) Get(c core.Context, ownerUid int64) (*extmodels.BusinessProfile, error) {
	if ownerUid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	profile := &extmodels.BusinessProfile{}
	has, err := ownerDB(ownerUid).NewSession(c).Where("owner_uid=?", ownerUid).Get(profile)

	if err != nil {
		return nil, err
	} else if !has {
		return &extmodels.BusinessProfile{OwnerUid: ownerUid}, nil
	}

	return profile, nil
}

// Save replaces a business's profile. Text that is too long is refused, never silently cut.
func (s *BusinessProfileService) Save(c core.Context, ownerUid int64, in ProfileInput) (*extmodels.BusinessProfile, error) {
	if ownerUid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	profile := &extmodels.BusinessProfile{
		OwnerUid: ownerUid, ReceiptName: trimmed(in.ReceiptName), Address: trimmed(in.Address),
		Phone: trimmed(in.Phone), Footer: trimmed(in.Footer), UpdatedUnixTime: nowUnix(),
	}

	if len(profile.ReceiptName) > 128 || len(profile.Address) > 255 || len(profile.Phone) > 32 || len(profile.Footer) > 255 {
		return nil, exterrs.ErrProfileFieldTooLong
	}

	err := ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		updated, err := sess.Where("owner_uid=?", ownerUid).Cols("receipt_name", "address", "phone", "footer", "updated_unix_time").Update(profile)

		if err != nil || updated > 0 {
			return err
		}

		_, err = sess.Insert(profile)

		return err
	})

	if err != nil {
		// two first saves can race on the insert; the loser finds the row now and only has to update it
		if existing, getErr := s.Get(c, ownerUid); getErr == nil && existing.UpdatedUnixTime > 0 {
			_, err = ownerDB(ownerUid).NewSession(c).Where("owner_uid=?", ownerUid).Cols("receipt_name", "address", "phone", "footer", "updated_unix_time").Update(profile)
		}

		if err != nil {
			return nil, err
		}
	}

	return profile, nil
}
