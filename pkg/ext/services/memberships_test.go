package extservices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

func TestMembership_InviteAcceptAndResolve(t *testing.T) {
	f := newFixture(t)
	staff := f.newUser(t, "staff")

	membership, err := Memberships.Invite(f.c, f.ownerUid, "staff@example.com", extmodels.RoleStaff)
	require.NoError(t, err)
	assert.Equal(t, extmodels.MembershipStatusPending, membership.Status)
	assert.Equal(t, staff, membership.StaffUid)

	_, err = Memberships.GetActive(f.c, f.ownerUid, staff)
	assert.Equal(t, exterrs.ErrNotAMemberOfBusiness, err, "a pending invitation gives no access")

	require.NoError(t, Memberships.Respond(f.c, staff, f.ownerUid, true))

	active, err := Memberships.GetActive(f.c, f.ownerUid, staff)
	require.NoError(t, err)
	assert.Equal(t, extmodels.RoleStaff, active.Role)

	assert.Equal(t, exterrs.ErrMembershipNotPending, Memberships.Respond(f.c, staff, f.ownerUid, true))
}

func TestMembership_DeclinedInviteGivesNoAccess(t *testing.T) {
	f := newFixture(t)
	staff := f.newUser(t, "staff")

	_, err := Memberships.Invite(f.c, f.ownerUid, "staff@example.com", extmodels.RoleManager)
	require.NoError(t, err)
	require.NoError(t, Memberships.Respond(f.c, staff, f.ownerUid, false))

	_, err = Memberships.GetActive(f.c, f.ownerUid, staff)
	assert.Equal(t, exterrs.ErrNotAMemberOfBusiness, err)

	list, err := Memberships.ListByStaff(f.c, staff)
	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestMembership_InviteRules(t *testing.T) {
	f := newFixture(t)
	f.newUser(t, "staff")

	_, err := Memberships.Invite(f.c, f.ownerUid, "owner@example.com", extmodels.RoleStaff)
	assert.Equal(t, exterrs.ErrCannotInviteSelf, err)

	_, err = Memberships.Invite(f.c, f.ownerUid, "nobody@example.com", extmodels.RoleStaff)
	assert.Equal(t, exterrs.ErrInviteeNotFound, err)

	_, err = Memberships.Invite(f.c, f.ownerUid, "staff@example.com", extmodels.RoleOwner)
	assert.Equal(t, exterrs.ErrRoleInvalid, err, "nobody can be invited as owner")
	_, err = Memberships.Invite(f.c, f.ownerUid, "staff@example.com", extmodels.Role(0))
	assert.Equal(t, exterrs.ErrRoleInvalid, err)
}

func TestMembership_ReinviteAfterRemovalAndRoleChange(t *testing.T) {
	f := newFixture(t)
	staff := f.newUser(t, "staff")

	_, err := Memberships.Invite(f.c, f.ownerUid, "staff@example.com", extmodels.RoleStaff)
	require.NoError(t, err)
	require.NoError(t, Memberships.Respond(f.c, staff, f.ownerUid, true))

	_, err = Memberships.Invite(f.c, f.ownerUid, "staff@example.com", extmodels.RoleManager)
	assert.Equal(t, exterrs.ErrMembershipAlreadyActive, err)

	require.NoError(t, Memberships.SetRole(f.c, f.ownerUid, staff, extmodels.RoleManager))
	active, err := Memberships.GetActive(f.c, f.ownerUid, staff)
	require.NoError(t, err)
	assert.Equal(t, extmodels.RoleManager, active.Role)
	assert.Equal(t, exterrs.ErrRoleInvalid, Memberships.SetRole(f.c, f.ownerUid, staff, extmodels.RoleOwner))

	require.NoError(t, Memberships.Remove(f.c, f.ownerUid, staff))
	_, err = Memberships.GetActive(f.c, f.ownerUid, staff)
	assert.Equal(t, exterrs.ErrNotAMemberOfBusiness, err, "removal takes access away immediately")

	again, err := Memberships.Invite(f.c, f.ownerUid, "staff@example.com", extmodels.RoleStaff)
	require.NoError(t, err)
	assert.Equal(t, extmodels.MembershipStatusPending, again.Status, "a removed member must accept again")
}

func TestMembership_ListsAreScopedAndHideRevoked(t *testing.T) {
	f := newFixture(t)
	a := f.newUser(t, "a")
	f.newUser(t, "b")
	otherOwner := f.newUser(t, "boss2")

	for _, email := range []string{"a@example.com", "b@example.com"} {
		_, err := Memberships.Invite(f.c, f.ownerUid, email, extmodels.RoleStaff)
		require.NoError(t, err)
	}

	_, err := Memberships.Invite(f.c, otherOwner, "a@example.com", extmodels.RoleManager)
	require.NoError(t, err)

	mine, err := Memberships.ListByOwner(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Len(t, mine, 2)

	forA, err := Memberships.ListByStaff(f.c, a)
	require.NoError(t, err)
	assert.Len(t, forA, 2, "a works for two businesses")

	require.NoError(t, Memberships.Remove(f.c, f.ownerUid, a))
	mine, err = Memberships.ListByOwner(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Len(t, mine, 1)
}

func TestMembership_UnknownMembershipErrors(t *testing.T) {
	f := newFixture(t)
	stranger := f.newUser(t, "stranger")

	assert.Equal(t, exterrs.ErrMembershipNotFound, Memberships.Respond(f.c, stranger, f.ownerUid, true))
	assert.Equal(t, exterrs.ErrMembershipNotFound, Memberships.Remove(f.c, f.ownerUid, stranger))
	assert.Equal(t, exterrs.ErrMembershipNotFound, Memberships.SetRole(f.c, f.ownerUid, stranger, extmodels.RoleStaff))
}

func TestAudit_RecordAndListNewestFirst(t *testing.T) {
	f := newFixture(t)

	for i := 0; i < 3; i++ {
		require.NoError(t, Audit.Record(f.c, &extmodels.AuditLog{OwnerUid: f.ownerUid, ActorUid: 9, Role: extmodels.RoleStaff, Method: "POST", Path: "/x", Status: 200}))
	}

	require.NoError(t, Audit.Record(f.c, &extmodels.AuditLog{OwnerUid: 12345, ActorUid: 9, Role: extmodels.RoleStaff, Method: "POST", Path: "/other", Status: 200}))

	entries, err := Audit.List(f.c, f.ownerUid, 2, 0)
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Greater(t, entries[0].AuditId, entries[1].AuditId)

	rest, err := Audit.List(f.c, f.ownerUid, 10, entries[1].AuditId)
	require.NoError(t, err)
	assert.Len(t, rest, 1)
}
