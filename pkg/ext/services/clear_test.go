package extservices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

func (f *fixture) count(t *testing.T, ownerUid int64, bean any) int64 {
	t.Helper()

	n, err := ownerDB(ownerUid).NewSession(f.c).Where("owner_uid=?", ownerUid).Count(bean)
	require.NoError(t, err)

	return n
}

// seedBusiness gives an owner one of everything
func (f *fixture) seedBusiness(t *testing.T, ownerUid int64) {
	t.Helper()

	item, err := Items.Add(f.c, ownerUid, ItemInput{Sku: "RICE", Name: "Rice", SalePrice: 1000, TrackStock: true})
	require.NoError(t, err)
	gone, err := Items.Add(f.c, ownerUid, ItemInput{Sku: "GONE", Name: "Gone", SalePrice: 1})
	require.NoError(t, err)
	require.NoError(t, Items.Delete(f.c, ownerUid, gone.ItemId))
	_, err = Locations.Add(f.c, ownerUid, "Shop 2")
	require.NoError(t, err)
	_, err = Stock.Receive(f.c, ownerUid, ownerUid, StockInput{ItemId: item.ItemId, Qty: 9000})
	require.NoError(t, err)
	customer, err := Customers.Add(f.c, ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)

	in := SaleInput{Lines: []SaleLineInput{{ItemId: item.ItemId, Qty: 2000}}, CustomerId: customer.CustomerId, AmountPaid: 500, PaymentAccountId: f.cashId, ReceivableAccountId: f.receivable, CategoryId: f.incomeCat}
	_, err = Sales.Create(f.c, ownerUid, ownerUid, in)
	require.NoError(t, err)
	_, err = Repayments.Add(f.c, ownerUid, ownerUid, RepaymentInput{CustomerId: customer.CustomerId, Amount: 300, PaymentAccountId: f.bankId, ReceivableAccountId: f.receivable, CategoryId: f.transferCat})
	require.NoError(t, err)
}

func TestClear_RemovesEveryBusinessRecord(t *testing.T) {
	f := newFixture(t)
	f.seedBusiness(t, f.ownerUid)

	require.Greater(t, f.count(t, f.ownerUid, &extmodels.Sale{}), int64(0))

	counts, err := Clear.ClearBusinessRecords(f.c, f.ownerUid, f.ownerUid, "203.0.113.5")
	require.NoError(t, err)

	assert.Equal(t, int64(2), counts.Items, "including the deleted one")
	assert.Equal(t, int64(2), counts.Locations)
	assert.Equal(t, int64(1), counts.Customers)
	assert.Equal(t, int64(1), counts.Sales)
	assert.Equal(t, int64(1), counts.SaleLines)
	assert.Equal(t, int64(1), counts.Repayments)
	assert.Equal(t, int64(1), counts.RepaymentAllocations)
	assert.Equal(t, int64(2), counts.StockMovements, "the stock received and the stock sold")

	for name, bean := range map[string]any{
		"items": &extmodels.Item{}, "locations": &extmodels.Location{}, "stock movements": &extmodels.StockMovement{}, "customers": &extmodels.Customer{},
		"sales": &extmodels.Sale{}, "sale lines": &extmodels.SaleLine{}, "repayments": &extmodels.Repayment{}, "allocations": &extmodels.RepaymentAllocation{},
	} {
		assert.Equal(t, int64(0), f.count(t, f.ownerUid, bean), name)
	}
}

func TestClear_LeavesOtherBusinessesAlone(t *testing.T) {
	f := newFixture(t)
	other := f.newUser(t, "other")
	f.seedBusiness(t, f.ownerUid)

	// the other business needs its own accounts and categories, so seed it directly
	_, err := Items.Add(f.c, other, ItemInput{Sku: "THEIRS", Name: "Theirs", SalePrice: 5, TrackStock: true})
	require.NoError(t, err)
	_, err = Customers.Add(f.c, other, CustomerInput{Name: "Their customer"})
	require.NoError(t, err)
	_, err = Locations.List(f.c, other)
	require.NoError(t, err)

	_, err = Clear.ClearBusinessRecords(f.c, f.ownerUid, f.ownerUid, "")
	require.NoError(t, err)

	assert.Equal(t, int64(1), f.count(t, other, &extmodels.Item{}))
	assert.Equal(t, int64(1), f.count(t, other, &extmodels.Customer{}))
	assert.Equal(t, int64(1), f.count(t, other, &extmodels.Location{}))
}

