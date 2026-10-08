package extservices

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// Regression tests for the defects found in code review.

func TestReview_SkusAreNeverTruncated(t *testing.T) {
	f := newFixture(t)
	prefix := strings.Repeat("A", 40)

	first, err := Items.Add(f.c, f.ownerUid, ItemInput{Sku: prefix + "1", Name: "one"})
	require.NoError(t, err)
	second, err := Items.Add(f.c, f.ownerUid, ItemInput{Sku: prefix + "2", Name: "two"})
	require.NoError(t, err, "two skus that only differ after 32 characters are different skus")

	assert.Equal(t, prefix+"1", first.Sku)
	assert.Equal(t, prefix+"2", second.Sku)

	_, err = Items.Add(f.c, f.ownerUid, ItemInput{Sku: strings.Repeat("B", 65), Name: "too long"})
	assert.Equal(t, exterrs.ErrItemFieldTooLong, err, "too long is an error, not a silent cut")
}

func TestReview_StockQuantitiesAreBounded(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)

	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 1 << 62})
	assert.Equal(t, exterrs.ErrStockTooLarge, err)
	_, err = Stock.Adjust(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: -(1 << 62)})
	assert.Equal(t, exterrs.ErrStockTooLarge, err)
	_, err = Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 1000, UnitCost: 1 << 62})
	assert.Equal(t, exterrs.ErrStockTooLarge, err)

	// stock close to the cap (seeded directly: reaching it through the API would take a thousand requests)
	main, err := Locations.Resolve(f.c, f.ownerUid, 0)
	require.NoError(t, err)
	_, err = ownerDB(f.ownerUid).NewSession(f.c).Insert(&extmodels.StockMovement{
		OwnerUid: f.ownerUid, ItemId: item.ItemId, LocationId: main.LocationId, QtyChange: maxStockOnHand - 1000, Reason: extmodels.StockReasonOpening, RefType: "",
	})
	require.NoError(t, err)

	_, err = Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 2000})
	assert.Equal(t, exterrs.ErrStockTooLarge, err, "one more receipt than the cap allows is refused")

	_, err = Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 1000})
	require.NoError(t, err, "exactly reaching the cap is fine")
	assert.Equal(t, maxStockOnHand, f.stock(t, item.ItemId, 0))
}

func TestReview_TwoLocationsCannotShareAnActiveName(t *testing.T) {
	f := newFixture(t)

	_, err := Locations.List(f.c, f.ownerUid) // creates "Main"
	require.NoError(t, err)

	// what a racing request would try: a second active "Main" row; the unique index must reject it
	_, err = ownerDB(f.ownerUid).NewSession(f.c).Insert(&extmodels.Location{OwnerUid: f.ownerUid, Name: DefaultLocationName, IsDefault: true})
	assert.Error(t, err)

	locations, err := Locations.List(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Len(t, locations, 1)
}

func TestReview_ADeletedLocationsNameCanBeUsedAgain(t *testing.T) {
	f := newFixture(t)

	shop, err := Locations.Add(f.c, f.ownerUid, "Shop")
	require.NoError(t, err)
	require.NoError(t, Locations.Delete(f.c, f.ownerUid, shop.LocationId))

	again, err := Locations.Add(f.c, f.ownerUid, "Shop")
	require.NoError(t, err)
	assert.NotEqual(t, shop.LocationId, again.LocationId)
}

func TestReview_DeleteChecksStockInsideTheTransaction(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)
	shop, err := Locations.Add(f.c, f.ownerUid, "Shop")
	require.NoError(t, err)

	_, err = Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, LocationId: shop.LocationId, Qty: 1000})
	require.NoError(t, err)
	assert.Equal(t, exterrs.ErrLocationHasStock, Locations.Delete(f.c, f.ownerUid, shop.LocationId))

	// stock cannot be added to a location that was deleted
	_, err = Stock.Adjust(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, LocationId: shop.LocationId, Qty: -1000})
	require.NoError(t, err)
	require.NoError(t, Locations.Delete(f.c, f.ownerUid, shop.LocationId))

	_, err = Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, LocationId: shop.LocationId, Qty: 1000})
	assert.Equal(t, exterrs.ErrLocationNotFound, err)
}

