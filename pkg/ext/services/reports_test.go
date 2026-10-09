package extservices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

func (f *fixture) receive(t *testing.T, itemId int64, locationId int64, qty int64) {
	t.Helper()

	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: itemId, LocationId: locationId, Qty: qty})
	require.NoError(t, err)
}

func TestReport_StockValueUsesCostAndSalePrice(t *testing.T) {
	f := newFixture(t)
	rice := f.newItem(t, "RICE", 1500, true)    // cost 750, sells 1500 per unit
	oil := f.newItem(t, "OIL", 4000, true)      // cost 2000, sells 4000
	f.newItem(t, "EMPTY", 900, true)            // no stock: left out
	service := f.newItem(t, "SVC", 5000, false) // services never have stock
	_ = service

	f.receive(t, rice.ItemId, 0, 2500) // 2.5 units
	f.receive(t, oil.ItemId, 0, 4000)  // 4 units

	report, err := Reports.StockValue(f.c, f.ownerUid, 0)
	require.NoError(t, err)
	require.Len(t, report.Rows, 2)

	// oil: 4 x 2000 = 8000 cost, 4 x 4000 = 16000 retail; rice: 2.5 x 750 = 1875, 2.5 x 1500 = 3750
	assert.Equal(t, "OIL", report.Rows[0].Item.Sku, "the most valuable stock comes first")
	assert.Equal(t, int64(8000), report.Rows[0].CostValue)
	assert.Equal(t, int64(16000), report.Rows[0].RetailValue)
	assert.Equal(t, int64(1875), report.Rows[1].CostValue)
	assert.Equal(t, int64(3750), report.Rows[1].RetailValue)
	assert.Equal(t, int64(9875), report.TotalCostValue)
	assert.Equal(t, int64(19750), report.TotalRetailValue)
}

func TestReport_StockValueCanBeLimitedToOneLocation(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true) // cost 500
	main, err := Locations.Resolve(f.c, f.ownerUid, 0)
	require.NoError(t, err)
	shop, err := Locations.Add(f.c, f.ownerUid, "Shop 2")
	require.NoError(t, err)

	f.receive(t, item.ItemId, main.LocationId, 3000)
	f.receive(t, item.ItemId, shop.LocationId, 1000)

	all, err := Reports.StockValue(f.c, f.ownerUid, 0)
	require.NoError(t, err)
	require.Len(t, all.Rows, 1)
	assert.Equal(t, int64(4000), all.Rows[0].Qty)
	assert.Equal(t, int64(2000), all.TotalCostValue)
	assert.Len(t, all.Rows[0].Locations, 2, "the breakdown by location is given when nothing is filtered")

	oneShop, err := Reports.StockValue(f.c, f.ownerUid, shop.LocationId)
	require.NoError(t, err)
	require.Len(t, oneShop.Rows, 1)
	assert.Equal(t, int64(1000), oneShop.Rows[0].Qty)
	assert.Equal(t, int64(500), oneShop.TotalCostValue)
	assert.Empty(t, oneShop.Rows[0].Locations)
}

func TestReport_StockValueIsEmptyWithoutStock(t *testing.T) {
	f := newFixture(t)

	report, err := Reports.StockValue(f.c, f.ownerUid, 0)
	require.NoError(t, err)
	assert.Empty(t, report.Rows)
	assert.Equal(t, int64(0), report.TotalCostValue)
}

func TestReport_BusinessesDoNotSeeEachOthersStock(t *testing.T) {
	f := newFixture(t)
	other := f.newUser(t, "other")
	item := f.newItem(t, "A1", 1000, true)
	f.receive(t, item.ItemId, 0, 1000)

	report, err := Reports.StockValue(f.c, other, 0)
	require.NoError(t, err)
	assert.Empty(t, report.Rows)
}

func TestReport_LowStock(t *testing.T) {
	f := newFixture(t)

	mk := func(sku string, reorder int64, stock int64) {
		item, err := Items.Add(f.c, f.ownerUid, ItemInput{Sku: sku, Name: sku, SalePrice: 100, ReorderLevel: reorder, TrackStock: true})
		require.NoError(t, err)

		if stock > 0 {
			f.receive(t, item.ItemId, 0, stock)
		}
	}

	mk("LOW", 5000, 3000) // 3 left, wants 5: short by 2
	mk("AT", 5000, 5000)  // exactly at the level counts as low
	mk("OK", 5000, 9000)  // fine
	mk("OUT", 4000, 0)    // nothing left, short by 4
	mk("NOLEVEL", 0, 0)   // no reorder level set: never listed

	rows, err := Reports.LowStock(f.c, f.ownerUid, 0)
	require.NoError(t, err)

	var skus []string
	for _, row := range rows {
		skus = append(skus, row.Item.Sku)
	}

	assert.Equal(t, []string{"OUT", "LOW", "AT"}, skus, "most urgent first, and only items that reached their level")
	assert.Equal(t, int64(4000), rows[0].Shortfall)
	assert.Equal(t, int64(2000), rows[1].Shortfall)
	assert.Equal(t, int64(0), rows[2].Shortfall)
}

