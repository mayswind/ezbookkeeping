package extservices

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// DefaultLocationName is the name of the location created automatically for single-store businesses
const DefaultLocationName = "Main"

// LocationService manages stores and warehouses
type LocationService struct{}

// Locations is the location service singleton
var Locations = &LocationService{}

// List returns all locations of a business, creating the default one on first use.
// Two requests racing on a brand-new business cannot create two defaults: the unique (owner, name) index rejects
// the loser, which then simply reads what the winner created.
func (s *LocationService) List(c core.Context, ownerUid int64) ([]*extmodels.Location, error) {
	if ownerUid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	locations, err := s.listActive(c, ownerUid)

	if err != nil || len(locations) > 0 {
		return locations, err
	}

	now := nowUnix()
	insertErr := ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		count, err := sess.Where("owner_uid=? AND deleted=?", ownerUid, false).Count(&extmodels.Location{})

		if err != nil || count > 0 {
			return err
		}

		_, err = sess.Insert(&extmodels.Location{OwnerUid: ownerUid, Name: DefaultLocationName, IsDefault: true, CreatedUnixTime: now, UpdatedUnixTime: now})

		return err
	})

	locations, err = s.listActive(c, ownerUid)

	if err != nil {
		return nil, err
	} else if len(locations) == 0 {
		return nil, insertErr
	}

	return locations, nil
}

func (s *LocationService) listActive(c core.Context, ownerUid int64) ([]*extmodels.Location, error) {
	locations := make([]*extmodels.Location, 0)
	err := ownerDB(ownerUid).NewSession(c).Where("owner_uid=? AND deleted=?", ownerUid, false).OrderBy("location_id").Find(&locations)

	return locations, err
}

// Resolve returns the location with the id, or the default location when the id is 0
func (s *LocationService) Resolve(c core.Context, ownerUid int64, locationId int64) (*extmodels.Location, error) {
	locations, err := s.List(c, ownerUid)

	if err != nil {
		return nil, err
	}

	for _, location := range locations {
		if (locationId == 0 && location.IsDefault) || location.LocationId == locationId {
			return location, nil
		}
	}

	if locationId == 0 {
		return locations[0], nil
	}

	return nil, exterrs.ErrLocationNotFound
}

// Add creates a location
func (s *LocationService) Add(c core.Context, ownerUid int64, name string) (*extmodels.Location, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, exterrs.ErrLocationNameIsEmpty
	}

	existing, err := s.List(c, ownerUid)

	if err != nil {
		return nil, err
	}

	for _, location := range existing {
		if strings.EqualFold(location.Name, name) {
			return nil, exterrs.ErrLocationNameExists
		}
	}

	now := nowUnix()
	location := &extmodels.Location{OwnerUid: ownerUid, Name: name, CreatedUnixTime: now, UpdatedUnixTime: now}
	_, err = ownerDB(ownerUid).NewSession(c).Insert(location)

	return location, err
}

// Modify renames a location
func (s *LocationService) Modify(c core.Context, ownerUid int64, locationId int64, name string) (*extmodels.Location, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, exterrs.ErrLocationNameIsEmpty
	}

	existing, err := s.List(c, ownerUid)

	if err != nil {
		return nil, err
	}

	var target *extmodels.Location

	for _, location := range existing {
		if location.LocationId == locationId {
			target = location
		} else if strings.EqualFold(location.Name, name) {
			return nil, exterrs.ErrLocationNameExists
		}
	}

	if target == nil {
		return nil, exterrs.ErrLocationNotFound
	}

	target.Name = name
	target.UpdatedUnixTime = nowUnix()
	_, err = ownerDB(ownerUid).NewSession(c).ID(locationId).Cols("name", "updated_unix_time").Update(target)

	return target, err
}

// Delete removes an empty location, never the last one. The location row is locked and the stock re-checked inside
// the same database transaction, so stock cannot arrive in a location while it is being deleted.
func (s *LocationService) Delete(c core.Context, ownerUid int64, locationId int64) error {
	if _, err := s.List(c, ownerUid); err != nil { // makes sure the default location exists
		return err
	}

	now := nowUnix()

	return ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		target := &extmodels.Location{}
		has, err := lockRows(sess).Where("owner_uid=? AND location_id=? AND deleted=?", ownerUid, locationId, false).Get(target)

		if err != nil {
			return err
		} else if !has {
			return exterrs.ErrLocationNotFound
		}

		others := make([]*extmodels.Location, 0)

		if err = sess.Where("owner_uid=? AND deleted=? AND location_id<>?", ownerUid, false, locationId).OrderBy("location_id").Find(&others); err != nil {
			return err
		}

		if len(others) == 0 {
			return exterrs.ErrCannotDeleteLastLocation
		}

		// stock of one item never goes below zero, so a zero total means every item is at zero
		stock, err := sess.Where("owner_uid=? AND location_id=?", ownerUid, locationId).SumInt(&extmodels.StockMovement{}, "qty_change")

		if err != nil {
			return err
		} else if stock != 0 {
			return exterrs.ErrLocationHasStock
		}

		wasDefault := target.IsDefault
		target.Deleted, target.IsDefault, target.DeletedUnixTime = true, false, now

		if _, err = sess.ID(locationId).Cols("deleted", "is_default", "deleted_unix_time").Update(target); err != nil {
			return err
		}

		if wasDefault {
			others[0].IsDefault = true
			_, err = sess.ID(others[0].LocationId).Cols("is_default").Update(others[0])
		}

		return err
	})
}
