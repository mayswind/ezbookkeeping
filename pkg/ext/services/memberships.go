package extservices

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
)

// MembershipService manages who may work inside a business
type MembershipService struct{}

// Memberships is the membership service singleton
var Memberships = &MembershipService{}

func isAssignableRole(role extmodels.Role) bool {
	return role == extmodels.RoleManager || role == extmodels.RoleStaff
}

// Invite invites an existing user, found by email, to work in the owner's business
func (s *MembershipService) Invite(c core.Context, ownerUid int64, email string, role extmodels.Role) (*extmodels.Membership, error) {
	if ownerUid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if !isAssignableRole(role) {
		return nil, exterrs.ErrRoleInvalid
	}

	invitee, err := services.Users.GetUserByEmail(c, strings.TrimSpace(email))

	if err == errs.ErrUserNotFound {
		return nil, exterrs.ErrInviteeNotFound
	} else if err != nil {
		return nil, err
	}

	if invitee.Uid == ownerUid {
		return nil, exterrs.ErrCannotInviteSelf
	}

	now := nowUnix()
	var result *extmodels.Membership

	err = globalDB().DoTransaction(c, func(sess *xormSession) error {
		existing := &extmodels.Membership{}
		has, err := sess.Where("owner_uid=? AND staff_uid=?", ownerUid, invitee.Uid).Get(existing)

		if err != nil {
			return err
		}

		if has {
			if existing.Status == extmodels.MembershipStatusActive {
				return exterrs.ErrMembershipAlreadyActive
			}

			existing.Role = role
			existing.Status = extmodels.MembershipStatusPending
			existing.UpdatedUnixTime = now
			_, err = sess.ID(existing.MembershipId).Cols("role", "status", "updated_unix_time").Update(existing)
			result = existing

			return err
		}

		membership := &extmodels.Membership{
			OwnerUid:        ownerUid,
			StaffUid:        invitee.Uid,
			Role:            role,
			Status:          extmodels.MembershipStatusPending,
			CreatedUnixTime: now,
			UpdatedUnixTime: now,
		}
		_, err = sess.Insert(membership)
		result = membership

		return err
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Respond lets the invited user accept or decline an invitation
func (s *MembershipService) Respond(c core.Context, staffUid int64, ownerUid int64, accept bool) error {
	membership, err := s.get(c, ownerUid, staffUid)

	if err != nil {
		return err
	}

	if membership.Status != extmodels.MembershipStatusPending {
		return exterrs.ErrMembershipNotPending
	}

	if accept {
		membership.Status = extmodels.MembershipStatusActive
	} else {
		membership.Status = extmodels.MembershipStatusRevoked
	}

	membership.UpdatedUnixTime = nowUnix()
	_, err = globalDB().NewSession(c).ID(membership.MembershipId).Cols("status", "updated_unix_time").Update(membership)

	return err
}

// Remove ends a membership, used by the owner to remove somebody and by a member to leave
func (s *MembershipService) Remove(c core.Context, ownerUid int64, staffUid int64) error {
	membership, err := s.get(c, ownerUid, staffUid)

	if err != nil {
		return err
	}

	membership.Status = extmodels.MembershipStatusRevoked
	membership.UpdatedUnixTime = nowUnix()
	_, err = globalDB().NewSession(c).ID(membership.MembershipId).Cols("status", "updated_unix_time").Update(membership)

	return err
}

// SetRole changes the role of an active member
func (s *MembershipService) SetRole(c core.Context, ownerUid int64, staffUid int64, role extmodels.Role) error {
	if !isAssignableRole(role) {
		return exterrs.ErrRoleInvalid
	}

	membership, err := s.get(c, ownerUid, staffUid)

	if err != nil {
		return err
	}

	if membership.Status == extmodels.MembershipStatusRevoked {
		return exterrs.ErrMembershipNotFound
	}

	membership.Role = role
	membership.UpdatedUnixTime = nowUnix()
	_, err = globalDB().NewSession(c).ID(membership.MembershipId).Cols("role", "updated_unix_time").Update(membership)

	return err
}

// ListByOwner returns the pending and active members of a business
func (s *MembershipService) ListByOwner(c core.Context, ownerUid int64) ([]*extmodels.Membership, error) {
	memberships := make([]*extmodels.Membership, 0)
	err := globalDB().NewSession(c).Where("owner_uid=? AND status<>?", ownerUid, extmodels.MembershipStatusRevoked).OrderBy("membership_id").Find(&memberships)

	return memberships, err
}

// ListAllByOwner returns every membership of a business including removed members, so records written by
// somebody who has since left can still show their name
func (s *MembershipService) ListAllByOwner(c core.Context, ownerUid int64) ([]*extmodels.Membership, error) {
	memberships := make([]*extmodels.Membership, 0)
	err := globalDB().NewSession(c).Where("owner_uid=? AND status<>?", ownerUid, extmodels.MembershipStatusPending).OrderBy("membership_id").Find(&memberships)

	return memberships, err
}

// ListByStaff returns the businesses a user has been invited to or works in
func (s *MembershipService) ListByStaff(c core.Context, staffUid int64) ([]*extmodels.Membership, error) {
	memberships := make([]*extmodels.Membership, 0)
	err := globalDB().NewSession(c).Where("staff_uid=? AND status<>?", staffUid, extmodels.MembershipStatusRevoked).OrderBy("membership_id").Find(&memberships)

	return memberships, err
}

// GetActive returns the active membership of a user in a business
func (s *MembershipService) GetActive(c core.Context, ownerUid int64, staffUid int64) (*extmodels.Membership, error) {
	membership := &extmodels.Membership{}
	has, err := globalDB().NewSession(c).Where("owner_uid=? AND staff_uid=? AND status=?", ownerUid, staffUid, extmodels.MembershipStatusActive).Get(membership)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, exterrs.ErrNotAMemberOfBusiness
	}

	return membership, nil
}

func (s *MembershipService) get(c core.Context, ownerUid int64, staffUid int64) (*extmodels.Membership, error) {
	membership := &extmodels.Membership{}
	has, err := globalDB().NewSession(c).Where("owner_uid=? AND staff_uid=?", ownerUid, staffUid).Get(membership)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, exterrs.ErrMembershipNotFound
	}

	return membership, nil
}
