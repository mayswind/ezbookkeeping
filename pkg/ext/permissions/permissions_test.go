package extperm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

func TestAllowed_StaffCanSellButNotManage(t *testing.T) {
	assert.True(t, Allowed(extmodels.RoleStaff, "POST", "/ext/sales/add.json"))
	assert.True(t, Allowed(extmodels.RoleStaff, "GET", "/accounts/list.json"))
	assert.True(t, Allowed(extmodels.RoleStaff, "POST", "/ext/repayments/add.json"))

	assert.False(t, Allowed(extmodels.RoleStaff, "POST", "/ext/sales/void.json"))
	assert.False(t, Allowed(extmodels.RoleStaff, "POST", "/ext/items/add.json"))
	assert.False(t, Allowed(extmodels.RoleStaff, "POST", "/ext/stock/adjust.json"))
	assert.False(t, Allowed(extmodels.RoleStaff, "POST", "/transactions/delete.json"))
	assert.False(t, Allowed(extmodels.RoleStaff, "POST", "/transactions/add.json"), "staff must not post raw expenses or transfers")
	assert.False(t, Allowed(extmodels.RoleStaff, "GET", "/transactions/statistics.json"))
	assert.False(t, Allowed(extmodels.RoleStaff, "GET", "/ext/stock/movements/list.json"))
}

func TestAllowed_ManagerCanOperateButNotExportOrDelete(t *testing.T) {
	assert.True(t, Allowed(extmodels.RoleManager, "POST", "/transactions/add.json"))
	assert.True(t, Allowed(extmodels.RoleManager, "POST", "/ext/items/add.json"))
	assert.True(t, Allowed(extmodels.RoleManager, "POST", "/ext/sales/void.json"))
	assert.True(t, Allowed(extmodels.RoleManager, "POST", "/transactions/modify.json"))
	assert.True(t, Allowed(extmodels.RoleManager, "GET", "/transactions/statistics.json"))
	assert.True(t, Allowed(extmodels.RoleManager, "GET", "/ext/stock/movements/list.json"))
	assert.True(t, Allowed(extmodels.RoleManager, "POST", "/accounts/add.json"))

	assert.False(t, Allowed(extmodels.RoleManager, "POST", "/accounts/delete.json"))
	assert.False(t, Allowed(extmodels.RoleManager, "POST", "/data/clear/all.json"))
	assert.False(t, Allowed(extmodels.RoleManager, "GET", "/data/export.csv"))
	assert.False(t, Allowed(extmodels.RoleManager, "POST", "/transactions/import.json"))
	assert.False(t, Allowed(extmodels.RoleManager, "POST", "/transaction/categories/delete.json"))
	assert.False(t, Allowed(extmodels.RoleManager, "POST", "/tokens/generate/api.json"))
}

func TestAllowed_UnknownRoutesAreDeniedByDefault(t *testing.T) {
	assert.False(t, Allowed(extmodels.RoleManager, "POST", "/some/new/upstream/endpoint.json"))
	assert.False(t, Allowed(extmodels.RoleStaff, "GET", "/some/new/upstream/endpoint.json"))
	assert.False(t, Allowed(extmodels.Role(0), "GET", "/accounts/list.json"))
	assert.True(t, Allowed(extmodels.RoleOwner, "POST", "/some/new/upstream/endpoint.json"))
}

func TestAllowed_MethodMustMatch(t *testing.T) {
	assert.False(t, Allowed(extmodels.RoleStaff, "GET", "/ext/sales/add.json"))
	assert.False(t, Allowed(extmodels.RoleManager, "POST", "/accounts/list.json"))
}

func TestIsSelfOnlyPath(t *testing.T) {
	assert.True(t, IsSelfOnlyPath("/users/profile/update.json"))
	assert.True(t, IsSelfOnlyPath("/tokens/list.json"))
	assert.True(t, IsSelfOnlyPath("/ext/staff/invite.json"))
	assert.True(t, IsSelfOnlyPath("/ext/me/businesses.json"))
	assert.False(t, IsSelfOnlyPath("/accounts/list.json"))
	assert.False(t, IsSelfOnlyPath("/ext/sales/add.json"))
}

func TestRelativePath(t *testing.T) {
	assert.Equal(t, "/accounts/list.json", RelativePath("/api/v1/accounts/list.json"))
}
