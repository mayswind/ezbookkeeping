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

// List returns all locations of a business, creating the default one on first use
func (s *LocationService) List(c core.Context, ownerUid int64) ([]*extmodels.Location, error) {
	if ownerUid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	locations := make([]*extmodels.Location, 0)
	err := ownerDB(ownerUid).NewSession(c).Where("owner_uid=? AND deleted=?", ownerUid, false).OrderBy("location_id").Find(&locations)

	if err != nil {
		return nil, err
	}

	if len(locations) > 0 {
		return locations, nil
	}

	now := nowUnix()
	location := &extmodels.Location{OwnerUid: ownerUid, Name: DefaultLocationName, IsDefault: true, CreatedUnixTime: now, UpdatedUnixTime: now}
	_, err = ownerDB(ownerUid).NewSession(c).Insert(location)

	if err != nil {
		return nil, err
	}

	return []*extmodels.Location{location}, nil
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

// Delete removes an empty location, never the last one
func (s *LocationService) Delete(c core.Context, ownerUid int64, locationId int64) error {
	existing, err := s.List(c, ownerUid)

	if err != nil {
		return err
	}

	var target, replacement *extmodels.Location

	for _, location := range existing {
		if location.LocationId == locationId {
			target = location
		} else if replacement == nil {
			replacement = location
		}
	}

	if target == nil {
		return exterrs.ErrLocationNotFound
	}

	if replacement == nil {
		return exterrs.ErrCannotDeleteLastLocation
	}

	levels, err := Stock.Levels(c, ownerUid, 0, locationId)

	if err != nil {
		return err
	}

	for _, level := range levels {
		if level.Qty != 0 {
			return exterrs.ErrLocationHasStock
		}
	}

	now := nowUnix()
	wasDefault := target.IsDefault

	return ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		target.Deleted = true
		target.IsDefault = false
		target.DeletedUnixTime = now
		_, err := sess.ID(locationId).Cols("deleted", "is_default", "deleted_unix_time").Update(target)

		if err != nil {
			return err
		}

		if wasDefault {
			replacement.IsDefault = true
			_, err = sess.ID(replacement.LocationId).Cols("is_default").Update(replacement)
		}

		return err
	})
}