func TestReport_LowStockPerLocation(t *testing.T) {
	f := newFixture(t)
	main, err := Locations.Resolve(f.c, f.ownerUid, 0)
	require.NoError(t, err)
	shop, err := Locations.Add(f.c, f.ownerUid, "Shop 2")
	require.NoError(t, err)

	item, err := Items.Add(f.c, f.ownerUid, ItemInput{Sku: "A", Name: "A", SalePrice: 100, ReorderLevel: 5000, TrackStock: true})
	require.NoError(t, err)
	f.receive(t, item.ItemId, main.LocationId, 9000)
	f.receive(t, item.ItemId, shop.LocationId, 1000)

	all, err := Reports.LowStock(f.c, f.ownerUid, 0)
	require.NoError(t, err)
	assert.Empty(t, all, "10 in total is above the level")

	atShop, err := Reports.LowStock(f.c, f.ownerUid, shop.LocationId)
	require.NoError(t, err)
	require.Len(t, atShop, 1)
	assert.Equal(t, int64(4000), atShop[0].Shortfall, "but the shop alone is running low")
}

// ---- receivables

const day = int64(24 * 60 * 60)

func (f *fixture) creditSaleAt(t *testing.T, customerId int64, itemId int64, qty int64, saleTime int64) *SaleDetail {
	t.Helper()

	in := f.saleInput(SaleLineInput{ItemId: itemId, Qty: qty})
	in.CustomerId = customerId
	in.Time = saleTime
	detail, err := Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)

	return detail
}

func TestReport_ReceivablesAgesDebtFromTheSaleDate(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "SVC", 1000, false) // 1.00 per unit
	ada, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)
	bola, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Bola"})
	require.NoError(t, err)

	now := int64(1_800_000_000)
	f.creditSaleAt(t, ada.CustomerId, item.ItemId, 1000, now-5*day)    // owes 1000, current
	f.creditSaleAt(t, ada.CustomerId, item.ItemId, 2000, now-30*day)   // owes 2000, still current (30 days)
	f.creditSaleAt(t, ada.CustomerId, item.ItemId, 3000, now-45*day)   // owes 3000, 31-60
	f.creditSaleAt(t, bola.CustomerId, item.ItemId, 4000, now-75*day)  // owes 4000, 61-90
	f.creditSaleAt(t, bola.CustomerId, item.ItemId, 5000, now-200*day) // owes 5000, over 90

	report, err := Reports.Receivables(f.c, f.ownerUid, now)
	require.NoError(t, err)
	require.Len(t, report.Rows, 2)

	assert.Equal(t, int64(15000), report.TotalOutstanding)
	assert.Equal(t, int64(3000), report.Current)
	assert.Equal(t, int64(3000), report.Days31To60)
	assert.Equal(t, int64(4000), report.Days61To90)
	assert.Equal(t, int64(5000), report.Over90)

	// largest debt first: Bola 9000, Ada 6000
	assert.Equal(t, "Bola", report.Rows[0].Customer.Name)
	assert.Equal(t, int64(9000), report.Rows[0].Outstanding)
	assert.Equal(t, int64(4000), report.Rows[0].Days61To90)
	assert.Equal(t, int64(5000), report.Rows[0].Over90)
	assert.Equal(t, now-200*day, report.Rows[0].OldestSaleTime)
	assert.Equal(t, 2, report.Rows[0].OpenSales)

	assert.Equal(t, "Ada", report.Rows[1].Customer.Name)
	assert.Equal(t, int64(6000), report.Rows[1].Outstanding)
	assert.Equal(t, int64(3000), report.Rows[1].Current)
	assert.Equal(t, int64(3000), report.Rows[1].Days31To60)
	assert.Equal(t, 3, report.Rows[1].OpenSales)
}

func TestReport_ReceivablesReflectRepaymentsAndVoids(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "SVC", 1000, false)
	ada, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)
	bola, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Bola"})
	require.NoError(t, err)

	now := nowUnix()
	f.creditSaleAt(t, ada.CustomerId, item.ItemId, 4000, now-40*day)
	voided := f.creditSaleAt(t, bola.CustomerId, item.ItemId, 7000, now-10*day)
	require.NoError(t, Sales.Void(f.c, f.ownerUid, f.ownerUid, voided.Sale.SaleId))

	_, err = Repayments.Add(f.c, f.ownerUid, f.ownerUid, f.repay(ada.CustomerId, 1500, 0))
	require.NoError(t, err)

	report, err := Reports.Receivables(f.c, f.ownerUid, now)
	require.NoError(t, err)
	require.Len(t, report.Rows, 1, "a customer whose only sale was voided owes nothing")
	assert.Equal(t, int64(2500), report.Rows[0].Outstanding)
	assert.Equal(t, int64(2500), report.Days31To60)
	assert.Equal(t, int64(2500), report.TotalOutstanding)

	// agrees with the customer balances used elsewhere
	balances, err := Customers.Balances(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Equal(t, report.TotalOutstanding, balances[ada.CustomerId])
}

