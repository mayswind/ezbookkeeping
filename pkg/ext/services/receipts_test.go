package extservices

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

func TestBusinessProfile_DefaultsToEmpty(t *testing.T) {
	f := newFixture(t)

	profile, err := BusinessProfiles.Get(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Equal(t, f.ownerUid, profile.OwnerUid)
	assert.Empty(t, profile.ReceiptName)
	assert.Empty(t, profile.Footer)
}

func TestBusinessProfile_SaveAndReplace(t *testing.T) {
	f := newFixture(t)

	saved, err := BusinessProfiles.Save(f.c, f.ownerUid, ProfileInput{ReceiptName: "  Ada Stores ", Address: "12 Market Road", Phone: "0800 000", Footer: "Thank you"})
	require.NoError(t, err)
	assert.Equal(t, "Ada Stores", saved.ReceiptName, "surrounding spaces are trimmed")

	got, err := BusinessProfiles.Get(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Equal(t, "12 Market Road", got.Address)
	assert.Equal(t, "Thank you", got.Footer)

	_, err = BusinessProfiles.Save(f.c, f.ownerUid, ProfileInput{ReceiptName: "Ada Stores Ltd"})
	require.NoError(t, err)

	got, err = BusinessProfiles.Get(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Equal(t, "Ada Stores Ltd", got.ReceiptName)
	assert.Empty(t, got.Address, "a save replaces the whole profile, so cleared fields stay cleared")

	count, err := ownerDB(f.ownerUid).NewSession(f.c).Where("owner_uid=?", f.ownerUid).Count(&extmodels.BusinessProfile{})
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "one profile per business")
}

func TestBusinessProfile_TooLongIsRefusedNotCut(t *testing.T) {
	f := newFixture(t)

	_, err := BusinessProfiles.Save(f.c, f.ownerUid, ProfileInput{ReceiptName: strings.Repeat("x", 129)})
	assert.Equal(t, exterrs.ErrProfileFieldTooLong, err)
	_, err = BusinessProfiles.Save(f.c, f.ownerUid, ProfileInput{Phone: strings.Repeat("1", 33)})
	assert.Equal(t, exterrs.ErrProfileFieldTooLong, err)

	got, err := BusinessProfiles.Get(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Empty(t, got.ReceiptName, "a refused save changes nothing")
}

func TestBusinessProfile_BusinessesAreSeparate(t *testing.T) {
	f := newFixture(t)
	other := f.newUser(t, "other")

	_, err := BusinessProfiles.Save(f.c, f.ownerUid, ProfileInput{ReceiptName: "Mine"})
	require.NoError(t, err)

	got, err := BusinessProfiles.Get(f.c, other)
	require.NoError(t, err)
	assert.Empty(t, got.ReceiptName)
}

func TestRepayment_GetReturnsTheAllocations(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "SVC", 1000, false)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)
	older := f.creditSale(t, customer.CustomerId, item.ItemId, 3000, 0)
	newer := f.creditSale(t, customer.CustomerId, item.ItemId, 2000, 0)

	made, err := Repayments.Add(f.c, f.ownerUid, f.ownerUid, f.repay(customer.CustomerId, 3500, 0))
	require.NoError(t, err)

	got, err := Repayments.Get(f.c, f.ownerUid, made.Repayment.RepaymentId)
	require.NoError(t, err)
	assert.Equal(t, int64(3500), got.Repayment.Amount)
	require.Len(t, got.Allocations, 2)
	assert.Equal(t, older.Sale.SaleId, got.Allocations[0].SaleId)
	assert.Equal(t, int64(3000), got.Allocations[0].Amount)
	assert.Equal(t, newer.Sale.SaleId, got.Allocations[1].SaleId)
	assert.Equal(t, int64(500), got.Allocations[1].Amount)
}

func TestRepayment_GetIsLimitedToTheBusiness(t *testing.T) {
	f := newFixture(t)
	other := f.newUser(t, "other")
	item := f.newItem(t, "SVC", 1000, false)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)
	f.creditSale(t, customer.CustomerId, item.ItemId, 1000, 0)
	made, err := Repayments.Add(f.c, f.ownerUid, f.ownerUid, f.repay(customer.CustomerId, 500, 0))
	require.NoError(t, err)

	_, err = Repayments.Get(f.c, other, made.Repayment.RepaymentId)
	assert.Equal(t, exterrs.ErrRepaymentNotFound, err)
	_, err = Repayments.Get(f.c, f.ownerUid, 999999)
	assert.Equal(t, exterrs.ErrRepaymentNotFound, err)
}
