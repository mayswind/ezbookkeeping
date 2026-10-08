package extmw

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

func memberOf(ownerUid int64, role extmodels.Role) MembershipLookup {
	return func(requested int64) (*extmodels.Membership, error) {
		if requested != ownerUid {
			return nil, exterrs.ErrNotAMemberOfBusiness
		}

		return &extmodels.Membership{OwnerUid: ownerUid, Role: role, Status: extmodels.MembershipStatusActive}, nil
	}
}

func notMember(int64) (*extmodels.Membership, error) { return nil, exterrs.ErrNotAMemberOfBusiness }

func TestDecide_NoHeaderMeansOwnDataEvenForMembers(t *testing.T) {
	decision, err := Decide(5, "", "POST", "/transactions/add.json", memberOf(10, extmodels.RoleStaff))

	assert.Nil(t, err)
	assert.False(t, decision.Delegated)
}

func TestDecide_StaffCanSellInTheirBusiness(t *testing.T) {
	decision, err := Decide(5, "10", "POST", "/ext/sales/add.json", memberOf(10, extmodels.RoleStaff))

	assert.Nil(t, err)
	assert.True(t, decision.Delegated)
	assert.Equal(t, int64(10), decision.OwnerUid)
	assert.Equal(t, extmodels.RoleStaff, decision.Role)
}

func TestDecide_StrangersAreRejected(t *testing.T) {
	_, err := Decide(5, "10", "GET", "/accounts/list.json", notMember)
	assert.Equal(t, exterrs.ErrNotAMemberOfBusiness, err)

	_, err = Decide(5, "10", "GET", "/accounts/list.json", memberOf(99, extmodels.RoleManager))
	assert.Equal(t, exterrs.ErrNotAMemberOfBusiness, err, "membership in one business gives no access to another")
}

func TestDecide_RoleLimitsAreEnforced(t *testing.T) {
	_, err := Decide(5, "10", "POST", "/ext/sales/void.json", memberOf(10, extmodels.RoleStaff))
	assert.Equal(t, exterrs.ErrNotPermittedInBusiness, err)

	_, err = Decide(5, "10", "GET", "/data/export.csv", memberOf(10, extmodels.RoleManager))
	assert.Equal(t, exterrs.ErrNotPermittedInBusiness, err, "even managers cannot export the owner's data")

	_, err = Decide(5, "10", "POST", "/data/clear/all.json", memberOf(10, extmodels.RoleManager))
	assert.Equal(t, exterrs.ErrNotPermittedInBusiness, err)

	_, err = Decide(5, "10", "POST", "/brand/new/upstream/route.json", memberOf(10, extmodels.RoleManager))
	assert.Equal(t, exterrs.ErrNotPermittedInBusiness, err, "unknown routes are closed by default")
}

func TestDecide_IdentityRoutesNeverTouchTheOwner(t *testing.T) {
	for _, path := range []string{"/users/profile/update.json", "/tokens/generate/api.json", "/users/2fa/disable.json", "/ext/staff/invite.json"} {
		decision, err := Decide(5, "10", "POST", path, memberOf(10, extmodels.RoleManager))

		assert.Nil(t, err, path)
		assert.False(t, decision.Delegated, "%s must act on the caller's own account, not the owner's", path)
	}
}

func TestDecide_OwnIdInHeaderIsHarmless(t *testing.T) {
	decision, err := Decide(10, "10", "POST", "/data/clear/all.json", notMember)

	assert.Nil(t, err)
	assert.False(t, decision.Delegated)
}

func TestDecide_InvalidHeader(t *testing.T) {
	for _, value := range []string{"abc", "-3", "0", "1.5", "10; DROP TABLE x"} {
		_, err := Decide(5, value, "GET", "/accounts/list.json", memberOf(10, extmodels.RoleStaff))
		assert.Equal(t, exterrs.ErrBusinessIdInvalid, err, value)
	}
}

func TestDecide_LookupFailureDoesNotGrantAccess(t *testing.T) {
	_, err := Decide(5, "10", "GET", "/accounts/list.json", func(int64) (*extmodels.Membership, error) {
		return nil, errors.New("database exploded")
	})

	assert.Equal(t, errs.ErrOperationFailed, err)
}