func TestReview_RepaymentSeesEveryOpenSaleNotJustTheNewest200(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "SVC", 1000, false)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)

	const count = 205
	var oldest *SaleDetail

	for i := 0; i < count; i++ {
		detail := f.creditSale(t, customer.CustomerId, item.ItemId, 1000, 0) // owes 1000 each

		if i == 0 {
			oldest = detail
		}
	}

	owed, err := Customers.Outstanding(f.c, f.ownerUid, customer.CustomerId)
	require.NoError(t, err)
	require.Equal(t, int64(count*1000), owed)

	// the oldest sale is beyond the newest 200 and must still be payable on its own
	_, err = Repayments.Add(f.c, f.ownerUid, f.ownerUid, f.repay(customer.CustomerId, 1000, oldest.Sale.SaleId))
	require.NoError(t, err)

	// and the whole remaining balance must be repayable in one go
	result, err := Repayments.Add(f.c, f.ownerUid, f.ownerUid, f.repay(customer.CustomerId, (count-1)*1000, 0))
	require.NoError(t, err)
	assert.Len(t, result.Allocations, count-1)

	owed, err = Customers.Outstanding(f.c, f.ownerUid, customer.CustomerId)
	require.NoError(t, err)
	assert.Equal(t, int64(0), owed)
}

func TestReview_VoidCanBeRetriedAfterAPartialFailure(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)
	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 10 * 1000})
	require.NoError(t, err)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)
	sale := f.creditSale(t, customer.CustomerId, item.ItemId, 4*1000, 1000)

	// simulate a crash after the sale was flipped to voided but before its transactions were removed
	require.NoError(t, Sales.voidState(f.c, f.ownerUid, f.ownerUid, sale.Sale))
	require.Equal(t, int64(10*1000), f.stock(t, item.ItemId, 0))
	require.Equal(t, int64(3000), f.accountBalance(t, f.receivable), "the books still carry the sale")

	require.NoError(t, Sales.Void(f.c, f.ownerUid, f.ownerUid, sale.Sale.SaleId), "a retry finishes the job")

	assert.Equal(t, int64(0), f.accountBalance(t, f.cashId))
	assert.Equal(t, int64(0), f.accountBalance(t, f.receivable))
	assert.Equal(t, int64(10*1000), f.stock(t, item.ItemId, 0), "stock is restocked exactly once")

	got, err := Sales.Get(f.c, f.ownerUid, sale.Sale.SaleId)
	require.NoError(t, err)
	assert.Zero(t, got.Sale.PaidTransactionId)
	assert.Zero(t, got.Sale.CreditTransactionId)

	assert.Equal(t, exterrs.ErrSaleAlreadyVoided, Sales.Void(f.c, f.ownerUid, f.ownerUid, sale.Sale.SaleId))
}

func TestReview_VoidRestocksTheDefaultLocationWhenTheSaleLocationWasDeleted(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)
	shop, err := Locations.Add(f.c, f.ownerUid, "Shop 2")
	require.NoError(t, err)
	main, err := Locations.Resolve(f.c, f.ownerUid, 0)
	require.NoError(t, err)

	_, err = Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, LocationId: shop.LocationId, Qty: 2 * 1000})
	require.NoError(t, err)

	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 2 * 1000})
	in.LocationId = shop.LocationId
	in.AmountPaid = 2000
	sale, err := Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)

	require.NoError(t, Locations.Delete(f.c, f.ownerUid, shop.LocationId), "the shop is empty after the sale")
	require.NoError(t, Sales.Void(f.c, f.ownerUid, f.ownerUid, sale.Sale.SaleId))

	assert.Equal(t, int64(2*1000), f.stock(t, item.ItemId, main.LocationId), "stock must not land in a deleted location")
}

func TestReview_ConcurrentSalesOfTheLastUnitCannotBothSucceed(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "LAST", 1000, true)
	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 1000})
	require.NoError(t, err)

	// sequential here (SQLite serializes writers); on MySQL/PostgreSQL the row locks give the same outcome under concurrency
	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 1000})
	in.AmountPaid = 1000

	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)
	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	assert.Equal(t, exterrs.ErrInsufficientStock, err)
	assert.Equal(t, int64(0), f.stock(t, item.ItemId, 0))
}