func TestClear_KeepsSettingsTeamAndAccountability(t *testing.T) {
	f := newFixture(t)
	f.newUser(t, "sam")
	f.seedBusiness(t, f.ownerUid)

	m, err := Memberships.Invite(f.c, f.ownerUid, "sam@example.com", extmodels.RoleStaff)
	require.NoError(t, err)
	require.NoError(t, Memberships.Respond(f.c, m.StaffUid, f.ownerUid, true))
	_, err = BusinessProfiles.Save(f.c, f.ownerUid, ProfileInput{ReceiptName: "Ada Stores"})
	require.NoError(t, err)
	_, err = UserSettings.SetBusinessFeatures(f.c, f.ownerUid, true)
	require.NoError(t, err)
	require.NoError(t, Terms.Accept(f.c, f.ownerUid, "v1", ""))
	require.NoError(t, Audit.Record(f.c, &extmodels.AuditLog{OwnerUid: f.ownerUid, ActorUid: m.StaffUid, Role: extmodels.RoleStaff, Method: "POST", Path: "/x", Status: 200, Action: "sales.add"}))

	_, err = Clear.ClearBusinessRecords(f.c, f.ownerUid, f.ownerUid, "198.51.100.1")
	require.NoError(t, err)

	members, err := Memberships.ListByOwner(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Len(t, members, 1, "the team stays")

	profile, err := BusinessProfiles.Get(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Equal(t, "Ada Stores", profile.ReceiptName, "receipt details stay")

	setting, configured, err := UserSettings.Get(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.True(t, configured)
	assert.True(t, setting.BusinessFeatures, "the person's settings stay")

	version, _, err := Terms.Latest(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Equal(t, "v1", version, "the acceptance stays")

	entries, err := Audit.List(f.c, f.ownerUid, 10, 0)
	require.NoError(t, err)
	require.Len(t, entries, 2, "the earlier entry stays, and the clearing itself is added")
	assert.Equal(t, "data.clear_all", entries[0].Action)
	assert.Equal(t, f.ownerUid, entries[0].ActorUid)
	assert.Equal(t, "198.51.100.1", entries[0].ClientIp)
	assert.Equal(t, "sales.add", entries[1].Action)
}

func TestClear_IsRepeatableAndTheBusinessCanStartAgain(t *testing.T) {
	f := newFixture(t)
	f.seedBusiness(t, f.ownerUid)

	_, err := Clear.ClearBusinessRecords(f.c, f.ownerUid, f.ownerUid, "")
	require.NoError(t, err)

	again, err := Clear.ClearBusinessRecords(f.c, f.ownerUid, f.ownerUid, "")
	require.NoError(t, err, "clearing an empty business is fine")
	assert.Equal(t, ClearedCounts{}, *again, "and removes nothing")

	// a fresh start: the default location comes back and an old sku can be used again
	locations, err := Locations.List(f.c, f.ownerUid)
	require.NoError(t, err)
	require.Len(t, locations, 1)
	assert.Equal(t, DefaultLocationName, locations[0].Name)

	_, err = Items.Add(f.c, f.ownerUid, ItemInput{Sku: "RICE", Name: "Rice again", TrackStock: true})
	require.NoError(t, err)

	levels, err := Stock.Levels(f.c, f.ownerUid, 0, 0)
	require.NoError(t, err)
	assert.Empty(t, levels, "no stock is left behind")

	owed, err := Customers.Balances(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Empty(t, owed, "nobody owes anything")
}

func TestClear_RejectsAnInvalidUser(t *testing.T) {
	f := newFixture(t)

	_, err := Clear.ClearBusinessRecords(f.c, 0, 0, "")
	assert.Error(t, err)
}