func TestReport_ReceivablesIgnoreWalkInSalesAndPaidSales(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "SVC", 1000, false)

	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 1000})
	in.AmountPaid = 1000
	_, err := Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)

	report, err := Reports.Receivables(f.c, f.ownerUid, nowUnix())
	require.NoError(t, err)
	assert.Empty(t, report.Rows)
	assert.Equal(t, int64(0), report.TotalOutstanding)
}

// ---- attribution

func TestStampText(t *testing.T) {
	assert.Equal(t, "Sale #5 · by Sam", StampText("Sale #5", "Sam", 255))
	assert.Equal(t, "by Sam", StampText("", "Sam", 255), "no text: just the mark")
	assert.Equal(t, "Sale #5", StampText("Sale #5", "", 255), "no name, no mark")
	assert.Equal(t, "Sale #5", StampText("Sale #5", "   ", 255))
	assert.Equal(t, "Sale #5 · by Sam", StampText("  Sale #5  ", "  Sam ", 255))
}

func TestStampText_NeverExceedsTheLimitAndKeepsTheMark(t *testing.T) {
	long := ""
	for i := 0; i < 400; i++ {
		long += "x"
	}

	got := StampText(long, "Sam", 255)
	assert.Len(t, []rune(got), 255)
	assert.Contains(t, got, " · by Sam", "the original text is shortened, never the mark")

	multibyte := StampText("日本語日本語日本語日本語", "Zoë", 12)
	assert.LessOrEqual(t, len([]rune(multibyte)), 12)
	assert.Contains(t, multibyte, "by Zoë")
}

func TestAttribution_StaffSalesAreMarkedOnTheTransactions(t *testing.T) {
	f := newFixture(t)
	staff := f.newUser(t, "sam")
	item := f.newItem(t, "A1", 1000, false)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)

	commentOf := func(id int64) string {
		tx := &models.Transaction{}
		has, err := ownerDB(f.ownerUid).NewSession(f.c).Where("uid=? AND transaction_id=?", f.ownerUid, id).Get(tx)
		require.NoError(t, err)
		require.True(t, has)

		return tx.Comment
	}

	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 2000})
	in.CustomerId = customer.CustomerId
	in.AmountPaid = 500
	in.Note = "weekend"
	bySam, err := Sales.Create(f.c, f.ownerUid, staff, in)
	require.NoError(t, err)

	assert.Equal(t, "Sale #"+itoa(bySam.Sale.SaleId)+": weekend · by sam", commentOf(bySam.Sale.PaidTransactionId), "sam is the nickname the fixture gives the user")
	assert.Equal(t, "Sale #"+itoa(bySam.Sale.SaleId)+" (on credit): weekend · by sam", commentOf(bySam.Sale.CreditTransactionId))

	byOwner, err := Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)
	assert.Equal(t, "Sale #"+itoa(byOwner.Sale.SaleId)+": weekend", commentOf(byOwner.Sale.PaidTransactionId), "the owner's own entries carry no mark")

	// repayments too
	_, err = Repayments.Add(f.c, f.ownerUid, staff, f.repay(customer.CustomerId, 100, 0))
	require.NoError(t, err)
	repayments, err := Repayments.List(f.c, f.ownerUid, customer.CustomerId, 0, 0)
	require.NoError(t, err)
	require.Len(t, repayments, 1)
	assert.Contains(t, commentOf(repayments[0].TransactionId), "· by sam")
	assert.Equal(t, staff, repayments[0].ActorUid)
}

func TestPeople_OwnerFirstThenEveryoneWhoWorkedThere(t *testing.T) {
	f := newFixture(t)
	f.newUser(t, "manager")
	f.newUser(t, "staff")
	manager := f.newUser(t, "gone")

	for _, email := range []string{"manager@example.com", "staff@example.com", "gone@example.com"} {
		role := extmodels.RoleStaff
		if email == "manager@example.com" {
			role = extmodels.RoleManager
		}

		m, err := Memberships.Invite(f.c, f.ownerUid, email, role)
		require.NoError(t, err)
		require.NoError(t, Memberships.Respond(f.c, m.StaffUid, f.ownerUid, true))
	}

	require.NoError(t, Memberships.Remove(f.c, f.ownerUid, manager))

	people, err := People(f.c, f.ownerUid)
	require.NoError(t, err)
	require.Len(t, people, 4)

	assert.Equal(t, f.ownerUid, people[0].Uid)
	assert.Equal(t, extmodels.RoleOwner, people[0].Role)
	assert.Equal(t, "owner", people[0].Name)
	assert.True(t, people[0].Active)

	byName := map[string]*Person{}
	for _, p := range people {
		byName[p.Name] = p
	}

	assert.True(t, byName["manager"].Active)
	assert.Equal(t, extmodels.RoleManager, byName["manager"].Role)
	assert.False(t, byName["gone"].Active, "a removed person still has a name for the records they wrote")
}

func itoa(n int64) string {
	return fmtInt(n)
}
